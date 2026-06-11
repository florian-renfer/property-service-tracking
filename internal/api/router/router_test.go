package router

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNew_WiresHealthRoutes(t *testing.T) {
	t.Parallel()

	healthHandler := &stubHealthHandler{}
	handler := New(Dependencies{
		HealthHandler: healthHandler,
	})

	paths := []string{"/api/v1/health", "/api/v1/health/"}
	for _, path := range paths {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		resp := httptest.NewRecorder()
		handler.ServeHTTP(resp, req)
		if resp.Code != http.StatusNoContent {
			t.Fatalf("GET %s status = %d, want %d", path, resp.Code, http.StatusNoContent)
		}
	}

	if healthHandler.calls != len(paths) {
		t.Fatalf("health handler calls = %d, want %d", healthHandler.calls, len(paths))
	}
}

type stubHealthHandler struct {
	calls int
}

func (h *stubHealthHandler) Get(w http.ResponseWriter, _ *http.Request) {
	h.calls++
	w.WriteHeader(http.StatusNoContent)
}
