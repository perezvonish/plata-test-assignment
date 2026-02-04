package rest

import (
	"net/http"
	"perezvonish/plata-test-assignment/internal/adapters/incoming/rest/middlewares"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RouterInitParams struct {
	Pool *pgxpool.Pool

	JobChannel chan<- uuid.UUID
}

func newRouter(params RouterInitParams) http.Handler {
	mux := http.NewServeMux()

	container := newContainer(ContainerInitParams{
		pool:       params.Pool,
		JobChannel: params.JobChannel,
	})
	container.registerAllModules(mux)

	return middlewares.Use(mux)
}
