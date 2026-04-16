package auth

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/liebeSonne/shortlink/internal/auth"
)

func TestNewAuthMiddleware(t *testing.T) {
	defaultCode := http.StatusOK

	type when struct {
		getTokenString string
		getTokenErr    error
		parseTokenData auth.Token
		parseTokenErr  error
	}
	type want struct {
		code         int
		contextToken *auth.Token
	}
	testCases := []struct {
		name string
		when when
		want want
	}{
		{
			"error on get cookies token",
			when{getTokenErr: errors.New("error")},
			want{http.StatusInternalServerError, nil},
		},
		{
			"empty cookies token",
			when{getTokenString: ""},
			want{defaultCode, nil},
		},
		{
			"error on parse token",
			when{getTokenString: "111", parseTokenErr: errors.New("error")},
			want{defaultCode, nil},
		},
		{
			"context token",
			when{getTokenString: "111", parseTokenData: auth.Token{UserID: "user1"}},
			want{defaultCode, &auth.Token{UserID: "user1"}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tokenService := new(mockTokenService)
			tokenService.On("Parse", mock.Anything).Return(tc.when.parseTokenData, tc.when.parseTokenErr)

			cookieService := new(mockService)
			cookieService.On("GetAuthToken", mock.Anything).Return(tc.when.getTokenString, tc.when.getTokenErr)

			contextToken := auth.Token{}
			existContextToken := false
			h := new(mockHandler)
			h.On("ServeHTTP", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
				w := args.Get(0).(http.ResponseWriter)
				w.WriteHeader(defaultCode)
				_, err := w.Write([]byte("ok"))
				require.NoError(t, err)

				r := args.Get(1).(*http.Request)
				ctx := r.Context()
				contextToken, existContextToken = auth.GetTokenFromContext(ctx)
			}).Return()

			handler := NewAuthMiddleware(h, tokenService, cookieService)

			srv := httptest.NewServer(handler)
			defer srv.Close()

			client := resty.New()

			req := client.R()
			req.Method = http.MethodGet
			req.URL = srv.URL + "/"
			req.SetDoNotParseResponse(true)

			resp, err := req.Send()
			require.NoError(t, err)

			require.Equal(t, tc.want.code, resp.StatusCode(), fmt.Sprintf("expected status code %d but got %d with body: %s", tc.want.code, resp.StatusCode(), string(resp.Body())))

			require.Equal(t, tc.want.contextToken != nil, existContextToken)
			if tc.want.contextToken != nil {
				require.Equal(t, tc.want.contextToken.UserID, contextToken.UserID)
			}
		})
	}
}
