package handler

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/liebeSonne/shortlink/internal/model"
	"github.com/liebeSonne/shortlink/internal/service"
)

func TestShortLinkHandler_HandleGet(t *testing.T) {
	link1 := "https://localhost/123"

	type on struct {
		id string
	}
	type want struct {
		code     int
		location *string
	}
	type when struct {
		err  error
		link *string
	}
	testCases := []struct {
		name string
		on   on
		want want
		when when
	}{
		{
			"empty id",
			on{""},
			want{http.StatusBadRequest, nil},
			when{nil, nil},
		},
		{
			"redirect",
			on{"id1"},
			want{http.StatusTemporaryRedirect, &link1},
			when{nil, &link1},
		},
		{
			"provider error",
			on{"id1"},
			want{http.StatusInternalServerError, &link1},
			when{errors.New("some provider error"), &link1},
		},
		{
			"not found",
			on{"id1"},
			want{http.StatusNotFound, nil},
			when{nil, nil},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			provider := new(mockProvider)
			if tc.when.err == nil && tc.when.link != nil {
				item := &model.ShortLink{ID: tc.on.id, URL: *tc.when.link}
				provider.On("Find", mock.Anything, tc.on.id).Return(item, tc.when.err)
			} else {
				provider.On("Find", mock.Anything, tc.on.id).Return(nil, tc.when.err)
			}

			urlAddress := "http://localhost:8080"
			handler := NewShortLinkHandler(new(mockService), provider, urlAddress)

			r := chi.NewRouter()
			r.Get("/", handler.HandleGet)
			r.Get("/{id}", handler.HandleGet)

			srv := httptest.NewServer(r)
			defer srv.Close()

			client := resty.New()
			client.SetRedirectPolicy(resty.NoRedirectPolicy())

			resp, err := client.R().
				Get(srv.URL + "/" + tc.on.id)

			require.Equal(t, tc.want.code, resp.StatusCode(), fmt.Sprintf("expected status code %d but got %d with body: %s", tc.want.code, resp.StatusCode(), string(resp.Body())))

			if tc.want.code == http.StatusTemporaryRedirect && tc.want.location != nil {
				require.Error(t, err)
				require.Contains(t, err.Error(), "auto redirect is disabled")

				locationURL := resp.Header().Get("Location")
				assert.Equal(t, *tc.want.location, locationURL)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestShortLinkHandler_HandleCreate(t *testing.T) {
	id1 := "id1"
	link1 := "https://localhost/123"
	urlAddress := "http://localhost:8080"

	type on struct {
		link string
	}
	type when struct {
		id  string
		err error
	}
	type want struct {
		code     int
		linkPath string
	}
	testCases := []struct {
		name string
		on   on
		want want
		when when
	}{
		{
			"success create",
			on{link1},
			want{http.StatusCreated, "/" + id1},
			when{id1, nil},
		},
		{
			"error in service",
			on{link1},
			want{http.StatusInternalServerError, ""},
			when{"", errors.New("some service error")},
		},
		{
			"error empty url",
			on{link1},
			want{http.StatusBadRequest, ""},
			when{"", service.ErrEmptyURL},
		},
		{
			"error invalid url",
			on{link1},
			want{http.StatusBadRequest, ""},
			when{"", service.ErrInvalidURL},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var item *model.ShortLink
			if tc.when.err == nil {
				item = &model.ShortLink{ID: tc.when.id, URL: link1}
			}
			s := new(mockService)
			s.On("Create", mock.Anything, tc.on.link).Return(item, tc.when.err)

			handler := NewShortLinkHandler(s, new(mockProvider), urlAddress)

			r := chi.NewRouter()
			r.Post("/", handler.HandleCreate)

			srv := httptest.NewServer(r)
			defer srv.Close()

			client := resty.New()
			client.SetRedirectPolicy(resty.NoRedirectPolicy())

			resp, err := client.R().
				SetBody(tc.on.link).
				Post(srv.URL + "/")

			require.NoError(t, err)

			require.Equal(t, tc.want.code, resp.StatusCode(), fmt.Sprintf("expected status code %d but got %d with body: %s", tc.want.code, resp.StatusCode(), string(resp.Body())))
			if tc.want.linkPath != "" {
				wantURL := urlAddress + tc.want.linkPath
				assert.Equal(t, wantURL, string(resp.Body()))
			}
		})
	}
}

func TestShortLinkHandler_HandleCreateShorten(t *testing.T) {
	id1 := "id1"
	link1 := "https://localhost/123"
	urlAddress := "http://localhost:8080"

	type on struct {
		body string
	}
	type when struct {
		id  string
		err error
	}
	type want struct {
		code int
		body string
	}
	testCases := []struct {
		name string
		on   on
		want want
		when when
	}{
		{
			"success create",
			on{fmt.Sprintf(`{"url": "%s"}`, link1)},
			want{http.StatusCreated, fmt.Sprintf(`{"result": "%s/%s"}`, urlAddress, id1)},
			when{id1, nil},
		},
		{
			"error in service",
			on{fmt.Sprintf(`{"url": "%s"}`, link1)},
			want{http.StatusInternalServerError, ""},
			when{"", errors.New("some service error")},
		},
		{
			"error empty url",
			on{fmt.Sprintf(`{"url": "%s"}`, link1)},
			want{http.StatusBadRequest, ""},
			when{"", service.ErrEmptyURL},
		},
		{
			"error invalid url",
			on{fmt.Sprintf(`{"url": "%s"}`, link1)},
			want{http.StatusBadRequest, ""},
			when{"", service.ErrInvalidURL},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var item *model.ShortLink
			if tc.when.err == nil {
				item = &model.ShortLink{ID: tc.when.id, URL: link1}
			}
			s := new(mockService)
			s.On("Create", mock.Anything, mock.Anything).Return(item, tc.when.err)

			handler := NewShortLinkHandler(s, new(mockProvider), urlAddress)

			r := chi.NewRouter()
			r.Post("/api/shorten", handler.HandleCreateShorten)

			srv := httptest.NewServer(r)
			defer srv.Close()

			client := resty.New()
			client.SetRedirectPolicy(resty.NoRedirectPolicy())

			req := client.R()
			req.Method = resty.MethodPost
			req.URL = srv.URL + "/api/shorten"

			if len(tc.on.body) > 0 {
				req.SetHeader("Content-Type", "application/json")
				req.SetBody(tc.on.body)
			}

			resp, err := req.Send()
			require.NoError(t, err)

			require.Equal(t, tc.want.code, resp.StatusCode(), fmt.Sprintf("expected status code %d but got %d with body: %s", tc.want.code, resp.StatusCode(), string(resp.Body())))

			if tc.want.body != "" {
				assert.JSONEq(t, tc.want.body, string(resp.Body()))
				assert.Equal(t, "application/json", resp.Header().Get("Content-Type"))
			}
		})
	}
}
