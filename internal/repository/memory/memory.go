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
		linksMap:        make(map[string]model.ShortLink),
		urlToIDMap:      make(map[string]string),
		userIDToIDs:     make(map[uuid.UUID]map[string]bool),
		shortIDToUserID: make(map[string]uuid.UUID),
	}
}

type memoryShortLinkRepository struct {
	linksMap        map[string]model.ShortLink
	urlToIDMap      map[string]string
	userIDToIDs     map[uuid.UUID]map[string]bool
	shortIDToUserID map[string]uuid.UUID
	mu              sync.RWMutex
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

	if idMap, ok := s.userIDToIDs[userID]; ok {
		for id := range idMap {
			if idMap[id] {
				if link, ok := s.linksMap[id]; ok {
					result = append(result, link)
				}
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
			s.userIDToIDs[*userID] = make(map[string]bool)
		}
		s.userIDToIDs[*userID][shortLink.ID] = true
		s.shortIDToUserID[shortLink.ID] = *userID
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
				s.userIDToIDs[*userID] = make(map[string]bool)
			}
			s.userIDToIDs[*userID][shortLink.ID] = true
			s.shortIDToUserID[shortLink.ID] = *userID
		}
	}
	return nil
}

func (s *memoryShortLinkRepository) DeleteByShortIDs(ctx context.Context, shortIDs []string, userID *uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, shortID := range shortIDs {
		if shortLink, ok := s.linksMap[shortID]; ok {
			isUserLink := s.isUserLink(shortLink, userID)
			if !isUserLink {
				continue
			}

			delete(s.urlToIDMap, shortLink.URL)
			delete(s.linksMap, shortID)
			delete(s.shortIDToUserID, shortID)
			if userID != nil {
				if userIDsMap, ok := s.userIDToIDs[*userID]; ok {
					if _, exist := userIDsMap[shortLink.ID]; exist {
						s.userIDToIDs[*userID][shortLink.ID] = false
					}
				}
			}
		}
	}
	return nil
}

func (s *memoryShortLinkRepository) isUserLink(shortLink model.ShortLink, userID *uuid.UUID) bool {
	isUserLink := false

	linkUserID, hasUser := s.shortIDToUserID[shortLink.ID]
	if hasUser && userID != nil && linkUserID == *userID {
		isUserLink = true
	}
	if !hasUser && userID == nil {
		isUserLink = true
	}

	return isUserLink
}
