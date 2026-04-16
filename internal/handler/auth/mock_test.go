package auth

import (
	"net/http"

	"github.com/stretchr/testify/mock"

	"github.com/liebeSonne/shortlink/internal/auth"
)

type mockHandler struct {
	mock.Mock
}

func (m *mockHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	m.Called(writer, request)
}

type mockService struct {
	mock.Mock
}

func (m *mockService) SetAuthToken(tokenString string, w http.ResponseWriter, r *http.Request) error {
	args := m.Called(tokenString, w, r)
	return args.Error(0)
}

func (m *mockService) GetAuthToken(r *http.Request) (string, error) {
	args := m.Called(r)
	return args.String(0), args.Error(1)
}

type mockTokenService struct {
	mock.Mock
}

func (m *mockTokenService) Create(tokenData auth.Token) (string, error) {
	args := m.Called(tokenData)
	return args.String(0), args.Error(1)
}

func (m *mockTokenService) Parse(tokenString string) (auth.Token, error) {
	args := m.Called(tokenString)
	return args.Get(0).(auth.Token), args.Error(1)
}
