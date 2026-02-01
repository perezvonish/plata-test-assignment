package job_worker

import (
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"
)

type Worker interface {
	Start() error
	Stop() error
}

type JobWorker struct {
	Name string

	jobChannel chan<- uuid.UUID
	logger     io.Writer

	startedAt time.Time
}

func (w *JobWorker) Start() error {
	startTime := time.Now()

	w.logger.Write([]byte(fmt.Sprintf("\nStarted new job worker %s at %v", w.Name, startTime)))

	return nil
}

func (w *JobWorker) Stop() error {
	return nil
}

type WorkerInitParams struct {
	Name string

	JobChannel chan<- uuid.UUID

	Logger io.Writer
}

func NewJobWorker(params WorkerInitParams) Worker {
	return &JobWorker{
		Name:       params.Name,
		jobChannel: params.JobChannel,
		logger:     params.Logger,
	}
}
