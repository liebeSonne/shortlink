package cookie

import (
	"errors"
	"net/http"
)

// Service - интерфейс сервиса управления авторизацией через cookie.
type Service interface {
	// SetAuthToken - установка токена авторизации в cookie.
	SetAuthToken(tokenString string, w http.ResponseWriter, r *http.Request) error
	// GetAuthToken - получение токена авторизации из cookie.
	GetAuthToken(r *http.Request) (string, error)
}

// NewService - создание экземпляра сервиса управления авторизацией.
func NewService(
	tokenKey string,
) Service {
	return &cookieServiceImpl{
		tokenKey: tokenKey,
	}
}

type cookieServiceImpl struct {
	tokenKey string
}

func (s *cookieServiceImpl) SetAuthToken(tokenString string, w http.ResponseWriter, r *http.Request) error {
	cookie := &http.Cookie{
		Name:  s.tokenKey,
		Value: tokenString,
	}
	http.SetCookie(w, cookie)
	r.AddCookie(cookie)
	return nil
}

func (s *cookieServiceImpl) GetAuthToken(r *http.Request) (string, error) {
	cookie, err := r.Cookie(s.tokenKey)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return "", nil
		}
		return "", err
	}

	if cookie == nil || cookie.Value == "" {
		return "", nil
	}

	return cookie.Value, nil
}
