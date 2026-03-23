package repository

import (
	"io"

	"github.com/liebeSonne/shortlink/internal/model"
)

type ShortLinkRepository interface {
	Get(id string) (model.ShortLink, error)
	Store(shortLink model.ShortLink) error
}

type ShortLinkRepositoryWithCloser interface {
	ShortLinkRepository
	io.Closer
}
