package usecases

import (
	"context"
	"fmt"
	"perezvonish/plata-test-assignment/internal/adapters/incoming/rest/response"
	"perezvonish/plata-test-assignment/internal/domain/currency"
	"perezvonish/plata-test-assignment/internal/domain/job"
)

type QuoteGetLatestUsecaseParams struct {
	From currency.Currency
	To   currency.Currency
}

type QuoteGetLatestUsecase interface {
	Execute(ctx context.Context, params QuoteGetLatestUsecaseParams) (response.Job, error)
}

type QuoteGetLatestUsecaseImpl struct {
	jobRepository job.Repository
}

func (q QuoteGetLatestUsecaseImpl) Execute(ctx context.Context, params QuoteGetLatestUsecaseParams) (response.Job, error) {
	lastJob, err := q.jobRepository.GetLatestByCurrencyPair(ctx, job.GetLatestByCurrencyPairParams{
		From: string(params.From),
		To:   string(params.To),
	})

	if err != nil {
		return response.Job{}, fmt.Errorf("failed to get latest job: %w", err)
	}

	if lastJob == nil {
		return response.Job{}, fmt.Errorf("no jobs found for currency pair %s/%s", params.From, params.To)
	}

	res := response.Job{
		Id:             lastJob.Id,
		Quote:          nil,
		Status:         lastJob.Status,
		IdempotencyKey: lastJob.IdempotencyKey,
		CreatedAt:      lastJob.CreatedAt,
	}

	if lastJob.Quote != nil {
		res.Quote = &response.Quote{
			Id:           lastJob.Quote.Id,
			FromCurrency: lastJob.Quote.FromCurrency,
			ToCurrency:   lastJob.Quote.ToCurrency,
			Price:        lastJob.GetConvertedPrice(lastJob.PriceE8Rate),
		}
	}

	return res, nil
}

type QuoteGetLatestUsecaseInitParams struct {
	JobRepository job.Repository
}

func NewQuoteGetLatestUsecase(params QuoteGetLatestUsecaseInitParams) QuoteGetLatestUsecase {
	return QuoteGetLatestUsecaseImpl{
		jobRepository: params.JobRepository,
	}
}
