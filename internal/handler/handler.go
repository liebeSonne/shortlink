package handler

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/liebeSonne/shortlink/internal/model"
	"github.com/liebeSonne/shortlink/internal/service"
)

type ShortLinkHandler interface {
	Handle(w http.ResponseWriter, r *http.Request)
}

func NewShortLinkHandler(
	service service.ShortLinkService,
	provider model.ShortLinkProvider,
) ShortLinkHandler {
	return &shortLinkHandler{
		service:  service,
		provider: provider,
	}
}

type shortLinkHandler struct {
	service  service.ShortLinkService
	provider model.ShortLinkProvider
}

func (h *shortLinkHandler) Handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleGet(w, r)
	case http.MethodPost:
		h.handlePost(w, r)
	default:
		w.WriteHeader(http.StatusNotAcceptable)
	}
}

func (h *shortLinkHandler) handleGet(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path
	id, _ = strings.CutPrefix(id, "/")

	err := validateShortLinkID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	itemPtr, err := h.provider.Find(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if itemPtr == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	url := (*itemPtr).URL()

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *shortLinkHandler) handlePost(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	link := string(body)

	err = validateLink(link)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shortLink, err := h.service.Create(link)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	url := fmt.Sprintf("%s://%s/%s", r.URL.Scheme, r.Host, shortLink.ID())

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(url))
}
