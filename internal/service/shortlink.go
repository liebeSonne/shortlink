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
	var id model.ShortLinkID

	try := 0
	for {
		shortID := s.generator.GenerateID(model.ShortLinkSize)
		id = model.ShortLinkID(shortID)

		_, err := s.repository.Get(id)

		if !errors.Is(err, model.ErrNotFound) {
			try++
			if try > maxTry {
				return nil, ErrTooManyAttempts
			}
			continue
		}
		break
	}

	item := model.NewShortLink(id, url)
	err := s.repository.Store(item)
	if err != nil {
		return item, err
	}
	return item, nil
}
