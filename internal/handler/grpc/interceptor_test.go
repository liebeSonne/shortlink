package grpc

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/liebeSonne/shortlink/internal/auth"
	"github.com/liebeSonne/shortlink/internal/logger"
)

func TestAuthUnaryInterceptor(t *testing.T) {
	userID := uuid.New()
	validToken := auth.Token{UserID: userID.String()}
	parseErr := errors.New("parse error")

	type given struct {
		mdValues []string
	}
	type want struct {
		tokenInCtx bool
		userID     uuid.UUID
		parseCall  bool
		logError   bool
	}
	type when struct {
		parseToken auth.Token
		parseErr   error
	}
	testCases := []struct {
		name  string
		given given
		want  want
		when  when
	}{
		{
			"no metadata",
			given{nil},
			want{false, uuid.UUID{}, false, false},
			when{},
		},
		{
			"no token header",
			given{[]string{}},
			want{false, uuid.UUID{}, false, false},
			when{},
		},
		{
			"empty token header",
			given{[]string{""}},
			want{false, uuid.UUID{}, false, false},
			when{},
		},
		{
			"raw token parsed success",
			given{[]string{"rawtoken"}},
			want{true, userID, true, false},
			when{validToken, nil},
		},
		{
			"bearer token parsed success",
			given{[]string{"Bearer mytoken"}},
			want{true, userID, true, false},
			when{validToken, nil},
		},
		{
			"token parse error",
			given{[]string{"badtoken"}},
			want{false, uuid.UUID{}, true, true},
			when{auth.Token{}, parseErr},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tokenService := auth.NewMockTokenService(t)
			l := logger.NewMockLogger(t)

			if tc.want.parseCall {
				tokenService.EXPECT().Parse(mock.Anything).Return(tc.when.parseToken, tc.when.parseErr)
			}
			if tc.want.logError {
				l.EXPECT().Errorf(mock.Anything, mock.Anything)
			}

			var capturedCtx context.Context
			handler := func(ctx context.Context, req interface{}) (interface{}, error) {
				capturedCtx = ctx
				return nil, nil
			}

			interceptor := AuthUnaryInterceptor("Authorization", tokenService, l)

			ctx := context.Background()
			if tc.given.mdValues != nil {
				md := metadata.New(map[string]string{})
				md["Authorization"] = tc.given.mdValues
				ctx = metadata.NewIncomingContext(ctx, md)
			}

			_, err := interceptor(ctx, nil, &grpc.UnaryServerInfo{}, handler)

			require.NoError(t, err)
			require.NotNil(t, capturedCtx)

			if tc.want.tokenInCtx {
				gotUserID, ok := auth.GetUserIDFromContext(capturedCtx)
				require.True(t, ok)
				assert.Equal(t, tc.want.userID, gotUserID)
			} else {
				_, ok := auth.GetUserIDFromContext(capturedCtx)
				assert.False(t, ok)
			}
		})
	}
}
