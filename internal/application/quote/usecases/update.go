package usecases

import (
	"context"
	"fmt"
	"perezvonish/plata-test-assignment/internal/adapters/incoming/rest/middlewares"
	"perezvonish/plata-test-assignment/internal/domain/currency"
	"perezvonish/plata-test-assignment/internal/domain/job"
	"perezvonish/plata-test-assignment/internal/domain/quote"

	"github.com/google/uuid"
)

type CreateParams struct {
	FromCurrency   currency.Currency `json:"from_currency"`
	ToCurrency     currency.Currency `json:"to_currency"`
	IdempotencyKey string            `json:"idempotency_key"`
}

type QuoteUpdateUsecaseInput struct {
	FromCurrency currency.Currency
	ToCurrency   currency.Currency
}

type QuoteUpdateUsecase interface {
	Execute(ctx context.Context, input QuoteUpdateUsecaseInput) (uuid.UUID, error)
}

type QuoteUpdateUsecaseImpl struct {
	jobRepository job.Repository
	jobChannel    chan<- uuid.UUID

	quoteRepository quote.Repository
}

func (u *QuoteUpdateUsecaseImpl) Execute(ctx context.Context, input QuoteUpdateUsecaseInput) (uuid.UUID, error) {
	idempotencyKey, ok := ctx.Value(middlewares.IdempotencyKey).(string)
	if !ok {
		if uid, ok := ctx.Value(middlewares.IdempotencyKey).(uuid.UUID); ok {
			idempotencyKey = uid.String()
		} else {
			return uuid.Nil, fmt.Errorf("idempotency key is missing in context")
		}
	}

	var (
		dbQuote *quote.Quote
		err     error
	)

	dbQuote, err = u.quoteRepository.GetByPair(ctx, quote.GetByPairParams{
		FromCurrency: input.FromCurrency,
		ToCurrency:   input.ToCurrency,
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to find quote for pair: %w", err)
	}

	if dbQuote == nil {
		dbQuote, err = u.quoteRepository.Save(ctx, quote.SaveParams{
			FromCurrency: input.FromCurrency,
			ToCurrency:   input.ToCurrency,
		})
		if err != nil {
			return uuid.Nil, fmt.Errorf("failed to save quote: %w", err)
		}
	}
	existingJob, err := u.jobRepository.GetByIdempotencyKey(ctx, idempotencyKey)
	if err != nil {
		return uuid.Nil, err
	}
	if existingJob != nil {
		return existingJob.Id, nil
	}

	updateID, err := u.jobRepository.Save(ctx, job.CreateParams{
		QuoteID:        dbQuote.Id,
		IdempotencyKey: idempotencyKey,
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create job: %w", err)
	}

	u.jobChannel <- updateID

	return updateID, nil
}

type QuoteUpdateUsecaseInitParams struct {
	JobRepository job.Repository
	JobChannel    chan<- uuid.UUID

	QuoteRepository quote.Repository
}

func NewQuoteUpdateUsecase(params QuoteUpdateUsecaseInitParams) QuoteUpdateUsecase {
	return &QuoteUpdateUsecaseImpl{
		jobRepository:   params.JobRepository,
		jobChannel:      params.JobChannel,
		quoteRepository: params.QuoteRepository,
	}
}
