package compress

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func TestHandleWithCompressor(t *testing.T) {
	generateResponse := func(requestBody string) string {
		return "Response: " + requestBody
	}

	type on struct {
		method  string
		path    string
		body    string
		headers map[string]string
	}
	type when struct {
		encodings    []Encoding
		contentTypes *[]string
	}
	type want struct {
		code      int
		body      string
		header    http.Header
		noHeader  http.Header
		encodings []Encoding
	}
	testCases := []struct {
		name string
		on   on
		when when
		want want
	}{
		// gzip compress
		{"compress gzip when ok",
			on{http.MethodGet, "/1", "123", map[string]string{
				"Accept-Encoding": "gzip",
				"Content-Type":    "application/json; charset=utf-8",
			}},
			when{[]Encoding{GzipEncoding}, nil},
			want{http.StatusOK, generateResponse("123"), http.Header{
				"Content-Encoding": []string{"gzip"},
			}, http.Header{}, []Encoding{GzipEncoding}},
		},
		{"compress gzip when error",
			on{http.MethodGet, "/1", "123", map[string]string{
				"Accept-Encoding": "gzip",
				"Content-Type":    "application/json; charset=utf-8",
			}},
			when{[]Encoding{GzipEncoding}, nil},
			want{http.StatusInternalServerError, "", http.Header{}, http.Header{
				"Content-Encoding": []string{"gzip"},
			}, nil},
		},
		{"compress gzip when unset content types",
			on{http.MethodGet, "/1", "123", map[string]string{
				"Accept-Encoding": "gzip",
				"Content-Type":    "application/json; charset=utf-8",
			}},
			when{[]Encoding{GzipEncoding}, nil},
			want{http.StatusOK, generateResponse("123"), http.Header{
				"Content-Encoding": []string{"gzip"},
			}, http.Header{}, []Encoding{GzipEncoding}},
		},
		{"compress gzip when empty set content types",
			on{http.MethodGet, "/1", "123", map[string]string{
				"Accept-Encoding": "gzip",
				"Content-Type":    "application/json; charset=utf-8",
			}},
			when{[]Encoding{GzipEncoding}, &[]string{}},
			want{http.StatusOK, generateResponse("123"), http.Header{}, http.Header{
				"Content-Encoding": []string{"gzip"},
			}, nil},
		},
		{"compress gzip when not set request content types",
			on{http.MethodGet, "/1", "123", map[string]string{
				"Accept-Encoding": "gzip",
				"Content-Type":    "application/json; charset=utf-8",
			}},
			when{[]Encoding{GzipEncoding}, &[]string{"test/html"}},
			want{http.StatusOK, generateResponse("123"), http.Header{}, http.Header{
				"Content-Encoding": []string{"gzip"},
			}, nil},
		},
		// deflate
		{"compress deflate when ok",
			on{http.MethodGet, "/1", "123", map[string]string{
				"Accept-Encoding": "deflate",
				"Content-Type":    "application/json; charset=utf-8",
			}},
			when{[]Encoding{DeflateEncoding}, nil},
			want{http.StatusOK, generateResponse("123"), http.Header{
				"Content-Encoding": []string{"deflate"},
			}, http.Header{}, []Encoding{DeflateEncoding}},
		},
		{"compress deflate when error",
			on{http.MethodGet, "/1", "123", map[string]string{
				"Accept-Encoding": "deflate",
				"Content-Type":    "application/json; charset=utf-8",
			}},
			when{[]Encoding{DeflateEncoding}, nil},
			want{http.StatusInternalServerError, "", http.Header{}, http.Header{
				"Content-Encoding": []string{"deflate"},
			}, nil},
		},
		{"compress deflate when unset content types",
			on{http.MethodGet, "/1", "123", map[string]string{
				"Accept-Encoding": "deflate",
				"Content-Type":    "application/json; charset=utf-8",
			}},
			when{[]Encoding{DeflateEncoding}, nil},
			want{http.StatusOK, generateResponse("123"), http.Header{
				"Content-Encoding": []string{"deflate"},
			}, http.Header{}, []Encoding{DeflateEncoding}},
		},
		{"compress deflate when empty set content types",
			on{http.MethodGet, "/1", "123", map[string]string{
				"Accept-Encoding": "deflate",
				"Content-Type":    "application/json; charset=utf-8",
			}},
			when{[]Encoding{DeflateEncoding}, &[]string{}},
			want{http.StatusOK, generateResponse("123"), http.Header{}, http.Header{
				"Content-Encoding": []string{"deflate"},
			}, nil},
		},
		{"compress deflate when not set request content types",
			on{http.MethodGet, "/1", "123", map[string]string{
				"Accept-Encoding": "deflate",
				"Content-Type":    "application/json; charset=utf-8",
			}},
			when{[]Encoding{DeflateEncoding}, &[]string{"test/html"}},
			want{http.StatusOK, generateResponse("123"), http.Header{}, http.Header{
				"Content-Encoding": []string{"deflate"},
			}, nil},
		},
		// request encoding != config encoding
		{"no compress when request encoding != config encodings",
			on{http.MethodGet, "/1", "123", map[string]string{
				"Accept-Encoding": "gzip",
				"Content-Type":    "application/json; charset=utf-8",
			}},
			when{[]Encoding{DeflateEncoding}, nil},
			want{http.StatusOK, generateResponse("123"), http.Header{}, http.Header{
				"Content-Encoding": []string{"gzip", "deflate"},
			}, nil},
		},
		// gzip and deflate
		{"compress gzip and deflate as deflate when ok",
			on{http.MethodGet, "/1", "123", map[string]string{
				"Accept-Encoding": "deflate, gzip",
				"Content-Type":    "application/json; charset=utf-8",
			}},
			when{[]Encoding{GzipEncoding, DeflateEncoding}, nil},
			want{http.StatusOK, generateResponse("123"), http.Header{
				"Content-Encoding": []string{"deflate"},
			}, http.Header{
				"Content-Encoding": []string{"gzip"},
			}, []Encoding{DeflateEncoding, GzipEncoding}},
		},
		{"compress gzip and deflate as gzip when ok",
			on{http.MethodGet, "/1", "123", map[string]string{
				"Accept-Encoding": "deflate, gzip",
				"Content-Type":    "application/json; charset=utf-8",
			}},
			when{[]Encoding{DeflateEncoding, GzipEncoding}, nil},
			want{http.StatusOK, generateResponse("123"), http.Header{
				"Content-Encoding": []string{"gzip"},
			}, http.Header{
				"Content-Encoding": []string{"deflate"},
			}, []Encoding{GzipEncoding, DeflateEncoding}},
		},
		// decompress
		{"gzip encoding request",
			on{http.MethodPost, "/", string(mustCompressGzip(t, []byte("123"))), map[string]string{
				"Content-Type":     "application/json; charset=utf-8",
				"Content-Encoding": "gzip",
			}},
			when{[]Encoding{GzipEncoding}, nil},
			want{http.StatusOK, "", http.Header{}, http.Header{}, nil},
		},
		{"gzip encoding request when no gzip encoding in config",
			on{http.MethodPost, "/", string(mustCompressGzip(t, []byte("123"))), map[string]string{
				"Content-Type":     "application/json; charset=utf-8",
				"Content-Encoding": "gzip",
			}},
			when{[]Encoding{DeflateEncoding}, nil},
			want{http.StatusOK, "", http.Header{}, http.Header{}, nil},
		},
		{"deflate encoding request",
			on{http.MethodPost, "/", string(mustCompressDeflate(t, []byte("123"))), map[string]string{
				"Content-Type":     "application/json; charset=utf-8",
				"Content-Encoding": "deflate",
			}},
			when{[]Encoding{DeflateEncoding}, nil},
			want{http.StatusOK, "", http.Header{}, http.Header{}, nil},
		},
		{"deflate encoding request when no deflate encoding in config",
			on{http.MethodPost, "/", string(mustCompressDeflate(t, []byte("123"))), map[string]string{
				"Content-Type":     "application/json; charset=utf-8",
				"Content-Encoding": "deflate",
			}},
			when{[]Encoding{GzipEncoding}, nil},
			want{http.StatusOK, "", http.Header{}, http.Header{}, nil},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			h := new(mockHandler)
			h.On("ServeHTTP", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
				w := args.Get(0).(http.ResponseWriter)
				w.WriteHeader(tc.want.code)
				_, err := w.Write([]byte(generateResponse(tc.on.body)))
				require.NoError(t, err)
			}).Return()

			router, err := NewCompressorMiddleware(h, CompressorConfig{
				Encodings:    tc.when.encodings,
				ContentTypes: tc.when.contentTypes,
			})
			require.NoError(t, err)

			srv := httptest.NewServer(router)
			defer srv.Close()

			client := resty.New()

			req := client.R()
			req.Method = tc.on.method
			req.URL = srv.URL + tc.on.path
			req.Body = tc.on.body
			req.SetHeaders(tc.on.headers)
			req.SetDoNotParseResponse(true)

			resp, err := req.Send()
			require.NoError(t, err)

			require.Equal(t, tc.want.code, resp.StatusCode(), fmt.Sprintf("expected status code %d but got %d with body: %s", tc.want.code, resp.StatusCode(), string(resp.Body())))

			defer resp.RawBody().Close()

			body, err := io.ReadAll(resp.RawBody())
			require.NoError(t, err)

			if tc.want.body != "" {
				for _, encoding := range tc.want.encodings {
					switch encoding {
					case GzipEncoding:
						responseBody, err := decompressGzip(body)
						require.NoError(t, err)
						body = responseBody
					case DeflateEncoding:
						responseBody, err := decompressDeflate(body)
						require.NoError(t, err)
						body = responseBody
					}
				}
				assert.Equal(t, tc.want.body, string(body))
			}

			for k, v := range tc.want.header {
				values := resp.Header().Values(k)
				require.True(t, len(values) >= len(v), fmt.Sprintf("expected header `%s` to contain %d values but got %d in %v", k, len(v), len(values), values))
				for _, value := range v {
					assert.True(t, slices.ContainsFunc(values, func(s string) bool {
						return strings.Contains(s, value)
					}))
				}
			}
			for k, v := range tc.want.noHeader {
				values := resp.Header().Values(k)
				for _, value := range v {
					assert.False(t, slices.ContainsFunc(values, func(s string) bool {
						return strings.Contains(s, value)
					}))
				}
			}
		})
	}
}

func mustCompressGzip(t *testing.T, data []byte) []byte {
	r, err := compressGzip(data)
	require.NoError(t, err)
	return r
}

func mustCompressDeflate(t *testing.T, data []byte) []byte {
	r, err := compressDeflate(data)
	require.NoError(t, err)
	return r
}

func compressGzip(data []byte) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	zw := gzip.NewWriter(buf)
	_, err := zw.Write(data)
	if err != nil {
		return nil, err
	}
	defer zw.Close()
	return buf.Bytes(), nil
}

func compressDeflate(data []byte) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	zw := zlib.NewWriter(buf)
	_, err := zw.Write(data)
	if err != nil {
		return nil, err
	}
	defer zw.Close()
	return buf.Bytes(), nil
}

func decompressGzip(data []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var b bytes.Buffer
	_, err = b.ReadFrom(r)
	if err != nil {
		return nil, fmt.Errorf("failed decompress data: %v", err)
	}

	return b.Bytes(), nil
}

func decompressDeflate(data []byte) ([]byte, error) {
	r, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var b bytes.Buffer
	_, err = b.ReadFrom(r)
	if err != nil {
		return nil, fmt.Errorf("failed decompress data: %v", err)
	}

	return b.Bytes(), nil
}
