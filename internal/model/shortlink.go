package model

const ShortLinkSize = 8

type ShortLink interface {
	ID() string
	URL() string
}

func NewShortLink(
	id string,
	url string,
) (ShortLink, error) {
	err := validateShortLinkID(id)
	if err != nil {
		return nil, err
	}
	err = validateLink(url)
	if err != nil {
		return nil, err
	}
	return &shortLink{
		id:  id,
		url: url,
	}, nil
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
	Get(id string) (ShortLink, error)
	Store(shortLink ShortLink) error
}

type ShortLinkProvider interface {
	Get(id string) (ShortLink, error)
}
