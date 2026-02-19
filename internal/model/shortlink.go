package model

import "errors"

var ErrNotFound = errors.New("short link not found")

const ShortLinkSize = 8

type ShortLinkID string

type ShortLink interface {
	ID() ShortLinkID
	URL() string
}

func NewShortLink(
	id ShortLinkID,
	url string,
) ShortLink {
	return &shortLink{
		id:  id,
		url: url,
	}
}

func (s *shortLink) ID() ShortLinkID {
	return s.id
}
func (s *shortLink) URL() string {
	return s.url
}

type shortLink struct {
	id  ShortLinkID
	url string
}

type ShortLinkRepository interface {
	Get(id ShortLinkID) (ShortLink, error)
	Store(shortLink ShortLink) error
}

type ShortLinkProvider interface {
	Find(id string) (*ShortLink, error)
}
