package health

import (
	"net/http"
)

var (
	PingRoute = "GET /v1/health/ping"
)

type Controller struct {
	handler *Handler
}

func NewController() *Controller {
	handler := NewHandler()

	return &Controller{
		handler: &handler,
	}
}

func (c *Controller) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc(PingRoute, c.handler.Ping)
}
