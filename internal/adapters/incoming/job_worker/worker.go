package job_worker

import (
	"context"
	"fmt"
	"io"
	"perezvonish/plata-test-assignment/internal/application/quote/usecases"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Worker interface {
	Start(ctx context.Context)
}

type JobWorker struct {
	Name string

	jobChannel <-chan uuid.UUID
	logger     io.Writer
	wg         *sync.WaitGroup

	processQuoteUpdateUsecase usecases.ProcessQuoteUpdateJobUsecase

	startedAt time.Time
}

func (w *JobWorker) Start(ctx context.Context) {
	w.wg.Add(1)
	w.startedAt = time.Now()

	go func() {
		defer w.wg.Done()
		w.logger.Write([]byte(fmt.Sprintf("[%s] Started at %v\n", w.Name, w.startedAt)))

		for {
			select {
			case <-ctx.Done():
				w.logger.Write([]byte(fmt.Sprintf("[%s] Stopping by context...\n", w.Name)))
				return

			case jobID, ok := <-w.jobChannel:
				if !ok {
					w.logger.Write([]byte(fmt.Sprintf("[%s] Stopping: channel closed\n", w.Name)))
					return
				}

				w.processJob(ctx, jobID)
			}
		}
	}()
}

func (w *JobWorker) processJob(ctx context.Context, id uuid.UUID) {
	w.logger.Write([]byte(fmt.Sprintf("[%s] Processing job: %s\n", w.Name, id)))

	err := w.processQuoteUpdateUsecase.Execute(ctx, usecases.ProcessQuoteUpdateJobUsecaseParams{
		Id: id,
	})
	if err != nil {
		w.logger.Write([]byte(fmt.Sprintf("[%s] ERROR processing job %s: %v\n", w.Name, id, err)))
		return
	}

	w.logger.Write([]byte(fmt.Sprintf("[%s] Successfully finished job: %s\n", w.Name, id)))
}

type WorkerInitParams struct {
	Name                      string
	JobChannel                <-chan uuid.UUID
	Logger                    io.Writer
	Wg                        *sync.WaitGroup
	ProcessQuoteUpdateUsecase usecases.ProcessQuoteUpdateJobUsecase
}

func NewJobWorker(params WorkerInitParams) Worker {
	return &JobWorker{
		Name:                      params.Name,
		jobChannel:                params.JobChannel,
		logger:                    params.Logger,
		wg:                        params.Wg,
		processQuoteUpdateUsecase: params.ProcessQuoteUpdateUsecase,
	}
}
