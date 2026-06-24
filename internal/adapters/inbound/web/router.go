package web

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

// Dependencies contains HTTP adapter dependencies.
type Dependencies struct {
	HealthHandler   HealthEndpoint
	PropertyHandler PropertyEndpoint
}

// HealthEndpoint is the inbound HTTP contract for health checks.
type HealthEndpoint interface {
	Get(http.ResponseWriter, *http.Request)
}

// PropertyEndpoint is the inbound HTTP contract for property management.
type PropertyEndpoint interface {
	Create(http.ResponseWriter, *http.Request)
	List(http.ResponseWriter, *http.Request)
	Get(http.ResponseWriter, *http.Request)
	Rename(http.ResponseWriter, *http.Request)
}

// NewRouter builds the HTTP router.
func NewRouter(deps Dependencies) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.ClientIPFromRemoteAddr)
	r.Use(middleware.Logger)
	r.Use(middleware.StripSlashes)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Route("/api/v1", func(v1 chi.Router) {
		v1.Get("/health", deps.HealthHandler.Get)
		if deps.PropertyHandler != nil {
			v1.Route("/properties", func(properties chi.Router) {
				properties.Post("/", deps.PropertyHandler.Create)
				properties.Get("/", deps.PropertyHandler.List)
				properties.Get("/{id}", deps.PropertyHandler.Get)
				properties.Patch("/{id}", deps.PropertyHandler.Rename)
			})
		}
	})

	return r
}
