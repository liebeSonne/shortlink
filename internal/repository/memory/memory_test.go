package memory

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/liebeSonne/shortlink/internal/model"
)

func TestShortLinkRepository_Find(t *testing.T) {
	userID1 := uuid.New()

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
		{
			"found by id when created by user",
			on{"id2"},
			when{[]userItems{{[]model.ShortLink{{ID: "id1", URL: "url1"}, {ID: "id2", URL: "url2"}}, &userID1}}},
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
	userID1 := uuid.New()

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
			when{[]userItems{{[]model.ShortLink{{ID: "id1", URL: "url1"}, {ID: "id1", URL: "url2"}, {ID: "id2", URL: "url2"}}, nil}}},
			want{&model.ShortLink{ID: "id2", URL: "url2"}, nil},
		},
		{
			"found by url when created by user",
			on{"url2"},
			when{[]userItems{{[]model.ShortLink{{ID: "id1", URL: "url1"}, {ID: "id2", URL: "url2"}}, &userID1}}},
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

func TestShortLinkRepository_FindByUserID(t *testing.T) {
	userID1 := uuid.New()
	userID2 := uuid.New()

	type on struct {
		userID uuid.UUID
	}
	type want struct {
		items []model.ShortLink
		err   error
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
			on{userID1},
			when{[]userItems{}},
			want{nil, nil},
		},
		{
			"not found when no user items",
			on{userID1},
			when{[]userItems{
				{[]model.ShortLink{{ID: "id1", URL: "url1"}}, nil},
				{[]model.ShortLink{{ID: "id2", URL: "url2"}}, &userID2},
			}},
			want{nil, nil},
		},
		{
			"found by userID",
			on{userID1},
			when{[]userItems{
				{[]model.ShortLink{{ID: "id1", URL: "url1"}, {ID: "id2", URL: "url2"}}, nil},
				{[]model.ShortLink{{ID: "id3", URL: "url3"}, {ID: "id4", URL: "url4"}}, &userID1},
				{[]model.ShortLink{{ID: "id5", URL: "url5"}, {ID: "id6", URL: "url6"}}, &userID2},
			}},
			want{[]model.ShortLink{{ID: "id3", URL: "url3"}, {ID: "id4", URL: "url4"}}, nil},
		},
		{
			"found last ids by userID",
			on{userID1},
			when{[]userItems{
				{[]model.ShortLink{{ID: "id1", URL: "url1"}, {ID: "id2", URL: "url2"}, {ID: "id2", URL: "url2"}}, nil},
				{[]model.ShortLink{{ID: "id1", URL: "url11"}, {ID: "id1", URL: "url12"}, {ID: "id2", URL: "url12"}, {ID: "id3", URL: "url13"}, {ID: "id3", URL: "url133"}}, &userID1},
				{[]model.ShortLink{{ID: "id1", URL: "url21"}, {ID: "id1", URL: "url22"}, {ID: "id2", URL: "url22"}}, nil},
				{[]model.ShortLink{{ID: "id2", URL: "url222"}}, &userID2},
			}},
			want{[]model.ShortLink{{ID: "id1", URL: "url22"}, {ID: "id2", URL: "url222"}, {ID: "id3", URL: "url133"}}, nil},
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
			items, err := repo.FindByUserID(t.Context(), tc.on.userID)
			if tc.want.err != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.want.err)
				return
			}

			require.NoError(t, err)
			require.Len(t, items, len(tc.want.items))
			for _, wantItem := range tc.want.items {
				for _, item := range items {
					if item.ID == wantItem.ID {
						assert.Equal(t, wantItem.ID, item.ID)
						assert.Equal(t, wantItem.URL, item.URL)
					}
				}
			}
		})
	}
}

func TestShortLinkRepository_Store(t *testing.T) {
	userID1 := uuid.New()

	testCases := []struct {
		name   string
		items  []model.ShortLink
		userID *uuid.UUID
		err    error
	}{
		{"correct store items", []model.ShortLink{{ID: "id1", URL: "url1"}, {ID: "id2", URL: "url2"}}, nil, nil},
		{"correct store with eq id", []model.ShortLink{{ID: "id1", URL: "url1"}, {ID: "id1", URL: "url2"}}, nil, nil},
		{"correct store with eq url", []model.ShortLink{{ID: "id2", URL: "url2"}, {ID: "id1", URL: "url2"}}, nil, nil},
		{"correct store by user", []model.ShortLink{{ID: "id1", URL: "url1"}, {ID: "id2", URL: "url2"}}, &userID1, nil},
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
	userID1 := uuid.New()

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
		{"correct store all by user", []model.ShortLink{{ID: "id1", URL: "url1"}, {ID: "id2", URL: "url2"}}, &userID1, nil},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := NewMemoryShortLinkRepository()
			err := repo.StoreAll(t.Context(), tc.items, tc.userID)
			assert.ErrorIs(t, err, tc.err)
		})
	}
}
