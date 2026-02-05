package app

import (
	"perezvonish/plata-test-assignment/internal/domain/job"
	infraJob "perezvonish/plata-test-assignment/internal/infrastructure/database/postgres/job"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type JobContainer struct {
	Repository job.Repository

	UpdateChannel chan uuid.UUID
}

type initJobContainerParams struct {
	Pool *pgxpool.Pool
}

func initJobContainer(params initJobContainerParams) *JobContainer {
	module := &JobContainer{
		UpdateChannel: make(chan uuid.UUID),
	}

	repository := infraJob.NewRepository(params.Pool)
	module.Repository = repository

	return module
}
