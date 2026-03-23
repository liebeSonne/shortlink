package repository

import (
	"sync"

	"github.com/liebeSonne/shortlink/internal/model"
)

func NewMemoryShortLinkRepository() ShortLinkRepository {
	return &memoryShortLinkRepository{
		linksMap: make(map[string]model.ShortLink),
	}
}

type memoryShortLinkRepository struct {
	linksMap map[string]model.ShortLink
	mu       sync.RWMutex
}

func (s *memoryShortLinkRepository) Find(shortID string) (*model.ShortLink, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	link, ok := s.linksMap[shortID]
	if !ok {
		return nil, nil
	}
	return &link, nil
}

func (s *memoryShortLinkRepository) Store(shortLink model.ShortLink) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.linksMap[shortLink.ID] = shortLink
	return nil
}
