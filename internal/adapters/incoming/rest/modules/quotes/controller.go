package quotes

import (
	"net/http"
	"perezvonish/plata-test-assignment/internal/adapters/incoming/rest/middlewares"
	"perezvonish/plata-test-assignment/internal/app"
)

var (
	updateRoute   = "PATCH /v1/quotes"
	getByUpdateId = "GET /v1/quotes/update-task/{id}"
	getLatest     = "GET /v1/quotes/latest"
)

type Controller struct {
	handler *Handler
}

func newController(app *app.Container) *Controller {
	handler := newHandler(HandlerInitParams{
		UpdateUsecase:        app.Quote.UpdateUsecase,
		GetByUpdateIdUsecase: app.Quote.GetByUpdateIdUsecase,
		GetLatestUsecase:     app.Quote.GetLatestUsecase,
	})

	return &Controller{
		handler: handler,
	}
}

func (c *Controller) RegisterRoutes(mux *http.ServeMux) {
	updateHandler := middlewares.Idempotency(http.HandlerFunc(c.handler.Update))

	mux.Handle(updateRoute, updateHandler)
	mux.HandleFunc(getByUpdateId, c.handler.GetByUpdateId)
	mux.HandleFunc(getLatest, c.handler.GetLatest)
}
