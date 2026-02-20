package repository

import (
	"github.com/liebeSonne/shortlink/internal/model"
)

func NewMemoryShortLinkRepository() model.ShortLinkRepository {
	return &shortLinkRepository{
		linksMap: make(map[string]model.ShortLink),
	}
}

type shortLinkRepository struct {
	linksMap map[string]model.ShortLink
}

func (s *shortLinkRepository) Get(id string) (*model.ShortLink, error) {
	link, ok := s.linksMap[id]
	if !ok {
		return nil, nil
	}
	return &link, nil
}

func (s *shortLinkRepository) Store(shortLink model.ShortLink) error {
	s.linksMap[shortLink.ID()] = shortLink
	return nil
}
