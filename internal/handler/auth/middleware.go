package auth

import (
	"net/http"

	"github.com/liebeSonne/shortlink/internal/auth"
	"github.com/liebeSonne/shortlink/internal/handler/cookie"
	"github.com/liebeSonne/shortlink/internal/logger"
)

func NewAuthMiddleware(
	next http.Handler,
	tokenService auth.TokenService,
	cookieService cookie.Service,
	logger logger.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		tokenString, err := cookieService.GetAuthToken(r)
		if err != nil {
			logger.Errorf("get cookie auth token error: %w", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if tokenString != "" {
			tokenData, err := tokenService.Parse(tokenString)
			if err == nil {
				ctx = auth.CreateTokenContext(ctx, tokenData)
			} else {
				logger.Errorf("parse token error: %w", err)
			}
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
