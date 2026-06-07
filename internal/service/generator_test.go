package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShortIDGenerator_GenerateID(t *testing.T) {
	testCases := []struct {
		name string
		size uint
	}{
		{"size 0", 0},
		{"size > 0", 10},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			generator := NewShortIDGenerator()
			id := generator.GenerateID(tc.size)
			assert.Len(t, id, int(tc.size))
		})
	}
}

func BenchmarkShortIDGenerator_GenerateID(b *testing.B) {
	generator := NewShortIDGenerator()
	for i := 0; i < b.N; i++ {
		generator.GenerateID(16)
	}
}
