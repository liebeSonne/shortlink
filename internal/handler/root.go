package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// RootRouter - интерфейс корневого роутера.
type RootRouter interface {
	// Router - возвращает роутер.
	Router() chi.Router
}

// NewRootRouter - создание экземпляра корневого роутера.
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

	r.Mount("/debug", middleware.Profiler())

	r.Get("/ping", h.databaseHandler.HandlePing)
	r.Get("/{id}", h.shortLinkHandler.HandleGet)
	r.Post("/", h.shortLinkHandler.HandleCreate)
	r.Post("/api/shorten", h.shortLinkHandler.HandleCreateShorten)
	r.Post("/api/shorten/batch", h.shortLinkHandler.HandleCreateShortenBatch)
	r.Get("/api/user/urls", h.shortLinkHandler.HandleGetUserUrls)
	r.Delete("/api/user/urls", h.shortLinkHandler.HandleDeleteUrls)

	return r
}
