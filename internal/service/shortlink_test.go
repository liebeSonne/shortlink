package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/liebeSonne/shortlink/internal/model"
	"github.com/liebeSonne/shortlink/internal/repository"
)

func TestShortLinkService_Create(t *testing.T) {
	id1 := "id1"
	url1 := "https://localhost/1"

	testCases := []struct {
		name       string
		items      []model.ShortLink
		url        string
		generateID string
		err        error
	}{
		{
			"success",
			[]model.ShortLink{},
			url1,
			id1,
			nil,
		},
		{
			"empty generated id",
			[]model.ShortLink{},
			url1,
			"",
			model.ErrEmptyID,
		},
		{
			"empty url",
			[]model.ShortLink{},
			"",
			id1,
			model.ErrEmptyURL,
		},
		{
			"invalid url",
			[]model.ShortLink{},
			"invalid",
			id1,
			model.ErrInvalidURL,
		},
		{
			"err too many generate attempts",
			[]model.ShortLink{testCreateShortLink(t, id1, url1)},
			url1,
			id1,
			ErrTooManyAttempts,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := repository.NewMemoryShortLinkRepository()
			for _, item := range tc.items {
				err := repo.Store(item)
				require.NoError(t, err)
			}
			generator := mockOneIDGenerator{generateID: tc.generateID}
			service := NewShortLinkService(repo, generator)
			item, err := service.Create(tc.url)
			if tc.err != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tc.err)
			} else {
				require.NotNil(t, item)
				assert.Equal(t, tc.url, item.URL())
				assert.Equal(t, tc.generateID, item.ID())
			}
		})
	}
}

type mockOneIDGenerator struct {
	generateID string
}

func (m mockOneIDGenerator) GenerateID(_ uint) string {
	return m.generateID
}

func testCreateShortLink(t *testing.T, id string, url string) model.ShortLink {
	item, err := model.NewShortLink(id, url)
	require.NoError(t, err)
	return item
}
