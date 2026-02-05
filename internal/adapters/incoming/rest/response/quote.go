package response

import (
	"perezvonish/plata-test-assignment/internal/domain/currency"

	"github.com/google/uuid"
)

type Quote struct {
	Id uuid.UUID `json:"id"`

	FromCurrency currency.Currency `json:"fromCurrency"`
	ToCurrency   currency.Currency `json:"toCurrency"`

	Price float64 `json:"price"`
}
