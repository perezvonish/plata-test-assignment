package quotes

import (
	"net/http"
	"perezvonish/plata-test-assignment/internal/app"
)

type Module struct {
	controller *Controller
}

func NewModule(app *app.Container) *Module {
	controller := newController(app)

	return &Module{
		controller: controller,
	}
}

func (m *Module) Register(mux *http.ServeMux) {
	m.controller.RegisterRoutes(mux)
}
