package database

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

//go:generate go run github.com/gojuno/minimock/v3/cmd/minimock -i Database -o ../../mocks/mock_database.go -g

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
	db, err := sql.Open("pgx", d.dataSourceName)
	if err != nil {
		return fmt.Errorf("opening database connection: %w", err)
	}

	defer func() {
		err = db.Close()
		if err != nil {
			fmt.Printf("error closing database connection: %v", err)
		}
	}()

	err = db.PingContext(ctx)
	if err != nil {
		return fmt.Errorf("pinging database: %w", err)
	}

	return nil
}
