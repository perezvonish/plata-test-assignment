package response

import (
	"perezvonish/plata-test-assignment/internal/domain/job"
	"time"

	"github.com/google/uuid"
)

type Job struct {
	Id uuid.UUID `json:"id"`

	Quote  *Quote     `json:"quote"`
	Status job.Status `json:"status"`

	IdempotencyKey string    `json:"idempotencyKey"`
	CreatedAt      time.Time `json:"createdAt"`
}
