package handler

import (
	"fmt"
	"io"
	"net/http"

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

func (h *shortLinkHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[1:]

	err := validateShortLinkID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	itemPtr, err := h.provider.Get(model.ShortLinkID(id))
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

func (h *shortLinkHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
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
