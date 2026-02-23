package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRootHandler_Handle(t *testing.T) {
	codeGetResult := http.StatusOK
	codePostResult := http.StatusOK
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

type mockShortLinkHandler struct {
	mock.Mock
}

func (m *mockShortLinkHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}
func (m *mockShortLinkHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}
