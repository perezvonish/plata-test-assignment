package rest

import (
	"net/http"
	"perezvonish/plata-test-assignment/internal/adapters/incoming/rest/modules/health"
)

type Module interface {
	Register(mux *http.ServeMux)
}

type Container struct {
	modules []Module
}

func newContainer() *Container {
	healthModule := health.NewModule()

	return &Container{
		modules: []Module{
			healthModule,
		},
	}
}

func (c *Container) registerAllModules(mux *http.ServeMux) {
	for _, m := range c.modules {
		m.Register(mux)
	}
}
