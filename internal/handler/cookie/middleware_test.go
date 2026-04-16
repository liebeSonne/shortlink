package cookie

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/liebeSonne/shortlink/internal/auth"
)

func TestNewAuthCookieMiddleware(t *testing.T) {
	userID1 := uuid.New()
	defaultCode := http.StatusOK

	type when struct {
		getTokenString    string
		getTokenErr       error
		setTokenErr       error
		parseTokenData    auth.Token
		parseTokenErr     error
		createTokenString string
		createTokenErr    error
	}
	type want struct {
		code int
	}
	testCases := []struct {
		name string
		when when
		want want
	}{
		{
			"get cookie token error",
			when{getTokenErr: errors.New("error")},
			want{http.StatusInternalServerError},
		},
		{
			"parse not empty cookie token string error",
			when{getTokenString: "111", parseTokenErr: errors.New("error")},
			want{http.StatusInternalServerError},
		},
		{
			"set new token on parse not empty cookie token string when error is token is not valid",
			when{getTokenString: "111", parseTokenErr: auth.ErrTokenIsNotValid, createTokenString: "222"},
			want{defaultCode},
		},
		{
			"create token error when not empty cookie token string",
			when{getTokenString: "111", parseTokenErr: auth.ErrTokenIsNotValid, createTokenErr: errors.New("error")},
			want{http.StatusInternalServerError},
		},
		{
			"set auth token error when not empty cookie token string",
			when{getTokenString: "111", parseTokenErr: auth.ErrTokenIsNotValid, createTokenString: "222", setTokenErr: errors.New("error")},
			want{http.StatusInternalServerError},
		},
		{
			"empty user id in parsed token",
			when{getTokenString: "111", parseTokenData: auth.Token{UserID: ""}},
			want{http.StatusUnauthorized},
		},
		{
			"not empty user id in parsed token",
			when{getTokenString: "111", parseTokenData: auth.Token{UserID: "user1"}},
			want{defaultCode},
		},
		{
			"set new token on empty cookie token string",
			when{getTokenString: "", createTokenString: "222"},
			want{defaultCode},
		},
		{
			"create token error when empty cookie token string",
			when{getTokenString: "", createTokenErr: errors.New("error")},
			want{http.StatusInternalServerError},
		},
		{
			"set auth token error when empty cookie token string",
			when{getTokenString: "", createTokenString: "222", setTokenErr: errors.New("error")},
			want{http.StatusInternalServerError},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tokenService := new(mockTokenService)
			tokenService.On("Parse", mock.Anything).Return(tc.when.parseTokenData, tc.when.parseTokenErr)
			tokenService.On("Create", mock.Anything).Return(tc.when.createTokenString, tc.when.createTokenErr)

			cookieService := new(mockService)
			cookieService.On("GetAuthToken", mock.Anything).Return(tc.when.getTokenString, tc.when.getTokenErr)
			lastSetTokenString := ""
			cookieService.On("SetAuthToken", mock.Anything, mock.Anything, mock.Anything).Return(tc.when.setTokenErr).Run(func(args mock.Arguments) {
				lastSetTokenString = args.String(0)
			})

			userService := new(mockUserService)
			userService.On("NextID", mock.Anything).Return(userID1, nil)

			h := new(mockHandler)
			h.On("ServeHTTP", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
				w := args.Get(0).(http.ResponseWriter)
				w.WriteHeader(defaultCode)
				_, err := w.Write([]byte("ok"))
				require.NoError(t, err)
			}).Return()

			handler := NewAuthCookieMiddleware(h, tokenService, cookieService, userService)

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

			if tc.when.createTokenErr == nil && tc.when.setTokenErr == nil {
				assert.Equal(t, tc.when.createTokenString, lastSetTokenString)
			}
		})
	}
}
