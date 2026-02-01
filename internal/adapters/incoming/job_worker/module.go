package job_worker

import (
	"io"
	"perezvonish/plata-test-assignment/internal/shared/config"
)

type Module struct {
	pool *WorkerPool
}

type ModuleInitParams struct {
	Config *config.Config
	Logger io.Writer
}

func NewModule(params ModuleInitParams) *Module {
	pool := NewWorkerPool(PoolInitParams{
		WorkerCount: params.Config.JobWorker.WorkerCount,
		Logger:      params.Logger,
	})

	return &Module{
		pool: pool,
	}
}

func (m *Module) Start() error {
	err := m.pool.Start()
	if err != nil {
		return err
	}
}
