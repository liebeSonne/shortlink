package database

import (
	"context"
	"database/sql"

	"github.com/liebeSonne/shortlink/internal/model"
	"github.com/liebeSonne/shortlink/internal/repository"
)

func NewShortLinkRepository(
	db *sql.DB,
) repository.ShortLinkRepository {
	return &shortLinkRepository{
		db: db,
	}
}

type shortLinkRepository struct {
	db *sql.DB
}

func (r *shortLinkRepository) Find(ctx context.Context, shortID string) (*model.ShortLink, error) {
	//TODO implement me
	panic("implement me")
}

func (r *shortLinkRepository) Store(ctx context.Context, shortLink model.ShortLink) error {
	//TODO implement me
	panic("implement me")
}
