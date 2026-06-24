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

	"github.com/florian-renfer/property-service-tracking/internal/adapters/inbound/web"
	"github.com/florian-renfer/property-service-tracking/internal/adapters/outbound/memory"
	"github.com/florian-renfer/property-service-tracking/internal/app"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

func main() {
	// Establish database connection and run migrations
	dbpool := connect()
	defer dbpool.Close()

	// Router dependencies
	healthHandler := web.NewHealthHandler()

	propertyRepo := memory.NewPropertyRepository()
	propertyService, err := app.NewPropertyService(propertyRepo)
	if err != nil {
		log.Fatal(err)
	}
	propertyHandler := web.NewPropertyHandler(propertyService)

	// The HTTP Server
	server := &http.Server{Addr: "0.0.0.0:8080", Handler: web.NewRouter(web.Dependencies{
		HealthHandler:   healthHandler,
		PropertyHandler: propertyHandler,
	})}

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

func connect() *pgxpool.Pool {
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

	// Wrap pgxpool with stdlib interface for migrate
	db := stdlib.OpenDBFromPool(dbpool)
	defer db.Close()

	// Create postgres driver from pooled connection
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Migration driver setup failed: %v\n", err)
		os.Exit(1)
	}

	// Run migrations using pooled connection
	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"postgres",
		driver,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Migration setup failed: %v\n", err)
		os.Exit(1)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		fmt.Fprintf(os.Stderr, "Migration failed: %v\n", err)
		os.Exit(1)
	}

	return dbpool
}
