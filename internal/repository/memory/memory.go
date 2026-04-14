package memory

import (
	"context"
	"sync"

	"github.com/google/uuid"

	"github.com/liebeSonne/shortlink/internal/model"
	"github.com/liebeSonne/shortlink/internal/repository"
)

func NewMemoryShortLinkRepository() repository.ShortLinkRepository {
	return &memoryShortLinkRepository{
		linksMap:    make(map[string]model.ShortLink),
		urlToIDMap:  make(map[string]string),
		userIDToIDs: make(map[uuid.UUID][]string),
	}
}

type memoryShortLinkRepository struct {
	linksMap    map[string]model.ShortLink
	urlToIDMap  map[string]string
	userIDToIDs map[uuid.UUID][]string
	mu          sync.RWMutex
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

func (s *memoryShortLinkRepository) FindByUserID(_ context.Context, userID uuid.UUID) ([]model.ShortLink, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]model.ShortLink, 0)

	if ids, ok := s.userIDToIDs[userID]; ok {
		for _, id := range ids {
			if link, ok := s.linksMap[id]; ok {
				result = append(result, link)
			}
		}
	}

	return result, nil
}

func (s *memoryShortLinkRepository) Store(_ context.Context, shortLink model.ShortLink, userID *uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.linksMap[shortLink.ID] = shortLink
	s.urlToIDMap[shortLink.URL] = shortLink.ID
	if userID != nil {
		if _, ok := s.userIDToIDs[*userID]; !ok {
			s.userIDToIDs[*userID] = make([]string, 0)
		}
		s.userIDToIDs[*userID] = append(s.userIDToIDs[*userID], shortLink.ID)
	}
	return nil
}

func (s *memoryShortLinkRepository) StoreAll(_ context.Context, shortLinks []model.ShortLink, userID *uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, shortLink := range shortLinks {
		s.linksMap[shortLink.ID] = shortLink
		s.urlToIDMap[shortLink.URL] = shortLink.ID
		if userID != nil {
			if _, ok := s.userIDToIDs[*userID]; !ok {
				s.userIDToIDs[*userID] = make([]string, 0)
			}
			s.userIDToIDs[*userID] = append(s.userIDToIDs[*userID], shortLink.ID)
		}
	}
	return nil
}
