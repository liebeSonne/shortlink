package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type RootRouter interface {
	Router() chi.Router
}

func NewRootRouter(
	shortLinkHandler ShortLinkHandler,
	enableLogs bool,
) RootRouter {
	return &rootHandler{
		shortLinkHandler: shortLinkHandler,
		enableLogs:       enableLogs,
	}
}

type rootHandler struct {
	shortLinkHandler ShortLinkHandler
	enableLogs       bool
}

func (h *rootHandler) Router() chi.Router {
	r := chi.NewRouter()

	if h.enableLogs {
		r.Use(middleware.Logger)
	}

	r.Get("/{id}", h.shortLinkHandler.HandleGet)
	r.Post("/", h.shortLinkHandler.HandleCreate)
	r.Post("/api/shorten", h.shortLinkHandler.HandleCreateShorten)

	return r
}
