package handler

import (
	"context"
	"net/http"

	"github.com/liebeSonne/shortlink/internal/model"
	"github.com/stretchr/testify/mock"
)

type mockShortLinkHandler struct {
	mock.Mock
}

func (m *mockShortLinkHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}
func (m *mockShortLinkHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

func (m *mockShortLinkHandler) HandleCreateShorten(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

type mockService struct {
	mock.Mock
}

func (m *mockService) Create(url string) (*model.ShortLink, error) {
	args := m.Called(url)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ShortLink), args.Error(1)
}

type mockProvider struct {
	mock.Mock
}

func (m *mockProvider) Find(id string) (*model.ShortLink, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ShortLink), args.Error(1)
}

type mockDatabaseHandler struct {
	mock.Mock
}

func (m *mockDatabaseHandler) HandlePing(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

type mockDatabase struct {
	mock.Mock
}

func (m *mockDatabase) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

type mockLogger struct {
	mock.Mock
}

func (l *mockLogger) Debugf(format string, args ...interface{}) {
	l.Called(format, args)
}
func (l *mockLogger) Infof(format string, args ...interface{}) {
	l.Called(format, args)
}
func (l *mockLogger) Warnf(format string, args ...interface{}) {
	l.Called(format, args)
}
func (l *mockLogger) Errorf(format string, args ...interface{}) {
	l.Called(format, args)
}
func (l *mockLogger) Fatalf(format string, args ...interface{}) {
	l.Called(format, args)
}
func (l *mockLogger) Panicf(format string, args ...interface{}) {
	l.Called(format, args)
}
func (l *mockLogger) Debugw(msg string, keysAndValues ...interface{}) {
	l.Called(msg, keysAndValues)
}
func (l *mockLogger) Infow(msg string, keysAndValues ...interface{}) {
	l.Called(msg, keysAndValues)
}
func (l *mockLogger) Warnw(msg string, keysAndValues ...interface{}) {
	l.Called(msg, keysAndValues)
}
func (l *mockLogger) Errorw(msg string, keysAndValues ...interface{}) {
	l.Called(msg, keysAndValues)
}
func (l *mockLogger) Fatalw(msg string, keysAndValues ...interface{}) {
	l.Called(msg, keysAndValues)
}
func (l *mockLogger) Panicw(msg string, keysAndValues ...interface{}) {
	l.Called(msg, keysAndValues)
}
