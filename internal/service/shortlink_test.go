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

	type on struct {
		url string
	}
	type when struct {
		items      []model.ShortLink
		generateID string
	}
	type want struct {
		err error
	}
	testCases := []struct {
		name string
		on   on
		when when
		want want
	}{
		{
			"success",
			on{url1},
			when{[]model.ShortLink{}, id1},
			want{nil},
		},
		{
			"empty generated id",
			on{url1},
			when{[]model.ShortLink{}, ""},
			want{model.ErrEmptyID},
		},
		{
			"empty url",
			on{""},
			when{[]model.ShortLink{}, id1},
			want{model.ErrEmptyURL},
		},
		{
			"invalid url",
			on{"invalid"},
			when{[]model.ShortLink{}, id1},
			want{model.ErrInvalidURL},
		},
		{
			"err too many generate attempts",
			on{url1},
			when{[]model.ShortLink{testCreateShortLink(t, id1, url1)}, id1},
			want{ErrTooManyAttempts},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := repository.NewMemoryShortLinkRepository()
			for _, item := range tc.when.items {
				err := repo.Store(item)
				require.NoError(t, err)
			}
			generator := mockOneIDGenerator{generateID: tc.when.generateID}
			service := NewShortLinkService(repo, generator)
			item, err := service.Create(tc.on.url)
			if tc.want.err != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tc.want.err)
			} else {
				require.NotNil(t, item)
				assert.Equal(t, tc.on.url, item.URL())
				assert.Equal(t, tc.when.generateID, item.ID())
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
