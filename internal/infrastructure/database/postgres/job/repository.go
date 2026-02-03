package job

import (
	"context"
	"errors"
	"time"

	"perezvonish/plata-test-assignment/internal/domain/job"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RepositoryImpl struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) job.Repository {
	return &RepositoryImpl{
		pool: pool,
	}
}

func (r *RepositoryImpl) GetById(ctx context.Context, id uuid.UUID) (*job.Job, error) {
	query := `
       SELECT id, quote_id, status, retry_count, idempotency_key, created_at, updated_at 
       FROM update_jobs WHERE id = $1`

	var m Model
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&m.Id, &m.QuoteId, &m.Status, &m.RetryCount, &m.IdempotencyKey, &m.CreatedAt, &m.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return m.MapToDomain(), nil
}

func (r *RepositoryImpl) GetByIdempotencyKey(ctx context.Context, key string) (*job.Job, error) {
	query := `
       SELECT id, quote_id, status, retry_count, idempotency_key, created_at, updated_at 
       FROM update_jobs WHERE idempotency_key = $1`

	var m Model
	err := r.pool.QueryRow(ctx, query, key).Scan(
		&m.Id, &m.QuoteId, &m.Status, &m.RetryCount, &m.IdempotencyKey, &m.CreatedAt, &m.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return m.MapToDomain(), nil
}

func (r *RepositoryImpl) UpdateStatus(ctx context.Context, params job.UpdateStatusParams) error {
	query := `
       UPDATE update_jobs 
       SET status = $1, 
           updated_at = $2,
           retry_count = CASE WHEN $1 = $3 THEN retry_count + 1 ELSE retry_count END
       WHERE id = $4`

	_, err := r.pool.Exec(ctx, query,
		string(params.Status),
		time.Now(),
		string(job.StatusFailure),
		params.Id,
	)
	return err
}
