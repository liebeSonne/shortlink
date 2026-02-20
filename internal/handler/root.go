package handler

import "net/http"

type RootHandler interface {
	Handle(w http.ResponseWriter, r *http.Request)
}

func NewRootHandler(
	shortLinkHandler ShortLinkHandler,
) RootHandler {
	return &rootHandler{
		shortLinkHandler: shortLinkHandler,
	}
}

type rootHandler struct {
	shortLinkHandler ShortLinkHandler
}

func (h *rootHandler) Handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.shortLinkHandler.HandleGet(w, r)
	case http.MethodPost:
		h.shortLinkHandler.HandleCreate(w, r)
	default:
		w.WriteHeader(http.StatusNotAcceptable)
	}
}
