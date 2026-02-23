package handler

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/liebeSonne/shortlink/internal/model"
	"github.com/liebeSonne/shortlink/internal/service"
)

type ShortLinkHandler interface {
	HandleGet(w http.ResponseWriter, r *http.Request)
	HandleCreate(w http.ResponseWriter, r *http.Request)
}

func NewShortLinkHandler(
	service service.ShortLinkService,
	provider model.ShortLinkProvider,
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
	provider   model.ShortLinkProvider
	urlAddress string
}

func (h *shortLinkHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if id == "" {
		http.Error(w, "empty id", http.StatusBadRequest)
		return
	}

	item, err := h.provider.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if item == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	url := item.URL()

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *shortLinkHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	link := string(body)

	shortLink, err := h.service.Create(link)
	if err != nil {
		if errors.Is(err, model.ErrInvalidURL) || errors.Is(err, model.ErrEmptyURL) || errors.Is(err, model.ErrEmptyID) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if shortLink == nil {
		http.Error(w, "not created", http.StatusInternalServerError)
		return
	}

	url := h.createShortLinkURL(shortLink.ID())

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(url))
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
}

func (h *shortLinkHandler) createShortLinkURL(id string) string {
	return fmt.Sprintf("%s/%s", h.urlAddress, id)
}
