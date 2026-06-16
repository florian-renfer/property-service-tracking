// Package rest defines http handlers.
package rest

import (
	"net/http"

	"github.com/florian-renfer/property-service-tracking/internal/interface/rest/auth"
	"github.com/go-chi/render"
)

// MeHandler serves the authenticated user response.
type MeHandler struct{}

// NewMeHandler creates a handler for the current authenticated user.
func NewMeHandler() *MeHandler {
	return &MeHandler{}
}

// Get returns the authenticated user.
func (h *MeHandler) Get(w http.ResponseWriter, r *http.Request) {
	user, ok := auth.UserFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthenticated", http.StatusUnauthorized)
		return
	}

	render.JSON(w, r, user)
}
