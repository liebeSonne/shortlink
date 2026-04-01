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
	databaseHandler DatabaseHandler,
	enableLogs bool,
) RootRouter {
	return &rootHandler{
		shortLinkHandler: shortLinkHandler,
		databaseHandler:  databaseHandler,
		enableLogs:       enableLogs,
	}
}

type rootHandler struct {
	shortLinkHandler ShortLinkHandler
	databaseHandler  DatabaseHandler
	enableLogs       bool
}

func (h *rootHandler) Router() chi.Router {
	r := chi.NewRouter()

	if h.enableLogs {
		r.Use(middleware.Logger)
	}

	r.Get("/ping", h.databaseHandler.HandlePing)
	r.Get("/{id}", h.shortLinkHandler.HandleGet)
	r.Post("/", h.shortLinkHandler.HandleCreate)
	r.Post("/api/shorten", h.shortLinkHandler.HandleCreateShorten)

	return r
}
