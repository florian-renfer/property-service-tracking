package app

import (
	"context"
	"errors"
	"time"

	"github.com/florian-renfer/property-service-tracking/internal/app/ports"
	"github.com/florian-renfer/property-service-tracking/internal/domain"

	"github.com/google/uuid"
)

// ErrPropertyRepositoryNil indicates the property service has no repository.
var ErrPropertyRepositoryNil = errors.New("property repository cannot be nil")

type (
	// PropertyService coordinates property use cases.
	PropertyService struct {
		repo ports.PropertyRepository
	}

	// CreatePropertyInput contains data required to create a property.
	CreatePropertyInput struct {
		Title     string
		CreatedBy uuid.UUID
	}

	// RenamePropertyInput contains data required to rename a property.
	RenamePropertyInput struct {
		ID        uuid.UUID
		Title     string
		UpdatedBy uuid.UUID
	}

	// PropertyOutput is the application-facing property representation.
	PropertyOutput struct {
		ID        uuid.UUID
		Title     string
		CreatedBy uuid.UUID
		CreatedAt time.Time
		UpdatedBy uuid.UUID
		UpdatedAt time.Time
	}
)

// NewPropertyService constructs a property application service.
func NewPropertyService(repo ports.PropertyRepository) (*PropertyService, error) {
	if repo == nil {
		return nil, ErrPropertyRepositoryNil
	}

	return &PropertyService{repo: repo}, nil
}

// CreateProperty creates and persists a property.
func (s *PropertyService) CreateProperty(ctx context.Context, input CreatePropertyInput) (PropertyOutput, error) {
	now := time.Now().UTC()
	property, err := domain.NewProperty(uuid.New(), input.Title, input.CreatedBy, now, input.CreatedBy, now)
	if err != nil {
		return PropertyOutput{}, err
	}

	if err := s.repo.Create(ctx, property); err != nil {
		return PropertyOutput{}, err
	}

	return newPropertyOutput(property), nil
}

// ListProperties returns all properties.
func (s *PropertyService) ListProperties(ctx context.Context) ([]PropertyOutput, error) {
	properties, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	outputs := make([]PropertyOutput, 0, len(properties))
	for _, property := range properties {
		outputs = append(outputs, newPropertyOutput(property))
	}

	return outputs, nil
}

// GetProperty returns a property by ID.
func (s *PropertyService) GetProperty(ctx context.Context, id uuid.UUID) (PropertyOutput, error) {
	if id == uuid.Nil {
		return PropertyOutput{}, domain.ErrPropertyIDNil
	}

	property, err := s.repo.Find(ctx, id)
	if err != nil {
		return PropertyOutput{}, err
	}

	return newPropertyOutput(property), nil
}

// RenameProperty updates a property's title and audit metadata.
func (s *PropertyService) RenameProperty(ctx context.Context, input RenamePropertyInput) (PropertyOutput, error) {
	if input.ID == uuid.Nil {
		return PropertyOutput{}, domain.ErrPropertyIDNil
	}

	property, err := s.repo.Find(ctx, input.ID)
	if err != nil {
		return PropertyOutput{}, err
	}

	if err := property.Rename(input.Title, input.UpdatedBy, time.Now().UTC()); err != nil {
		return PropertyOutput{}, err
	}

	if err := s.repo.Update(ctx, property); err != nil {
		return PropertyOutput{}, err
	}

	return newPropertyOutput(property), nil
}

func newPropertyOutput(property domain.Property) PropertyOutput {
	return PropertyOutput{
		ID:        property.ID(),
		Title:     property.Title(),
		CreatedBy: property.CreatedBy(),
		CreatedAt: property.CreatedAt(),
		UpdatedBy: property.UpdatedBy(),
		UpdatedAt: property.UpdatedAt(),
	}
}
