package job

import (
	"context"

	"github.com/google/uuid"
)

type UpdateStatusParams struct {
	Id     uuid.UUID
	Status Status
}

type Repository interface {
	GetById(ctx context.Context, id uuid.UUID) (*Job, error)
	GetByIdempotencyKey(ctx context.Context, key string) (*Job, error)
	UpdateStatus(ctx context.Context, params UpdateStatusParams) error
}
