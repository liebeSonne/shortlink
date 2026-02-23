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
				mockItem := new(mockShortLink)
				mockItem.On("ID").Return(tc.on.id).On("URL").Return(*tc.when.link)
				provider.On("Get", tc.on.id).Return(mockItem, tc.when.err)
			} else {
				provider.On("Get", tc.on.id).Return(nil, tc.when.err)
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
			when{"", model.ErrEmptyURL},
		},
		{
			"error empty id",
			on{link1},
			want{http.StatusBadRequest, ""},
			when{"", model.ErrEmptyID},
		},
		{
			"error invalid url",
			on{link1},
			want{http.StatusBadRequest, ""},
			when{"", model.ErrInvalidURL},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var item model.ShortLink
			if tc.when.err == nil {
				mockItem := new(mockShortLink)
				mockItem.On("ID").Return(tc.when.id).On("URL").Return(link1)
				item = mockItem
			}
			service := new(mockService)
			service.On("Create", tc.on.link).Return(item, tc.when.err)

			urlAddress := "http://localhost:8080"
			handler := NewShortLinkHandler(service, new(mockProvider), urlAddress)

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

type mockService struct {
	mock.Mock
}

func (m *mockService) Create(url string) (model.ShortLink, error) {
	args := m.Called(url)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(model.ShortLink), args.Error(1)
}

type mockProvider struct {
	mock.Mock
}

func (m *mockProvider) Get(id string) (model.ShortLink, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(model.ShortLink), args.Error(1)
}

type mockShortLink struct {
	mock.Mock
}

func (m *mockShortLink) ID() string {
	args := m.Called()
	return args.String(0)
}

func (m *mockShortLink) URL() string {
	args := m.Called()
	return args.String(0)
}
