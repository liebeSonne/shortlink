package memory

import (
	"context"
	"sync"

	"github.com/liebeSonne/shortlink/internal/model"
	"github.com/liebeSonne/shortlink/internal/repository"
)

func NewMemoryShortLinkRepository() repository.ShortLinkRepository {
	return &memoryShortLinkRepository{
		linksMap:   make(map[string]model.ShortLink),
		urlToIDMap: make(map[string]string),
	}
}

type memoryShortLinkRepository struct {
	linksMap   map[string]model.ShortLink
	urlToIDMap map[string]string
	mu         sync.RWMutex
}

func (s *memoryShortLinkRepository) Find(_ context.Context, shortID string) (*model.ShortLink, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	link, ok := s.linksMap[shortID]
	if !ok {
		return nil, nil
	}
	return &link, nil
}

func (s *memoryShortLinkRepository) FindByURL(_ context.Context, url string) (*model.ShortLink, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if id, ok := s.urlToIDMap[url]; ok {
		link, ok := s.linksMap[id]
		if !ok {
			return nil, nil
		}
		return &link, nil
	}
	return nil, nil
}

func (s *memoryShortLinkRepository) Store(_ context.Context, shortLink model.ShortLink) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.linksMap[shortLink.ID] = shortLink
	s.urlToIDMap[shortLink.URL] = shortLink.ID
	return nil
}

func (s *memoryShortLinkRepository) StoreAll(_ context.Context, shortLinks []model.ShortLink) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, shortLink := range shortLinks {
		s.linksMap[shortLink.ID] = shortLink
		s.urlToIDMap[shortLink.URL] = shortLink.ID
	}
	return nil
}
