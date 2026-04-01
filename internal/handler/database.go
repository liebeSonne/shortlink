package handler

import (
	"net/http"
)

type DatabaseHandler interface {
	HandlePing(w http.ResponseWriter, r *http.Request)
}

func NewDatabaseHandler() DatabaseHandler {
	return &databaseHandler{}
}

type databaseHandler struct {
}

func (h *databaseHandler) HandlePing(w http.ResponseWriter, _ *http.Request) {
	// TODO
	isPing := true

	if isPing {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
