package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/liebeSonne/shortlink/internal/auth"
	"github.com/liebeSonne/shortlink/internal/handler/audit"
	"github.com/liebeSonne/shortlink/internal/logger"
	"github.com/liebeSonne/shortlink/internal/model"
	"github.com/liebeSonne/shortlink/internal/provider"
	"github.com/liebeSonne/shortlink/internal/service"
)

func ExampleShortLinkHandler_handleGet() {
	id := "id1"
	link := "https://example1.com"
	urlAddress := "http://localhost:8080"

	// mocks
	t := &testing.T{}
	p := provider.NewMockShortLinkProvider(t)
	p.EXPECT().Find(mock.Anything, mock.Anything).Return(&model.ShortLink{ID: id, URL: link}, nil).Maybe()
	d := service.NewMockShortLinkDeleter(t)
	l := logger.NewMockLogger(t)
	s := service.NewMockShortLinkService(t)
	ap := audit.NewMockPublisher(t)
	ap.EXPECT().Notify(mock.Anything).Return().Maybe()

	handler := NewShortLinkHandler(s, p, urlAddress, d, l, ap)

	r := chi.NewRouter()
	r.Get("/{id}", handler.HandleGet)

	req := httptest.NewRequest(http.MethodGet, "/"+id, nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	fmt.Println("Status:", rec.Code)
	fmt.Println("Location:", rec.Header().Get("Location"))

	// Output:
	// Status: 307
	// Location: https://example1.com
}

func ExampleShortLinkHandler_handleCreate() {
	id := "id1"
	link := "https://example1.com"
	urlAddress := "http://localhost:8080"

	// mocks
	t := &testing.T{}
	p := provider.NewMockShortLinkProvider(t)
	p.EXPECT().FindByURL(mock.Anything, mock.Anything).Return(nil, nil).Maybe()
	d := service.NewMockShortLinkDeleter(t)
	l := logger.NewMockLogger(t)
	s := service.NewMockShortLinkService(t)
	s.EXPECT().Create(mock.Anything, mock.Anything, mock.Anything).Return(&model.ShortLink{ID: id, URL: link}, nil)
	ap := audit.NewMockPublisher(t)
	ap.EXPECT().Notify(mock.Anything).Return().Maybe()

	handler := NewShortLinkHandler(s, p, urlAddress, d, l, ap)

	r := chi.NewRouter()
	r.Post("/", handler.HandleCreate)

	reqBody := strings.NewReader(link)
	req := httptest.NewRequest(http.MethodPost, "/", reqBody)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	fmt.Println("Status:", rec.Code)
	fmt.Println("Body:", rec.Body.String())

	// Output:
	// Status: 201
	// Body: http://localhost:8080/id1
}

func ExampleShortLinkHandler_handleCreateShorten() {
	id := "id1"
	link := "https://example1.com"
	urlAddress := "http://localhost:8080"

	// mocks
	t := &testing.T{}
	p := provider.NewMockShortLinkProvider(t)
	p.EXPECT().FindByURL(mock.Anything, mock.Anything).Return(nil, nil).Maybe()
	d := service.NewMockShortLinkDeleter(t)
	l := logger.NewMockLogger(t)
	s := service.NewMockShortLinkService(t)
	s.EXPECT().Create(mock.Anything, mock.Anything, mock.Anything).Return(&model.ShortLink{ID: id, URL: link}, nil)
	ap := audit.NewMockPublisher(t)
	ap.EXPECT().Notify(mock.Anything).Return().Maybe()

	handler := NewShortLinkHandler(s, p, urlAddress, d, l, ap)

	r := chi.NewRouter()
	r.Post("/api/shorten", handler.HandleCreateShorten)

	reqBody := strings.NewReader(fmt.Sprintf(`{"url": "%s"}`, link))
	req := httptest.NewRequest(http.MethodPost, "/api/shorten", reqBody)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	fmt.Println("Status:", rec.Code)
	fmt.Println("Body:", rec.Body.String())

	// Output:
	// Status: 201
	// Body: {"result":"http://localhost:8080/id1"}
}

func ExampleShortLinkHandler_handleCreateShortenBatch() {
	id := "id1"
	link := "https://example1.com"
	urlAddress := "http://localhost:8080"

	// mocks
	t := &testing.T{}
	p := provider.NewMockShortLinkProvider(t)
	p.EXPECT().FindByURL(mock.Anything, mock.Anything).Return(nil, nil).Maybe()
	d := service.NewMockShortLinkDeleter(t)
	l := logger.NewMockLogger(t)
	s := service.NewMockShortLinkService(t)
	s.EXPECT().CreateBatch(mock.Anything, mock.Anything, mock.Anything).Return([]service.OutputShortLinkData{
		{CorrelationID: "c1", ShortLink: model.ShortLink{ID: id, URL: link}},
	}, nil)

	ap := audit.NewMockPublisher(t)
	ap.EXPECT().Notify(mock.Anything).Return().Maybe()

	handler := NewShortLinkHandler(s, p, urlAddress, d, l, ap)

	r := chi.NewRouter()
	r.Post("/api/shorten/batch", handler.HandleCreateShortenBatch)

	reqBody := strings.NewReader(fmt.Sprintf(`[{"correlation_id": "c1", "original_url": "%s"}]`, link))
	req := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", reqBody)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	fmt.Println("Status:", rec.Code)
	fmt.Println("Body:", rec.Body.String())

	// Output:
	// Status: 201
	// Body: [{"correlation_id":"c1","short_url":"http://localhost:8080/id1"}]
}

func ExampleShortLinkHandler_handleGetUserUrls() {
	urlAddress := "http://localhost:8080"
	userID := uuid.New()

	// mocks
	t := &testing.T{}
	p := provider.NewMockShortLinkProvider(t)
	p.EXPECT().FindByUserID(mock.Anything, mock.Anything).Return([]model.ShortLink{
		{ID: "id1", URL: "https://example1.com"},
		{ID: "id2", URL: "https://example2.com"},
	}, nil).Maybe()
	d := service.NewMockShortLinkDeleter(t)
	l := logger.NewMockLogger(t)
	s := service.NewMockShortLinkService(t)
	ap := audit.NewMockPublisher(t)
	ap.EXPECT().Notify(mock.Anything).Return().Maybe()

	handler := NewShortLinkHandler(s, p, urlAddress, d, l, ap)

	r := chi.NewRouter()
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = auth.CreateTokenContext(ctx, auth.Token{UserID: userID.String()})
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})
	r.Get("/api/user/urls", handler.HandleGetUserUrls)

	req := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	fmt.Println("Status:", rec.Code)
	fmt.Println("Body:", rec.Body.String())

	// Output:
	// Status: 200
	// Body: [{"short_url":"http://localhost:8080/id1","original_url":"https://example1.com"},{"short_url":"http://localhost:8080/id2","original_url":"https://example2.com"}]
}

func ExampleShortLinkHandler_handleDeleteUrls() {
	urlAddress := "http://localhost:8080"
	userID := uuid.New()

	// mocks
	t := &testing.T{}
	p := provider.NewMockShortLinkProvider(t)
	d := service.NewMockShortLinkDeleter(t)
	d.EXPECT().Add(mock.Anything).Return(nil).Maybe()
	l := logger.NewMockLogger(t)
	s := service.NewMockShortLinkService(t)

	ap := audit.NewMockPublisher(t)
	ap.EXPECT().Notify(mock.Anything).Return().Maybe()

	handler := NewShortLinkHandler(s, p, urlAddress, d, l, ap)

	r := chi.NewRouter()
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = auth.CreateTokenContext(ctx, auth.Token{UserID: userID.String()})
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})
	r.Delete("/api/user/urls", handler.HandleDeleteUrls)

	reqBody := strings.NewReader(`["id1","id2"]`)
	req := httptest.NewRequest(http.MethodDelete, "/api/user/urls", reqBody)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	fmt.Println("Status:", rec.Code)

	// Output:
	// Status: 202
}
