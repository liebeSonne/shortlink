package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/liebeSonne/shortlink/internal/model"
	"github.com/stretchr/testify/mock"
)

type mockOneIDGenerator struct {
	mock.Mock
}

func (m *mockOneIDGenerator) GenerateID(size uint) string {
	args := m.Called(size)
	if rf, ok := args.Get(0).(func(uint) string); ok {
		return rf(size)
	}
	return args.String(0)
}

type mockShortLinkRepository struct {
	mock.Mock
}

func (m *mockShortLinkRepository) Find(ctx context.Context, shortID string) (*model.ShortLink, error) {
	args := m.Called(ctx, shortID)
	return args.Get(0).(*model.ShortLink), args.Error(1)
}
func (m *mockShortLinkRepository) FindByURL(ctx context.Context, url string) (*model.ShortLink, error) {
	args := m.Called(ctx, url)
	return args.Get(0).(*model.ShortLink), args.Error(1)
}
func (m *mockShortLinkRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]model.ShortLink, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]model.ShortLink), args.Error(1)
}
func (m *mockShortLinkRepository) Store(ctx context.Context, shortLink model.ShortLink, userID *uuid.UUID) error {
	args := m.Called(ctx, shortLink, userID)
	return args.Error(0)
}
func (m *mockShortLinkRepository) StoreAll(ctx context.Context, shortLinks []model.ShortLink, userID *uuid.UUID) error {
	args := m.Called(ctx, shortLinks, userID)
	return args.Error(0)
}
func (m *mockShortLinkRepository) DeleteByShortIDs(ctx context.Context, shortIDs []string, userID *uuid.UUID) error {
	args := m.Called(ctx, shortIDs, userID)
	return args.Error(0)
}
