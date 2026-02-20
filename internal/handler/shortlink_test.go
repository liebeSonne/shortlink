package handler

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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
			provider := &mockProvider{url: tc.when.link, err: tc.when.err}
			handler := NewShortLinkHandler(new(mockService), provider)

			request := httptest.NewRequest(http.MethodGet, "/"+tc.on.id, nil)
			w := httptest.NewRecorder()
			handler.HandleGet(w, request)
			res := w.Result()

			require.Equal(t, tc.want.code, res.StatusCode)

			if tc.want.code == http.StatusTemporaryRedirect && tc.want.location != nil {
				location := w.Header().Get("Location")
				assert.Equal(t, *tc.want.location, location)
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
		code int
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
			want{http.StatusCreated},
			when{id1, nil},
		},
		{
			"error in service",
			on{link1},
			want{http.StatusInternalServerError},
			when{"", errors.New("some service error")},
		},
		{
			"error empty url",
			on{link1},
			want{http.StatusBadRequest},
			when{"", model.ErrEmptyURL},
		},
		{
			"error empty id",
			on{link1},
			want{http.StatusBadRequest},
			when{"", model.ErrEmptyID},
		},
		{
			"error invalid url",
			on{link1},
			want{http.StatusBadRequest},
			when{"", model.ErrInvalidURL},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			service := &mockService{id: tc.when.id, err: tc.when.err}
			handler := NewShortLinkHandler(service, new(mockProvider))

			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tc.on.link))
			request.Host = "localhost"
			request.URL.Scheme = "http"
			request.TLS = nil
			w := httptest.NewRecorder()
			handler.HandleCreate(w, request)
			res := w.Result()

			require.Equal(t, tc.want.code, res.StatusCode)

			if tc.want.code == http.StatusCreated {
				defer res.Body.Close()
				resBody, err := io.ReadAll(res.Body)

				require.NoError(t, err)
				wantURL := fmt.Sprintf("%s://%s/%s", request.URL.Scheme, request.Host, tc.when.id)
				assert.Equal(t, wantURL, string(resBody))
			}
		})
	}
}

type mockService struct {
	id  string
	err error
}

func (m *mockService) Create(url string) (model.ShortLink, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &mockShortLink{id: m.id, url: url}, nil
}

type mockProvider struct {
	url *string
	err error
}

func (m *mockProvider) Get(id string) (*model.ShortLink, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.url != nil {
		item := testNewMockShortLink(id, *m.url)
		return &item, nil
	}
	return nil, nil
}

func testNewMockShortLink(id, url string) model.ShortLink {
	return &mockShortLink{id: id, url: url}
}

type mockShortLink struct {
	id  string
	url string
}

func (m *mockShortLink) ID() string {
	return m.id
}

func (m *mockShortLink) URL() string {
	return m.url
}
