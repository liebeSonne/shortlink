package filestorage

import (
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

	type itemData struct {
		id  string
		url string
	}
	type on struct {
		url    string
		userID *uuid.UUID
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
			on{url1, nil},
			when{[]userItems{}},
			want{nil, nil},
		},
		{
			"not found when empty url",
			on{"", nil},
			when{[]userItems{{[]itemData{{id1, url1}}, nil}}},
			want{nil, nil},
		},
		{
			"found by url",
			on{url2, nil},
			when{[]userItems{{[]itemData{{id1, url1}, {id2, url2}}, nil}}},
			want{&itemData{id2, url2}, nil},
		},
		{
			"found first by url",
			on{url1, nil},
			when{[]userItems{{[]itemData{{id1, url1}, {id2, url2}, {id2, url1}}, nil}}},
			want{&itemData{id1, url1}, nil},
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

func TestFileShortLinkRepository_Store(t *testing.T) {
	id1 := "id1"
	id2 := "id2"
	url1 := "https://example1.com"
	url2 := "https://example2.com"

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

	testCases := []struct {
		name   string
		items  []model.ShortLink
		userID *uuid.UUID
		err    error
	}{
		{"correct store all items", []model.ShortLink{{ID: id1, URL: url1}, {ID: id2, URL: url2}}, nil, nil},
		{"correct store all with eq id", []model.ShortLink{{ID: id1, URL: url1}, {ID: id1, URL: url2}}, nil, nil},
		{"correct store all with eq url", []model.ShortLink{{ID: id2, URL: url2}, {ID: id1, URL: url2}}, nil, nil},
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
