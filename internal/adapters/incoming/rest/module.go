package rest

import (
	"net/http"
	swaggerDocs "perezvonish/plata-test-assignment/internal/adapters/incoming/rest/modules/docs"
	"perezvonish/plata-test-assignment/internal/adapters/incoming/rest/modules/health"
	"perezvonish/plata-test-assignment/internal/adapters/incoming/rest/modules/quotes"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Module interface {
	Register(mux *http.ServeMux)
}

type Container struct {
	modules []Module
}

type ContainerInitParams struct {
	pool *pgxpool.Pool

	JobChannel chan<- uuid.UUID
}

func newContainer(params ContainerInitParams) *Container {
	docs := swaggerDocs.NewModule()
	healthModule := health.NewModule()
	quoteModule := quotes.NewModule(quotes.ModuleInitParams{
		Pool:       params.pool,
		JobChannel: params.JobChannel,
	})

	return &Container{
		modules: []Module{
			docs,
			healthModule,
			quoteModule,
		},
	}
}

func (c *Container) registerAllModules(mux *http.ServeMux) {
	for _, m := range c.modules {
		m.Register(mux)
	}
}
