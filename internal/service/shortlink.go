package service

import (
	"context"
	"errors"
	"github.com/avast/retry-go"

	"github.com/liebeSonne/shortlink/internal/model"
	"github.com/liebeSonne/shortlink/internal/repository"
)

const ShortLinkSize = 8
const DefaultMaxAttemptsToGenerateUniqueID = 5

var ErrTooManyAttempts = errors.New("too many attempts to generate unique short id")

var errIDAlreadyExists = errors.New("id already exists")
var errEmptyID = errors.New("empty ID")

type InputShortLinkData struct {
	CorrelationID string
	URL           string
}

type OutputShortLinkData struct {
	CorrelationID string
	ShortLink     model.ShortLink
}

type ShortLinkService interface {
	Create(ctx context.Context, url string) (*model.ShortLink, error)
	CreateBatch(ctx context.Context, urlsData []InputShortLinkData) ([]OutputShortLinkData, error)
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

func (s *shortLinkService) Create(ctx context.Context, url string) (*model.ShortLink, error) {
	err := validateLink(url)
	if err != nil {
		return nil, err
	}

	id, err := s.nextID(ctx, []string{})
	if err != nil {
		return nil, err
	}

	item := model.ShortLink{ID: id, URL: url}

	err = s.repository.Store(ctx, item)
	if err != nil {
		return &item, err
	}

	return &item, nil
}

func (s *shortLinkService) CreateBatch(ctx context.Context, urlsData []InputShortLinkData) ([]OutputShortLinkData, error) {
	for _, urlData := range urlsData {
		err := validateLink(urlData.URL)
		if err != nil {
			return nil, err
		}
	}

	outputURLsData := make([]OutputShortLinkData, 0, len(urlsData))
	exceptions := []string{}

	for _, urlData := range urlsData {
		id, err := s.nextID(ctx, exceptions)
		if err != nil {
			return nil, err
		}

		exceptions = append(exceptions, id)

		item := model.ShortLink{ID: id, URL: urlData.URL}

		outputURLsData = append(outputURLsData, OutputShortLinkData{
			CorrelationID: urlData.CorrelationID,
			ShortLink:     item,
		})
	}

	items := make([]model.ShortLink, 0, len(outputURLsData))
	for _, urlData := range outputURLsData {
		items = append(items, urlData.ShortLink)
	}

	err := s.repository.StoreAll(ctx, items)
	if err != nil {
		return nil, err
	}

	return outputURLsData, nil
}

func (s *shortLinkService) nextID(ctx context.Context, exceptions []string) (string, error) {
	var nextID *string
	var err error

	_ = retry.Do(
		func() error {
			id := s.generator.GenerateID(ShortLinkSize)

			if id == "" {
				return errEmptyID
			}

			for _, exception := range exceptions {
				if exception == id {
					return errIDAlreadyExists
				}
			}

			item, err1 := s.repository.Find(ctx, id)
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
