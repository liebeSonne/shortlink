package auth

import (
	"context"

	"github.com/google/uuid"
)

var tokenKey struct{}

func CreateTokenContext(ctx context.Context, token Token) context.Context {
	return context.WithValue(ctx, tokenKey, token)
}

func GetTokenFromContext(ctx context.Context) (Token, bool) {
	token, ok := ctx.Value(tokenKey).(Token)
	if !ok {
		return Token{}, false
	}

	return token, true
}

func GetUserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	tokenPtr, ok := GetTokenFromContext(ctx)
	if !ok || tokenPtr.UserID == "" {
		return uuid.UUID{}, false
	}

	userID, err := uuid.Parse(tokenPtr.UserID)
	if err != nil {
		return uuid.UUID{}, false
	}

	return userID, true
}
