package web_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/florian-renfer/property-service-tracking/internal/adapters/inbound/web"
)

func TestNewRouterWiresHealthRoutes(t *testing.T) {
	t.Parallel()

	handler := web.NewRouter(web.Dependencies{
		HealthHandler: web.NewHealthHandler(),
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
