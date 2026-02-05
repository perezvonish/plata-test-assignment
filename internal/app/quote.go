package app

import (
	"perezvonish/plata-test-assignment/internal/application/quote/usecases"
	"perezvonish/plata-test-assignment/internal/domain/quote"
	infraQuote "perezvonish/plata-test-assignment/internal/infrastructure/database/postgres/quote"

	"github.com/jackc/pgx/v5/pgxpool"
)

type QuoteContainer struct {
	Repository quote.Repository

	UpdateUsecase        usecases.QuoteUpdateUsecase
	GetByUpdateIdUsecase usecases.QuoteGetByUpdateIdUsecase
	GetLatestUsecase     usecases.QuoteGetLatestUsecase
}

type initQuoteContainerParams struct {
	Pool *pgxpool.Pool

	JobContainer *JobContainer
}

func initQuoteContainer(params initQuoteContainerParams) *QuoteContainer {
	repository := infraQuote.NewRepository(params.Pool)
	updateUsecase := usecases.NewQuoteUpdateUsecase(usecases.QuoteUpdateUsecaseInitParams{
		JobRepository:   params.JobContainer.Repository,
		JobChannel:      params.JobContainer.UpdateChannel,
		QuoteRepository: repository,
	})
	getByUpdateIdUsecase := usecases.NewQuoteGetByUpdateIdUsecase(usecases.QuoteGetByUpdateIdUsecaseInitParams{JobRepository: params.JobContainer.Repository})
	getLatestUsecase := usecases.NewQuoteGetLatestUsecase(usecases.QuoteGetLatestUsecaseInitParams{
		JobRepository: params.JobContainer.Repository,
	})

	return &QuoteContainer{
		Repository:           repository,
		UpdateUsecase:        updateUsecase,
		GetByUpdateIdUsecase: getByUpdateIdUsecase,
		GetLatestUsecase:     getLatestUsecase,
	}
}
