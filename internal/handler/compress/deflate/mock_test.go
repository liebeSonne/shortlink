package deflate

import (
	"net/http"

	"github.com/stretchr/testify/mock"
)

type mockResponseWriter struct {
	mock.Mock
}

func (m *mockResponseWriter) Header() http.Header {
	args := m.Called()
	return args[0].(http.Header)
}

func (m *mockResponseWriter) Write(bytes []byte) (int, error) {
	args := m.Called(bytes)
	return args.Int(0), args.Error(1)
}

func (m *mockResponseWriter) WriteHeader(statusCode int) {
	m.Called(statusCode)
}

type mockHandler struct {
	mock.Mock
}

func (m *mockHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	m.Called(writer, request)
}
