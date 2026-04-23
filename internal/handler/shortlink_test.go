package handler

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/liebeSonne/shortlink/internal/auth"
	"github.com/liebeSonne/shortlink/internal/logger"
	"github.com/liebeSonne/shortlink/internal/model"
	"github.com/liebeSonne/shortlink/internal/provider"
	"github.com/liebeSonne/shortlink/internal/repository"
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
			want{http.StatusGone, nil},
			when{nil, nil},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := provider.NewMockShortLinkProvider(t)
			if tc.when.err == nil && tc.when.link != nil {
				item := &model.ShortLink{ID: tc.on.id, URL: *tc.when.link}
				p.EXPECT().Find(mock.Anything, tc.on.id).Return(item, tc.when.err)
			} else {
				p.EXPECT().Find(mock.Anything, tc.on.id).Return(nil, tc.when.err).Maybe()
			}

			d := service.NewMockShortLinkDeleter(t)

			l := logger.NewMockLogger(t)
			l.EXPECT().Errorf(mock.Anything, mock.Anything).Maybe()

			s := service.NewMockShortLinkService(t)

			urlAddress := "http://localhost:8080"
			handler := NewShortLinkHandler(s, p, urlAddress, d, l)

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
		createItem *model.ShortLink
		createErr  error
		findItem   *model.ShortLink
		findErr    error
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
			when{&model.ShortLink{ID: id1, URL: link1}, nil, nil, nil},
		},
		{
			"error in service",
			on{link1},
			want{http.StatusInternalServerError, ""},
			when{nil, errors.New("some service error"), nil, nil},
		},
		{
			"error empty url",
			on{link1},
			want{http.StatusBadRequest, ""},
			when{nil, service.ErrEmptyURL, nil, nil},
		},
		{
			"error invalid url",
			on{link1},
			want{http.StatusBadRequest, ""},
			when{nil, service.ErrInvalidURL, nil, nil},
		},
		{
			"conflict unique url",
			on{link1},
			want{http.StatusConflict, "/" + id1},
			when{nil, repository.NewErrConflictURL(link1, errors.New("conflict error")), &model.ShortLink{ID: id1, URL: link1}, nil},
		},
		{
			"conflict unique url and not found",
			on{link1},
			want{http.StatusInternalServerError, ""},
			when{nil, repository.NewErrConflictURL(link1, errors.New("conflict error")), nil, nil},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := service.NewMockShortLinkService(t)
			s.EXPECT().Create(mock.Anything, mock.Anything, mock.Anything).Return(tc.when.createItem, tc.when.createErr)

			p := provider.NewMockShortLinkProvider(t)
			p.EXPECT().FindByURL(mock.Anything, mock.Anything).Return(tc.when.findItem, tc.when.findErr).Maybe()

			d := service.NewMockShortLinkDeleter(t)

			l := logger.NewMockLogger(t)
			l.EXPECT().Errorf(mock.Anything, mock.Anything).Maybe()

			handler := NewShortLinkHandler(s, p, urlAddress, d, l)

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
		createItem *model.ShortLink
		createErr  error
		findItem   *model.ShortLink
		findErr    error
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
			when{&model.ShortLink{ID: id1, URL: link1}, nil, nil, nil},
		},
		{
			"error in service",
			on{fmt.Sprintf(`{"url": "%s"}`, link1)},
			want{http.StatusInternalServerError, ""},
			when{nil, errors.New("some service error"), nil, nil},
		},
		{
			"error empty url",
			on{fmt.Sprintf(`{"url": "%s"}`, link1)},
			want{http.StatusBadRequest, ""},
			when{nil, service.ErrEmptyURL, nil, nil},
		},
		{
			"error invalid url",
			on{fmt.Sprintf(`{"url": "%s"}`, link1)},
			want{http.StatusBadRequest, ""},
			when{nil, service.ErrInvalidURL, nil, nil},
		},
		{
			"conflict unique url",
			on{fmt.Sprintf(`{"url": "%s"}`, link1)},
			want{http.StatusConflict, fmt.Sprintf(`{"result": "%s/%s"}`, urlAddress, id1)},
			when{nil, repository.NewErrConflictURL(link1, errors.New("conflict error")), &model.ShortLink{ID: id1, URL: link1}, nil},
		},
		{
			"conflict unique url and not found",
			on{fmt.Sprintf(`{"url": "%s"}`, link1)},
			want{http.StatusInternalServerError, ""},
			when{nil, repository.NewErrConflictURL(link1, errors.New("conflict error")), nil, nil},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := service.NewMockShortLinkService(t)
			s.EXPECT().Create(mock.Anything, mock.Anything, mock.Anything).Return(tc.when.createItem, tc.when.createErr)

			p := provider.NewMockShortLinkProvider(t)
			p.EXPECT().FindByURL(mock.Anything, mock.Anything).Return(tc.when.findItem, tc.when.findErr).Maybe()

			d := service.NewMockShortLinkDeleter(t)

			l := logger.NewMockLogger(t)
			l.EXPECT().Errorf(mock.Anything, mock.Anything).Maybe()

			handler := NewShortLinkHandler(s, p, urlAddress, d, l)

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

func TestShortLinkHandler_HandleCreateShortenBatch(t *testing.T) {
	correlationID1 := "correlationID1"
	correlationID2 := "correlationID2"
	id1 := "id1"
	id2 := "id2"
	link1 := "https://localhost/111"
	link2 := "https://localhost/222"
	urlAddress := "http://localhost:8080"

	type on struct {
		body string
	}
	type when struct {
		outputs []service.OutputShortLinkData
		err     error
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
			"success create one link",
			on{fmt.Sprintf(`[{"correlation_id": "%s", "original_url": "%s"}]`, correlationID1, link1)},
			want{http.StatusCreated, fmt.Sprintf(`[{"correlation_id": "%s", "short_url": "%s/%s"}]`, correlationID1, urlAddress, id1)},
			when{[]service.OutputShortLinkData{{CorrelationID: correlationID1, ShortLink: model.ShortLink{ID: id1, URL: link1}}}, nil},
		},
		{
			"success create many link",
			on{fmt.Sprintf(`[{"correlation_id": "%s", "original_url": "%s"},{"correlation_id": "%s", "original_url": "%s"}]`, correlationID1, link1, correlationID2, link2)},
			want{http.StatusCreated, fmt.Sprintf(`[{"correlation_id": "%s", "short_url": "%s/%s"}, {"correlation_id": "%s", "short_url": "%s/%s"}]`, correlationID1, urlAddress, id1, correlationID2, urlAddress, id2)},
			when{[]service.OutputShortLinkData{
				{CorrelationID: correlationID1, ShortLink: model.ShortLink{ID: id1, URL: link1}},
				{CorrelationID: correlationID2, ShortLink: model.ShortLink{ID: id2, URL: link2}},
			}, nil},
		},
		{
			"error in service",
			on{fmt.Sprintf(`[{"correlation_id": "%s", "original_url": "%s"}]`, correlationID1, link1)},
			want{http.StatusInternalServerError, ""},
			when{nil, errors.New("some service error")},
		},
		{
			"error empty url",
			on{fmt.Sprintf(`[{"correlation_id": "%s", "original_url": "%s"}]`, correlationID1, link1)},
			want{http.StatusBadRequest, ""},
			when{nil, service.ErrEmptyURL},
		},
		{
			"error invalid url",
			on{fmt.Sprintf(`[{"correlation_id": "%s", "original_url": "%s"}]`, correlationID1, link1)},
			want{http.StatusBadRequest, ""},
			when{nil, service.ErrInvalidURL},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := service.NewMockShortLinkService(t)
			s.EXPECT().CreateBatch(mock.Anything, mock.Anything, mock.Anything).Return(tc.when.outputs, tc.when.err)

			d := service.NewMockShortLinkDeleter(t)

			l := logger.NewMockLogger(t)
			l.EXPECT().Errorf(mock.Anything, mock.Anything).Maybe()

			p := provider.NewMockShortLinkProvider(t)

			handler := NewShortLinkHandler(s, p, urlAddress, d, l)

			r := chi.NewRouter()
			r.Post("/api/shorten/batch", handler.HandleCreateShortenBatch)

			srv := httptest.NewServer(r)
			defer srv.Close()

			client := resty.New()
			client.SetRedirectPolicy(resty.NoRedirectPolicy())

			req := client.R()
			req.Method = resty.MethodPost
			req.URL = srv.URL + "/api/shorten/batch"

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

func TestShortLinkHandler_HandleGetUserUrls(t *testing.T) {
	urlAddress := "http://localhost:8080"
	userID1 := uuid.New()

	type on struct {
		userID *uuid.UUID
	}
	type want struct {
		code int
		body string
	}
	type when struct {
		items []model.ShortLink
		err   error
	}
	testCases := []struct {
		name string
		on   on
		want want
		when when
	}{
		{
			"no user",
			on{nil},
			want{http.StatusUnauthorized, ""},
			when{nil, nil},
		},
		{
			"empty items",
			on{&userID1},
			want{http.StatusNoContent, ""},
			when{nil, nil},
		},
		{
			"provider error",
			on{&userID1},
			want{http.StatusInternalServerError, ""},
			when{nil, errors.New("some provider error")},
		},
		{
			"found items",
			on{&userID1},
			want{http.StatusOK, fmt.Sprintf(`[{"short_url": "%s/id1", "original_url": "https://example1.com"},{"short_url": "%s/id2", "original_url": "https://example2.com"}]`, urlAddress, urlAddress)},
			when{[]model.ShortLink{{ID: "id1", URL: "https://example1.com"}, {ID: "id2", URL: "https://example2.com"}}, nil},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := provider.NewMockShortLinkProvider(t)
			if tc.on.userID != nil {
				p.EXPECT().FindByUserID(mock.Anything, *(tc.on.userID)).Return(tc.when.items, tc.when.err)
			}

			d := service.NewMockShortLinkDeleter(t)

			l := logger.NewMockLogger(t)
			l.EXPECT().Errorf(mock.Anything, mock.Anything).Maybe()

			s := service.NewMockShortLinkService(t)

			handler := NewShortLinkHandler(s, p, urlAddress, d, l)

			r := chi.NewRouter()
			r.Get("/api/user/urls", handler.HandleGetUserUrls)

			middleware := func(next http.Handler) http.HandlerFunc {
				return func(w http.ResponseWriter, re *http.Request) {
					ctx := re.Context()
					if tc.on.userID != nil {
						ctx = auth.CreateTokenContext(ctx, auth.Token{UserID: tc.on.userID.String()})
					}
					next.ServeHTTP(w, re.WithContext(ctx))
				}
			}
			var router http.Handler = r
			router = middleware(router)

			srv := httptest.NewServer(router)
			defer srv.Close()

			client := resty.New()
			client.SetRedirectPolicy(resty.NoRedirectPolicy())

			resp, err := client.R().
				Get(srv.URL + "/api/user/urls")

			require.NoError(t, err)

			require.Equal(t, tc.want.code, resp.StatusCode(), fmt.Sprintf("expected status code %d but got %d with body: %s", tc.want.code, resp.StatusCode(), string(resp.Body())))

			if tc.want.body != "" {
				assert.JSONEq(t, tc.want.body, string(resp.Body()))
			}

			if tc.want.code == http.StatusNoContent || tc.want.code == http.StatusOK {
				assert.Equal(t, "application/json", resp.Header().Get("Content-Type"))
			}
		})
	}
}

func TestShortLinkHandler_HandleDeleteUrls(t *testing.T) {
	urlAddress := "http://localhost:8080"
	userID1 := uuid.New()

	type on struct {
		userID *uuid.UUID
		body   string
	}
	type want struct {
		code int
	}
	type when struct {
		deleteErr error
	}
	testCases := []struct {
		name string
		on   on
		want want
		when when
	}{
		{
			"empty ids by user",
			on{&userID1, "[]"},
			want{http.StatusAccepted},
			when{nil},
		},
		{
			"empty ids by no user",
			on{nil, "[]"},
			want{http.StatusUnauthorized},
			when{nil},
		},
		{
			"not empty ids by user",
			on{&userID1, `["id1","id2","id3"]`},
			want{http.StatusAccepted},
			when{nil},
		},
		{
			"not empty ids by no user",
			on{nil, `["id1","id2","id3"]`},
			want{http.StatusUnauthorized},
			when{nil},
		},
		{
			"invalid ids",
			on{&userID1, `"not array of ids"`},
			want{http.StatusInternalServerError},
			when{nil},
		},
		{
			"delete error",
			on{&userID1, `["id1","id2","id3"]`},
			want{http.StatusInternalServerError},
			when{errors.New("some server error")},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := service.NewMockShortLinkDeleter(t)
			d.EXPECT().Add(mock.Anything).Return(tc.when.deleteErr).Maybe()

			l := logger.NewMockLogger(t)
			l.EXPECT().Errorf(mock.Anything, mock.Anything).Maybe()

			s := service.NewMockShortLinkService(t)
			p := provider.NewMockShortLinkProvider(t)

			handler := NewShortLinkHandler(s, p, urlAddress, d, l)

			r := chi.NewRouter()
			r.Delete("/api/user/urls", handler.HandleDeleteUrls)

			middleware := func(next http.Handler) http.HandlerFunc {
				return func(w http.ResponseWriter, re *http.Request) {
					ctx := re.Context()
					if tc.on.userID != nil {
						ctx = auth.CreateTokenContext(ctx, auth.Token{UserID: tc.on.userID.String()})
					}
					next.ServeHTTP(w, re.WithContext(ctx))
				}
			}
			var router http.Handler = r
			router = middleware(router)

			srv := httptest.NewServer(router)
			defer srv.Close()

			client := resty.New()
			client.SetRedirectPolicy(resty.NoRedirectPolicy())

			req := client.R()
			req.Method = resty.MethodDelete
			req.URL = srv.URL + "/api/user/urls"

			if len(tc.on.body) > 0 {
				req.SetHeader("Content-Type", "application/json")
				req.SetBody(tc.on.body)
			}

			resp, err := req.Send()
			require.NoError(t, err)

			require.Equal(t, tc.want.code, resp.StatusCode(), fmt.Sprintf("expected status code %d but got %d with body: %s", tc.want.code, resp.StatusCode(), string(resp.Body())))
		})
	}
}
