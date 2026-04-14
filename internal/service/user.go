package service

import "github.com/google/uuid"

type UserService interface {
	NextID() uuid.UUID
}

func NewUserService() UserService {
	return &userServiceImpl{}
}

type userServiceImpl struct {
}

func (s *userServiceImpl) NextID() uuid.UUID {
	return uuid.New()
}
