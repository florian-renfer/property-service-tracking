// Package router provides API router composition.
package router

import (
	"net/http"

	v1 "github.com/florian-renfer/property-service-tracking/internal/api/v1"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

// Dependencies contains dependencies for route wiring.
type Dependencies struct {
	HealthHandler  v1.HealthEndpoint
	MeHandler      v1.MeEndpoint
	AuthMiddleware v1.AuthMiddleware
}

// New builds the full HTTP router.
func New(deps Dependencies) http.Handler {
	r := chi.NewRouter()

	// Ingress
	r.Use(middleware.RequestID)
	r.Use(middleware.ClientIPFromRemoteAddr)
	r.Use(middleware.Logger)
	r.Use(middleware.StripSlashes)
	r.Use(middleware.Recoverer)

	// Egress
	r.Use(middleware.Compress(5))
	r.Use(render.SetContentType(render.ContentTypeJSON))

	v1.RegisterRoutes(r, v1.Dependencies{
		HealthHandler:  deps.HealthHandler,
		MeHandler:      deps.MeHandler,
		AuthMiddleware: deps.AuthMiddleware,
	})

	return r
}
