package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"slices"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/liebeSonne/shortlink/internal/model"
	"github.com/liebeSonne/shortlink/internal/repository"
)

const chunkSize = 500

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
		SELECT short_id, url 
		FROM short_link 
		WHERE short_id = $1 AND deleted_at IS NULL 
		LIMIT 1
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
		SELECT short_id, url 
		FROM short_link 
		WHERE url = $1 AND deleted_at IS NULL
		LIMIT 1
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

func (r *shortLinkRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]model.ShortLink, error) {
	const sqlQuery = `
		SELECT short_id, url 
		FROM short_link 
		WHERE user_id = $1 AND deleted_at IS NULL
	`

	rows, err := r.pool.Query(ctx, sqlQuery, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("error on query: %w", err)
	}

	shortLinks := make([]model.ShortLink, 0)
	for rows.Next() {
		var shortLink model.ShortLink
		err := rows.Scan(&shortLink.ID, &shortLink.URL)
		if err != nil {
			return nil, fmt.Errorf("error on scan row: %w", err)
		}
		shortLinks = append(shortLinks, shortLink)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("error on scan rows: %w", rows.Err())
	}

	return shortLinks, nil
}

func (r *shortLinkRepository) Store(ctx context.Context, shortLink model.ShortLink, userID *uuid.UUID) error {
	return r.StoreAll(ctx, []model.ShortLink{shortLink}, userID)
}

func (r *shortLinkRepository) StoreAll(ctx context.Context, shortLinks []model.ShortLink, userID *uuid.UUID) error {
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
		INSERT INTO short_link (short_id, url, user_id) VALUES ($1, $2, $3)
	`

	stmtName := "insert_short_link"
	_, err = tx.Conn().Prepare(ctx, stmtName, sqlQuery)
	if err != nil {
		return fmt.Errorf("error on prepare statement: %w", err)
	}

	insertErrors := make([]error, 0)
	for _, shortLink := range shortLinks {
		_, err := tx.Exec(ctx, stmtName, shortLink.ID, shortLink.URL, userID)
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

func (r *shortLinkRepository) DeleteByShortIDs(ctx context.Context, shortIDs []string, userID *uuid.UUID) error {
	if len(shortIDs) == 0 {
		return nil
	}

	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})

	const sqlQuery = `
		UPDATE short_link 
		SET deleted_at = NOW() 
	    WHERE 
	        short_id = ANY($1) 
	        AND user_id = $2 
			AND deleted_at IS NULL
	`

	stmtName := "update_short_link"
	_, err = tx.Conn().Prepare(ctx, stmtName, sqlQuery)
	if err != nil {
		return fmt.Errorf("error on prepare statement: %w", err)
	}

	chunkErrors := make([]error, 0)
	for chunkIDs := range slices.Chunk(shortIDs, chunkSize) {
		_, err := tx.Exec(ctx, stmtName, pq.Array(chunkIDs), userID)
		if err != nil {
			chunkErrors = append(chunkErrors, err)
		}
	}

	err = errors.Join(chunkErrors...)
	if err != nil {
		return fmt.Errorf("error on update: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("error on commit transaction: %w", err)
	}

	return nil
}
