package usecases

import (
	"context"
	"perezvonish/plata-test-assignment/internal/adapters/incoming/rest/response"
	"perezvonish/plata-test-assignment/internal/application/quote"
	"perezvonish/plata-test-assignment/internal/domain/job"

	"github.com/google/uuid"
)

type QuoteGetByUpdateIdUsecaseInput struct {
	UpdateId uuid.UUID
}

type QuoteGetByUpdateIdUsecase interface {
	Execute(ctx context.Context, input QuoteGetByUpdateIdUsecaseInput) (response.Job, error)
}

type QuoteGetByUpdateIdUsecaseImpl struct {
	jobRepository job.Repository
}

func (q *QuoteGetByUpdateIdUsecaseImpl) Execute(ctx context.Context, input QuoteGetByUpdateIdUsecaseInput) (response.Job, error) {
	currentJob, err := q.jobRepository.GetByUpdateId(ctx, input.UpdateId)
	if err != nil {
		return response.Job{}, quote.ErrorWhileFindingJob
	}

	if currentJob == nil {
		return response.Job{}, quote.ErrorNotFoundJob
	}

	res := response.Job{
		Id:             currentJob.Id,
		Quote:          nil,
		Status:         currentJob.Status,
		IdempotencyKey: currentJob.IdempotencyKey,
		CreatedAt:      currentJob.CreatedAt,
	}

	if currentJob.Quote != nil {
		res.Quote = &response.Quote{
			Id:           currentJob.Quote.Id,
			FromCurrency: currentJob.Quote.FromCurrency,
			ToCurrency:   currentJob.Quote.ToCurrency,
			Price:        currentJob.GetConvertedPrice(currentJob.PriceE8Rate),
		}
	}

	return res, nil
}

type QuoteGetByUpdateIdUsecaseInitParams struct {
	JobRepository job.Repository
}

func NewQuoteGetByUpdateIdUsecase(params QuoteGetByUpdateIdUsecaseInitParams) QuoteGetByUpdateIdUsecase {
	return &QuoteGetByUpdateIdUsecaseImpl{
		jobRepository: params.JobRepository,
	}
}
