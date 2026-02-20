package model

const ShortLinkSize = 8

type ShortLink interface {
	ID() string
	URL() string
}

func NewShortLink(
	id string,
	url string,
) ShortLink {
	return &shortLink{
		id:  id,
		url: url,
	}
}

func (s *shortLink) ID() string {
	return s.id
}
func (s *shortLink) URL() string {
	return s.url
}

type shortLink struct {
	id  string
	url string
}

type ShortLinkRepository interface {
	Get(id string) (*ShortLink, error)
	Store(shortLink ShortLink) error
}

type ShortLinkProvider interface {
	Get(id string) (*ShortLink, error)
}
