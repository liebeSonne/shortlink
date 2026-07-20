package grpc

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/liebeSonne/shortlink/internal/auth"
	"github.com/liebeSonne/shortlink/internal/logger"
)

func AuthUnaryInterceptor(
	tokenKey string,
	tokenService auth.TokenService,
	logger logger.Logger,
) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		var tokenString string
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			values := md.Get(tokenKey)
			if len(values) > 0 {
				authHeader := values[0]
				if strings.HasPrefix(authHeader, "Bearer ") {
					tokenString = strings.TrimPrefix(authHeader, "Bearer ")
				} else {
					tokenString = authHeader
				}
			}
		}

		if tokenString != "" {
			tokenData, err := tokenService.Parse(tokenString)
			if err == nil {
				ctx = auth.CreateTokenContext(ctx, tokenData)
			} else {
				logger.Errorf("parse token error: %v", err)
			}
		}

		return handler(ctx, req)
	}
}
