package service

import "github.com/stretchr/testify/mock"

type mockOneIDGenerator struct {
	mock.Mock
}

func (m *mockOneIDGenerator) GenerateID(size uint) string {
	args := m.Called(size)
	if rf, ok := args.Get(0).(func(uint) string); ok {
		return rf(size)
	}
	return args.String(0)
}
