package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var ErrTokenIsNotValid = errors.New("token is not valid")

type Claims struct {
	jwt.RegisteredClaims
	AuthToken
}

type AuthToken struct {
	UserID string
}

type Service interface {
	Create(tokenData AuthToken) (string, error)
	Parse(tokenString string) (AuthToken, error)
}

func NewService(
	secretKey string,
	tokenExpiry time.Duration,
) Service {
	return &tokenServiceImpl{
		secretKey:   secretKey,
		tokenExpiry: tokenExpiry,
	}
}

type tokenServiceImpl struct {
	secretKey   string
	tokenExpiry time.Duration
}

func (s *tokenServiceImpl) Create(tokenData AuthToken) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.tokenExpiry)),
		},
		AuthToken: tokenData,
	})

	tokenString, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *tokenServiceImpl) Parse(tokenString string) (AuthToken, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(s.secretKey), nil
	})
	if err != nil {
		return AuthToken{}, fmt.Errorf("parse token error: %w", err)
	}

	if !token.Valid {
		return AuthToken{}, ErrTokenIsNotValid
	}

	return claims.AuthToken, nil
}
