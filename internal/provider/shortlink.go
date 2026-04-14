package provider

import (
	"context"

	"github.com/google/uuid"

	"github.com/liebeSonne/shortlink/internal/model"
)

type ShortLinkProvider interface {
	Find(ctx context.Context, shortID string) (*model.ShortLink, error)
	FindByURL(ctx context.Context, url string) (*model.ShortLink, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]model.ShortLink, error)
}
