package web

import (
	"net/http"

	"github.com/go-chi/render"
)

// HealthHandler serves health responses.
type HealthHandler struct{}

// NewHealthHandler creates a health endpoint handler.
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Get returns service health.
func (h *HealthHandler) Get(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, map[string]string{
		"status": "ok",
	})
}
