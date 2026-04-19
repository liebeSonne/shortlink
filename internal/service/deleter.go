package service

import (
	"context"
	"github.com/google/uuid"

	"github.com/liebeSonne/shortlink/internal/logger"
)

const deleterChanSize = 100

type InputDelete struct {
	IDs    []string
	UserID *uuid.UUID
}

type ShortLinkDeleter interface {
	Add(input InputDelete) error
}

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

	go instance.flush()

	return instance
}

type deleter struct {
	ctx     context.Context
	logger  logger.Logger
	handler func(input InputDelete) error

	inputCh chan InputDelete
}

func (s *deleter) Add(input InputDelete) error {
	s.inputCh <- input
	return nil
}

func (s *deleter) flush() {
	for {
		select {
		case <-s.ctx.Done():
			return
		case input := <-s.inputCh:
			if len(input.IDs) == 0 {
				continue
			}
			err := s.handler(input)
			if err != nil {
				s.logger.Errorf("failed to handle: %w", err)
				continue
			}
		}
	}
}
