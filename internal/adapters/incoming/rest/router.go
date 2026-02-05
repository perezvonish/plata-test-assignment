package rest

import (
	"net/http"
	"perezvonish/plata-test-assignment/internal/adapters/incoming/rest/middlewares"
	"perezvonish/plata-test-assignment/internal/app"
)

func newRouter(app *app.Container) http.Handler {
	mux := http.NewServeMux()

	container := newContainer(app)
	container.registerAllModules(mux)

	return middlewares.Use(mux)
}
