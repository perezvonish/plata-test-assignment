package job_worker

import (
	"context"
	"io"
	"perezvonish/plata-test-assignment/internal/app"
	"perezvonish/plata-test-assignment/internal/shared/config"
)

type Module struct {
	pool *WorkerPool
}

type ModuleInitParams struct {
	Config *config.Config

	Logger       io.Writer
	AppContainer *app.Container
}

func NewModule(params ModuleInitParams) *Module {
	pool := NewWorkerPool(PoolInitParams{
		WorkerCount:               params.Config.JobWorker.WorkerCount,
		Logger:                    params.Logger,
		ConsumerChannel:           params.AppContainer.Job.UpdateChannel,
		ProcessQuoteUpdateUsecase: params.AppContainer.ExternalExchange.ProcessQuoteUpdateJobUsecase,
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
