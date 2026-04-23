package cookie

import (
	"errors"
	"net/http"

	"github.com/liebeSonne/shortlink/internal/auth"
	"github.com/liebeSonne/shortlink/internal/logger"
	"github.com/liebeSonne/shortlink/internal/service"
)

var ErrInvalidTokenUserID = errors.New("invalid token user id")

func NewAuthCookieMiddleware(
	next http.Handler,
	tokenService auth.TokenService,
	cookieService Service,
	userService service.UserService,
	logger logger.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := cookieService.GetAuthToken(r)
		if err != nil {
			logger.Errorf("get cookie auth token error: %w", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		hasValidCookie := false

		if tokenString != "" {
			tokenData, err := tokenService.Parse(tokenString)
			if err != nil {
				logger.Errorf("parse token error: %w", err)
				if !errors.Is(err, auth.ErrTokenIsNotValid) {
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
			} else {
				hasValidCookie = true
				if hasValidCookie && tokenData.UserID == "" {
					http.Error(w, ErrInvalidTokenUserID.Error(), http.StatusUnauthorized)
					return
				}
			}
		}

		if !hasValidCookie {
			userID := userService.NextID()
			tokenData := auth.Token{
				UserID: userID.String(),
			}
			tokenString, err = tokenService.Create(tokenData)
			if err != nil {
				logger.Errorf("create token error: %w", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			err = cookieService.SetAuthToken(tokenString, w, r)
			if err != nil {
				logger.Errorf("set cookie auth token error: %w", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
		}

		next.ServeHTTP(w, r)
	}
}
