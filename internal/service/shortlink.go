package service

import (
	"errors"

	"github.com/liebeSonne/shortlink/internal/model"
	"github.com/liebeSonne/shortlink/internal/repository"
)

const maxTry = 1000

var ErrTooManyAttempts = errors.New("too many attempts")

type ShortLinkService interface {
	Create(url string) (model.ShortLink, error)
}

func NewShortLinkService(
	repository repository.ShortLinkRepository,
	generator ShortIDGenerator,
) ShortLinkService {
	return &shortLinkService{
		repository: repository,
		generator:  generator,
	}
}

type shortLinkService struct {
	repository repository.ShortLinkRepository
	generator  ShortIDGenerator
}

func (s *shortLinkService) Create(url string) (model.ShortLink, error) {
	id, err := s.nextID()
	if err != nil {
		return nil, err
	}

	item, err := model.NewShortLink(id, url)
	if err != nil {
		return nil, err
	}

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

		item, err := s.repository.Get(id)
		if err != nil {
			return "", err
		}

		if item == nil {
			return id, nil
		}

		try++
		if try > maxTry {
			return id, ErrTooManyAttempts
		}
		continue
	}
}
