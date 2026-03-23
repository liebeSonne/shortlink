package memory

import "github.com/stretchr/testify/mock"

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
