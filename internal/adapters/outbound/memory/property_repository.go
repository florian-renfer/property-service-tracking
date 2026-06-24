package memory

import (
	"context"
	"sync"

	"github.com/florian-renfer/property-service-tracking/internal/app/ports"
	"github.com/florian-renfer/property-service-tracking/internal/domain"
	"github.com/google/uuid"
)

// PropertyRepository stores properties in process memory.
type PropertyRepository struct {
	mu         sync.RWMutex
	properties map[uuid.UUID]domain.Property
	order      []uuid.UUID
}

var _ ports.PropertyRepository = (*PropertyRepository)(nil)

// NewPropertyRepository creates an empty in-memory property repository.
func NewPropertyRepository() *PropertyRepository {
	return &PropertyRepository{
		properties: make(map[uuid.UUID]domain.Property),
	}
}

// Create stores a property.
func (r *PropertyRepository) Create(_ context.Context, property domain.Property) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.properties[property.ID()]; !exists {
		r.order = append(r.order, property.ID())
	}
	r.properties[property.ID()] = property
	return nil
}

// FindAll returns all stored properties in insertion order.
func (r *PropertyRepository) FindAll(_ context.Context) ([]domain.Property, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	properties := make([]domain.Property, 0, len(r.order))
	for _, id := range r.order {
		properties = append(properties, r.properties[id])
	}
	return properties, nil
}

// Find returns a stored property by ID.
func (r *PropertyRepository) Find(_ context.Context, id uuid.UUID) (domain.Property, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	property, ok := r.properties[id]
	if !ok {
		return domain.Property{}, ports.ErrPropertyNotFound
	}
	return property, nil
}

// Update replaces a stored property.
func (r *PropertyRepository) Update(_ context.Context, property domain.Property) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.properties[property.ID()]; !ok {
		return ports.ErrPropertyNotFound
	}
	r.properties[property.ID()] = property
	return nil
}
