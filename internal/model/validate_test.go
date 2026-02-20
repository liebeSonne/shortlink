package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_validateShortLinkID(t *testing.T) {
	testCases := []struct {
		name string
		id   string
		err  error
	}{
		{"valid", "123", nil},
		{"empty", "", ErrEmptyID},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateShortLinkID(tc.id)
			assert.ErrorIs(t, err, tc.err)
		})
	}
}

func Test_validateLink(t *testing.T) {
	testCases := []struct {
		name string
		id   string
		err  error
	}{
		{"valid", "https://github.com/some/path?q=123", nil},
		{"valid with path and params", "https://github.com/some/path?q=123", nil},
		{"valid with port", "https://github.com:8080", nil},
		{"valid with ip and port", "http://127.0.0.1:8080", nil},
		{"empty", "", ErrEmptyURL},
		{"invalid format", "://github.com", ErrInvalidURL},
		{"empty host", "https://", ErrInvalidURL},
		{"without schema", "github.com", ErrInvalidURL},
		{"empty host", "/some/path?q=123", ErrInvalidURL},
		{"invalid", "invalid", ErrInvalidURL},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateLink(tc.id)
			if tc.err == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tc.err)
			}
		})
	}
}
