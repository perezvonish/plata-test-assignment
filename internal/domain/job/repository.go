package job

import (
	"context"

	"github.com/google/uuid"
)

type UpdateStatusParams struct {
	Id     uuid.UUID
	Status Status
}

type CreateParams struct {
	QuoteID        uuid.UUID `json:"quote_id"`
	IdempotencyKey string    `json:"idempotency_key"`
}

type UpdatePriceParams struct {
	Id          uuid.UUID
	PriceE8Rate int64
}

type Repository interface {
	GetById(ctx context.Context, id uuid.UUID) (*Job, error)
	GetByUpdateId(ctx context.Context, updateId uuid.UUID) (*Job, error)
	GetByIdempotencyKey(ctx context.Context, key string) (*Job, error)
	UpdateStatus(ctx context.Context, params UpdateStatusParams) error
	UpdatePrice(ctx context.Context, params UpdatePriceParams) error
	Save(ctx context.Context, params CreateParams) (uuid.UUID, error)
}
