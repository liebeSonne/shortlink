package compress

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewEncodingMiddleware(t *testing.T) {
	type on struct {
		encoding     Encoding
		contentTypes *[]string
		h            http.Handler
	}
	type want struct {
		err error
	}
	testCases := []struct {
		name string
		on   on
		want want
	}{
		{"gzip", on{GzipEncoding, &[]string{"application/json", "text/html"}, new(mockHandler)}, want{nil}},
		{"gzip with empty content types", on{GzipEncoding, &[]string{}, new(mockHandler)}, want{nil}},
		{"gzip with unset content types", on{GzipEncoding, nil, new(mockHandler)}, want{nil}},
		{"deflate", on{DeflateEncoding, &[]string{"application/json", "text/html"}, new(mockHandler)}, want{nil}},
		{"deflate with empty content types", on{DeflateEncoding, &[]string{}, new(mockHandler)}, want{nil}},
		{"deflate with unset content types", on{DeflateEncoding, nil, new(mockHandler)}, want{nil}},
		{"unknown", on{-1, &[]string{"application/json", "text/html"}, new(mockHandler)}, want{ErrUnknownEncoding}},
		{"unknown with empty content types", on{-1, &[]string{}, new(mockHandler)}, want{ErrUnknownEncoding}},
		{"unknown with unset content types", on{-1, nil, new(mockHandler)}, want{ErrUnknownEncoding}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			h, err := NewEncodingMiddleware(tc.on.h, tc.on.encoding, tc.on.contentTypes)
			if tc.want.err != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tc.want.err)
				return
			}

			require.NotNil(t, h)
		})
	}
}
