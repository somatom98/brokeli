package integration

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type PostgresContainer struct {
	*postgres.PostgresContainer
	DSN string
	DB  *sql.DB
}

func SetupPostgres(ctx context.Context) (*PostgresContainer, error) {
	dbName := "brokeli_test"
	dbUser := "postgres"
	dbPassword := "postgres"

	container, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(15*time.Second)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start container: %w", err)
	}

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("failed to get connection string: %w", err)
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	return &PostgresContainer{
		PostgresContainer: container,
		DSN:               dsn,
		DB:                db,
	}, nil
}

func (c *PostgresContainer) Close(ctx context.Context) error {
	if c.DB != nil {
		c.DB.Close()
	}
	return c.Terminate(ctx)
}
