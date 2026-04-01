package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/require"
)

func TestDatabaseHandler_HandlePing(t *testing.T) {
	type want struct {
		code int
	}
	testCases := []struct {
		name string
		want want
	}{
		{
			"200",
			want{http.StatusOK},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handler := NewDatabaseHandler()

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
