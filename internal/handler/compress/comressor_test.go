package compress

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewCompressorMiddleware(t *testing.T) {
	type on struct {
		cfg CompressorConfig
		h   http.Handler
	}
	type want struct {
		err error
	}
	testCases := []struct {
		name string
		on   on
		want want
	}{
		{"empty encoding and empty content types", on{CompressorConfig{[]Encoding{}, &[]string{}}, new(mockHandler)}, want{nil}},
		{"empty encoding and unset content types", on{CompressorConfig{[]Encoding{}, nil}, new(mockHandler)}, want{nil}},
		{"empty encoding and not empty content types", on{CompressorConfig{[]Encoding{}, &[]string{"application/json", "text/html"}}, new(mockHandler)}, want{nil}},
		{"not empty encoding and empty content types", on{CompressorConfig{[]Encoding{GzipEncoding}, &[]string{}}, new(mockHandler)}, want{nil}},
		{"not empty encoding and unset content types", on{CompressorConfig{[]Encoding{GzipEncoding}, nil}, new(mockHandler)}, want{nil}},
		{"not empty encoding and not empty content types", on{CompressorConfig{[]Encoding{GzipEncoding}, &[]string{"application/json", "text/html"}}, new(mockHandler)}, want{nil}},
		{"gzip", on{CompressorConfig{[]Encoding{GzipEncoding}, &[]string{"application/json", "text/html"}}, new(mockHandler)}, want{nil}},
		{"gzip with unset content types", on{CompressorConfig{[]Encoding{GzipEncoding}, nil}, new(mockHandler)}, want{nil}},
		{"deflate", on{CompressorConfig{[]Encoding{DeflateEncoding}, &[]string{"application/json", "text/html"}}, new(mockHandler)}, want{nil}},
		{"deflate with unset content types", on{CompressorConfig{[]Encoding{DeflateEncoding}, nil}, new(mockHandler)}, want{nil}},
		{"gzip and deflate", on{CompressorConfig{[]Encoding{GzipEncoding, DeflateEncoding}, &[]string{"application/json", "text/html"}}, new(mockHandler)}, want{nil}},
		{"unknown encoding", on{CompressorConfig{[]Encoding{-1}, &[]string{"application/json", "text/html"}}, new(mockHandler)}, want{ErrUnknownEncoding}},
		{"unknown one of encoding", on{CompressorConfig{[]Encoding{GzipEncoding, -1}, &[]string{"application/json", "text/html"}}, new(mockHandler)}, want{ErrUnknownEncoding}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			h, err := NewCompressorMiddleware(tc.on.h, tc.on.cfg)
			if tc.want.err != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tc.want.err)
				return
			}

			require.NotNil(t, h)
		})
	}
}
