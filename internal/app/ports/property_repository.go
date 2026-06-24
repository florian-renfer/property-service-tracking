package ports

import (
	"context"

	"github.com/florian-renfer/property-service-tracking/internal/domain"

	"github.com/google/uuid"
)

// PropertyRepository defines persistence operations for properties.
type PropertyRepository interface {
	Create(ctx context.Context, property domain.Property) error
	FindAll(ctx context.Context) ([]domain.Property, error)
	Find(ctx context.Context, id uuid.UUID) (domain.Property, error)
	Update(ctx context.Context, property domain.Property) error
}
