package quotes

import (
	"net/http"
	"perezvonish/plata-test-assignment/internal/adapters/incoming/rest/middlewares"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	updateRoute   = "PATCH /v1/quotes"
	getByUpdateId = "GET /v1/quotes/update-task/{id}"
	getLatest     = "GET /v1/quotes/latest"
)

type Controller struct {
	handler *Handler
}

type ControllerInitParams struct {
	Pool *pgxpool.Pool

	JobChannel chan<- uuid.UUID
}

func newController(params ControllerInitParams) *Controller {
	handler := newHandler(HandlerInitParams{
		Pool:       params.Pool,
		JobChannel: params.JobChannel,
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
