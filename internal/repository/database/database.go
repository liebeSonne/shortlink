package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Database interface {
	Ping(ctx context.Context) error
}

func NewDatabase(
	dataSourceName string,
) Database {
	return &database{
		dataSourceName: dataSourceName,
	}
}

type database struct {
	dataSourceName string
}

func (d *database) Ping(ctx context.Context) error {
	pool, err := pgxpool.New(ctx, d.dataSourceName)
	if err != nil {
		return fmt.Errorf("erro on pgxpool.New: %w", err)
	}

	defer pool.Close()

	err = pool.Ping(ctx)
	if err != nil {
		return fmt.Errorf("pinging database: %w", err)
	}

	return nil
}
