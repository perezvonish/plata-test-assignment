package job_worker

import (
	"context"
	"fmt"
	"io"
	"perezvonish/plata-test-assignment/internal/application/quote/usecases"
	"sync"

	"github.com/google/uuid"
)

type WorkerPool struct {
	workerCount int
	workers     []Worker

	logger io.Writer
	wg     sync.WaitGroup

	consumeChannel <-chan uuid.UUID
}

type PoolInitParams struct {
	WorkerCount int

	ConsumerChannel <-chan uuid.UUID
	Logger          io.Writer

	ProcessQuoteUpdateUsecase usecases.ProcessQuoteUpdateJobUsecase
}

func NewWorkerPool(params PoolInitParams) *WorkerPool {
	var workers []Worker

	wg := sync.WaitGroup{}

	for i := 0; i < params.WorkerCount; i++ {
		w := NewJobWorker(WorkerInitParams{
			Name:                      fmt.Sprintf("worker-%d", i),
			JobChannel:                params.ConsumerChannel,
			Logger:                    params.Logger,
			Wg:                        &wg,
			ProcessQuoteUpdateUsecase: params.ProcessQuoteUpdateUsecase,
		})
		workers = append(workers, w)
	}

	return &WorkerPool{
		workerCount:    params.WorkerCount,
		workers:        workers,
		wg:             wg,
		logger:         params.Logger,
		consumeChannel: params.ConsumerChannel,
	}
}

func (wp *WorkerPool) Start(ctx context.Context) {
	wp.logger.Write([]byte(fmt.Sprintf("\n============= STARTING WORKER POOL (count: %d) =============\n", wp.workerCount)))

	for _, w := range wp.workers {
		w.Start(ctx)
	}
}

func (wp *WorkerPool) Stop() {
	wp.logger.Write([]byte("\n============= STOPPING WORKER POOL =============\n"))
	wp.wg.Wait()
	wp.logger.Write([]byte("\n============= WORKER POOL STOPPED =============\n"))
}
