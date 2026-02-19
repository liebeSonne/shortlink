package repository

import (
	"github.com/liebeSonne/shortlink/internal/model"
)

func NewMemoryShortLinkRepository() model.ShortLinkRepository {
	return &shortLinkRepository{
		linksMap: make(map[model.ShortLinkID]model.ShortLink),
	}
}

type shortLinkRepository struct {
	linksMap map[model.ShortLinkID]model.ShortLink
}

func (s *shortLinkRepository) Get(id model.ShortLinkID) (model.ShortLink, error) {
	link, ok := s.linksMap[id]
	if !ok {
		return nil, model.ErrNotFound
	}
	return link, nil
}

func (s *shortLinkRepository) Store(shortLink model.ShortLink) error {
	s.linksMap[shortLink.ID()] = shortLink
	return nil
}
