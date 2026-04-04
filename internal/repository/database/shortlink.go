package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

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
	const sqlQuery = `
		SELECT short_id, url FROM short_link WHERE short_id = $1
	`

	var shortLink model.ShortLink
	row := r.db.QueryRowContext(ctx, sqlQuery, shortID)
	err := row.Scan(&shortLink.ID, &shortLink.URL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("error on scan row: %w", err)
	}

	return &shortLink, nil
}

func (r *shortLinkRepository) Store(ctx context.Context, shortLink model.ShortLink) error {
	const sqlQuery = `
		INSERT INTO short_link (short_id, url) VALUES ($1, $2)
	`

	_, err := r.db.ExecContext(ctx, sqlQuery, shortLink.ID, shortLink.URL)
	if err != nil {
		return fmt.Errorf("error on insert: %w", err)
	}

	return nil
}
