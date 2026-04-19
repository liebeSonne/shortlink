package repository

import (
	"context"
	"io"

	"github.com/google/uuid"

	"github.com/liebeSonne/shortlink/internal/model"
)

type ShortLinkRepository interface {
	Find(ctx context.Context, shortID string) (*model.ShortLink, error)
	FindByURL(ctx context.Context, url string) (*model.ShortLink, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]model.ShortLink, error)
	Store(ctx context.Context, shortLink model.ShortLink, userID *uuid.UUID) error
	StoreAll(ctx context.Context, shortLinks []model.ShortLink, userID *uuid.UUID) error
	DeleteByShortIDs(ctx context.Context, shortIDs []string, userID *uuid.UUID) error
}

type ShortLinkRepositoryWithCloser interface {
	ShortLinkRepository
	io.Closer
}
