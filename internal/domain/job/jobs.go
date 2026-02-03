package job

import (
	"time"

	"github.com/google/uuid"
)

type Job struct {
	Id uuid.UUID

	QuoteId uuid.UUID
	Status  Status

	RetryCount int

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
