package provider

import "github.com/liebeSonne/shortlink/internal/model"

type ShortLinkProvider interface {
	Get(id string) (model.ShortLink, error)
}
