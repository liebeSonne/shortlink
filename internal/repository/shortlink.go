package repository

import (
	"github.com/liebeSonne/shortlink/internal/model"
)

func NewMemoryShortLinkRepository() model.ShortLinkRepository {
	return &memoryShortLinkRepository{
		linksMap: make(map[string]model.ShortLink),
	}
}

type memoryShortLinkRepository struct {
	linksMap map[string]model.ShortLink
}

func (s *memoryShortLinkRepository) Get(id string) (model.ShortLink, error) {
	link, ok := s.linksMap[id]
	if !ok {
		return nil, nil
	}
	return link, nil
}

func (s *memoryShortLinkRepository) Store(shortLink model.ShortLink) error {
	s.linksMap[shortLink.ID()] = shortLink
	return nil
}
