package deflate

import (
	"bytes"
	"compress/zlib"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/liebeSonne/shortlink/internal/logger"
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
			l := logger.NewMockLogger(t)
			h := NewDeflateHandlerMiddleware(tc.on.h, tc.on.contentTypes, l)
			require.NotNil(t, h)
		})
	}
}

func TestDeflateHandlerMiddleware_Handle(t *testing.T) {
	responseBody := []byte("test response body content")

	type on struct {
		acceptEncoding *string
		contentType    *string
		contentTypes   *[]string
		handlerBody    []byte
	}
	type want struct {
		contentEncoding *string
		body            []byte
	}
	testCases := []struct {
		name string
		on   on
		want want
	}{
		{
			"compress response when accept-encoding deflate",
			on{
				acceptEncoding: strPtr("deflate"),
				contentTypes:   nil,
				handlerBody:    responseBody,
			},
			want{
				contentEncoding: strPtr("deflate"),
				body:            responseBody,
			},
		},
		{
			"compress response when accept-encoding contains deflate",
			on{
				acceptEncoding: strPtr("gzip, deflate, br"),
				contentTypes:   nil,
				handlerBody:    responseBody,
			},
			want{
				contentEncoding: strPtr("deflate"),
				body:            responseBody,
			},
		},
		{
			"no compress when accept-encoding not deflate",
			on{
				acceptEncoding: strPtr("gzip"),
				contentTypes:   nil,
				handlerBody:    responseBody,
			},
			want{
				contentEncoding: nil,
				body:            responseBody,
			},
		},
		{
			"no compress when accept-encoding empty",
			on{
				acceptEncoding: nil,
				contentTypes:   nil,
				handlerBody:    responseBody,
			},
			want{
				contentEncoding: nil,
				body:            responseBody,
			},
		},
		{
			"compress when content-type matches",
			on{
				acceptEncoding: strPtr("deflate"),
				contentType:    strPtr("application/json"),
				contentTypes:   ptrSlice([]string{"application/json"}),
				handlerBody:    responseBody,
			},
			want{
				contentEncoding: strPtr("deflate"),
				body:            responseBody,
			},
		},
		{
			"no compress when content-type not allowed",
			on{
				acceptEncoding: strPtr("deflate"),
				contentType:    strPtr("application/json"),
				contentTypes:   ptrSlice([]string{"text/html"}),
				handlerBody:    responseBody,
			},
			want{
				contentEncoding: nil,
				body:            responseBody,
			},
		},
		{
			"compress when content-types unset",
			on{
				acceptEncoding: strPtr("deflate"),
				contentTypes:   nil,
				handlerBody:    responseBody,
			},
			want{
				contentEncoding: strPtr("deflate"),
				body:            responseBody,
			},
		},
		{
			"no compress when content-types empty",
			on{
				acceptEncoding: strPtr("deflate"),
				contentTypes:   ptrSlice([]string{}),
				handlerBody:    responseBody,
			},
			want{
				contentEncoding: nil,
				body:            responseBody,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockHandler := new(mockHandler)
			mockHandler.On("ServeHTTP", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
				w := args.Get(0).(http.ResponseWriter)
				w.WriteHeader(http.StatusOK)
				w.Write(tc.on.handlerBody)
			}).Once()

			l := logger.NewMockLogger(t)
			l.EXPECT().Errorf(mock.Anything, mock.Anything).Maybe()

			middleware := NewDeflateHandlerMiddleware(mockHandler, tc.on.contentTypes, l)

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tc.on.acceptEncoding != nil {
				req.Header.Set("Accept-Encoding", *tc.on.acceptEncoding)
			}
			if tc.on.contentType != nil {
				req.Header.Set("Content-Type", *tc.on.contentType)
			}

			rec := httptest.NewRecorder()
			middleware.ServeHTTP(rec, req)

			mockHandler.AssertExpectations(t)

			if tc.want.contentEncoding != nil {
				assert.Equal(t, *tc.want.contentEncoding, rec.Header().Get("Content-Encoding"))

				dec, err := zlib.NewReader(bytes.NewReader(rec.Body.Bytes()))
				require.NoError(t, err)
				defer dec.Close()

				decompressed, err := io.ReadAll(dec)
				require.NoError(t, err)
				assert.Equal(t, tc.want.body, decompressed)
			} else {
				assert.Empty(t, rec.Header().Get("Content-Encoding"))
				assert.Equal(t, tc.want.body, rec.Body.Bytes())
			}
		})
	}
}

func TestDeflateHandlerMiddleware_RequestDecompress(t *testing.T) {
	compressedData := compressData(t, []byte("test request body"))

	type on struct {
		contentEncoding *string
		body            []byte
	}
	type want struct {
		requestBody []byte
	}
	testCases := []struct {
		name string
		on   on
		want want
	}{
		{
			"decompress request when content-encoding deflate",
			on{
				contentEncoding: strPtr("deflate"),
				body:            compressedData,
			},
			want{
				requestBody: []byte("test request body"),
			},
		},
		{
			"no decompress when content-encoding not deflate",
			on{
				contentEncoding: strPtr("gzip"),
				body:            []byte("raw body"),
			},
			want{
				requestBody: nil,
			},
		},
		{
			"no decompress when content-encoding empty",
			on{
				contentEncoding: nil,
				body:            []byte("raw body"),
			},
			want{
				requestBody: nil,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockHandler := new(mockHandler)
			mockHandler.On("ServeHTTP", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
				r := args.Get(1).(*http.Request)
				if tc.want.requestBody != nil {
					bodyBytes, _ := io.ReadAll(r.Body)
					assert.Equal(t, tc.want.requestBody, bodyBytes)
				}
			}).Once()

			l := logger.NewMockLogger(t)
			l.EXPECT().Errorf(mock.Anything, mock.Anything).Maybe()

			middleware := NewDeflateHandlerMiddleware(mockHandler, nil, l)

			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(tc.on.body))
			if tc.on.contentEncoding != nil {
				req.Header.Set("Content-Encoding", *tc.on.contentEncoding)
			}

			rec := httptest.NewRecorder()
			middleware.ServeHTTP(rec, req)

			mockHandler.AssertExpectations(t)
		})
	}
}

func TestDeflateHandlerMiddleware_ErrorCases(t *testing.T) {
	testCases := []struct {
		name string
		on   struct {
			contentEncoding *string
			acceptEncoding  *string
		}
	}{
		{
			"error creating deflate reader",
			struct {
				contentEncoding *string
				acceptEncoding  *string
			}{
				contentEncoding: strPtr("deflate"),
				acceptEncoding:  nil,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockHandler := new(mockHandler)

			l := logger.NewMockLogger(t)
			l.EXPECT().Errorf(mock.Anything, mock.Anything).Once()

			middleware := NewDeflateHandlerMiddleware(mockHandler, nil, l)

			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte("invalid zlib data")))
			if tc.on.contentEncoding != nil {
				req.Header.Set("Content-Encoding", *tc.on.contentEncoding)
			}
			if tc.on.acceptEncoding != nil {
				req.Header.Set("Accept-Encoding", *tc.on.acceptEncoding)
			}

			rec := httptest.NewRecorder()
			middleware.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		})
	}
}

func strPtr(s string) *string {
	return &s
}

func ptrSlice(s []string) *[]string {
	return &s
}

func compressData(t *testing.T, data []byte) []byte {
	var buf bytes.Buffer
	w, err := zlib.NewWriterLevel(&buf, zlib.BestSpeed)
	require.NoError(t, err)
	_, err = w.Write(data)
	require.NoError(t, err)
	err = w.Close()
	require.NoError(t, err)
	return buf.Bytes()
}
