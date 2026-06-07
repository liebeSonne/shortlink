package provider

import (
	"context"

	"github.com/google/uuid"

	"github.com/liebeSonne/shortlink/internal/model"
)

// ShortLinkProvider - интерфейс провайдера данных сокращенных ссылок
type ShortLinkProvider interface {
	// Find - поиск данных сокращенной ссылки по названию сокращенной ссылки.
	Find(ctx context.Context, shortID string) (*model.ShortLink, error)
	// FindByURL - поиск данных сокращенной ссылки по ссылке.
	FindByURL(ctx context.Context, url string) (*model.ShortLink, error)
	// FindByUserID - поиск данных сокращенных ссылок по идентификатору пользователя.
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]model.ShortLink, error)
}
