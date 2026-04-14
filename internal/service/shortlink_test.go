package service

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/liebeSonne/shortlink/internal/model"
	"github.com/liebeSonne/shortlink/internal/repository/memory"
)

func TestShortLinkService_Create(t *testing.T) {
	type on struct {
		url    string
		userID *uuid.UUID
	}
	type userItems struct {
		items  []model.ShortLink
		userID *uuid.UUID
	}
	type when struct {
		userItems   []userItems
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
			on{"https://github.com/shortlink/?q=123", nil},
			when{[]userItems{}, "id1", 2},
			want{nil},
		},
		{
			"empty generated id",
			on{"https://localhost/1", nil},
			when{[]userItems{}, "", 2},
			want{ErrTooManyAttempts},
		},
		{
			"empty url",
			on{"", nil},
			when{[]userItems{}, "id1", 2},
			want{ErrEmptyURL},
		},
		{
			"invalid url",
			on{"invalid", nil},
			when{[]userItems{}, "id1", 2},
			want{ErrInvalidURL},
		},
		{
			"err too many generate attempts",
			on{"https://localhost/1", nil},
			when{
				[]userItems{{[]model.ShortLink{{ID: "id1", URL: "https://localhost/1"}}, nil}},
				"id1",
				2,
			},
			want{ErrTooManyAttempts},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := memory.NewMemoryShortLinkRepository()
			for _, userItem := range tc.when.userItems {
				for _, item := range userItem.items {
					err := repo.Store(t.Context(), item, userItem.userID)
					require.NoError(t, err)
				}
			}

			generator := new(mockOneIDGenerator)
			generator.On("GenerateID", mock.Anything).Return(tc.when.generateID)

			service := NewShortLinkService(repo, generator, tc.when.maxAttempts)
			item, err := service.Create(t.Context(), tc.on.url, tc.on.userID)
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

func TestShortLinkService_CreateBatch(t *testing.T) {
	correlationID1 := "cid1"
	correlationID2 := "cid2"
	link1 := "https://github.com/shortlink/?q=123"
	link2 := "https://localhost/1"
	id1 := "id1"
	id2 := "id2"
	invalidLink := "invalid"

	type on struct {
		inputData []InputShortLinkData
		userID    *uuid.UUID
	}
	type userItems struct {
		items  []model.ShortLink
		userID *uuid.UUID
	}
	type when struct {
		userItems   []userItems
		generateIDs []string
		maxAttempts uint
	}
	type want struct {
		outputData []OutputShortLinkData
		err        error
	}
	testCases := []struct {
		name string
		on   on
		when when
		want want
	}{
		{
			"crete one valid url",
			on{[]InputShortLinkData{
				{correlationID1, link1},
			}, nil},
			when{[]userItems{}, []string{id1}, 2},
			want{[]OutputShortLinkData{
				{correlationID1, model.ShortLink{ID: id1, URL: link1}},
			}, nil},
		},
		{
			"crete many valid url",
			on{[]InputShortLinkData{
				{CorrelationID: correlationID1, URL: link1},
				{CorrelationID: correlationID2, URL: link2},
			}, nil},
			when{[]userItems{}, []string{id1, id2}, 2},
			want{[]OutputShortLinkData{
				{correlationID1, model.ShortLink{ID: id1, URL: link1}},
				{correlationID2, model.ShortLink{ID: id2, URL: link2}},
			}, nil},
		},
		{
			"empty generated id",
			on{[]InputShortLinkData{
				{correlationID1, link1},
			}, nil},
			when{[]userItems{}, []string{""}, 2},
			want{nil, ErrTooManyAttempts},
		},
		{
			"empty url",
			on{[]InputShortLinkData{
				{correlationID1, ""},
			}, nil},
			when{[]userItems{}, []string{id1}, 2},
			want{nil, ErrEmptyURL},
		},
		{
			"invalid url",
			on{[]InputShortLinkData{
				{correlationID1, invalidLink},
			}, nil},
			when{[]userItems{}, []string{id1}, 2},
			want{nil, ErrInvalidURL},
		},
		{
			"err too many generate attempts",
			on{[]InputShortLinkData{
				{correlationID1, link1},
			}, nil},
			when{
				[]userItems{{[]model.ShortLink{{ID: id1, URL: link2}}, nil}},
				[]string{id1},
				2,
			},
			want{nil, ErrTooManyAttempts},
		},
		{
			"err too many generate attempts in many items",
			on{[]InputShortLinkData{
				{CorrelationID: correlationID1, URL: link1},
				{CorrelationID: correlationID2, URL: link2},
			}, nil},
			when{
				[]userItems{},
				[]string{id1, id1},
				2,
			},
			want{nil, ErrTooManyAttempts},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := memory.NewMemoryShortLinkRepository()
			for _, userItem := range tc.when.userItems {
				for _, item := range userItem.items {
					err := repo.Store(t.Context(), item, userItem.userID)
					require.NoError(t, err)
				}
			}

			lastGenerateIndex := 0
			generator := new(mockOneIDGenerator)
			generator.On("GenerateID", mock.Anything).Return(func(_ uint) string {
				if lastGenerateIndex >= len(tc.when.generateIDs) {
					return ""
				}
				id := tc.when.generateIDs[lastGenerateIndex]
				lastGenerateIndex++
				return id
			})

			service := NewShortLinkService(repo, generator, tc.when.maxAttempts)
			outputItems, err := service.CreateBatch(t.Context(), tc.on.inputData, tc.on.userID)
			if tc.want.err != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.want.err)
				return
			}

			require.NoError(t, err)
			require.Len(t, outputItems, len(tc.want.outputData))

			for _, expectOutput := range tc.want.outputData {
				exist := false
				for _, outputItem := range outputItems {
					if expectOutput.CorrelationID == outputItem.CorrelationID {
						exist = true
						assert.Equal(t, expectOutput.CorrelationID, outputItem.CorrelationID)
						assert.Equal(t, expectOutput.ShortLink.URL, outputItem.ShortLink.URL)
						assert.Equal(t, expectOutput.ShortLink.ID, outputItem.ShortLink.ID)
					}
				}
				assert.True(t, exist)
			}
		})
	}
}
