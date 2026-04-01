package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/liebeSonne/shortlink/internal/model"
	"github.com/liebeSonne/shortlink/internal/repository/memory"
)

func TestShortLinkService_Create(t *testing.T) {
	type itemData struct {
		id  string
		url string
	}
	type on struct {
		url string
	}
	type when struct {
		items       []itemData
		generateID  string
		maxAttempts uint
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
			"valid url",
			on{"https://github.com/shortlink/?q=123"},
			when{[]itemData{}, "id1", 2},
			want{nil},
		},
		{
			"empty generated id",
			on{"https://localhost/1"},
			when{[]itemData{}, "", 2},
			want{ErrTooManyAttempts},
		},
		{
			"empty url",
			on{""},
			when{[]itemData{}, "id1", 2},
			want{ErrEmptyURL},
		},
		{
			"invalid url",
			on{"invalid"},
			when{[]itemData{}, "id1", 2},
			want{ErrInvalidURL},
		},
		{
			"err too many generate attempts",
			on{"https://localhost/1"},
			when{
				[]itemData{{"id1", "https://localhost/1"}},
				"id1",
				2,
			},
			want{ErrTooManyAttempts},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := memory.NewMemoryShortLinkRepository()
			for _, item := range tc.when.items {
				shortLink := model.ShortLink{ID: item.id, URL: item.url}
				err := repo.Store(shortLink)
				require.NoError(t, err)
			}

			generator := new(mockOneIDGenerator)
			generator.On("GenerateID", mock.Anything).Return(tc.when.generateID)

			service := NewShortLinkService(repo, generator, tc.when.maxAttempts)
			item, err := service.Create(tc.on.url)
			if tc.want.err != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.want.err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, item)
			assert.Equal(t, tc.on.url, item.URL)
			assert.Equal(t, tc.when.generateID, item.ID)
		})
	}
}

type mockOneIDGenerator struct {
	mock.Mock
}

func (m *mockOneIDGenerator) GenerateID(size uint) string {
	args := m.Called(size)
	return args.String(0)
}
