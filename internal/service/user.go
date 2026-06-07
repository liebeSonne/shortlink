package service

import "github.com/google/uuid"

// UserService - интерфейс сервиса управления пользователями.
type UserService interface {
	// NextID - вернуть следующий идентификатор пользователя.
	NextID() uuid.UUID
}

// NewUserService - создание экземпляра сервиса управления пользователями.
func NewUserService() UserService {
	return &userServiceImpl{}
}

type userServiceImpl struct {
}

func (s *userServiceImpl) NextID() uuid.UUID {
	return uuid.New()
}
