package docs

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

var (
	getSwagger = "GET /swagger/"
)

type Controller struct{}

func newController() *Controller {
	return &Controller{}
}

func (c *Controller) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), // Ссылка на JSON-спецификацию
	))
}
