package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"perezvonish/plata-test-assignment/internal/shared/config"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ConnectInitParams struct {
	Config *config.Config
}

func (p ConnectInitParams) GetDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		p.Config.Postgres.Username,
		p.Config.Postgres.Password,
		p.Config.Postgres.Host,
		p.Config.Postgres.Port,
		p.Config.Postgres.DatabaseName,
	)
}

func ConnectWithRetry(params ConnectInitParams) (*pgxpool.Pool, error) {
	var pool *pgxpool.Pool
	var err error

	dsn := params.GetDSN()
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
			fmt.Println("Connected to Postgres!")

			if err := Migrate(params); err != nil {
				pool.Close()
				return nil, fmt.Errorf("migration error: %w", err)
			}

			return pool, nil
		}

		fmt.Printf("Postgres not ready: %v. Retrying in %v...\n", err, time.Duration(i)*time.Second)
		time.Sleep(time.Duration(i*2) * time.Second)
	}

	return nil, ErrorOutOfRetryCounts
}

func Migrate(params ConnectInitParams) error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	migrationsDir := filepath.Join(wd, "internal", "infrastructure", "database", "postgres", "migrations")

	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		return fmt.Errorf("migrations directory not found at: %s", migrationsDir)
	}

	migrationPath := fmt.Sprintf("file://%s", filepath.ToSlash(migrationsDir))

	dsn := params.GetDSN()
	return RunMigrations(dsn, migrationPath)
}

func RunMigrations(dbUrl string, migrationsPath string) error {
	m, err := migrate.New(migrationsPath, dbUrl)
	if err != nil {
		return fmt.Errorf("could not create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("could not run up migrations: %w", err)
	}

	log.Println("Migrations applied successfully!")
	return nil
}
