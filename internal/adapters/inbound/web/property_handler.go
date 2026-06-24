package web

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/florian-renfer/property-service-tracking/internal/app"
	"github.com/florian-renfer/property-service-tracking/internal/app/ports"
	"github.com/florian-renfer/property-service-tracking/internal/domain"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

const principalIDHeader = "X-Principal-ID"

// PropertyHandler serves property endpoints.
type PropertyHandler struct {
	service *app.PropertyService
}

type (
	propertyRequest struct {
		Title string `json:"title"`
	}

	propertyResponse struct {
		ID        uuid.UUID `json:"id"`
		Title     string    `json:"title"`
		CreatedBy uuid.UUID `json:"created_by"`
		CreatedAt string    `json:"created_at"`
		UpdatedBy uuid.UUID `json:"updated_by"`
		UpdatedAt string    `json:"updated_at"`
	}

	errorResponse struct {
		Error string `json:"error"`
	}
)

// NewPropertyHandler creates a property REST handler.
func NewPropertyHandler(service *app.PropertyService) *PropertyHandler {
	return &PropertyHandler{service: service}
}

// Create creates a property.
func (h *PropertyHandler) Create(w http.ResponseWriter, r *http.Request) {
	principalID, ok := principalID(w, r)
	if !ok {
		return
	}

	var req propertyRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	property, err := h.service.CreateProperty(r.Context(), app.CreatePropertyInput{
		Title:     req.Title,
		CreatedBy: principalID,
	})
	if err != nil {
		writePropertyError(w, r, err)
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, newPropertyResponse(property))
}

// List returns all properties.
func (h *PropertyHandler) List(w http.ResponseWriter, r *http.Request) {
	properties, err := h.service.ListProperties(r.Context())
	if err != nil {
		writePropertyError(w, r, err)
		return
	}

	response := make([]propertyResponse, 0, len(properties))
	for _, property := range properties {
		response = append(response, newPropertyResponse(property))
	}

	render.JSON(w, r, response)
}

// Get returns a property.
func (h *PropertyHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, ok := propertyID(w, r)
	if !ok {
		return
	}

	property, err := h.service.GetProperty(r.Context(), id)
	if err != nil {
		writePropertyError(w, r, err)
		return
	}

	render.JSON(w, r, newPropertyResponse(property))
}

// Rename renames a property.
func (h *PropertyHandler) Rename(w http.ResponseWriter, r *http.Request) {
	principalID, ok := principalID(w, r)
	if !ok {
		return
	}
	id, ok := propertyID(w, r)
	if !ok {
		return
	}

	var req propertyRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	property, err := h.service.RenameProperty(r.Context(), app.RenamePropertyInput{
		ID:        id,
		Title:     req.Title,
		UpdatedBy: principalID,
	})
	if err != nil {
		writePropertyError(w, r, err)
		return
	}

	render.JSON(w, r, newPropertyResponse(property))
}

func propertyID(w http.ResponseWriter, r *http.Request) (uuid.UUID, bool) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil || id == uuid.Nil {
		writeError(w, r, http.StatusBadRequest, "invalid property id")
		return uuid.Nil, false
	}
	return id, true
}

func principalID(w http.ResponseWriter, r *http.Request) (uuid.UUID, bool) {
	id, err := uuid.Parse(r.Header.Get(principalIDHeader))
	if err != nil || id == uuid.Nil {
		writeError(w, r, http.StatusUnauthorized, "missing or invalid principal")
		return uuid.Nil, false
	}
	return id, true
}

func decodeJSON(w http.ResponseWriter, r *http.Request, dst any) bool {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(dst); err != nil {
		writeError(w, r, http.StatusBadRequest, "invalid json")
		return false
	}
	return true
}

func writePropertyError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, ports.ErrPropertyNotFound):
		writeError(w, r, http.StatusNotFound, err.Error())
	case errors.Is(err, domain.ErrPropertyIDNil),
		errors.Is(err, domain.ErrPropertyTitleBlank),
		errors.Is(err, domain.ErrPropertyTitleTooLong),
		errors.Is(err, domain.ErrPropertyCreatedByNil),
		errors.Is(err, domain.ErrPropertyCreatedAtZero),
		errors.Is(err, domain.ErrPropertyUpdatedByNil),
		errors.Is(err, domain.ErrPropertyUpdatedAtZero),
		errors.Is(err, domain.ErrPropertyUpdatedAtBeforeCreatedAt):
		writeError(w, r, http.StatusBadRequest, err.Error())
	default:
		writeError(w, r, http.StatusInternalServerError, "internal server error")
	}
}

func writeError(w http.ResponseWriter, r *http.Request, status int, message string) {
	render.Status(r, status)
	render.JSON(w, r, errorResponse{Error: message})
}

func newPropertyResponse(property app.PropertyOutput) propertyResponse {
	return propertyResponse{
		ID:        property.ID,
		Title:     property.Title,
		CreatedBy: property.CreatedBy,
		CreatedAt: property.CreatedAt.Format(time.RFC3339Nano),
		UpdatedBy: property.UpdatedBy,
		UpdatedAt: property.UpdatedAt.Format(time.RFC3339Nano),
	}
}
