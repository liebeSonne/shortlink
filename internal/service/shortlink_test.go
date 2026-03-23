package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/liebeSonne/shortlink/internal/model"
	"github.com/liebeSonne/shortlink/internal/repository"
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
		items      []itemData
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
			on{"https://localhost/1"},
			when{[]itemData{}, "id1"},
			want{nil},
		},
		{
			"empty generated id",
			on{"https://localhost/1"},
			when{[]itemData{}, ""},
			want{model.ErrEmptyID},
		},
		{
			"empty url",
			on{""},
			when{[]itemData{}, "id1"},
			want{model.ErrEmptyURL},
		},
		{
			"invalid url",
			on{"invalid"},
			when{[]itemData{}, "id1"},
			want{model.ErrInvalidURL},
		},
		{
			"err too many generate attempts",
			on{"https://localhost/1"},
			when{
				[]itemData{{"id1", "https://localhost/1"}},
				"id1",
			},
			want{ErrTooManyAttempts},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := repository.NewMemoryShortLinkRepository()
			for _, item := range tc.when.items {
				mockItem := new(mockShortLink)
				mockItem.On("ID").Return(item.id).On("URL").Return(item.url)

				err := repo.Store(mockItem)
				require.NoError(t, err)
			}

			generator := new(mockOneIDGenerator)
			generator.On("GenerateID", mock.Anything).Return(tc.when.generateID)

			service := NewShortLinkService(repo, generator)
			item, err := service.Create(tc.on.url)
			if tc.want.err != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.want.err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, item)
			assert.Equal(t, tc.on.url, item.URL())
			assert.Equal(t, tc.when.generateID, item.ID())
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
