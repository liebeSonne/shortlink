package main

import (
	"context"
	"fmt"

	"github.com/liebeSonne/shortlink/internal/config"
	internalio "github.com/liebeSonne/shortlink/internal/io"
	"github.com/liebeSonne/shortlink/internal/repository/database"
)

func createDatabase(cfg config.Config) database.Database {
	databaseDSN := ""
	if cfg.DatabaseDSN != nil {
		databaseDSN = *cfg.DatabaseDSN
	}
	return database.NewDatabase(databaseDSN)
}

func initDatabaseClient(
	ctx context.Context,
	cfg config.Config,
	closer *internalio.MultiCloser,
) (*database.Client, error) {
	if cfg.DatabaseDSN != nil && *cfg.DatabaseDSN != "" {
		client, err := createDatabaseClient(ctx, *cfg.DatabaseDSN, closer)
		if err != nil {
			return nil, err
		}
		return &client, nil
	}
	return nil, nil
}

func createDatabaseClient(
	ctx context.Context,
	databaseDSN string,
	closer *internalio.MultiCloser,
) (database.Client, error) {
	client, err := database.NewClient(ctx, databaseDSN)
	if err != nil {
		return nil, fmt.Errorf("error init db client: %w", err)
	}

	if closer != nil {
		closer.AddCloser(internalio.CloserFunc(
			func() error {
				return client.Close()
			},
		))
	}

	return client, nil
}
