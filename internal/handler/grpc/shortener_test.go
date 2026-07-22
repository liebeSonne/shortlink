package grpc

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	pb "github.com/liebeSonne/shortlink/api/proto"
	"github.com/liebeSonne/shortlink/internal/auth"
	"github.com/liebeSonne/shortlink/internal/handler/audit"
	"github.com/liebeSonne/shortlink/internal/logger"
	"github.com/liebeSonne/shortlink/internal/model"
	"github.com/liebeSonne/shortlink/internal/provider"
	"github.com/liebeSonne/shortlink/internal/repository"
	"github.com/liebeSonne/shortlink/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestShortenerGRPCServer_ShortenURL(t *testing.T) {
	id1 := "id1"
	link1 := "https://localhost/123"
	urlAddress := "http://localhost:8080"
	userID1 := uuid.New()

	type on struct {
		link   string
		userID *uuid.UUID
	}
	type want struct {
		result   string
		grpcCode codes.Code
		hasError bool
	}
	type when struct {
		createItem *model.ShortLink
		createErr  error
		findItem   *model.ShortLink
		findErr    error
	}
	testCases := []struct {
		name string
		on   on
		want want
		when when
	}{
		{
			"success",
			on{link1, nil},
			want{urlAddress + "/" + id1, codes.OK, false},
			when{&model.ShortLink{ID: id1, URL: link1}, nil, nil, nil},
		},
		{
			"success with user",
			on{link1, &userID1},
			want{urlAddress + "/" + id1, codes.OK, false},
			when{&model.ShortLink{ID: id1, URL: link1}, nil, nil, nil},
		},
		{
			"service error",
			on{link1, nil},
			want{"", codes.Internal, true},
			when{nil, errors.New("some service error"), nil, nil},
		},
		{
			"invalid url error",
			on{link1, nil},
			want{"", codes.InvalidArgument, true},
			when{nil, service.ErrInvalidURL, nil, nil},
		},
		{
			"empty url error",
			on{link1, nil},
			want{"", codes.InvalidArgument, true},
			when{nil, service.ErrEmptyURL, nil, nil},
		},
		{
			"conflict found existing",
			on{link1, nil},
			want{urlAddress + "/" + id1, codes.OK, false},
			when{nil, repository.NewErrConflictURL(link1, errors.New("conflict")), &model.ShortLink{ID: id1, URL: link1}, nil},
		},
		{
			"conflict not found",
			on{link1, nil},
			want{"", codes.Internal, true},
			when{nil, repository.NewErrConflictURL(link1, errors.New("conflict")), nil, nil},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := service.NewMockShortLinkService(t)
			s.EXPECT().Create(mock.Anything, mock.Anything, mock.Anything).Return(tc.when.createItem, tc.when.createErr)

			p := provider.NewMockShortLinkProvider(t)
			var isConflict bool
			if tc.when.createErr != nil {
				var conflictErr *repository.ErrConflictURL
				isConflict = errors.As(tc.when.createErr, &conflictErr) && conflictErr.URL == tc.on.link
			}
			if isConflict {
				p.EXPECT().FindByURL(mock.Anything, mock.Anything).Return(tc.when.findItem, tc.when.findErr)
			}

			l := logger.NewMockLogger(t)
			l.EXPECT().Errorf(mock.Anything, mock.Anything).Maybe()

			mockAuditPublisher := audit.NewMockPublisher(t)
			mockAuditPublisher.EXPECT().Notify(mock.Anything).Maybe()

			server := NewShortenerGRPCServer(s, p, urlAddress, l, mockAuditPublisher)

			ctx := context.Background()
			if tc.on.userID != nil {
				ctx = auth.CreateTokenContext(ctx, auth.Token{UserID: tc.on.userID.String()})
			}

			req := &pb.URLShortenRequest{}
			req.SetUrl(tc.on.link)

			resp, err := server.ShortenURL(ctx, req)

			if tc.want.hasError {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				assert.Equal(t, tc.want.grpcCode, st.Code())
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.want.result, resp.GetResult())
			}
		})
	}
}

func TestShortenerGRPCServer_ExpandURL(t *testing.T) {
	id1 := "id1"
	link1 := "https://localhost/123"
	userID1 := uuid.New()

	type on struct {
		id     string
		userID *uuid.UUID
	}
	type want struct {
		result   string
		grpcCode codes.Code
		hasError bool
	}
	type when struct {
		item *model.ShortLink
		err  error
	}
	testCases := []struct {
		name string
		on   on
		want want
		when when
	}{
		{
			"empty id",
			on{"", nil},
			want{"", codes.InvalidArgument, true},
			when{nil, nil},
		},
		{
			"success",
			on{id1, nil},
			want{link1, codes.OK, false},
			when{&model.ShortLink{ID: id1, URL: link1}, nil},
		},
		{
			"success with user",
			on{id1, &userID1},
			want{link1, codes.OK, false},
			when{&model.ShortLink{ID: id1, URL: link1}, nil},
		},
		{
			"provider error",
			on{id1, nil},
			want{"", codes.Internal, true},
			when{nil, errors.New("some provider error")},
		},
		{
			"not found",
			on{id1, nil},
			want{"", codes.NotFound, true},
			when{nil, nil},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := service.NewMockShortLinkService(t)

			p := provider.NewMockShortLinkProvider(t)
			if tc.on.id != "" {
				p.EXPECT().Find(mock.Anything, tc.on.id).Return(tc.when.item, tc.when.err)
			}

			l := logger.NewMockLogger(t)
			l.EXPECT().Errorf(mock.Anything, mock.Anything).Maybe()

			mockAuditPublisher := audit.NewMockPublisher(t)
			mockAuditPublisher.EXPECT().Notify(mock.Anything).Run(func(event audit.Event) {
				if !tc.want.hasError {
					assert.Equal(t, tc.want.result, event.URL)
					assert.Equal(t, audit.ActionFollow, event.Action)
				}
			}).Maybe()

			urlAddress := "http://localhost:8080"
			server := NewShortenerGRPCServer(s, p, urlAddress, l, mockAuditPublisher)

			ctx := context.Background()
			if tc.on.userID != nil {
				ctx = auth.CreateTokenContext(ctx, auth.Token{UserID: tc.on.userID.String()})
			}

			req := &pb.URLExpandRequest{}
			req.SetId(tc.on.id)

			resp, err := server.ExpandURL(ctx, req)

			if tc.want.hasError {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				assert.Equal(t, tc.want.grpcCode, st.Code())
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.want.result, resp.GetResult())
			}
		})
	}
}

func TestShortenerGRPCServer_ListUserURLs(t *testing.T) {
	urlAddress := "http://localhost:8080"
	userID1 := uuid.New()

	type on struct {
		userID *uuid.UUID
	}
	type want struct {
		grpcCode     codes.Code
		hasError     bool
		urlCount     int
		shortURLs    []string
		originalURLS []string
	}
	type when struct {
		items []model.ShortLink
		err   error
	}
	testCases := []struct {
		name string
		on   on
		want want
		when when
	}{
		{
			"unauthenticated",
			on{nil},
			want{codes.Unauthenticated, true, 0, nil, nil},
			when{nil, nil},
		},
		{
			"success",
			on{&userID1},
			want{codes.OK, false, 2, []string{urlAddress + "/id1", urlAddress + "/id2"}, []string{"https://example1.com", "https://example2.com"}},
			when{[]model.ShortLink{{ID: "id1", URL: "https://example1.com"}, {ID: "id2", URL: "https://example2.com"}}, nil},
		},
		{
			"provider error",
			on{&userID1},
			want{codes.Internal, true, 0, nil, nil},
			when{nil, errors.New("some provider error")},
		},
		{
			"no content",
			on{&userID1},
			want{codes.NotFound, true, 0, nil, nil},
			when{nil, nil},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := service.NewMockShortLinkService(t)

			p := provider.NewMockShortLinkProvider(t)
			if tc.on.userID != nil {
				p.EXPECT().FindByUserID(mock.Anything, *(tc.on.userID)).Return(tc.when.items, tc.when.err)
			}

			l := logger.NewMockLogger(t)
			l.EXPECT().Errorf(mock.Anything, mock.Anything).Maybe()

			mockAuditPublisher := audit.NewMockPublisher(t)

			server := NewShortenerGRPCServer(s, p, urlAddress, l, mockAuditPublisher)

			ctx := context.Background()
			if tc.on.userID != nil {
				ctx = auth.CreateTokenContext(ctx, auth.Token{UserID: tc.on.userID.String()})
			}

			resp, err := server.ListUserURLs(ctx, &pb.UserURLsRequest{})

			if tc.want.hasError {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				assert.Equal(t, tc.want.grpcCode, st.Code())
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.want.urlCount, len(resp.GetUrl()))
				if len(tc.want.shortURLs) > 0 {
					gotShortURLs := make([]string, 0, len(resp.GetUrl()))
					for _, item := range resp.GetUrl() {
						gotShortURLs = append(gotShortURLs, item.GetShortUrl())
					}
					assert.Equal(t, tc.want.shortURLs, gotShortURLs)
				}
				if len(tc.want.originalURLS) > 0 {
					gotURLs := make([]string, 0, len(resp.GetUrl()))
					for _, item := range resp.GetUrl() {
						gotURLs = append(gotURLs, item.GetOriginalUrl())
					}
					assert.Equal(t, tc.want.originalURLS, gotURLs)
				}
			}
		})
	}
}
