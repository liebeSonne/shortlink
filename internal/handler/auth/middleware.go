package auth

import (
	"net/http"

	"github.com/liebeSonne/shortlink/internal/auth"
	"github.com/liebeSonne/shortlink/internal/handler/cookie"
)

func NewAuthMiddleware(
	next http.Handler,
	tokenService auth.TokenService,
	cookieService cookie.Service,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		tokenString, err := cookieService.GetAuthToken(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if tokenString != "" {
			tokenData, err := tokenService.Parse(tokenString)
			if err == nil {
				ctx = auth.CreateTokenContext(ctx, tokenData)
			}
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
