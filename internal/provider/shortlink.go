package provider

import (
	"context"
	"github.com/liebeSonne/shortlink/internal/model"
)

type ShortLinkProvider interface {
	Find(ctx context.Context, shortID string) (*model.ShortLink, error)
}
