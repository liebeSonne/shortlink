package repository

import (
	"context"
	"io"

	"github.com/google/uuid"

	"github.com/liebeSonne/shortlink/internal/model"
)

// ShortLinkRepository - интерфейс репозитория сокращенных ссылок.
type ShortLinkRepository interface {
	// Find - поиск данных сокращенной ссылки по названию сокращенной ссылки.
	Find(ctx context.Context, shortID string) (*model.ShortLink, error)
	// FindByURL - поиск данных сокращенной ссылки по ссылке.
	FindByURL(ctx context.Context, url string) (*model.ShortLink, error)
	// FindByUserID - поиск данных сокращенных ссылок по идентификатору пользователя.
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]model.ShortLink, error)
	// Store - сохранение данных сокращенной ссылки.
	Store(ctx context.Context, shortLink model.ShortLink, userID *uuid.UUID) error
	// StoreAll - сохранение данных массива сокращенных ссылок.
	StoreAll(ctx context.Context, shortLinks []model.ShortLink, userID *uuid.UUID) error
	// DeleteByShortIDs - удаление сокращенных ссылок.
	DeleteByShortIDs(ctx context.Context, shortIDs []string, userID *uuid.UUID) error
	// Stats - статистика
	Stats(ctx context.Context) (model.StatsData, error)
}

// ShortLinkRepositoryWithCloser - интерфейс репозитория сокращенных ссылок с методом закрытия.
type ShortLinkRepositoryWithCloser interface {
	ShortLinkRepository
	io.Closer
}
