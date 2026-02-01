package database

import (
	"context"
	"fmt"
	"perezvonish/plata-test-assignment/internal/shared/config"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ConnectInitParams struct {
	Config *config.Config
}

func ConnectWithRetry(params ConnectInitParams) (*pgxpool.Pool, error) {
	var pool *pgxpool.Pool
	var err error

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		params.Config.Postgres.Username,
		params.Config.Postgres.Password,
		params.Config.Postgres.Host,
		params.Config.Postgres.Port,
		params.Config.Postgres.DatabaseName,
	)

	reconnectAttempts := params.Config.Postgres.ConnectRetryCount

	for i := 1; i <= reconnectAttempts; i++ {
		fmt.Printf("Attempting to connect to Postgres (attempt %d/%d)...\n", i, reconnectAttempts)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		pool, err = pgxpool.New(ctx, dsn)
		if err == nil {
			err = pool.Ping(ctx)
		}
		cancel()

		if err == nil {
			fmt.Printf("Connected to Postgres!")
			return pool, nil
		}

		fmt.Printf("Postgres not ready: %v. Retrying in %v...\n", err, time.Duration(i)*time.Second)
		time.Sleep(time.Duration(i+i+1) * time.Second)
	}

	return nil, ErrorOutOfRetryCounts
}
