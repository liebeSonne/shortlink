package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/liebeSonne/shortlink/internal/auth"
	"github.com/liebeSonne/shortlink/internal/handler/cookie"
	"github.com/liebeSonne/shortlink/internal/model"
	"github.com/liebeSonne/shortlink/internal/provider"
	"github.com/liebeSonne/shortlink/internal/repository"
	"github.com/liebeSonne/shortlink/internal/service"
)

var ErrNotCreated = errors.New("short link not created")

type ShortLinkHandler interface {
	HandleGet(w http.ResponseWriter, r *http.Request)
	HandleCreate(w http.ResponseWriter, r *http.Request)
	HandleCreateShorten(w http.ResponseWriter, r *http.Request)
	HandleCreateShortenBatch(w http.ResponseWriter, r *http.Request)
	HandleGetUserUrls(w http.ResponseWriter, r *http.Request)
}

func NewShortLinkHandler(
	service service.ShortLinkService,
	provider provider.ShortLinkProvider,
	urlAddress string,
) ShortLinkHandler {
	return &shortLinkHandler{
		service:    service,
		provider:   provider,
		urlAddress: urlAddress,
	}
}

type shortLinkHandler struct {
	service       service.ShortLinkService
	provider      provider.ShortLinkProvider
	urlAddress    string
	cookieService cookie.Service
	tokenService  auth.TokenService
}

func (h *shortLinkHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if id == "" {
		http.Error(w, "empty id", http.StatusBadRequest)
		return
	}

	item, err := h.provider.Find(ctx, id)
	if err != nil {
		h.responseError(w, err)
		return
	}
	if item == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	url := item.URL

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *shortLinkHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.responseError(w, err)
		return
	}
	link := string(body)

	var userIDPtr *uuid.UUID
	if userID, hasUserID := auth.GetUserIDFromContext(ctx); hasUserID {
		userIDPtr = &userID
	}

	shortLink, status, err := h.createShortLink(ctx, link, userIDPtr)
	if err != nil {
		h.responseError(w, err)
		return
	}

	url := h.createShortLinkURL(shortLink.ID)

	w.WriteHeader(status)
	_, err = w.Write([]byte(url))
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
}

func (h *shortLinkHandler) HandleCreateShorten(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var request ShortenRequest
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&request)
	if err != nil {
		h.responseError(w, err)
		return
	}

	var userIDPtr *uuid.UUID
	if userID, hasUserID := auth.GetUserIDFromContext(ctx); hasUserID {
		userIDPtr = &userID
	}

	shortLink, status, err := h.createShortLink(ctx, request.URL, userIDPtr)
	if err != nil {
		h.responseError(w, err)
		return
	}

	url := h.createShortLinkURL(shortLink.ID)

	resp := ShortenResponse{
		Result: url,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	enc := json.NewEncoder(w)
	err = enc.Encode(resp)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
}

func (h *shortLinkHandler) HandleCreateShortenBatch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var request ShortenBatchRequest
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&request)
	if err != nil {
		h.responseError(w, err)
		return
	}

	inputs := make([]service.InputShortLinkData, 0, len(request))
	for _, item := range request {
		inputs = append(inputs, service.InputShortLinkData{
			CorrelationID: item.CorrelationID,
			URL:           item.OriginalURL,
		})
	}

	var userIDPtr *uuid.UUID
	if userID, hasUserID := auth.GetUserIDFromContext(ctx); hasUserID {
		userIDPtr = &userID
	}

	outputs, err := h.service.CreateBatch(ctx, inputs, userIDPtr)
	if err != nil {
		h.responseError(w, err)
		return
	}

	resp := make(ShortenBatchResponse, 0, len(outputs))
	for _, output := range outputs {
		url := h.createShortLinkURL(output.ShortLink.ID)
		resp = append(resp, ShortenBatchResponseItem{
			CorrelationID: output.CorrelationID,
			ShortURL:      url,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	enc := json.NewEncoder(w)
	err = enc.Encode(resp)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
}

func (h *shortLinkHandler) HandleGetUserUrls(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, hasUserID := auth.GetUserIDFromContext(ctx)

	if !hasUserID {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	items, err := h.provider.FindByUserID(ctx, userID)
	if err != nil {
		h.responseError(w, err)
		return
	}

	if len(items) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	resp := make(UserUrlsResponse, 0, len(items))
	for _, item := range items {
		url := h.createShortLinkURL(item.ID)
		resp = append(resp, UserUrlsResponseItem{
			ShortURL:    url,
			OriginalURL: item.URL,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	enc := json.NewEncoder(w)
	err = enc.Encode(resp)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
}

func (h *shortLinkHandler) createShortLink(ctx context.Context, link string, userID *uuid.UUID) (*model.ShortLink, int, error) {
	var shortLink *model.ShortLink
	status := http.StatusCreated

	shortLink, err := h.service.Create(ctx, link, userID)
	if err != nil {
		var conflictErr *repository.ErrConflictURL
		if errors.As(err, &conflictErr) && conflictErr.URL == link {
			shortLink, err = h.provider.FindByURL(ctx, link)
			status = http.StatusConflict
		}
	}
	if err != nil {
		return nil, status, err
	}

	if shortLink == nil {
		return nil, http.StatusInternalServerError, ErrNotCreated
	}

	return shortLink, status, nil
}

func (h *shortLinkHandler) createShortLinkURL(id string) string {
	return fmt.Sprintf("%s/%s", h.urlAddress, id)
}

func (h *shortLinkHandler) responseError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}
	if errors.Is(err, service.ErrInvalidURL) || errors.Is(err, service.ErrEmptyURL) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if errors.Is(err, ErrNotCreated) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Error(w, err.Error(), http.StatusInternalServerError)
}
