package job_worker

import (
	"context"
	"io"
	"perezvonish/plata-test-assignment/internal/application/quote/usecases"
	"perezvonish/plata-test-assignment/internal/shared/config"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Module struct {
	pool *WorkerPool
}

type ModuleInitParams struct {
	Pool   *pgxpool.Pool
	Config *config.Config

	Logger          io.Writer
	ConsumerChannel <-chan uuid.UUID
}

func NewModule(params ModuleInitParams) *Module {
	processQuoteUpdateUsecase := usecases.NewProcessQuoteUpdateUsecase(usecases.ProcessQuoteUpdateUsecaseInitParams{
		Pool:   params.Pool,
		Config: &params.Config.ExchangeApiConfig,
	})

	pool := NewWorkerPool(PoolInitParams{
		WorkerCount:               params.Config.JobWorker.WorkerCount,
		Logger:                    params.Logger,
		ConsumerChannel:           params.ConsumerChannel,
		ProcessQuoteUpdateUsecase: processQuoteUpdateUsecase,
	})

	return &Module{
		pool: pool,
	}
}

func (m *Module) StartWorkers(ctx context.Context) {
	m.pool.Start(ctx)
}

func (m *Module) StopWorkers() {
	m.pool.Stop()
}
