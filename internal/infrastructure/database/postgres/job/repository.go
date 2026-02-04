package job

import (
	"context"
	"errors"
	"fmt"
	"time"

	"perezvonish/plata-test-assignment/internal/domain/job"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
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

func (r *RepositoryImpl) Save(ctx context.Context, params job.CreateParams) (uuid.UUID, error) {
	jobId := uuid.New()
	now := time.Now().UTC()
	status := "pending"
	var retryCount int64 = 0

	const query = `
        INSERT INTO update_jobs (
            id, 
            quote_id, 
            status, 
            retry_count, 
            idempotency_key, 
            created_at, 
            updated_at
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        ON CONFLICT (idempotency_key) DO NOTHING
    `

	_, err := r.pool.Exec(ctx, query,
		jobId,
		params.QuoteID,
		status,
		retryCount,
		params.IdempotencyKey,
		now,
		now,
	)

	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to save job: %w", err)
	}

	return jobId, nil
}

func (r *RepositoryImpl) GetById(ctx context.Context, id uuid.UUID) (*job.Job, error) {
	query := `
		SELECT id, quote_id, status, retry_count, idempotency_key, created_at, updated_at
		FROM update_jobs
		WHERE id = $1
	`

	var m Model

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&m.Id,
		&m.QuoteId,
		&m.Status,
		&m.RetryCount,
		&m.IdempotencyKey,
		&m.CreatedAt,
		&m.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to scan job by id: %w", err)
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
	var pgID pgtype.UUID
	if err := pgID.Scan(params.Id.String()); err != nil {
		return fmt.Errorf("failed to scan uuid for pgx: %w", err)
	}

	query := `
       UPDATE update_jobs 
       SET status = $1::text, 
           updated_at = $2,
           retry_count = CASE 
               WHEN $1::text = $3::text THEN retry_count + 1 
               ELSE retry_count 
           END
       WHERE id = $4::uuid`

	_, err := r.pool.Exec(ctx, query,
		string(params.Status),
		time.Now().UTC(),
		string(job.StatusFailure),
		pgID,
	)

	if err != nil {
		return fmt.Errorf("update status repository error: %w", err)
	}

	return nil
}
