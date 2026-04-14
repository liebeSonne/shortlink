package cookie

import (
	"errors"
	"net/http"

	"github.com/liebeSonne/shortlink/internal/auth"
	"github.com/liebeSonne/shortlink/internal/service"
)

var ErrInvalidTokenUserID = errors.New("invalid token user id")

func NewAuthCookieMiddleware(
	next http.Handler,
	tokenService auth.TokenService,
	cookieService Service,
	userService service.UserService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := cookieService.GetAuthToken(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		hasValidCookie := false

		if tokenString != "" {
			tokenData, err := tokenService.Parse(tokenString)
			if err != nil {
				if !errors.Is(err, auth.ErrTokenIsNotValid) {
					http.Error(w, err.Error(), http.StatusInternalServerError)
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
			tokenString, err := tokenService.Create(tokenData)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			err = cookieService.SetAuthToken(tokenString, w)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		next.ServeHTTP(w, r)
	}
}
