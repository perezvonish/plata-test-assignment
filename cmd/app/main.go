package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"perezvonish/plata-test-assignment/internal/adapters/incoming/job_worker"
	"perezvonish/plata-test-assignment/internal/adapters/incoming/rest"
	"perezvonish/plata-test-assignment/internal/shared/config"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		panic(err)
	}

	//jobChannel := make(chan uuid.UUID, 100)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	jobWorkers := job_worker.NewModule(job_worker.ModuleInitParams{
		Config: cfg,
		Logger: os.Stdout,
	})

	if err := jobWorkers.Start(); err != nil {
		panic(err)
	}

	httpServer := rest.NewServer(ctx, *cfg)

	go httpServer.Start()

	<-ctx.Done()
	shutdown(httpServer)
}

func shutdown(httpServer *rest.Server) {
	fmt.Println("\nShutdown signal received...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := httpServer.Stop(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	fmt.Println("Application stopped gracefully")
}
