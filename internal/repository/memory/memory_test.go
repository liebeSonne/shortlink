package memory

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/liebeSonne/shortlink/internal/model"
)

func TestShortLinkRepository_Find(t *testing.T) {
	type on struct {
		id string
	}
	type want struct {
		item *model.ShortLink
		err  error
	}
	type userItems struct {
		items  []model.ShortLink
		userID *uuid.UUID
	}
	type when struct {
		userItems []userItems
	}
	testCases := []struct {
		name string
		on   on
		when when
		want want
	}{
		{
			"not found when no items",
			on{"id1"},
			when{[]userItems{}},
			want{nil, nil},
		},
		{
			"not found when empty id",
			on{""},
			when{[]userItems{{[]model.ShortLink{{ID: "id1", URL: "url1"}}, nil}}},
			want{nil, nil},
		},
		{
			"found by id",
			on{"id2"},
			when{[]userItems{{[]model.ShortLink{{ID: "id1", URL: "url1"}, {ID: "id2", URL: "url2"}}, nil}}},
			want{&model.ShortLink{ID: "id2", URL: "url2"}, nil},
		},
		{
			"found last by id",
			on{"id1"},
			when{[]userItems{{[]model.ShortLink{{ID: "id1", URL: "url1"}, {ID: "id2", URL: "url2"}, {ID: "id1", URL: "url2"}}, nil}}},
			want{&model.ShortLink{ID: "id1", URL: "url2"}, nil},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := NewMemoryShortLinkRepository()
			for _, userItem := range tc.when.userItems {
				for _, item := range userItem.items {
					err := repo.Store(t.Context(), item, userItem.userID)
					require.NoError(t, err)
				}
			}
			item, err := repo.Find(t.Context(), tc.on.id)
			if tc.want.err != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.want.err)
				return
			}

			require.NoError(t, err)
			if tc.want.item != nil {
				require.NotNil(t, item)
				assert.Equal(t, tc.want.item.ID, item.ID)
				assert.Equal(t, tc.want.item.URL, item.URL)
			}
		})
	}
}

func TestShortLinkRepository_FindByURL(t *testing.T) {
	type on struct {
		url string
	}
	type want struct {
		item *model.ShortLink
		err  error
	}
	type userItems struct {
		items  []model.ShortLink
		userID *uuid.UUID
	}
	type when struct {
		userItems []userItems
	}
	testCases := []struct {
		name string
		on   on
		when when
		want want
	}{
		{
			"not found when no items",
			on{"url1"},
			when{[]userItems{}},
			want{nil, nil},
		},
		{
			"not found when empty url",
			on{""},
			when{[]userItems{{[]model.ShortLink{{ID: "id1", URL: "url1"}}, nil}}},
			want{nil, nil},
		},
		{
			"found by url",
			on{"url2"},
			when{[]userItems{{[]model.ShortLink{{ID: "id1", URL: "url1"}, {ID: "id2", URL: "url2"}}, nil}}},
			want{&model.ShortLink{ID: "id2", URL: "url2"}, nil},
		},
		{
			"found first by url",
			on{"url2"},
			when{[]userItems{{[]model.ShortLink{{ID: "id1", URL: "url1"}, {ID: "id2", URL: "url2"}, {ID: "id2", URL: "url2"}}, nil}}},
			want{&model.ShortLink{ID: "id2", URL: "url2"}, nil},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := NewMemoryShortLinkRepository()
			for _, userItem := range tc.when.userItems {
				for _, item := range userItem.items {
					err := repo.Store(t.Context(), item, userItem.userID)
					require.NoError(t, err)
				}
			}
			item, err := repo.FindByURL(t.Context(), tc.on.url)
			if tc.want.err != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.want.err)
				return
			}

			require.NoError(t, err)
			if tc.want.item != nil {
				require.NotNil(t, item)
				assert.Equal(t, tc.want.item.ID, item.ID)
				assert.Equal(t, tc.want.item.URL, item.URL)
			}
		})
	}
}

func TestShortLinkRepository_Store(t *testing.T) {
	testCases := []struct {
		name   string
		items  []model.ShortLink
		userID *uuid.UUID
		err    error
	}{
		{"correct store items", []model.ShortLink{{ID: "id1", URL: "url1"}, {ID: "id2", URL: "url2"}}, nil, nil},
		{"correct store with eq id", []model.ShortLink{{ID: "id1", URL: "url1"}, {ID: "id1", URL: "url2"}}, nil, nil},
		{"correct store with eq url", []model.ShortLink{{ID: "id2", URL: "url2"}, {ID: "id1", URL: "url2"}}, nil, nil},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := NewMemoryShortLinkRepository()
			for _, item := range tc.items {
				err := repo.Store(t.Context(), item, tc.userID)
				assert.ErrorIs(t, err, tc.err)
			}
		})
	}
}

func TestShortLinkRepository_StoreAll(t *testing.T) {
	testCases := []struct {
		name   string
		items  []model.ShortLink
		userID *uuid.UUID
		err    error
	}{
		{"correct store all one items", []model.ShortLink{{ID: "id1", URL: "url1"}}, nil, nil},
		{"correct store all many items", []model.ShortLink{{ID: "id1", URL: "url1"}, {ID: "id2", URL: "url2"}}, nil, nil},
		{"correct store all with eq id", []model.ShortLink{{ID: "id1", URL: "url1"}, {ID: "id1", URL: "url2"}}, nil, nil},
		{"correct store all with eq url", []model.ShortLink{{ID: "id2", URL: "url2"}, {ID: "id1", URL: "url2"}}, nil, nil},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := NewMemoryShortLinkRepository()
			err := repo.StoreAll(t.Context(), tc.items, tc.userID)
			assert.ErrorIs(t, err, tc.err)
		})
	}
}
