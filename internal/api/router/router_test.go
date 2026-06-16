package router

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/florian-renfer/property-service-tracking/internal/interface/rest"
	"github.com/florian-renfer/property-service-tracking/internal/interface/rest/auth"
)

func TestNew_WiresHealthRoutes(t *testing.T) {
	t.Parallel()

	handler := New(Dependencies{
		HealthHandler:  rest.NewHandler(),
		MeHandler:      rest.NewMeHandler(),
		AuthMiddleware: passthroughAuthMiddleware{},
	})

	paths := []string{"/api/v1/health", "/api/v1/health/"}
	for _, path := range paths {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		resp := httptest.NewRecorder()
		handler.ServeHTTP(resp, req)
		if resp.Code != http.StatusOK {
			t.Fatalf("GET %s status = %d, want %d", path, resp.Code, http.StatusOK)
		}
		contentType := resp.Header().Get("Content-Type")
		if !strings.HasPrefix(contentType, "application/json") {
			t.Fatalf("GET %s content-type = %q, want application/json", path, contentType)
		}
		var got map[string]string
		if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
			t.Fatalf("GET %s decode body: %v", path, err)
		}
		if got["status"] != "ok" {
			t.Fatalf("GET %s body status = %q, want %q", path, got["status"], "ok")
		}
	}
}

func TestNew_WiresMeRoute(t *testing.T) {
	t.Parallel()

	handler := New(Dependencies{
		HealthHandler: rest.NewHandler(),
		MeHandler:     rest.NewMeHandler(),
		AuthMiddleware: injectUserAuthMiddleware{
			user: auth.AuthenticatedUser{
				Subject: "user-123",
				Email:   "user@example.com",
				Roles:   []string{"GLOBAL_ADMIN"},
			},
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/me", nil)
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("GET /api/v1/me status = %d, want %d", resp.Code, http.StatusOK)
	}

	var got auth.AuthenticatedUser
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if got.Subject != "user-123" || got.Email != "user@example.com" || len(got.Roles) != 1 || got.Roles[0] != "GLOBAL_ADMIN" {
		t.Fatalf("body = %#v", got)
	}
}

func TestNew_WiresMeRouteUnauthorized(t *testing.T) {
	t.Parallel()

	handler := New(Dependencies{
		HealthHandler:  rest.NewHandler(),
		MeHandler:      rest.NewMeHandler(),
		AuthMiddleware: passthroughAuthMiddleware{},
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/me", nil)
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)

	if resp.Code != http.StatusUnauthorized {
		t.Fatalf("GET /api/v1/me status = %d, want %d", resp.Code, http.StatusUnauthorized)
	}
}

type passthroughAuthMiddleware struct{}

func (passthroughAuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return next
}

type injectUserAuthMiddleware struct {
	user auth.AuthenticatedUser
}

func (m injectUserAuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r.WithContext(auth.ContextWithUser(r.Context(), m.user)))
	})
}
