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
	Create(url string) (*model.ShortLink, error)
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

func (s *shortLinkService) Create(url string) (*model.ShortLink, error) {
	err := validateLink(url)
	if err != nil {
		return nil, err
	}

	id, err := s.nextID()
	if err != nil {
		return nil, err
	}

	item := model.ShortLink{ID: id, URL: url}

	err = s.repository.Store(item)
	if err != nil {
		return &item, err
	}

	return &item, nil
}

func (s *shortLinkService) nextID() (string, error) {
	var nextID *string
	var err error

	var errIDAlreadyExists = errors.New("id already exists")
	var errEmptyID = errors.New("empty ID")

	_ = retry.Do(
		func() error {
			id := s.generator.GenerateID(ShortLinkSize)

			if id == "" {
				return errEmptyID
			}

			item, err1 := s.repository.Find(id)
			if err1 != nil {
				err = err1
				return err1
			}

			if item != nil {
				return errIDAlreadyExists
			}

			nextID = &id
			return nil
		},
		retry.Attempts(s.maxAttemptsToGenerateUniqueID),
		retry.RetryIf(func(retryErr error) bool {
			return errors.Is(retryErr, errIDAlreadyExists) || errors.Is(err, errEmptyID)
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
