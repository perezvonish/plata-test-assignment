package job

import (
	"perezvonish/plata-test-assignment/internal/domain/job"
	"time"

	"github.com/google/uuid"
)

type Model struct {
	Id uuid.UUID `db:"id"`

	QuoteId uuid.UUID `db:"quote_id"`
	Status  string    `db:"status"`

	RetryCount int `db:"retry_count"`

	IdempotencyKey string `db:"idempotency_key"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (m *Model) MapToDomain() *job.Job {
	return &job.Job{
		Id:             m.Id,
		QuoteId:        m.QuoteId,
		Status:         job.Status(m.Status),
		RetryCount:     m.RetryCount,
		IdempotencyKey: m.IdempotencyKey,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
}

func MapToModel(j *job.Job) *Model {
	return &Model{
		Id:             j.Id,
		QuoteId:        j.QuoteId,
		Status:         string(j.Status),
		RetryCount:     int(j.RetryCount),
		IdempotencyKey: j.IdempotencyKey,
		CreatedAt:      j.CreatedAt,
		UpdatedAt:      j.UpdatedAt,
	}
}
