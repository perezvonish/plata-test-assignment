package usecases

import (
	"context"
	"fmt"
	"perezvonish/plata-test-assignment/internal/application/quote/services"
	"perezvonish/plata-test-assignment/internal/domain/job"
	"perezvonish/plata-test-assignment/internal/domain/quote"

	"github.com/google/uuid"
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
		fmt.Println()
		currentJob.MarkAsFailure()
		_ = p.jobRepository.UpdateStatus(ctx, job.UpdateStatusParams{
			Id:     currentJob.Id,
			Status: currentJob.Status,
		})
		return fmt.Errorf("external api error: %w", err)
	}

	newPriceE8 := dbQuote.GetE8Price(rate)

	currentJob.UpdatePrice(newPriceE8)
	err = p.jobRepository.UpdatePrice(ctx, job.UpdatePriceParams{
		Id:          currentJob.Id,
		PriceE8Rate: currentJob.PriceE8Rate,
	})
	if err != nil {
		currentJob.MarkAsFailure()
		_ = p.jobRepository.UpdateStatus(ctx, job.UpdateStatusParams{
			Id:     currentJob.Id,
			Status: currentJob.Status,
		})
		return fmt.Errorf("failed to update price in job: %w", err)
	}

	err = p.quoteRepository.UpdatePrice(ctx, quote.UpdatePriceParams{
		Id:    dbQuote.Id,
		Price: newPriceE8,
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
	JobRepository   job.Repository
	QuoteRepository quote.Repository
	ExchangeService services.ExchangePrice
}

func NewProcessQuoteUpdateUsecase(params ProcessQuoteUpdateUsecaseInitParams) ProcessQuoteUpdateJobUsecase {
	return &ProcessQuoteUpdateJobUsecaseImpl{
		jobRepository:        params.JobRepository,
		quoteRepository:      params.QuoteRepository,
		exchangePriceService: params.ExchangeService,
	}
}
