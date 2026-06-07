package database

import (
	"context"
	"fmt"
	"io"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Client - интерфейс клиента к базе данных.
type Client interface {
	// Pool - возвращает пул соединений с базой данных.
	Pool() *pgxpool.Pool

	io.Closer
}

// NewClient - создание экземпляра клиента к базе данных.
func NewClient(
	ctx context.Context,
	dataSourceName string,
) (Client, error) {
	pool, err := pgxpool.New(ctx, dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("erro on pgxpool.New: %w", err)
	}

	return &client{
		pool: pool,
	}, nil
}

type client struct {
	pool *pgxpool.Pool
}

func (d *client) Pool() *pgxpool.Pool {
	return d.pool
}

func (d *client) Close() error {
	d.pool.Close()
	return nil
}
