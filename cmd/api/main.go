// Package main starts the API service.
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/florian-renfer/property-service-tracking/internal/api/router"
	"github.com/florian-renfer/property-service-tracking/internal/interface/rest"
)

func main() {
	// Router dependencies
	healthHandler := rest.NewHandler()

	// The HTTP Server
	server := &http.Server{Addr: "0.0.0.0:8080", Handler: router.New(router.Dependencies{HealthHandler: healthHandler})}

	// Create context that listens for the interrupt signal
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Run server in the background
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	// Listen for the interrupt signal
	<-ctx.Done()

	// Create shutdown context with 30-second timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Trigger graceful shutdown
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatal(err)
	}
}
