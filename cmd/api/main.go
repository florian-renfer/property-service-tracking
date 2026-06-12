// Package main starts the API service.
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/florian-renfer/property-service-tracking/internal/api/router"
	"github.com/florian-renfer/property-service-tracking/internal/interface/rest"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// Establish database connection
	connect()

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

func connect() {
	databaseURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?pool_max_conns=10&pool_max_conn_lifetime=30m",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)

	dbpool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	var greeting string
	err = dbpool.QueryRow(context.Background(), "select 'Hello, world!'").Scan(&greeting)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(greeting)
}
