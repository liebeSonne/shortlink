package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/liebeSonne/shortlink/internal/logger"
)

// deleterChanSize - размер канала удаляемых ссылок.
const deleterChanSize = 5

// deleterWorkersCount - количество обработчиков удаления сокращенных ссылок.
const deleterWorkersCount = 5

// InputDelete - набор данных для удаления сокращенных ссылок.
type InputDelete struct {
	IDs    []string
	UserID *uuid.UUID
}

// ShortLinkDeleter - интерфейс сервиса отложенного удаления сокращенных ссылок.
type ShortLinkDeleter interface {
	// Add - добавление в пул отложено удаляемых сокращенных ссылок.
	Add(input InputDelete) error
}

// NewShortLinkDeleter - создание экземпляра сервиса отложенного удаления сокращенных ссылок.
func NewShortLinkDeleter(
	ctx context.Context,
	logger logger.Logger,
	handler func(input InputDelete) error,
) ShortLinkDeleter {
	instance := &deleter{
		ctx:     ctx,
		logger:  logger,
		handler: handler,

		inputCh: make(chan InputDelete, deleterChanSize),
	}

	for i := 0; i < deleterWorkersCount; i++ {
		go instance.flush()
	}

	return instance
}

type deleter struct {
	ctx     context.Context
	logger  logger.Logger
	handler func(input InputDelete) error

	inputCh chan InputDelete
}

func (s *deleter) Add(input InputDelete) error {
	s.logger.Debugf("add to deleter, userID: %v, ids: %v", input.UserID, input.IDs)
	s.inputCh <- input
	return nil
}

func (s *deleter) flush() {
	for {
		select {
		case <-s.ctx.Done():
			s.logger.Debugf("flush deleter: context closed")
			return
		case input := <-s.inputCh:
			if len(input.IDs) == 0 {
				continue
			}
			s.logger.Debugf("handle deleter for userID: %v, ids: %v", input.UserID, input.IDs)
			err := s.handler(input)
			if err != nil {
				s.logger.Errorf("failed to handle: %v", err)
				continue
			}
		}
	}
}
