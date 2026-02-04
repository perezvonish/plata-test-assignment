package quotes

import (
	"context"
	"net/http"
	"perezvonish/plata-test-assignment/internal/domain/currency"
)

type ParamsValidator interface {
	ValidateUpdate(ctx context.Context, r *http.Request) (UpdateInput, error)
	ValidateGetByUpdateId(ctx context.Context, r *http.Request)
	ValidateGetLatest(ctx context.Context, r *http.Request)
}

type ParamsValidatorImpl struct{}

func (v ParamsValidatorImpl) ValidateUpdate(ctx context.Context, r *http.Request) (UpdateInput, error) {
	query := r.URL.Query()
	fromRaw := query.Get("from")
	toRaw := query.Get("to")

	if fromRaw == "" || toRaw == "" {
		return UpdateInput{}, ErrorFromAndToAreRequired
	}

	from := currency.Currency(fromRaw)
	to := currency.Currency(toRaw)

	if !currency.IsValid(from) || !currency.IsValid(to) {
		return UpdateInput{}, ErrorInvalidCurrency
	}

	if from == to {
		return UpdateInput{}, ErrorIdenticalCurrency
	}

	return UpdateInput{
		From: from,
		To:   to,
	}, nil
}

func (v ParamsValidatorImpl) ValidateGetByUpdateId(ctx context.Context, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (v ParamsValidatorImpl) ValidateGetLatest(ctx context.Context, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func newParamsValidator() ParamsValidator {
	return ParamsValidatorImpl{}
}
