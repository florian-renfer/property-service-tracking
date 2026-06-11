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

// Dependencies contains dependencies for v1 route wiring.
type Dependencies struct {
	HealthHandler HealthEndpoint
}

// RegisterRoutes registers all v1 endpoints.
func RegisterRoutes(r chi.Router, deps Dependencies) {
	r.Route("/api/v1", func(v1 chi.Router) {
		v1.Get("/health", deps.HealthHandler.Get)
	})
}
