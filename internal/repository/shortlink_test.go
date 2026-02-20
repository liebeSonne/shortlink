package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/liebeSonne/shortlink/internal/model"
)

func TestShortLinkRepository_Get(t *testing.T) {
	id1 := "id1"
	id2 := "id1"
	url1 := "https://localhost/1"
	url2 := "https://localhost/2"
	itemWithID1AndUrl1, err := model.NewShortLink(id1, url1)
	require.NoError(t, err)
	itemWithID1AndUrl2, err := model.NewShortLink(id1, url2)
	require.NoError(t, err)
	itemWithID2AndUrl2, err := model.NewShortLink(id2, url2)
	require.NoError(t, err)

	testCases := []struct {
		name  string
		items []model.ShortLink
		id    string
		want  *model.ShortLink
		err   error
	}{
		{
			name:  "not found when no items",
			items: []model.ShortLink{},
			id:    id1,
			want:  nil,
			err:   nil,
		},
		{
			name:  "not found when empty id",
			items: []model.ShortLink{itemWithID1AndUrl1},
			id:    "",
			want:  nil,
			err:   nil,
		},
		{
			name:  "found by id",
			items: []model.ShortLink{itemWithID1AndUrl1, itemWithID2AndUrl2},
			id:    id2,
			want:  &itemWithID2AndUrl2,
			err:   nil,
		},
		{
			name:  "found last by id",
			items: []model.ShortLink{itemWithID1AndUrl1, itemWithID2AndUrl2, itemWithID1AndUrl2},
			id:    id1,
			want:  &itemWithID1AndUrl2,
			err:   nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := NewMemoryShortLinkRepository()
			for _, item := range tc.items {
				err := repo.Store(item)
				assert.NoError(t, err)
			}
			itemPtr, err := repo.Get(tc.id)
			if tc.err != nil {
				assert.ErrorIs(t, err, tc.err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.want, itemPtr)
			}
		})
	}
}

func TestShortLinkRepository_Store(t *testing.T) {
	id1 := "id1"
	id2 := "id1"
	url1 := "https://localhost/1"
	url2 := "https://localhost/2"
	itemWithID1AndUrl1, err := model.NewShortLink(id1, url1)
	require.NoError(t, err)
	itemWithID1AndUrl2, err := model.NewShortLink(id1, url2)
	require.NoError(t, err)
	itemWithID2AndUrl2, err := model.NewShortLink(id2, url2)
	require.NoError(t, err)

	testCases := []struct {
		name  string
		items []model.ShortLink
		err   error
	}{
		{"correct store items", []model.ShortLink{itemWithID1AndUrl1, itemWithID2AndUrl2}, nil},
		{"correct store with eq id", []model.ShortLink{itemWithID1AndUrl1, itemWithID1AndUrl2}, nil},
		{"correct store with eq url", []model.ShortLink{itemWithID2AndUrl2, itemWithID1AndUrl2}, nil},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := NewMemoryShortLinkRepository()
			for _, item := range tc.items {
				err := repo.Store(item)
				assert.ErrorIs(t, err, tc.err)
			}
		})
	}
}
