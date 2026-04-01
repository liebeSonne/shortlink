package handler

import (
	"net/http"

	"github.com/liebeSonne/shortlink/internal/logger"
	"github.com/liebeSonne/shortlink/internal/repository/database"
)

type DatabaseHandler interface {
	HandlePing(w http.ResponseWriter, r *http.Request)
}

func NewDatabaseHandler(
	database database.Database,
	logger logger.Logger,
) DatabaseHandler {
	return &databaseHandler{
		database: database,
		logger:   logger,
	}
}

type databaseHandler struct {
	database database.Database
	logger   logger.Logger
}

func (h *databaseHandler) HandlePing(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	err := h.database.Ping(ctx)
	if err != nil {
		h.logger.Debugf("ping database error: %w", err)
	}
	isPing := err == nil

	if isPing {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
