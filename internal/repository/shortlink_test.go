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

	type on struct {
		id string
	}
	type want struct {
		item *model.ShortLink
		err  error
	}
	type when struct {
		items []model.ShortLink
	}
	testCases := []struct {
		name string
		on   on
		when when
		want want
	}{
		{
			"not found when no items",
			on{id1},
			when{[]model.ShortLink{}},
			want{nil, nil},
		},
		{
			"not found when empty id",
			on{""},
			when{[]model.ShortLink{itemWithID1AndUrl1}},
			want{nil, nil},
		},
		{
			"found by id",
			on{id2},
			when{[]model.ShortLink{itemWithID1AndUrl1, itemWithID2AndUrl2}},
			want{&itemWithID2AndUrl2, nil},
		},
		{
			"found last by id",
			on{id1},
			when{[]model.ShortLink{itemWithID1AndUrl1, itemWithID2AndUrl2, itemWithID1AndUrl2}},
			want{&itemWithID1AndUrl2, nil},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := NewMemoryShortLinkRepository()
			for _, item := range tc.when.items {
				err := repo.Store(item)
				assert.NoError(t, err)
			}
			itemPtr, err := repo.Get(tc.on.id)
			if tc.want.err != nil {
				assert.ErrorIs(t, err, tc.want.err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.want.item, itemPtr)
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
