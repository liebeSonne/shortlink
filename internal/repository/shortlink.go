package repository

import (
	"context"
	"io"

	"github.com/liebeSonne/shortlink/internal/model"
)

type ShortLinkRepository interface {
	Find(ctx context.Context, shortID string) (*model.ShortLink, error)
	Store(ctx context.Context, shortLink model.ShortLink) error
}

type ShortLinkRepositoryWithCloser interface {
	ShortLinkRepository
	io.Closer
}
