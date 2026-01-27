package rest

import (
	"net/http"
	"perezvonish/plata-test-assignment/internal/adapters/incoming/rest/middlewares"
)

func newRouter() http.Handler {
	mux := http.NewServeMux()

	container := newContainer()
	container.registerAllModules(mux)

	return middlewares.Use(mux)
}
