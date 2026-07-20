package deflate

import (
	"io"
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

type MockReadCloser struct {
	mock.Mock
}

func (m *MockReadCloser) Read(p []byte) (n int, err error) {
	args := m.Called(p)
	return args.Int(0), args.Error(1)
}

func (m *MockReadCloser) Close() error {
	args := m.Called()
	return args.Error(0)
}

var _ io.ReadCloser = (*MockReadCloser)(nil)
