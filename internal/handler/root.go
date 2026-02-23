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
) RootRouter {
	return &rootHandler{
		shortLinkHandler: shortLinkHandler,
	}
}

type rootHandler struct {
	shortLinkHandler ShortLinkHandler
}

func (h *rootHandler) Router() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/{id}", h.shortLinkHandler.HandleGet)
	r.Post("/", h.shortLinkHandler.HandleCreate)

	return r
}
