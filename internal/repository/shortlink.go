package repository

import (
	"sync"

	"github.com/liebeSonne/shortlink/internal/model"
)

func NewMemoryShortLinkRepository() model.ShortLinkRepository {
	return &memoryShortLinkRepository{
		linksMap: make(map[string]model.ShortLink),
	}
}

type memoryShortLinkRepository struct {
	linksMap map[string]model.ShortLink
	mu       sync.RWMutex
}

func (s *memoryShortLinkRepository) Get(id string) (model.ShortLink, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	link, ok := s.linksMap[id]
	if !ok {
		return nil, nil
	}
	return link, nil
}

func (s *memoryShortLinkRepository) Store(shortLink model.ShortLink) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.linksMap[shortLink.ID()] = shortLink
	return nil
}
