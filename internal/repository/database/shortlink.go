package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/liebeSonne/shortlink/internal/model"
	"github.com/liebeSonne/shortlink/internal/repository"
)

func NewShortLinkRepository(
	pool *pgxpool.Pool,
) repository.ShortLinkRepository {
	return &shortLinkRepository{
		pool: pool,
	}
}

type shortLinkRepository struct {
	pool *pgxpool.Pool
}

func (r *shortLinkRepository) Find(ctx context.Context, shortID string) (*model.ShortLink, error) {
	const sqlQuery = `
		SELECT short_id, url FROM short_link WHERE short_id = $1 LIMIT 1
	`

	var shortLink model.ShortLink
	row := r.pool.QueryRow(ctx, sqlQuery, shortID)
	err := row.Scan(&shortLink.ID, &shortLink.URL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("error on scan row: %w", err)
	}

	return &shortLink, nil
}

func (r *shortLinkRepository) FindByURL(ctx context.Context, url string) (*model.ShortLink, error) {
	const sqlQuery = `
		SELECT short_id, url FROM short_link WHERE url = $1 LIMIT 1
	`

	var shortLink model.ShortLink
	row := r.pool.QueryRow(ctx, sqlQuery, url)
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
	return r.StoreAll(ctx, []model.ShortLink{shortLink})
}

func (r *shortLinkRepository) StoreAll(ctx context.Context, shortLinks []model.ShortLink) error {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("error on begin transaction: %w", err)
	}

	defer func() {
		err = tx.Rollback(ctx)
		if err != nil {
			fmt.Printf("error on rollback transaction: %v\n", err)
		}
	}()

	const sqlQuery = `
		INSERT INTO short_link (short_id, url) VALUES ($1, $2)
	`

	stmtName := "insert_short_link"
	_, err = tx.Conn().Prepare(ctx, stmtName, sqlQuery)
	if err != nil {
		return fmt.Errorf("error on prepare statement: %w", err)
	}

	insertErrors := make([]error, 0)
	for _, shortLink := range shortLinks {
		_, err := tx.Exec(ctx, stmtName, shortLink.ID, shortLink.URL)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgerrcode.UniqueViolation == pgErr.Code {
				err = repository.NewErrConflictURL(shortLink.URL, err)
			}
			insertErrors = append(insertErrors, err)
		}
	}

	err = errors.Join(insertErrors...)
	if err != nil {
		return fmt.Errorf("error on insert: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("error on commit transaction: %w", err)
	}

	return nil
}
