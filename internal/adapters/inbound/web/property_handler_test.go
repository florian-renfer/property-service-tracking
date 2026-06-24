package web_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/florian-renfer/property-service-tracking/internal/adapters/inbound/web"
	"github.com/florian-renfer/property-service-tracking/internal/adapters/outbound/memory"
	"github.com/florian-renfer/property-service-tracking/internal/app"
	"github.com/florian-renfer/property-service-tracking/internal/app/ports"
	"github.com/florian-renfer/property-service-tracking/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var errUnexpectedRepository = errors.New("unexpected repository error")

func TestPropertyHandlerFlow(t *testing.T) {
	handler := newTestRouter(t)
	principalID := uuid.New()

	createResp := propertyResponse{}
	doJSON(t, handler, http.MethodPost, "/api/v1/properties", principalID, map[string]string{
		"title": "  Fronhofallee 40-42  ",
	}, http.StatusCreated, &createResp)
	assert.NotEqual(t, uuid.Nil, createResp.ID)
	assert.Equal(t, "Fronhofallee 40-42", createResp.Title)
	assert.Equal(t, principalID, createResp.CreatedBy)
	assert.Equal(t, principalID, createResp.UpdatedBy)
	assert.NotEmpty(t, createResp.CreatedAt)
	assert.Equal(t, createResp.CreatedAt, createResp.UpdatedAt)

	getResp := propertyResponse{}
	doJSON(t, handler, http.MethodGet, "/api/v1/properties/"+createResp.ID.String(), uuid.Nil, nil, http.StatusOK, &getResp)
	assert.Equal(t, createResp, getResp)

	listResp := []propertyResponse{}
	doJSON(t, handler, http.MethodGet, "/api/v1/properties", uuid.Nil, nil, http.StatusOK, &listResp)
	require.Len(t, listResp, 1)
	assert.Equal(t, createResp.ID, listResp[0].ID)

	updaterID := uuid.New()
	renameResp := propertyResponse{}
	doJSON(t, handler, http.MethodPatch, "/api/v1/properties/"+createResp.ID.String(), updaterID, map[string]string{
		"title": "Bahnhofstrasse 1",
	}, http.StatusOK, &renameResp)
	assert.Equal(t, createResp.ID, renameResp.ID)
	assert.Equal(t, "Bahnhofstrasse 1", renameResp.Title)
	assert.Equal(t, principalID, renameResp.CreatedBy)
	assert.Equal(t, createResp.CreatedAt, renameResp.CreatedAt)
	assert.Equal(t, updaterID, renameResp.UpdatedBy)
	assert.NotEqual(t, createResp.UpdatedAt, renameResp.UpdatedAt)
}

func TestPropertyHandlerUsesPrincipalHeaderForAuditFields(t *testing.T) {
	handler := newTestRouter(t)
	principalID := uuid.New()

	doRawJSON(t, handler, http.MethodPost, "/api/v1/properties", principalID, `{
		"title": "Fronhofallee 40-42",
		"created_by": "00000000-0000-0000-0000-000000000001"
	}`, http.StatusBadRequest, nil)

	createResp := propertyResponse{}
	doJSON(t, handler, http.MethodPost, "/api/v1/properties", principalID, map[string]string{
		"title": "Fronhofallee 40-42",
	}, http.StatusCreated, &createResp)
	assert.Equal(t, principalID, createResp.CreatedBy)
}

func TestPropertyHandlerPrincipalValidation(t *testing.T) {
	tests := []struct {
		name      string
		method    string
		path      string
		principal string
		body      string
	}{
		{
			name:      "create missing principal",
			method:    http.MethodPost,
			path:      "/api/v1/properties",
			principal: "",
			body:      `{"title":"Fronhofallee 40-42"}`,
		},
		{
			name:      "create invalid principal",
			method:    http.MethodPost,
			path:      "/api/v1/properties",
			principal: "not-a-uuid",
			body:      `{"title":"Fronhofallee 40-42"}`,
		},
		{
			name:      "rename missing principal",
			method:    http.MethodPatch,
			path:      "/api/v1/properties/" + uuid.NewString(),
			principal: "",
			body:      `{"title":"Bahnhofstrasse 1"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := newTestRouter(t)
			req := httptest.NewRequest(tt.method, tt.path, bytes.NewBufferString(tt.body))
			if tt.principal != "" {
				req.Header.Set("X-Principal-ID", tt.principal)
			}
			resp := httptest.NewRecorder()

			handler.ServeHTTP(resp, req)

			assert.Equal(t, http.StatusUnauthorized, resp.Code)
			assertErrorBody(t, resp, "missing or invalid principal")
		})
	}
}

func TestPropertyHandlerBadRequests(t *testing.T) {
	handler := newTestRouter(t)
	principalID := uuid.New()

	tests := []struct {
		name   string
		method string
		path   string
		body   string
	}{
		{
			name:   "invalid json",
			method: http.MethodPost,
			path:   "/api/v1/properties",
			body:   `{`,
		},
		{
			name:   "blank title",
			method: http.MethodPost,
			path:   "/api/v1/properties",
			body:   `{"title":"   "}`,
		},
		{
			name:   "invalid path uuid",
			method: http.MethodGet,
			path:   "/api/v1/properties/not-a-uuid",
		},
		{
			name:   "rename invalid json",
			method: http.MethodPatch,
			path:   "/api/v1/properties/" + uuid.NewString(),
			body:   `{`,
		},
		{
			name:   "rename invalid path uuid",
			method: http.MethodPatch,
			path:   "/api/v1/properties/not-a-uuid",
			body:   `{"title":"Bahnhofstrasse 1"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var principal uuid.UUID
			if tt.method == http.MethodPost || tt.method == http.MethodPatch {
				principal = principalID
			}

			doRawJSON(t, handler, tt.method, tt.path, principal, tt.body, http.StatusBadRequest, nil)
		})
	}
}

func TestPropertyHandlerNotFound(t *testing.T) {
	handler := newTestRouter(t)
	missingID := uuid.New()

	doJSON(t, handler, http.MethodGet, "/api/v1/properties/"+missingID.String(), uuid.Nil, nil, http.StatusNotFound, nil)
	doJSON(t, handler, http.MethodPatch, "/api/v1/properties/"+missingID.String(), uuid.New(), map[string]string{
		"title": "Bahnhofstrasse 1",
	}, http.StatusNotFound, nil)
}

func TestPropertyHandlerUnexpectedRepositoryError(t *testing.T) {
	handler := newTestRouterWithRepository(t, failingPropertyRepository{})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/properties", nil)
	resp := httptest.NewRecorder()

	handler.ServeHTTP(resp, req)

	require.Equal(t, http.StatusInternalServerError, resp.Code)
	assertErrorBody(t, resp, "internal server error")
}

type propertyResponse struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	CreatedBy uuid.UUID `json:"created_by"`
	CreatedAt string    `json:"created_at"`
	UpdatedBy uuid.UUID `json:"updated_by"`
	UpdatedAt string    `json:"updated_at"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func newTestRouter(t *testing.T) http.Handler {
	t.Helper()

	repo := memory.NewPropertyRepository()
	return newTestRouterWithRepository(t, repo)
}

func newTestRouterWithRepository(t *testing.T, repo ports.PropertyRepository) http.Handler {
	t.Helper()

	service, err := app.NewPropertyService(repo)
	require.NoError(t, err)

	return web.NewRouter(web.Dependencies{
		HealthHandler:   web.NewHealthHandler(),
		PropertyHandler: web.NewPropertyHandler(service),
	})
}

func doJSON(t *testing.T, handler http.Handler, method, path string, principalID uuid.UUID, body any, wantStatus int, dst any) {
	t.Helper()

	var rawBody string
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		require.NoError(t, err)
		rawBody = string(bodyBytes)
	}
	doRawJSON(t, handler, method, path, principalID, rawBody, wantStatus, dst)
}

func doRawJSON(t *testing.T, handler http.Handler, method, path string, principalID uuid.UUID, body string, wantStatus int, dst any) {
	t.Helper()

	req := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	if principalID != uuid.Nil {
		req.Header.Set("X-Principal-ID", principalID.String())
	}
	resp := httptest.NewRecorder()

	handler.ServeHTTP(resp, req)

	require.Equal(t, wantStatus, resp.Code, resp.Body.String())
	if dst != nil {
		require.NoError(t, json.NewDecoder(resp.Body).Decode(dst))
	}
}

func assertErrorBody(t *testing.T, resp *httptest.ResponseRecorder, want string) {
	t.Helper()

	var got errorResponse
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&got))
	assert.Equal(t, want, got.Error)
}

type failingPropertyRepository struct{}

func (failingPropertyRepository) Create(context.Context, domain.Property) error {
	return errUnexpectedRepository
}

func (failingPropertyRepository) FindAll(context.Context) ([]domain.Property, error) {
	return nil, errUnexpectedRepository
}

func (failingPropertyRepository) Find(context.Context, uuid.UUID) (domain.Property, error) {
	return domain.Property{}, errUnexpectedRepository
}

func (failingPropertyRepository) Update(context.Context, domain.Property) error {
	return errUnexpectedRepository
}
