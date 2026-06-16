// Package v1 wires versioned API endpoints.
package v1

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// HealthEndpoint is the inbound HTTP contract for health checks.
type HealthEndpoint interface {
	Get(http.ResponseWriter, *http.Request)
}

// MeEndpoint is the inbound HTTP contract for the authenticated user endpoint.
type MeEndpoint interface {
	Get(http.ResponseWriter, *http.Request)
}

// AuthMiddleware is the inbound HTTP contract for auth protection.
type AuthMiddleware interface {
	Authenticate(http.Handler) http.Handler
}

// Dependencies contains dependencies for v1 route wiring.
type Dependencies struct {
	HealthHandler  HealthEndpoint
	MeHandler      MeEndpoint
	AuthMiddleware AuthMiddleware
}

// RegisterRoutes registers all v1 endpoints.
func RegisterRoutes(r chi.Router, deps Dependencies) {
	r.Route("/api/v1", func(v1 chi.Router) {
		v1.Get("/health", deps.HealthHandler.Get)
		v1.Group(func(r chi.Router) {
			r.Use(deps.AuthMiddleware.Authenticate)
			r.Get("/me", deps.MeHandler.Get)
		})
	})
}
