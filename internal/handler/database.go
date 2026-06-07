package handler

import (
	"net/http"

	"github.com/liebeSonne/shortlink/internal/logger"
	"github.com/liebeSonne/shortlink/internal/repository/database"
)

// DatabaseHandler - интерфейс обработчика событий к базе данных.
type DatabaseHandler interface {
	// HandlePing - обработчик запроса проверки соединения с базой данных.
	HandlePing(w http.ResponseWriter, r *http.Request)
}

// NewDatabaseHandler - создание экземпляра обработчика запросов к базе данных.
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

// HandlePing    godoc
// @Summary      Проверка соединения с базой данных
// @Tags         database
// @Success      200
// @Failure      500
// @Router       /ping [get]
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
