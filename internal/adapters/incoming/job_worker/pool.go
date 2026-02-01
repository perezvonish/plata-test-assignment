package job_worker

import (
	"fmt"
	"io"
	"perezvonish/plata-test-assignment/internal/shared/utils"

	"github.com/google/uuid"
)

type WorkerPool struct {
	workerCount int64
	workers     []Worker

	logger io.Writer

	consumeChannel chan<- uuid.UUID
}

type PoolInitParams struct {
	WorkerCount     int64
	ConsumerChannel chan<- uuid.UUID
	Logger          io.Writer
}

func NewWorkerPool(params PoolInitParams) *WorkerPool {
	return &WorkerPool{
		workerCount:    params.WorkerCount,
		consumeChannel: params.ConsumerChannel,
		logger:         params.Logger,
	}
}

func (wp *WorkerPool) Start() error {
	var workers []Worker

	for i := int64(0); i < wp.workerCount; i++ {
		name := utils.GenerateUUID()
		workers = append(workers, NewJobWorker(WorkerInitParams{
			Name:       name,
			JobChannel: wp.consumeChannel,
			Logger:     wp.logger,
		}))
	}

	wp.workers = workers

	for _, worker := range wp.workers {
		err := worker.Start()
		if err != nil {
			return err
		}
	}

	wp.logger.Write([]byte(fmt.Sprintf("\n=============WORKER POOL STARTED =============\n")))

	return nil
}

func (wp *WorkerPool) Stop() error {
	return nil
}
