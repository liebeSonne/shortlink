package provider

import "github.com/liebeSonne/shortlink/internal/model"

type ShortLinkProvider interface {
	Find(shortID string) (*model.ShortLink, error)
}
