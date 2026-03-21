package handler

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

	"github.com/liebeSonne/shortlink/internal/handler/compress"
)

func TestRootHandler_Handle(t *testing.T) {
	codeGetResult := http.StatusOK
	codePostResult := http.StatusCreated
	getResponse := "get"
	postResponse := "post"

	mockHandler := new(mockShortLinkHandler)
	mockHandler.On("HandleGet", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		w := args.Get(0).(http.ResponseWriter)
		w.WriteHeader(codeGetResult)
		_, err := w.Write([]byte(getResponse))
		require.NoError(t, err)
	}).Return()
	mockHandler.On("HandleCreate", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		w := args.Get(0).(http.ResponseWriter)
		w.WriteHeader(codePostResult)
		_, err := w.Write([]byte(postResponse))
		require.NoError(t, err)
	}).Return()
	mockHandler.On("HandleCreateShorten", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		w := args.Get(0).(http.ResponseWriter)
		w.WriteHeader(codePostResult)
		_, err := w.Write([]byte(postResponse))
		require.NoError(t, err)
	}).Return()

	type want struct {
		code int
		body string
	}
	testCases := []struct {
		name   string
		method string
		path   string
		want   want
	}{
		{"get handler", http.MethodGet, "/123", want{codeGetResult, getResponse}},
		{"post handler", http.MethodPost, "/", want{codePostResult, postResponse}},
		{"post api shorten handler", http.MethodPost, "/api/shorten", want{codePostResult, postResponse}},
		{"not head handler", http.MethodHead, "/", want{http.StatusMethodNotAllowed, ""}},
		{"not acceptable pur", http.MethodPut, "/", want{http.StatusMethodNotAllowed, ""}},
		{"not acceptable patch", http.MethodPatch, "/", want{http.StatusMethodNotAllowed, ""}},
		{"not acceptable connect", http.MethodConnect, "/", want{http.StatusMethodNotAllowed, ""}},
		{"not acceptable delete", http.MethodDelete, "/", want{http.StatusMethodNotAllowed, ""}},
		{"not acceptable options", http.MethodOptions, "/", want{http.StatusMethodNotAllowed, ""}},
		{"not acceptable trace", http.MethodTrace, "/", want{http.StatusMethodNotAllowed, ""}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			router := NewRootRouter(mockHandler, false)

			srv := httptest.NewServer(router.Router())
			defer srv.Close()

			client := resty.New()

			req := client.R()
			req.Method = tc.method
			req.URL = srv.URL + tc.path

			resp, err := req.Send()
			require.NoError(t, err)

			require.Equal(t, tc.want.code, resp.StatusCode(), fmt.Sprintf("expected status code %d but got %d with body: %s", tc.want.code, resp.StatusCode(), string(resp.Body())))

			if tc.want.body != "" {
				assert.Equal(t, tc.want.body, string(resp.Body()))
			}
		})
	}
}

func TestRootHandlerWithCompressor(t *testing.T) {
	generateGetResponse := func(requestBody string) string {
		return "GetResponse: " + requestBody
	}
	generateCreateResponse := func(requestBody string) string {
		return "CreateResponse: " + requestBody
	}
	generateCreateShortenResponse := func(requestBody string) string {
		return "CreateShortenResponse: " + requestBody
	}

	type on struct {
		method  string
		path    string
		body    string
		headers map[string]string
	}
	type when struct {
		encodings    []compress.Encoding
		contentTypes *[]string
	}
	type want struct {
		code      int
		body      string
		header    http.Header
		noHeader  http.Header
		encodings []compress.Encoding
	}
	testCases := []struct {
		name string
		on   on
		when when
		want want
	}{
		// all methods
		{"get ok with compress gzip",
			on{http.MethodGet, "/1", "123", map[string]string{
				"Accept-Encoding": "gzip",
				"Content-Type":    "text/html; charset=utf-8",
			}},
			when{[]compress.Encoding{compress.GzipEncoding}, &[]string{"application/json", "text/html"}},
			want{http.StatusOK, generateGetResponse("123"), http.Header{
				"Content-Encoding": []string{"gzip"},
			}, http.Header{}, []compress.Encoding{compress.GzipEncoding}},
		},
		{"get error with compress gzip",
			on{http.MethodGet, "/1", "123", map[string]string{
				"Accept-Encoding": "gzip",
				"Content-Type":    "text/html; charset=utf-8",
			}},
			when{[]compress.Encoding{compress.GzipEncoding}, &[]string{"application/json", "text/html"}},
			want{http.StatusNotFound, "", http.Header{}, http.Header{
				"Content-Encoding": []string{"gzip"},
			}, nil},
		},
		{"create ok with compress gzip",
			on{http.MethodPost, "/", "123", map[string]string{
				"Accept-Encoding": "gzip",
				"Content-Type":    "text/html; charset=utf-8",
			}},
			when{[]compress.Encoding{compress.GzipEncoding}, &[]string{"application/json", "text/html"}},
			want{http.StatusOK, generateCreateResponse("123"), http.Header{
				"Content-Encoding": []string{"gzip"},
			}, http.Header{}, []compress.Encoding{compress.GzipEncoding}},
		},
		{"create error with compress gzip",
			on{http.MethodPost, "/", "123", map[string]string{
				"Accept-Encoding": "gzip",
				"Content-Type":    "text/html; charset=utf-8",
			}},
			when{[]compress.Encoding{compress.GzipEncoding}, &[]string{"application/json", "text/html"}},
			want{http.StatusBadRequest, "", http.Header{}, http.Header{
				"Content-Encoding": []string{"gzip"},
			}, nil},
		},
		{"create shorten ok with compress gzip",
			on{http.MethodPost, "/api/shorten", "123", map[string]string{
				"Accept-Encoding": "gzip",
				"Content-Type":    "application/json; charset=utf-8",
			}},
			when{[]compress.Encoding{compress.GzipEncoding}, &[]string{"application/json", "text/html"}},
			want{http.StatusOK, generateCreateShortenResponse("123"), http.Header{
				"Content-Encoding": []string{"gzip"},
			}, http.Header{}, []compress.Encoding{compress.GzipEncoding}},
		},
		{"create shorten error with compress gzip",
			on{http.MethodPost, "/api/shorten", "123", map[string]string{
				"Accept-Encoding": "gzip",
				"Content-Type":    "application/json; charset=utf-8",
			}},
			when{[]compress.Encoding{compress.GzipEncoding}, &[]string{"application/json", "text/html"}},
			want{http.StatusInternalServerError, "", http.Header{}, http.Header{
				"Content-Encoding": []string{"gzip"},
			}, nil},
		},
		// gzip compress
		{"compress gzip when ok",
			on{http.MethodGet, "/1", "123", map[string]string{
				"Accept-Encoding": "gzip",
				"Content-Type":    "application/json; charset=utf-8",
			}},
			when{[]compress.Encoding{compress.GzipEncoding}, nil},
			want{http.StatusOK, generateGetResponse("123"), http.Header{
				"Content-Encoding": []string{"gzip"},
			}, http.Header{}, []compress.Encoding{compress.GzipEncoding}},
		},
		{"compress gzip when error",
			on{http.MethodGet, "/1", "123", map[string]string{
				"Accept-Encoding": "gzip",
				"Content-Type":    "application/json; charset=utf-8",
			}},
			when{[]compress.Encoding{compress.GzipEncoding}, nil},
			want{http.StatusInternalServerError, "", http.Header{}, http.Header{
				"Content-Encoding": []string{"gzip"},
			}, nil},
		},
		{"compress gzip when unset content types",
			on{http.MethodGet, "/1", "123", map[string]string{
				"Accept-Encoding": "gzip",
				"Content-Type":    "application/json; charset=utf-8",
			}},
			when{[]compress.Encoding{compress.GzipEncoding}, nil},
			want{http.StatusOK, generateGetResponse("123"), http.Header{
				"Content-Encoding": []string{"gzip"},
			}, http.Header{}, []compress.Encoding{compress.GzipEncoding}},
		},
		{"compress gzip when empty set content types",
			on{http.MethodGet, "/1", "123", map[string]string{
				"Accept-Encoding": "gzip",
				"Content-Type":    "application/json; charset=utf-8",
			}},
			when{[]compress.Encoding{compress.GzipEncoding}, &[]string{}},
			want{http.StatusOK, generateGetResponse("123"), http.Header{}, http.Header{
				"Content-Encoding": []string{"gzip"},
			}, nil},
		},
		{"compress gzip when not set request content types",
			on{http.MethodGet, "/1", "123", map[string]string{
				"Accept-Encoding": "gzip",
				"Content-Type":    "application/json; charset=utf-8",
			}},
			when{[]compress.Encoding{compress.GzipEncoding}, &[]string{"test/html"}},
			want{http.StatusOK, generateGetResponse("123"), http.Header{}, http.Header{
				"Content-Encoding": []string{"gzip"},
			}, nil},
		},
		// deflate
		{"compress deflate when ok",
			on{http.MethodGet, "/1", "123", map[string]string{
				"Accept-Encoding": "deflate",
				"Content-Type":    "application/json; charset=utf-8",
			}},
			when{[]compress.Encoding{compress.DeflateEncoding}, nil},
			want{http.StatusOK, generateGetResponse("123"), http.Header{
				"Content-Encoding": []string{"deflate"},
			}, http.Header{}, []compress.Encoding{compress.DeflateEncoding}},
		},
		{"compress deflate when error",
			on{http.MethodGet, "/1", "123", map[string]string{
				"Accept-Encoding": "deflate",
				"Content-Type":    "application/json; charset=utf-8",
			}},
			when{[]compress.Encoding{compress.DeflateEncoding}, nil},
			want{http.StatusInternalServerError, "", http.Header{}, http.Header{
				"Content-Encoding": []string{"deflate"},
			}, nil},
		},
		{"compress deflate when unset content types",
			on{http.MethodGet, "/1", "123", map[string]string{
				"Accept-Encoding": "deflate",
				"Content-Type":    "application/json; charset=utf-8",
			}},
			when{[]compress.Encoding{compress.DeflateEncoding}, nil},
			want{http.StatusOK, generateGetResponse("123"), http.Header{
				"Content-Encoding": []string{"deflate"},
			}, http.Header{}, []compress.Encoding{compress.DeflateEncoding}},
		},
		{"compress deflate when empty set content types",
			on{http.MethodGet, "/1", "123", map[string]string{
				"Accept-Encoding": "deflate",
				"Content-Type":    "application/json; charset=utf-8",
			}},
			when{[]compress.Encoding{compress.DeflateEncoding}, &[]string{}},
			want{http.StatusOK, generateGetResponse("123"), http.Header{}, http.Header{
				"Content-Encoding": []string{"deflate"},
			}, nil},
		},
		{"compress deflate when not set request content types",
			on{http.MethodGet, "/1", "123", map[string]string{
				"Accept-Encoding": "deflate",
				"Content-Type":    "application/json; charset=utf-8",
			}},
			when{[]compress.Encoding{compress.DeflateEncoding}, &[]string{"test/html"}},
			want{http.StatusOK, generateGetResponse("123"), http.Header{}, http.Header{
				"Content-Encoding": []string{"deflate"},
			}, nil},
		},
		// request encoding != config encoding
		{"no compress when request encoding != config encodings",
			on{http.MethodGet, "/1", "123", map[string]string{
				"Accept-Encoding": "gzip",
				"Content-Type":    "application/json; charset=utf-8",
			}},
			when{[]compress.Encoding{compress.DeflateEncoding}, nil},
			want{http.StatusOK, generateGetResponse("123"), http.Header{}, http.Header{
				"Content-Encoding": []string{"gzip", "deflate"},
			}, nil},
		},
		// gzip and deflate
		{"compress gzip and deflate as deflate when ok",
			on{http.MethodGet, "/1", "123", map[string]string{
				"Accept-Encoding": "deflate, gzip",
				"Content-Type":    "application/json; charset=utf-8",
			}},
			when{[]compress.Encoding{compress.GzipEncoding, compress.DeflateEncoding}, nil},
			want{http.StatusOK, generateGetResponse("123"), http.Header{
				"Content-Encoding": []string{"deflate"},
			}, http.Header{
				"Content-Encoding": []string{"gzip"},
			}, []compress.Encoding{compress.DeflateEncoding, compress.GzipEncoding}},
		},
		{"compress gzip and deflate as gzip when ok",
			on{http.MethodGet, "/1", "123", map[string]string{
				"Accept-Encoding": "deflate, gzip",
				"Content-Type":    "application/json; charset=utf-8",
			}},
			when{[]compress.Encoding{compress.DeflateEncoding, compress.GzipEncoding}, nil},
			want{http.StatusOK, generateGetResponse("123"), http.Header{
				"Content-Encoding": []string{"gzip"},
			}, http.Header{
				"Content-Encoding": []string{"deflate"},
			}, []compress.Encoding{compress.GzipEncoding, compress.DeflateEncoding}},
		},
		// decompress
		{"gzip encoding request",
			on{http.MethodPost, "/", string(mustCompressGzip(t, []byte("123"))), map[string]string{
				"Content-Type":     "application/json; charset=utf-8",
				"Content-Encoding": "gzip",
			}},
			when{[]compress.Encoding{compress.GzipEncoding}, nil},
			want{http.StatusOK, "", http.Header{}, http.Header{}, nil},
		},
		{"gzip encoding request when no gzip encoding in config",
			on{http.MethodPost, "/", string(mustCompressGzip(t, []byte("123"))), map[string]string{
				"Content-Type":     "application/json; charset=utf-8",
				"Content-Encoding": "gzip",
			}},
			when{[]compress.Encoding{compress.DeflateEncoding}, nil},
			want{http.StatusOK, "", http.Header{}, http.Header{}, nil},
		},
		{"deflate encoding request",
			on{http.MethodPost, "/", string(mustCompressDeflate(t, []byte("123"))), map[string]string{
				"Content-Type":     "application/json; charset=utf-8",
				"Content-Encoding": "deflate",
			}},
			when{[]compress.Encoding{compress.DeflateEncoding}, nil},
			want{http.StatusOK, "", http.Header{}, http.Header{}, nil},
		},
		{"deflate encoding request when no deflate encoding in config",
			on{http.MethodPost, "/", string(mustCompressDeflate(t, []byte("123"))), map[string]string{
				"Content-Type":     "application/json; charset=utf-8",
				"Content-Encoding": "deflate",
			}},
			when{[]compress.Encoding{compress.GzipEncoding}, nil},
			want{http.StatusOK, "", http.Header{}, http.Header{}, nil},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockHandler := new(mockShortLinkHandler)
			mockHandler.On("HandleGet", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
				w := args.Get(0).(http.ResponseWriter)
				w.WriteHeader(tc.want.code)
				_, err := w.Write([]byte(generateGetResponse(tc.on.body)))
				require.NoError(t, err)
			}).Return()
			mockHandler.On("HandleCreate", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
				w := args.Get(0).(http.ResponseWriter)
				w.WriteHeader(tc.want.code)
				_, err := w.Write([]byte(generateCreateResponse(tc.on.body)))
				require.NoError(t, err)
			}).Return()
			mockHandler.On("HandleCreateShorten", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
				w := args.Get(0).(http.ResponseWriter)
				w.WriteHeader(tc.want.code)
				_, err := w.Write([]byte(generateCreateShortenResponse(tc.on.body)))
				require.NoError(t, err)
			}).Return()

			rootRouter := NewRootRouter(mockHandler, false)
			router := rootRouter.Router().(http.Handler)
			router, err := compress.NewCompressorMiddleware(router, compress.CompressorConfig{
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
					case compress.GzipEncoding:
						responseBody, err := decompressGzip(body)
						require.NoError(t, err)
						body = responseBody
					case compress.DeflateEncoding:
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
