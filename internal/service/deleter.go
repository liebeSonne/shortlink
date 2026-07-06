package service

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"

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
	Stop() error
}

// NewShortLinkDeleter - создание экземпляра сервиса отложенного удаления сокращенных ссылок.
func NewShortLinkDeleter(
	ctx context.Context,
	logger logger.Logger,
	handler func(input InputDelete) error,
) ShortLinkDeleter {
	workerCtx, cancel := context.WithCancel(ctx)
	eg, ctx := errgroup.WithContext(workerCtx)

	instance := &deleter{
		ctx:     ctx,
		logger:  logger,
		handler: handler,

		inputCh:    make(chan InputDelete, deleterChanSize),
		errorGroup: eg,
		cancel:     cancel,
	}

	instance.start()

	return instance
}

type deleter struct {
	ctx     context.Context
	logger  logger.Logger
	handler func(input InputDelete) error

	inputCh    chan InputDelete
	errorGroup *errgroup.Group
	cancel     context.CancelFunc
}

func (s *deleter) Add(input InputDelete) error {
	s.logger.Debugf("add to deleter, userID: %v, ids: %v", input.UserID, input.IDs)
	s.inputCh <- input
	return nil
}

func (s *deleter) Stop() error {
	s.cancel()

	close(s.inputCh)

	err := s.errorGroup.Wait()
	if err != nil {
		return fmt.Errorf("error waiting for deleter: %w", err)
	}

	return nil
}

func (s *deleter) start() {
	for i := 0; i < deleterWorkersCount; i++ {
		s.errorGroup.Go(func() error {
			s.worker()
			return nil
		})
	}
}

func (s *deleter) worker() {
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
