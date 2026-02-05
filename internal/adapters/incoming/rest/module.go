package rest

import (
	"net/http"
	swaggerDocs "perezvonish/plata-test-assignment/internal/adapters/incoming/rest/modules/docs"
	"perezvonish/plata-test-assignment/internal/adapters/incoming/rest/modules/health"
	"perezvonish/plata-test-assignment/internal/adapters/incoming/rest/modules/quotes"
	"perezvonish/plata-test-assignment/internal/app"
)

type Module interface {
	Register(mux *http.ServeMux)
}

type Container struct {
	modules []Module
}

func newContainer(app *app.Container) *Container {
	docs := swaggerDocs.NewModule()
	healthModule := health.NewModule()
	quoteModule := quotes.NewModule(app)

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
