package job

import (
	"perezvonish/plata-test-assignment/internal/domain/job"
	"perezvonish/plata-test-assignment/internal/infrastructure/database/postgres/quote"
	"time"

	"github.com/google/uuid"
)

type Model struct {
	Id uuid.UUID `db:"id"`

	QuoteId uuid.UUID    `db:"quote_id"`
	Quote   *quote.Model `db:"quote"`
	Status  string       `db:"status"`

	RetryCount int `db:"retry_count"`

	PriceE8Rate int64 `db:"price_e8_rate"`

	IdempotencyKey string `db:"idempotency_key"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (m *Model) MapToDomain() *job.Job {
	res := &job.Job{
		Id:             m.Id,
		QuoteId:        m.QuoteId,
		Status:         job.Status(m.Status),
		RetryCount:     m.RetryCount,
		PriceE8Rate:    m.PriceE8Rate,
		IdempotencyKey: m.IdempotencyKey,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}

	if m.Quote != nil {
		res.Quote = m.Quote.MapToDomain()
	}

	return res
}

func MapToModel(j *job.Job) *Model {
	return &Model{
		Id:             j.Id,
		QuoteId:        j.QuoteId,
		Status:         string(j.Status),
		RetryCount:     int(j.RetryCount),
		PriceE8Rate:    j.PriceE8Rate,
		IdempotencyKey: j.IdempotencyKey,
		CreatedAt:      j.CreatedAt,
		UpdatedAt:      j.UpdatedAt,
	}
}
