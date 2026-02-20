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
	}{
		{"empty", "", ""},
		{"empty id", "", "url"},
		{"empty url", "id", ""},
		{"uri", "id", "https://github.com/shortlink/"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			item := NewShortLink(tc.id, tc.url)
			assert.Equal(t, tc.id, item.ID())
			assert.Equal(t, tc.url, item.URL())
		})
	}
}
