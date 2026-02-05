package app

import (
	"perezvonish/plata-test-assignment/internal/application/quote/services"
	"perezvonish/plata-test-assignment/internal/application/quote/usecases"
	"perezvonish/plata-test-assignment/internal/shared/config"
)

type ExternalExchangeContainer struct {
	ExchangePriceService services.ExchangePrice

	ProcessQuoteUpdateJobUsecase usecases.ProcessQuoteUpdateJobUsecase
}

type initExternalExchangeContainerParams struct {
	Config         *config.ExchangeApiConfig
	JobContainer   *JobContainer
	QuoteContainer *QuoteContainer
}

func NewExternalExchangeContainer(params initExternalExchangeContainerParams) *ExternalExchangeContainer {
	exchangePriceService := services.NewExchangePrice(services.ExchangePriceInitParams{
		ExchangeApiConfig: params.Config,
	})

	processQuoteUpdateJobUsecase := usecases.NewProcessQuoteUpdateUsecase(usecases.ProcessQuoteUpdateUsecaseInitParams{
		JobRepository:   params.JobContainer.Repository,
		QuoteRepository: params.QuoteContainer.Repository,
		ExchangeService: exchangePriceService,
	})

	return &ExternalExchangeContainer{
		ExchangePriceService:         exchangePriceService,
		ProcessQuoteUpdateJobUsecase: processQuoteUpdateJobUsecase,
	}
}
