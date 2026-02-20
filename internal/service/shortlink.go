package service

import (
	"errors"

	"github.com/liebeSonne/shortlink/internal/model"
)

const maxTry = 1000

var ErrTooManyAttempts = errors.New("too many attempts")

type ShortLinkService interface {
	Create(url string) (model.ShortLink, error)
}

func NewShortLinkService(
	repository model.ShortLinkRepository,
	generator model.ShortIDGenerator,
) ShortLinkService {
	return &shortLinkService{
		repository: repository,
		generator:  generator,
	}
}

type shortLinkService struct {
	repository model.ShortLinkRepository
	generator  model.ShortIDGenerator
}

func (s *shortLinkService) Create(url string) (model.ShortLink, error) {
	id, err := s.nextID()
	if err != nil {
		return nil, err
	}

	item := model.NewShortLink(id, url)

	err = s.repository.Store(item)
	if err != nil {
		return item, err
	}

	return item, nil
}

func (s *shortLinkService) nextID() (string, error) {
	try := 0
	for {
		id := s.generator.GenerateID(model.ShortLinkSize)

		itemPtr, err := s.repository.Get(id)
		if err != nil {
			return "", err
		}

		if itemPtr == nil {
			return id, nil
		}

		try++
		if try > maxTry {
			return id, ErrTooManyAttempts
		}
		continue
	}
}
