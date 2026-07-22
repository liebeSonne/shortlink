package grpc

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/liebeSonne/shortlink/api/proto"
	"github.com/liebeSonne/shortlink/internal/auth"
	"github.com/liebeSonne/shortlink/internal/handler/audit"
	"github.com/liebeSonne/shortlink/internal/logger"
	"github.com/liebeSonne/shortlink/internal/model"
	"github.com/liebeSonne/shortlink/internal/provider"
	"github.com/liebeSonne/shortlink/internal/repository"
	"github.com/liebeSonne/shortlink/internal/service"
)

type ShortenerGRPCServer struct {
	pb.UnimplementedShortenerServiceServer

	service        service.ShortLinkService
	provider       provider.ShortLinkProvider
	urlAddress     string
	logger         logger.Logger
	auditPublisher audit.Publisher
}

func NewShortenerGRPCServer(
	service service.ShortLinkService,
	provider provider.ShortLinkProvider,
	urlAddress string,
	logger logger.Logger,
	auditPublisher audit.Publisher,
) *ShortenerGRPCServer {
	return &ShortenerGRPCServer{
		service:        service,
		provider:       provider,
		urlAddress:     urlAddress,
		logger:         logger,
		auditPublisher: auditPublisher,
	}
}

func (s *ShortenerGRPCServer) ShortenURL(ctx context.Context, r *pb.URLShortenRequest) (*pb.URLShortenResponse, error) {
	link := r.GetUrl()

	var userIDPtr *uuid.UUID
	if userID, ok := auth.GetUserIDFromContext(ctx); ok {
		userIDPtr = &userID
	}

	var shortLink *model.ShortLink
	shortLink, err := s.service.Create(ctx, link, userIDPtr)
	if err != nil {
		var conflictErr *repository.ErrConflictURL
		if errors.As(err, &conflictErr) && conflictErr.URL == link {
			shortLink, err = s.provider.FindByURL(ctx, link)
			if err != nil {
				s.logger.Errorf("find by url error: %v", err)
				return nil, status.Errorf(codes.Internal, "conflict")
			}
		}
	}
	if err != nil {
		return nil, s.translateError(err)
	}
	if shortLink == nil {
		return nil, status.Errorf(codes.Internal, "not created")
	}

	url := s.createShortLinkURL(shortLink.ID)

	s.sendAuditEvent(audit.ActionShorted, userIDPtr, shortLink.URL)

	var response pb.URLShortenResponse
	response.SetResult(url)

	return &response, nil
}

func (s *ShortenerGRPCServer) ExpandURL(ctx context.Context, r *pb.URLExpandRequest) (*pb.URLExpandResponse, error) {
	id := r.GetId()

	if id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "empty id")
	}

	var userIDPtr *uuid.UUID
	if userID, ok := auth.GetUserIDFromContext(ctx); ok {
		userIDPtr = &userID
	}

	item, err := s.provider.Find(ctx, id)
	if err != nil {
		return nil, s.translateError(err)
	}
	if item == nil {
		return nil, status.Error(codes.NotFound, codes.NotFound.String())
	}

	url := item.URL

	s.sendAuditEvent(audit.ActionFollow, userIDPtr, url)

	var response pb.URLExpandResponse
	response.SetResult(url)

	return &response, nil
}

func (s *ShortenerGRPCServer) ListUserURLs(ctx context.Context, _ *pb.UserURLsRequest) (*pb.UserURLsResponse, error) {
	userID, ok := auth.GetUserIDFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated")
	}

	items, err := s.provider.FindByUserID(ctx, userID)
	if err != nil {
		return nil, s.translateError(err)
	}

	if len(items) == 0 {
		return nil, status.Errorf(codes.NotFound, "no content")
	}

	var response pb.UserURLsResponse
	responseItems := make([]*pb.URLData, 0, len(items))

	for _, item := range items {
		url := s.createShortLinkURL(item.ID)
		var responseItem pb.URLData
		responseItem.SetOriginalUrl(item.URL)
		responseItem.SetShortUrl(url)
		responseItems = append(responseItems, &responseItem)
	}
	response.SetUrl(responseItems)

	return &response, nil
}

func (s *ShortenerGRPCServer) sendAuditEvent(action audit.Action, userIDPtr *uuid.UUID, link string) {
	var userIDStrPtr *string
	if userIDPtr != nil {
		userIDStr := userIDPtr.String()
		userIDStrPtr = &userIDStr
	}
	s.auditPublisher.Notify(audit.Event{
		Time:   time.Now(),
		Action: action,
		UserID: userIDStrPtr,
		URL:    link,
	})
}

func (s *ShortenerGRPCServer) translateError(err error) error {
	if err == nil {
		return nil
	}

	s.logger.Errorf("response error: %v", err)

	if errors.Is(err, service.ErrInvalidURL) || errors.Is(err, service.ErrEmptyURL) {
		return status.Error(codes.InvalidArgument, codes.InvalidArgument.String())
	}

	return status.Error(codes.Internal, codes.Internal.String())
}

func (s *ShortenerGRPCServer) createShortLinkURL(id string) string {
	return fmt.Sprintf("%s/%s", s.urlAddress, id)
}
