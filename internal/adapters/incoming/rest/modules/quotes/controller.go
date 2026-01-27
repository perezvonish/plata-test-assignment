package quotes

import "net/http"

var (
	updateRoute   = "PATCH /v1/quotes"
	getByUpdateId = "GET /v1/quotes/{updateId}"
	getLatest     = "GET /v1/quotes/latest"
)

type Controller struct {
	handler *Handler
}

func newController() *Controller {
	handler := newHandler()

	return &Controller{
		handler: handler,
	}
}

func (c *Controller) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc(updateRoute, c.handler.Update)
	mux.HandleFunc(getByUpdateId, c.handler.GetByUpdateId)
	mux.HandleFunc(getLatest, c.handler.GetLatest)
}
