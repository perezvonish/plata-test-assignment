package quote

import (
	"perezvonish/plata-test-assignment/internal/domain/currency"
	"perezvonish/plata-test-assignment/internal/domain/quote"
	"time"

	"github.com/google/uuid"
)

type Model struct {
	Id           uuid.UUID `db:"id"`
	FromCurrency string    `db:"from_currency"`
	ToCurrency   string    `db:"to_currency"`
	PriceE8Rate  int64     `db:"price_e8_rate"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

func (m *Model) MapToDomain() *quote.Quote {
	return &quote.Quote{
		Id:           m.Id,
		FromCurrency: currency.Currency(m.FromCurrency),
		ToCurrency:   currency.Currency(m.ToCurrency),
		PriceE8Rate:  m.PriceE8Rate,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

func MapToModel(q *quote.Quote) *Model {
	return &Model{
		Id:           q.Id,
		FromCurrency: string(q.FromCurrency),
		ToCurrency:   string(q.ToCurrency),
		PriceE8Rate:  q.PriceE8Rate,
		CreatedAt:    q.CreatedAt,
		UpdatedAt:    q.UpdatedAt,
	}
}
