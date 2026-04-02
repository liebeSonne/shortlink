package handler

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDatabaseHandler_HandlePing(t *testing.T) {
	type want struct {
		code int
	}
	type when struct {
		err error
	}
	testCases := []struct {
		name string
		when when
		want want
	}{
		{
			"200",
			when{nil},
			want{http.StatusOK},
		},
		{
			"500",
			when{errors.New("som ping database error")},
			want{http.StatusInternalServerError},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db := new(mockDatabase)
			db.On("Ping", mock.Anything).Return(tc.when.err)

			logger := new(mockLogger)
			logger.On("Debugf", mock.Anything, mock.Anything).Return()

			handler := NewDatabaseHandler(db, logger)

			r := chi.NewRouter()
			r.Get("/ping", handler.HandlePing)

			srv := httptest.NewServer(r)
			defer srv.Close()

			client := resty.New()
			client.SetRedirectPolicy(resty.NoRedirectPolicy())

			resp, err := client.R().
				Get(srv.URL + "/ping")

			require.Equal(t, tc.want.code, resp.StatusCode(), fmt.Sprintf("expected status code %d but got %d with body: %s", tc.want.code, resp.StatusCode(), string(resp.Body())))
			require.NoError(t, err)
		})
	}
}
