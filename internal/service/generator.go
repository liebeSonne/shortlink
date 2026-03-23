package service

import (
	"math/rand"
)

type ShortIDGenerator interface {
	GenerateID(size uint) string
}

func NewShortIDGenerator() ShortIDGenerator {
	return &shortIDGenerator{
		symbols: []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"),
	}
}

type shortIDGenerator struct {
	symbols []rune
}

func (s *shortIDGenerator) GenerateID(size uint) string {
	id := make([]rune, size)
	for i := range size {
		index := rand.Intn(len(s.symbols))
		id[i] = s.symbols[index]
	}
	return string(id)
}
