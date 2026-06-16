package rest

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/florian-renfer/property-service-tracking/internal/interface/rest/auth"
)

func TestMeHandlerGet(t *testing.T) {
	t.Parallel()

	handler := NewMeHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/me", nil)
	req = req.WithContext(auth.ContextWithUser(req.Context(), auth.AuthenticatedUser{
		Subject: "user-123",
		Email:   "user@example.com",
		Roles:   []string{"GLOBAL_ADMIN"},
	}))
	resp := httptest.NewRecorder()

	handler.Get(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.Code, http.StatusOK)
	}

	var got auth.AuthenticatedUser
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if got.Subject != "user-123" || got.Email != "user@example.com" || len(got.Roles) != 1 || got.Roles[0] != "GLOBAL_ADMIN" {
		t.Fatalf("body = %#v", got)
	}
}

func TestMeHandlerGetUnauthorized(t *testing.T) {
	t.Parallel()

	handler := NewMeHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/me", nil)
	resp := httptest.NewRecorder()

	handler.Get(resp, req)

	if resp.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", resp.Code, http.StatusUnauthorized)
	}
}
