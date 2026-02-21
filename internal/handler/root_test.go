package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRootHandler_Handle(t *testing.T) {
	mockHandler := &mockShortLinkHandler{
		t:            t,
		code:         http.StatusOK,
		getResponse:  "get",
		postResponse: "post",
	}

	type want struct {
		code int
	}
	testCases := []struct {
		name   string
		method string
		want   want
	}{
		{"get handler", http.MethodGet, want{mockHandler.code}},
		{"post handler", http.MethodPost, want{mockHandler.code}},
		{"not acceptable head", http.MethodHead, want{http.StatusNotAcceptable}},
		{"not acceptable pur", http.MethodPut, want{http.StatusNotAcceptable}},
		{"not acceptable patch", http.MethodPatch, want{http.StatusNotAcceptable}},
		{"not acceptable connect", http.MethodConnect, want{http.StatusNotAcceptable}},
		{"not acceptable delete", http.MethodDelete, want{http.StatusNotAcceptable}},
		{"not acceptable options", http.MethodOptions, want{http.StatusNotAcceptable}},
		{"not acceptable trace", http.MethodTrace, want{http.StatusNotAcceptable}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handler := NewRootHandler(mockHandler)

			request := httptest.NewRequest(tc.method, "/", strings.NewReader(""))
			w := httptest.NewRecorder()
			handler.Handle(w, request)
			res := w.Result()
			defer func() {
				err := res.Body.Close()
				require.NoError(t, err)
			}()

			require.Equal(t, tc.want.code, res.StatusCode)

			if tc.want.code != http.StatusNotAcceptable {
				resBody, err := io.ReadAll(res.Body)

				wantResponse := ""
				if tc.method == http.MethodGet {
					wantResponse = mockHandler.getResponse
				}
				if tc.method == http.MethodPost {
					wantResponse = mockHandler.postResponse
				}

				require.NoError(t, err)
				assert.Equal(t, wantResponse, string(resBody))
			}
		})
	}
}

type mockShortLinkHandler struct {
	t            *testing.T
	code         int
	getResponse  string
	postResponse string
}

func (m *mockShortLinkHandler) HandleGet(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(m.code)
	_, err := w.Write([]byte(m.getResponse))
	require.NoError(m.t, err)
}
func (m *mockShortLinkHandler) HandleCreate(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(m.code)
	_, err := w.Write([]byte(m.postResponse))
	require.NoError(m.t, err)
}
