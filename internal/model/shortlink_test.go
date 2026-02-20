package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewShortLink(t *testing.T) {
	testCases := []struct {
		name string
		id   string
		url  string
		err  error
	}{
		{"valid", "id", "https://github.com/shortlink/?q=123", nil},
		{"empty id", "", "url", ErrEmptyID},
		{"empty url", "id", "", ErrEmptyURL},
		{"invalid url format", "id", "://github.com/", ErrInvalidURL},
		{"invalid url format without schema and host", "id", "/shortlink/?q=123", ErrInvalidURL},
		{"empty url schema", "id", "github.com", ErrInvalidURL},
		{"empty url host", "id", "https://", ErrInvalidURL},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			item, err := NewShortLink(tc.id, tc.url)
			assert.ErrorIs(t, tc.err, err)
			if tc.err == nil {
				assert.Equal(t, tc.id, item.ID())
				assert.Equal(t, tc.url, item.URL())
			}
		})
	}
}
