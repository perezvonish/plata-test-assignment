package rest

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"perezvonish/plata-test-assignment/internal/app"
	"perezvonish/plata-test-assignment/internal/shared/config"
	"time"
)

type Server struct {
	formattedPort string

	httpServer *http.Server
}

type ServerInitParams struct {
	Ctx    context.Context
	Config *config.Config

	AppContainer *app.Container
}

func NewServer(params ServerInitParams) *Server {
	formattedPort := fmt.Sprintf(":%d", params.Config.Server.Port)

	router := newRouter(params.AppContainer)

	httpServer := &http.Server{
		Addr:         formattedPort,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		BaseContext: func(l net.Listener) context.Context {
			return params.Ctx
		},
	}

	return &Server{
		formattedPort: formattedPort,
		httpServer:    httpServer,
	}
}

func (s *Server) Start() {
	fmt.Println("Server is starting on " + s.formattedPort)
	if err := s.httpServer.ListenAndServe(); err != nil {
		fmt.Printf("Server failed: %s\n", err)
	}
}

func (s *Server) Stop(ctx context.Context) error {
	fmt.Printf("Server is stopping on %s\n", s.formattedPort)
	return s.httpServer.Shutdown(ctx)
}
