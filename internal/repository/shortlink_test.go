package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestShortLinkRepository_Get(t *testing.T) {
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
	type when struct {
		items []itemData
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
			when{[]itemData{}},
			want{nil, nil},
		},
		{
			"not found when empty id",
			on{""},
			when{[]itemData{{"id1", "url1"}}},
			want{nil, nil},
		},
		{
			"found by id",
			on{"id2"},
			when{[]itemData{{"id1", "url1"}, {"id2", "url2"}}},
			want{&itemData{"id2", "url2"}, nil},
		},
		{
			"found last by id",
			on{"id1"},
			when{[]itemData{{"id1", "url1"}, {"id2", "url2"}, {"id1", "url2"}}},
			want{&itemData{"id1", "url2"}, nil},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := NewMemoryShortLinkRepository()
			for _, item := range tc.when.items {
				mockItem := new(mockShortLink)
				mockItem.On("ID").Return(item.id).On("URL").Return(item.url)

				err := repo.Store(mockItem)
				require.NoError(t, err)
			}
			item, err := repo.Get(tc.on.id)
			if tc.want.err != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.want.err)
				return
			}

			require.NoError(t, err)
			if tc.want.item != nil {
				require.NotNil(t, item)
				assert.Equal(t, tc.want.item.id, item.ID())
				assert.Equal(t, tc.want.item.url, item.URL())
			}
		})
	}
}

func TestShortLinkRepository_Store(t *testing.T) {
	type itemData struct {
		id  string
		url string
	}
	testCases := []struct {
		name  string
		items []itemData
		err   error
	}{
		{"correct store items", []itemData{{"id1", "url1"}, {"id2", "url2"}}, nil},
		{"correct store with eq id", []itemData{{"id1", "url1"}, {"id1", "url2"}}, nil},
		{"correct store with eq url", []itemData{{"id2", "url2"}, {"id1", "url2"}}, nil},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := NewMemoryShortLinkRepository()
			for _, item := range tc.items {
				mockItem := new(mockShortLink)
				mockItem.On("ID").Return(item.id).On("URL").Return(item.url)

				err := repo.Store(mockItem)
				assert.ErrorIs(t, err, tc.err)
			}
		})
	}
}

type mockShortLink struct {
	mock.Mock
}

func (m *mockShortLink) ID() string {
	args := m.Called()
	return args.String(0)
}
func (m *mockShortLink) URL() string {
	args := m.Called()
	return args.String(0)
}
