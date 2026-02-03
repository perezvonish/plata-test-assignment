package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"perezvonish/plata-test-assignment/internal/adapters/incoming/job_worker"
	"perezvonish/plata-test-assignment/internal/adapters/incoming/rest"
	"perezvonish/plata-test-assignment/internal/infrastructure/database"
	"perezvonish/plata-test-assignment/internal/shared/config"
	"syscall"
	"time"

	"github.com/google/uuid"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		panic(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	dbPool, err := database.ConnectWithRetry(database.ConnectInitParams{
		Config: cfg,
	})
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer dbPool.Close()

	jobChan := make(chan uuid.UUID, 100)
	workerModule := job_worker.NewModule(job_worker.ModuleInitParams{
		Config:          cfg,
		Logger:          os.Stdout,
		ConsumerChannel: jobChan,
	})
	workerModule.StartWorkers(ctx)

	httpServer := rest.NewServer(ctx, *cfg)
	go httpServer.Start()

	<-ctx.Done()
	gracefulShutdown(httpServer, workerModule)
}

func gracefulShutdown(httpServer *rest.Server, workerModule *job_worker.Module) {
	fmt.Println("\nShutdown signal received...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := httpServer.Stop(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	workerModule.StopWorkers()

	fmt.Println("Application stopped gracefully")
}
