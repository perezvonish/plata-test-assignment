package app

import (
	"perezvonish/plata-test-assignment/internal/shared/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Container struct {
	Job              *JobContainer
	Quote            *QuoteContainer
	ExternalExchange *ExternalExchangeContainer
}

type ContainerInitParams struct {
	Config *config.Config
	Pool   *pgxpool.Pool
}

func NewContainer(params ContainerInitParams) *Container {
	jobContainer := initJobContainer(initJobContainerParams{
		Pool: params.Pool,
	})

	quoteContainer := initQuoteContainer(initQuoteContainerParams{
		Pool:         params.Pool,
		JobContainer: jobContainer,
	})

	externalExchangeContainer := NewExternalExchangeContainer(initExternalExchangeContainerParams{
		Config:         &params.Config.ExchangeApiConfig,
		JobContainer:   jobContainer,
		QuoteContainer: quoteContainer,
	})

	return &Container{
		Job:              jobContainer,
		Quote:            quoteContainer,
		ExternalExchange: externalExchangeContainer,
	}
}
