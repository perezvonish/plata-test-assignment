package health

import "net/http"

type Module struct {
	controller *Controller
}

func NewModule() *Module {
	controller := NewController()

	return &Module{
		controller: controller,
	}
}

func (m *Module) Register(mux *http.ServeMux) {
	m.controller.RegisterRoutes(mux)
}
