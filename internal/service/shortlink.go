package service

import (
	"errors"

	"github.com/avast/retry-go"

	"github.com/liebeSonne/shortlink/internal/model"
	"github.com/liebeSonne/shortlink/internal/repository"
)

const ShortLinkSize = 8
const DefaultMaxAttemptsToGenerateUniqueID = 5

var ErrTooManyAttempts = errors.New("too many attempts to generate unique short id")

type ShortLinkService interface {
	Create(url string) (model.ShortLink, error)
}

func NewShortLinkService(
	repository repository.ShortLinkRepository,
	generator ShortIDGenerator,
	maxAttemptsToGenerateUniqueID uint,
) ShortLinkService {
	return &shortLinkService{
		repository:                    repository,
		generator:                     generator,
		maxAttemptsToGenerateUniqueID: maxAttemptsToGenerateUniqueID,
	}
}

type shortLinkService struct {
	repository                    repository.ShortLinkRepository
	generator                     ShortIDGenerator
	maxAttemptsToGenerateUniqueID uint
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
	var nextID *string
	var err error

	var errIDAlreadyExists = errors.New("id already exists")

	_ = retry.Do(
		func() error {
			id := s.generator.GenerateID(ShortLinkSize)
			item, err := s.repository.Get(id)
			if err != nil {
				return err
			}

			if item != nil {
				return errIDAlreadyExists
			}

			nextID = &id
			return nil
		},
		retry.Attempts(s.maxAttemptsToGenerateUniqueID),
		retry.RetryIf(func(retryErr error) bool {
			return errors.Is(retryErr, errIDAlreadyExists)
		}),
	)

	if err != nil {
		return "", err
	}

	if nextID != nil {
		return *nextID, nil
	}

	return "", ErrTooManyAttempts
}
