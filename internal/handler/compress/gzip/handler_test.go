package gzip

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewDeflateHandlerMiddleware2(t *testing.T) {
	type on struct {
		h            http.Handler
		contentTypes *[]string
	}
	testCases := []struct {
		name string
		on   on
	}{
		{"mock handler", on{new(mockHandler), &[]string{"application/json", "text/html"}}},
		{"mock handler with empty content types", on{new(mockHandler), &[]string{}}},
		{"mock handler with unset content types", on{new(mockHandler), nil}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			h := NewGzipHandlerMiddleware(tc.on.h, tc.on.contentTypes)
			require.NotNil(t, h)
		})
	}
}
