package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRootHandler_Handle(t *testing.T) {
	codeGetResult := http.StatusOK
	codePostResult := http.StatusOK
	getResponse := "get"
	postResponse := "get"

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

	type want struct {
		code int
	}
	testCases := []struct {
		name   string
		method string
		want   want
	}{
		{"get handler", http.MethodGet, want{codeGetResult}},
		{"post handler", http.MethodPost, want{codePostResult}},
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
					wantResponse = getResponse
				}
				if tc.method == http.MethodPost {
					wantResponse = postResponse
				}

				require.NoError(t, err)
				assert.Equal(t, wantResponse, string(resBody))
			}
		})
	}
}

type mockShortLinkHandler struct {
	mock.Mock
}

func (m *mockShortLinkHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}
func (m *mockShortLinkHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}
