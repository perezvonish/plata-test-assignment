package quotes

import (
	"errors"
	"perezvonish/plata-test-assignment/internal/domain/currency"

	"github.com/google/uuid"
)

type UpdateInput struct {
	From currency.Currency `json:"from"`
	To   currency.Currency `json:"to"`
}

type UpdateOutput struct {
	UpdateId uuid.UUID `json:"update_id"`
}

var (
	ErrorFromAndToAreRequired = errors.New("parameters 'from' and 'to' are required")
	ErrorInvalidCurrency      = errors.New("invalid currency code provided")
	ErrorIdenticalCurrency    = errors.New("currencies are identical")
)
