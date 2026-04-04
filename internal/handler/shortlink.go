package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/liebeSonne/shortlink/internal/provider"
	"github.com/liebeSonne/shortlink/internal/service"
)

type ShortLinkHandler interface {
	HandleGet(w http.ResponseWriter, r *http.Request)
	HandleCreate(w http.ResponseWriter, r *http.Request)
	HandleCreateShorten(w http.ResponseWriter, r *http.Request)
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
	service    service.ShortLinkService
	provider   provider.ShortLinkProvider
	urlAddress string
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

	shortLink, err := h.service.Create(ctx, link)
	if err != nil {
		h.responseError(w, err)
		return
	}
	if shortLink == nil {
		http.Error(w, "not created", http.StatusInternalServerError)
		return
	}

	url := h.createShortLinkURL(shortLink.ID)

	w.WriteHeader(http.StatusCreated)
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

	shortLink, err := h.service.Create(ctx, request.URL)
	if err != nil {
		h.responseError(w, err)
		return
	}
	if shortLink == nil {
		http.Error(w, "not created", http.StatusInternalServerError)
		return
	}

	url := h.createShortLinkURL(shortLink.ID)

	resp := ShortenResponse{
		Result: url,
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
	http.Error(w, err.Error(), http.StatusInternalServerError)
}
