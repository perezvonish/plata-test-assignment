package job

import (
	"perezvonish/plata-test-assignment/internal/domain/quote"
	"time"

	"github.com/google/uuid"
)

type Job struct {
	Id uuid.UUID

	QuoteId uuid.UUID
	Quote   *quote.Quote
	Status  Status

	RetryCount int

	PriceE8Rate int64

	IdempotencyKey string

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (j *Job) MarkAsProcessing() {
	j.Status = StatusProcessing
	j.UpdatedAt = time.Now()
}

func (j *Job) MarkAsSuccess() {
	j.Status = StatusSuccess
	j.UpdatedAt = time.Now()
}

func (j *Job) MarkAsFailure() {
	j.Status = StatusFailure
	j.RetryCount++
	j.UpdatedAt = time.Now()
}

func (j *Job) UpdatePrice(price int64) {
	j.PriceE8Rate = price
	j.UpdatedAt = time.Now()
}

func (j *Job) GetConvertedPrice(value int64) float64 {
	return float64(value) / float64(quote.PricePrecision)
}
