package quote

import (
	"perezvonish/plata-test-assignment/internal/domain/currency"
	"time"

	"github.com/google/uuid"
)

type Quote struct {
	Id uuid.UUID

	FromCurrency currency.Currency
	ToCurrency   currency.Currency

	PriceE8Rate int64

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (q *Quote) UpdatePriceE8Rate(price int64) {
	q.PriceE8Rate = price
	q.UpdatedAt = time.Now()
}

func (q *Quote) GetE8Price(value float64) int64 {
	return int64(value * PricePrecision)
}

func (q *Quote) GetConvertedPrice(value int64) float64 {
	return float64(value) / float64(PricePrecision)
}
