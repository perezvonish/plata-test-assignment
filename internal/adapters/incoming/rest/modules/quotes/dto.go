package quotes

import (
	"errors"
	"perezvonish/plata-test-assignment/internal/adapters/incoming/rest/response"
	"perezvonish/plata-test-assignment/internal/domain/currency"

	"github.com/google/uuid"
)

type UpdateInput struct {
	From currency.Currency `json:"from"`
	To   currency.Currency `json:"to"`
}

type UpdateOutput struct {
	UpdateId uuid.UUID `json:"updateId"`
}

type GetByUpdateIdInput struct {
	Id uuid.UUID `json:"id"`
}
type GetByUpdateIdOutput struct {
	Job response.Job `json:"job"`
}

type GetLatestInput struct {
	From currency.Currency `json:"from"`
	To   currency.Currency `json:"to"`
}
type GetLatestOutput struct {
	Job response.Job `json:"job"`
}

var (
	ErrorFromAndToAreRequired = errors.New("parameters 'from' and 'to' are required")
	ErrorInvalidCurrency      = errors.New("invalid currency code provided")
	ErrorIdenticalCurrency    = errors.New("currencies are identical")
	ErrorNotPassedJobId       = errors.New("'job id' is required")
	ErrorNotValidJobId        = errors.New("'job id' isn't valid")
)
