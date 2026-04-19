package filestorage

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/liebeSonne/shortlink/internal/model"
)

func TestFileShortLinkRepository_Find(t *testing.T) {
	id1 := "id1"
	id2 := "id2"
	url1 := "https://example1.com"
	url2 := "https://example2.com"
	userID1 := uuid.New()

	type itemData struct {
		id  string
		url string
	}
	type on struct {
		id string
	}
	type want struct {
		item *itemData
		err  error
	}
	type userItems struct {
		items  []itemData
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
			on{id1},
			when{[]userItems{}},
			want{nil, nil},
		},
		{
			"not found when empty id",
			on{""},
			when{[]userItems{{[]itemData{{id1, url1}}, nil}}},
			want{nil, nil},
		},
		{
			"found by id",
			on{id2},
			when{[]userItems{{[]itemData{{id1, url1}, {id2, url2}}, nil}}},
			want{&itemData{id2, url2}, nil},
		},
		{
			"found first by id",
			on{id1},
			when{[]userItems{{[]itemData{{id1, url1}, {id2, url2}, {id1, url2}}, nil}}},
			want{&itemData{id1, url1}, nil},
		},
		{
			"found by id when created by user",
			on{id2},
			when{[]userItems{{[]itemData{{id1, url1}, {id2, url2}}, &userID1}}},
			want{&itemData{id2, url2}, nil},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tempDir := t.TempDir()
			filePath := filepath.Join(tempDir, "tmp-1.json")

			repo, err := NewFileShortLinkRepository(filePath)
			require.NoError(t, err)
			t.Cleanup(func() {
				err = repo.Close()
				require.NoError(t, err)
			})

			for _, userItem := range tc.when.userItems {
				for _, item := range userItem.items {
					shortLink := model.ShortLink{ID: item.id, URL: item.url}
					err := repo.Store(t.Context(), shortLink, userItem.userID)
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
				assert.Equal(t, tc.want.item.id, item.ID)
				assert.Equal(t, tc.want.item.url, item.URL)
			}
		})
	}
}

func TestFileShortLinkRepository_FindByURL(t *testing.T) {
	id1 := "id1"
	id2 := "id2"
	url1 := "https://example1.com"
	url2 := "https://example2.com"
	userID1 := uuid.New()

	type itemData struct {
		id  string
		url string
	}
	type on struct {
		url string
	}
	type want struct {
		item *itemData
		err  error
	}
	type userItems struct {
		items  []itemData
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
			on{url1},
			when{[]userItems{}},
			want{nil, nil},
		},
		{
			"not found when empty url",
			on{""},
			when{[]userItems{{[]itemData{{id1, url1}}, nil}}},
			want{nil, nil},
		},
		{
			"found by url",
			on{url2},
			when{[]userItems{{[]itemData{{id1, url1}, {id2, url2}}, nil}}},
			want{&itemData{id2, url2}, nil},
		},
		{
			"found first by url",
			on{url1},
			when{[]userItems{{[]itemData{{id1, url1}, {id2, url2}, {id2, url1}}, nil}}},
			want{&itemData{id1, url1}, nil},
		},
		{
			"found by url when created by user",
			on{url2},
			when{[]userItems{{[]itemData{{id1, url1}, {id2, url2}}, &userID1}}},
			want{&itemData{id2, url2}, nil},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tempDir := t.TempDir()
			filePath := filepath.Join(tempDir, "tmp-1.json")

			repo, err := NewFileShortLinkRepository(filePath)
			require.NoError(t, err)
			t.Cleanup(func() {
				err = repo.Close()
				require.NoError(t, err)
			})

			for _, userItem := range tc.when.userItems {
				for _, item := range userItem.items {
					shortLink := model.ShortLink{ID: item.id, URL: item.url}
					err := repo.Store(t.Context(), shortLink, userItem.userID)
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
				assert.Equal(t, tc.want.item.id, item.ID)
				assert.Equal(t, tc.want.item.url, item.URL)
			}
		})
	}
}

func TestFileShortLinkRepository_FindByUserID(t *testing.T) {
	id1 := "id1"
	id2 := "id2"
	id3 := "id3"
	id4 := "id4"
	url1 := "https://example1.com"
	url2 := "https://example2.com"
	url3 := "https://example3.com"
	url4 := "https://example4.com"
	url5 := "https://example5.com"
	url6 := "https://example6.com"
	userID1 := uuid.New()
	userID2 := uuid.New()

	type itemData struct {
		id  string
		url string
	}
	type on struct {
		userID uuid.UUID
	}
	type want struct {
		items []itemData
		err   error
	}
	type userItems struct {
		items  []itemData
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
				{[]itemData{{id1, url1}}, nil},
				{[]itemData{{id2, url2}}, &userID2},
			}},
			want{nil, nil},
		},
		{
			"found by userID",
			on{userID1},
			when{[]userItems{
				{[]itemData{{id1, url1}, {id2, url4}}, nil},
				{[]itemData{{id1, url2}, {id3, url5}}, &userID1},
				{[]itemData{{id1, url3}, {id4, url6}}, &userID2},
			}},
			want{[]itemData{{id1, url2}, {id3, url5}}, nil},
		},
		{
			"found all by userID",
			on{userID1},
			when{[]userItems{
				{[]itemData{{id1, url1}, {id2, url4}, {id2, url1}}, nil},
				{[]itemData{{id1, url2}, {id3, url5}, {id3, url2}}, &userID1},
				{[]itemData{{id1, url3}, {id4, url6}, {id4, url3}}, &userID2},
			}},
			want{[]itemData{{id1, url2}, {id3, url5}, {id3, url2}}, nil},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tempDir := t.TempDir()
			filePath := filepath.Join(tempDir, "tmp-1.json")

			repo, err := NewFileShortLinkRepository(filePath)
			require.NoError(t, err)
			t.Cleanup(func() {
				err = repo.Close()
				require.NoError(t, err)
			})

			for _, userItem := range tc.when.userItems {
				for _, item := range userItem.items {
					shortLink := model.ShortLink{ID: item.id, URL: item.url}
					err := repo.Store(t.Context(), shortLink, userItem.userID)
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
			countEqual := 0
			for _, wantItem := range tc.want.items {
				for _, item := range items {
					if item.ID == wantItem.id && item.URL == wantItem.url {
						countEqual++
					}
				}
			}
			assert.Equal(t, len(tc.want.items), countEqual)
		})
	}
}

func TestFileShortLinkRepository_Store(t *testing.T) {
	id1 := "id1"
	id2 := "id2"
	url1 := "https://example1.com"
	url2 := "https://example2.com"
	userID1 := uuid.New()

	type itemData struct {
		id  string
		url string
	}
	testCases := []struct {
		name   string
		items  []itemData
		userID *uuid.UUID
		err    error
	}{
		{"correct store items", []itemData{{id1, url1}, {id2, url2}}, nil, nil},
		{"correct store with eq id", []itemData{{id1, url1}, {id1, url2}}, nil, nil},
		{"correct store with eq url", []itemData{{id2, url2}, {id1, url2}}, nil, nil},
		{"correct store items by user", []itemData{{id1, url1}, {id2, url2}}, &userID1, nil},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tempDir := t.TempDir()
			filePath := filepath.Join(tempDir, "tmp-1.json")

			repo, err := NewFileShortLinkRepository(filePath)
			require.NoError(t, err)
			t.Cleanup(func() {
				err = repo.Close()
				require.NoError(t, err)
			})

			for _, item := range tc.items {
				shortLink := model.ShortLink{ID: item.id, URL: item.url}
				err := repo.Store(t.Context(), shortLink, tc.userID)
				assert.ErrorIs(t, err, tc.err)
			}
		})
	}
}

func TestFileShortLinkRepository_StoreAll(t *testing.T) {
	id1 := "id1"
	id2 := "id2"
	url1 := "https://example1.com"
	url2 := "https://example2.com"
	userID1 := uuid.New()

	testCases := []struct {
		name   string
		items  []model.ShortLink
		userID *uuid.UUID
		err    error
	}{
		{"correct store all items", []model.ShortLink{{ID: id1, URL: url1}, {ID: id2, URL: url2}}, nil, nil},
		{"correct store all with eq id", []model.ShortLink{{ID: id1, URL: url1}, {ID: id1, URL: url2}}, nil, nil},
		{"correct store all with eq url", []model.ShortLink{{ID: id2, URL: url2}, {ID: id1, URL: url2}}, nil, nil},
		{"correct store all items by user", []model.ShortLink{{ID: id1, URL: url1}, {ID: id2, URL: url2}}, &userID1, nil},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tempDir := t.TempDir()
			filePath := filepath.Join(tempDir, "tmp-1.json")

			repo, err := NewFileShortLinkRepository(filePath)
			require.NoError(t, err)
			t.Cleanup(func() {
				err = repo.Close()
				require.NoError(t, err)
			})

			err = repo.StoreAll(t.Context(), tc.items, tc.userID)
			assert.ErrorIs(t, err, tc.err)
		})
	}
}

func TestFileShortLinkRepository_DeleteByShortIDs(t *testing.T) {
	userID1 := uuid.New()
	userID2 := uuid.New()

	type on struct {
		shortIDs []string
		userID   *uuid.UUID
	}
	type want struct {
		deletedIDs    []string
		notDeletedIDs []string
		err           error
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
			"empty ids by no user when no items",
			on{[]string{}, nil},
			when{},
			want{[]string{}, nil, nil},
		},
		{
			"empty ids by user when no items",
			on{[]string{}, &userID1},
			when{},
			want{[]string{}, nil, nil},
		},
		{
			"empty ids by no user when has items",
			on{[]string{}, nil},
			when{
				[]userItems{
					{[]model.ShortLink{{ID: "id1", URL: "url1"}, {ID: "id2", URL: "url2"}}, nil},
					{[]model.ShortLink{{ID: "id3", URL: "url3"}, {ID: "id4", URL: "url4"}}, &userID1},
					{[]model.ShortLink{{ID: "id5", URL: "url5"}, {ID: "id6", URL: "url6"}}, &userID2},
				},
			},
			want{[]string{}, []string{"id1", "id2", "id3", "id4", "id5", "id6"}, nil},
		},
		{
			"empty ids by user when has items",
			on{[]string{}, &userID1},
			when{
				[]userItems{
					{[]model.ShortLink{{ID: "id1", URL: "url1"}, {ID: "id2", URL: "url2"}}, nil},
					{[]model.ShortLink{{ID: "id3", URL: "url3"}, {ID: "id4", URL: "url4"}}, &userID1},
					{[]model.ShortLink{{ID: "id5", URL: "url5"}, {ID: "id6", URL: "url6"}}, &userID2},
				},
			},
			want{[]string{}, []string{"id1", "id2", "id3", "id4", "id5", "id6"}, nil},
		},
		{
			"by user",
			on{[]string{"id1", "id3", "id5"}, &userID1},
			when{
				[]userItems{
					{[]model.ShortLink{{ID: "id1", URL: "url1"}, {ID: "id2", URL: "url2"}}, nil},
					{[]model.ShortLink{{ID: "id3", URL: "url3"}, {ID: "id4", URL: "url4"}}, &userID1},
					{[]model.ShortLink{{ID: "id5", URL: "url5"}, {ID: "id6", URL: "url6"}}, &userID2},
				},
			},
			want{[]string{"id3"}, []string{"id1", "id2", "id4", "id5", "id6"}, nil},
		},
		{
			"by no user",
			on{[]string{"id1", "id3", "id5"}, nil},
			when{
				[]userItems{
					{[]model.ShortLink{{ID: "id1", URL: "url1"}, {ID: "id2", URL: "url2"}}, nil},
					{[]model.ShortLink{{ID: "id3", URL: "url3"}, {ID: "id4", URL: "url4"}}, &userID1},
					{[]model.ShortLink{{ID: "id5", URL: "url5"}, {ID: "id6", URL: "url6"}}, &userID2},
				},
			},
			want{[]string{"id1"}, []string{"id2", "id3", "id4", "id5", "id6"}, nil},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := t.Context()
			tempDir := t.TempDir()
			filePath := filepath.Join(tempDir, "tmp-1.json")

			repo, err := NewFileShortLinkRepository(filePath)
			require.NoError(t, err)
			t.Cleanup(func() {
				err = repo.Close()
				require.NoError(t, err)
			})

			for _, userItem := range tc.when.userItems {
				for _, item := range userItem.items {
					err := repo.Store(ctx, item, userItem.userID)
					require.NoError(t, err)
				}
			}

			err = repo.DeleteByShortIDs(ctx, tc.on.shortIDs, tc.on.userID)
			if tc.want.err != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.want.err)
				return
			}

			require.NoError(t, err)

			for _, deletedID := range tc.want.deletedIDs {
				result, err := repo.Find(ctx, deletedID)
				require.NoError(t, err)
				assert.Nil(t, result, fmt.Sprintf("must be deleted id: %s", deletedID))
			}
			for _, notDeletedID := range tc.want.notDeletedIDs {
				result, err := repo.Find(ctx, notDeletedID)
				require.NoError(t, err)
				assert.NotNil(t, result, fmt.Sprintf("must be not deleted id: %s", notDeletedID))
			}
		})
	}
}
