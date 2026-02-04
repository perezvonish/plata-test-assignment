package usecases

import (
	"perezvonish/plata-test-assignment/internal/domain/job"

	"github.com/google/uuid"
)

type UpdateQuotesResponse struct {
	UpdateID uuid.UUID `json:"update_id"`
}

type UpdateStatusResponse struct {
	ID     string     `json:"id"`
	Status job.Status `json:"status"`
}
