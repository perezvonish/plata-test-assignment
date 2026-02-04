package quote

import (
	"context"
	"perezvonish/plata-test-assignment/internal/domain/currency"

	"github.com/google/uuid"
)

type GetByPairParams struct {
	FromCurrency currency.Currency
	ToCurrency   currency.Currency
}

type UpdatePriceParams struct {
	Id    uuid.UUID
	Price int64
}

type SaveParams struct {
	FromCurrency currency.Currency
	ToCurrency   currency.Currency
}

type Repository interface {
	GetById(ctx context.Context, id uuid.UUID) (*Quote, error)
	GetByPair(ctx context.Context, params GetByPairParams) (*Quote, error)
	UpdatePrice(ctx context.Context, params UpdatePriceParams) error
	Save(ctx context.Context, params SaveParams) (*Quote, error)
}
