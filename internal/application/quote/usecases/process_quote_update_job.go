package usecases

import (
	"context"
	"fmt"
	"perezvonish/plata-test-assignment/internal/application/quote/services"
	"perezvonish/plata-test-assignment/internal/domain/job"
	"perezvonish/plata-test-assignment/internal/domain/quote"
	infraJob "perezvonish/plata-test-assignment/internal/infrastructure/database/postgres/job"
	infraQuote "perezvonish/plata-test-assignment/internal/infrastructure/database/postgres/quote"
	"perezvonish/plata-test-assignment/internal/shared/config"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProcessQuoteUpdateJobUsecaseParams struct {
	Id uuid.UUID
}

type ProcessQuoteUpdateJobUsecase interface {
	Execute(ctx context.Context, params ProcessQuoteUpdateJobUsecaseParams) error
}

type ProcessQuoteUpdateJobUsecaseImpl struct {
	jobRepository   job.Repository
	quoteRepository quote.Repository

	exchangePriceService services.ExchangePrice
}

func (p *ProcessQuoteUpdateJobUsecaseImpl) Execute(ctx context.Context, params ProcessQuoteUpdateJobUsecaseParams) error {
	currentJob, err := p.jobRepository.GetById(ctx, params.Id)
	if err != nil || currentJob == nil {
		return fmt.Errorf("job not found or error: %w", err)
	}

	currentJob.MarkAsProcessing()
	err = p.jobRepository.UpdateStatus(ctx, job.UpdateStatusParams{
		Id:     currentJob.Id,
		Status: currentJob.Status,
	})
	if err != nil {
		return err
	}

	dbQuote, err := p.quoteRepository.GetById(ctx, currentJob.QuoteId)
	if err != nil {
		return err
	}

	rate, err := p.exchangePriceService.GetRate(ctx, services.GetRateParams{
		FromCurrency: dbQuote.FromCurrency,
		ToCurrency:   dbQuote.ToCurrency,
	})

	if err != nil {
		currentJob.MarkAsFailure()

		_ = p.jobRepository.UpdateStatus(ctx, job.UpdateStatusParams{
			Id:     currentJob.Id,
			Status: currentJob.Status,
		})
		return fmt.Errorf("external api error: %w", err)
	}

	err = p.quoteRepository.UpdatePrice(ctx, quote.UpdatePriceParams{
		Id:    dbQuote.Id,
		Price: dbQuote.GetE8Price(rate),
	})
	if err != nil {
		currentJob.MarkAsFailure()
		_ = p.jobRepository.UpdateStatus(ctx, job.UpdateStatusParams{
			Id:     currentJob.Id,
			Status: currentJob.Status,
		})
		return err
	}

	currentJob.MarkAsSuccess()
	return p.jobRepository.UpdateStatus(ctx, job.UpdateStatusParams{
		Id:     currentJob.Id,
		Status: currentJob.Status,
	})
}

type ProcessQuoteUpdateUsecaseInitParams struct {
	Pool   *pgxpool.Pool
	Config *config.ExchangeApiConfig
}

func NewProcessQuoteUpdateUsecase(params ProcessQuoteUpdateUsecaseInitParams) ProcessQuoteUpdateJobUsecase {
	jobRepository := infraJob.NewRepository(params.Pool)
	quoteRepository := infraQuote.NewRepository(params.Pool)
	exchangeService := services.NewExchangePrice(services.ExchangePriceInitParams{
		ExchangeApiConfig: params.Config,
	})

	return &ProcessQuoteUpdateJobUsecaseImpl{
		jobRepository:        jobRepository,
		quoteRepository:      quoteRepository,
		exchangePriceService: exchangeService,
	}
}
