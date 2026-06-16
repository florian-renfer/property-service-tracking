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
	restauth "github.com/florian-renfer/property-service-tracking/internal/interface/rest/auth"

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
	healthHandler := rest.NewHandler()
	meHandler := rest.NewMeHandler()

	authCtx, cancelAuth := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancelAuth()

	issuerURL := getEnvOrDefault(
		"KEYCLOAK_ISSUER_URL",
		fmt.Sprintf(
			"http://%s:%s/realms/property-service-tracking",
			getEnvOrDefault("KEYCLOAK_URL", "127.0.0.1"),
			getEnvOrDefault("KEYCLOAK_PORT", "8090"),
		),
	)

	authMiddleware, err := restauth.NewMiddleware(
		authCtx,
		issuerURL,
		getEnvOrDefault("KEYCLOAK_AUDIENCE", "property-service-api"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// The HTTP Server
	port := getEnvOrDefault("API_PORT", "4000")
	server := &http.Server{Addr: "0.0.0.0:" + port, Handler: router.New(router.Dependencies{HealthHandler: healthHandler, MeHandler: meHandler, AuthMiddleware: authMiddleware})}

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

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
