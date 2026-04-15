package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUserService_NextID(t *testing.T) {
	userService := NewUserService()

	userID1 := userService.NextID()
	userID2 := userService.NextID()

	require.NotEqual(t, userID1, userID2)
}
