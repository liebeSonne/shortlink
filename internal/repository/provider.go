package repository

import (
	"errors"

	"github.com/liebeSonne/shortlink/internal/model"
)

func NewRepositoryShortLinkProvider(
	repository model.ShortLinkRepository,
) model.ShortLinkProvider {
	return &repositoryShortLinkProvider{
		repository: repository,
	}
}

type repositoryShortLinkProvider struct {
	repository model.ShortLinkRepository
}

func (p *repositoryShortLinkProvider) Find(id string) (*model.ShortLink, error) {
	shortLinkID := model.ShortLinkID(id)
	item, err := p.repository.Get(shortLinkID)
	if errors.Is(err, model.ErrNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &item, nil
}
