package service

import (
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/liebeSonne/shortlink/internal/logger"
)

func BenchmarkDeleter_Add(b *testing.B) {
	l := logger.NewMockLogger(b)
	l.EXPECT().Errorf(mock.Anything, mock.Anything).Maybe()
	l.EXPECT().Debugf(mock.Anything, mock.Anything).Maybe()

	handler := func(input InputDelete) error {
		_ = input
		return nil
	}

	shortLinkDeleter := NewShortLinkDeleter(b.Context(), l, handler)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		input := InputDelete{}
		_ = shortLinkDeleter.Add(input)
	}
}
