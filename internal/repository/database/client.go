package database

import (
	"database/sql"
	"fmt"
	"io"
)

type Client interface {
	DB() *sql.DB

	io.Closer
}

func NewClient(
	dataSourceName string,
) (Client, error) {
	db, err := sql.Open("pgx", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("opening database connection: %w", err)
	}

	return &client{
		db: db,
	}, nil
}

type client struct {
	db *sql.DB
}

func (d *client) DB() *sql.DB {
	return d.db
}

func (d *client) Close() error {
	if d.db != nil {
		err := d.db.Close()
		if err != nil {
			return fmt.Errorf("error closing database connection: %w", err)
		}
		return nil
	}
	return nil
}
