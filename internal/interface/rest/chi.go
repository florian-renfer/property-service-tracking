// Package rest defines http handler.
package rest

import (
	"net/http"

	"github.com/go-chi/render"
)

// Handler serves health responses.
type Handler struct{}

// NewHandler creates a health endpoint handler.
func NewHandler() *Handler {
	return &Handler{}
}

// Get returns service health.
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, map[string]string{
		"status": "ok",
	})
}
