package quotes

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Module struct {
	controller *Controller
}

type ModuleInitParams struct {
	Pool *pgxpool.Pool

	JobChannel chan<- uuid.UUID
}

func NewModule(params ModuleInitParams) *Module {
	controller := newController(ControllerInitParams{
		Pool:       params.Pool,
		JobChannel: params.JobChannel,
	})

	return &Module{
		controller: controller,
	}
}

func (m *Module) Register(mux *http.ServeMux) {
	m.controller.RegisterRoutes(mux)
}
