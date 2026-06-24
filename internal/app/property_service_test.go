package app_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/florian-renfer/property-service-tracking/internal/app"
	"github.com/florian-renfer/property-service-tracking/internal/app/ports"
	"github.com/florian-renfer/property-service-tracking/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	errCreateProperty = errors.New("create property failed")
	errFindProperties = errors.New("find properties failed")
	errFindProperty   = errors.New("find property failed")
	errUpdateProperty = errors.New("update property failed")
	errPropertyMiss   = errors.New("property not found")
)

func TestNewPropertyService(t *testing.T) {
	t.Run("nil repository", func(t *testing.T) {
		service, err := app.NewPropertyService(nil)

		require.ErrorIs(t, err, app.ErrPropertyRepositoryNil)
		assert.Nil(t, service)
	})

	t.Run("valid repository", func(t *testing.T) {
		service, err := app.NewPropertyService(newMemoryPropertyRepository())

		require.NoError(t, err)
		assert.NotNil(t, service)
	})
}

func TestPropertyServiceCreateProperty(t *testing.T) {
	ctx := context.Background()
	createdBy := uuid.New()

	t.Run("creates and persists property", func(t *testing.T) {
		repo := newMemoryPropertyRepository()
		service := mustNewPropertyService(t, repo)
		before := time.Now().UTC()

		property, err := service.CreateProperty(ctx, app.CreatePropertyInput{
			Title:     "  Fronhofallee 40-42  ",
			CreatedBy: createdBy,
		})

		after := time.Now().UTC()
		require.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, property.ID)
		assert.Equal(t, "Fronhofallee 40-42", property.Title)
		assert.Equal(t, createdBy, property.CreatedBy)
		assert.Equal(t, createdBy, property.UpdatedBy)
		assert.WithinRange(t, property.CreatedAt, before, after)
		assert.Equal(t, property.CreatedAt, property.UpdatedAt)
		assert.Equal(t, 1, repo.createCalls)

		persisted, err := repo.Find(ctx, property.ID)
		require.NoError(t, err)
		assert.Equal(t, property.ID, persisted.ID())
		assert.Equal(t, property.Title, persisted.Title())
	})

	t.Run("domain validation error does not persist", func(t *testing.T) {
		tests := []struct {
			name    string
			input   app.CreatePropertyInput
			wantErr error
		}{
			{
				name: "blank title",
				input: app.CreatePropertyInput{
					Title:     "   ",
					CreatedBy: createdBy,
				},
				wantErr: domain.ErrPropertyTitleBlank,
			},
			{
				name: "nil creator",
				input: app.CreatePropertyInput{
					Title:     "Fronhofallee 40-42",
					CreatedBy: uuid.Nil,
				},
				wantErr: domain.ErrPropertyCreatedByNil,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				repo := newMemoryPropertyRepository()
				service := mustNewPropertyService(t, repo)

				property, err := service.CreateProperty(ctx, tt.input)

				require.ErrorIs(t, err, tt.wantErr)
				assert.Empty(t, property)
				assert.Equal(t, 0, repo.createCalls)
				assert.Empty(t, repo.properties)
			})
		}
	})

	t.Run("repository error is propagated", func(t *testing.T) {
		repo := newMemoryPropertyRepository()
		repo.createErr = errCreateProperty
		service := mustNewPropertyService(t, repo)

		property, err := service.CreateProperty(ctx, app.CreatePropertyInput{
			Title:     "Fronhofallee 40-42",
			CreatedBy: createdBy,
		})

		require.ErrorIs(t, err, errCreateProperty)
		assert.Empty(t, property)
		assert.Equal(t, 1, repo.createCalls)
		assert.Empty(t, repo.properties)
	})
}

func TestPropertyServiceListProperties(t *testing.T) {
	ctx := context.Background()

	t.Run("returns empty slice", func(t *testing.T) {
		service := mustNewPropertyService(t, newMemoryPropertyRepository())

		properties, err := service.ListProperties(ctx)

		require.NoError(t, err)
		assert.Empty(t, properties)
	})

	t.Run("returns stored properties", func(t *testing.T) {
		repo := newMemoryPropertyRepository()
		first := mustNewDomainProperty(t, "Fronhofallee 40-42")
		second := mustNewDomainProperty(t, "Bahnhofstrasse 1")
		repo.mustStore(t, first, second)
		service := mustNewPropertyService(t, repo)

		properties, err := service.ListProperties(ctx)

		require.NoError(t, err)
		require.Len(t, properties, 2)
		assert.Equal(t, first.ID(), properties[0].ID)
		assert.Equal(t, first.Title(), properties[0].Title)
		assert.Equal(t, second.ID(), properties[1].ID)
		assert.Equal(t, second.Title(), properties[1].Title)
		assert.Equal(t, 1, repo.findAllCalls)
	})

	t.Run("repository error is propagated", func(t *testing.T) {
		repo := newMemoryPropertyRepository()
		repo.findAllErr = errFindProperties
		service := mustNewPropertyService(t, repo)

		properties, err := service.ListProperties(ctx)

		require.ErrorIs(t, err, errFindProperties)
		assert.Nil(t, properties)
		assert.Equal(t, 1, repo.findAllCalls)
	})
}

func TestPropertyServiceGetProperty(t *testing.T) {
	ctx := context.Background()

	t.Run("returns stored property", func(t *testing.T) {
		repo := newMemoryPropertyRepository()
		stored := mustNewDomainProperty(t, "Fronhofallee 40-42")
		repo.mustStore(t, stored)
		service := mustNewPropertyService(t, repo)

		property, err := service.GetProperty(ctx, stored.ID())

		require.NoError(t, err)
		assertPropertyOutput(t, stored, property)
		assert.Equal(t, 1, repo.findCalls)
	})

	t.Run("nil id does not query repository", func(t *testing.T) {
		repo := newMemoryPropertyRepository()
		service := mustNewPropertyService(t, repo)

		property, err := service.GetProperty(ctx, uuid.Nil)

		require.ErrorIs(t, err, domain.ErrPropertyIDNil)
		assert.Empty(t, property)
		assert.Equal(t, 0, repo.findCalls)
	})

	t.Run("repository error is propagated", func(t *testing.T) {
		repo := newMemoryPropertyRepository()
		repo.findErr = errFindProperty
		service := mustNewPropertyService(t, repo)

		property, err := service.GetProperty(ctx, uuid.New())

		require.ErrorIs(t, err, errFindProperty)
		assert.Empty(t, property)
		assert.Equal(t, 1, repo.findCalls)
	})

	t.Run("missing property error is propagated", func(t *testing.T) {
		service := mustNewPropertyService(t, newMemoryPropertyRepository())

		property, err := service.GetProperty(ctx, uuid.New())

		require.ErrorIs(t, err, errPropertyMiss)
		assert.Empty(t, property)
	})
}

func TestPropertyServiceRenameProperty(t *testing.T) {
	ctx := context.Background()
	updaterID := uuid.New()

	t.Run("renames stored property", func(t *testing.T) {
		repo := newMemoryPropertyRepository()
		stored := mustNewDomainProperty(t, "Fronhofallee 40-42")
		repo.mustStore(t, stored)
		service := mustNewPropertyService(t, repo)
		before := time.Now().UTC()

		property, err := service.RenameProperty(ctx, app.RenamePropertyInput{
			ID:        stored.ID(),
			Title:     "  Bahnhofstrasse 1  ",
			UpdatedBy: updaterID,
		})

		after := time.Now().UTC()
		require.NoError(t, err)
		assert.Equal(t, stored.ID(), property.ID)
		assert.Equal(t, "Bahnhofstrasse 1", property.Title)
		assert.Equal(t, stored.CreatedBy(), property.CreatedBy)
		assert.Equal(t, stored.CreatedAt(), property.CreatedAt)
		assert.Equal(t, updaterID, property.UpdatedBy)
		assert.WithinRange(t, property.UpdatedAt, before, after)
		assert.Equal(t, 1, repo.findCalls)
		assert.Equal(t, 1, repo.updateCalls)

		persisted, err := repo.Find(ctx, stored.ID())
		require.NoError(t, err)
		assert.Equal(t, "Bahnhofstrasse 1", persisted.Title())
		assert.Equal(t, updaterID, persisted.UpdatedBy())
	})

	t.Run("validation error does not update repository", func(t *testing.T) {
		tests := []struct {
			name      string
			input     app.RenamePropertyInput
			wantErr   error
			wantFinds int
		}{
			{
				name: "nil id",
				input: app.RenamePropertyInput{
					ID:        uuid.Nil,
					Title:     "Bahnhofstrasse 1",
					UpdatedBy: updaterID,
				},
				wantErr: domain.ErrPropertyIDNil,
			},
			{
				name: "blank title",
				input: app.RenamePropertyInput{
					ID:        uuid.New(),
					Title:     "   ",
					UpdatedBy: updaterID,
				},
				wantErr:   domain.ErrPropertyTitleBlank,
				wantFinds: 1,
			},
			{
				name: "nil updater",
				input: app.RenamePropertyInput{
					ID:        uuid.New(),
					Title:     "Bahnhofstrasse 1",
					UpdatedBy: uuid.Nil,
				},
				wantErr:   domain.ErrPropertyUpdatedByNil,
				wantFinds: 1,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				repo := newMemoryPropertyRepository()
				stored := mustNewDomainPropertyWithID(t, tt.input.ID, "Fronhofallee 40-42")
				if tt.input.ID != uuid.Nil {
					repo.mustStore(t, stored)
				}
				service := mustNewPropertyService(t, repo)

				property, err := service.RenameProperty(ctx, tt.input)

				require.ErrorIs(t, err, tt.wantErr)
				assert.Empty(t, property)
				assert.Equal(t, tt.wantFinds, repo.findCalls)
				assert.Equal(t, 0, repo.updateCalls)
			})
		}
	})

	t.Run("find error is propagated", func(t *testing.T) {
		repo := newMemoryPropertyRepository()
		repo.findErr = errFindProperty
		service := mustNewPropertyService(t, repo)

		property, err := service.RenameProperty(ctx, app.RenamePropertyInput{
			ID:        uuid.New(),
			Title:     "Bahnhofstrasse 1",
			UpdatedBy: updaterID,
		})

		require.ErrorIs(t, err, errFindProperty)
		assert.Empty(t, property)
		assert.Equal(t, 1, repo.findCalls)
		assert.Equal(t, 0, repo.updateCalls)
	})

	t.Run("missing property error is propagated", func(t *testing.T) {
		service := mustNewPropertyService(t, newMemoryPropertyRepository())

		property, err := service.RenameProperty(ctx, app.RenamePropertyInput{
			ID:        uuid.New(),
			Title:     "Bahnhofstrasse 1",
			UpdatedBy: updaterID,
		})

		require.ErrorIs(t, err, errPropertyMiss)
		assert.Empty(t, property)
	})

	t.Run("update error is propagated", func(t *testing.T) {
		repo := newMemoryPropertyRepository()
		stored := mustNewDomainProperty(t, "Fronhofallee 40-42")
		repo.mustStore(t, stored)
		repo.updateErr = errUpdateProperty
		service := mustNewPropertyService(t, repo)

		property, err := service.RenameProperty(ctx, app.RenamePropertyInput{
			ID:        stored.ID(),
			Title:     "Bahnhofstrasse 1",
			UpdatedBy: updaterID,
		})

		require.ErrorIs(t, err, errUpdateProperty)
		assert.Empty(t, property)
		assert.Equal(t, 1, repo.findCalls)
		assert.Equal(t, 1, repo.updateCalls)

		persisted, err := repo.Find(ctx, stored.ID())
		require.NoError(t, err)
		assert.Equal(t, stored.Title(), persisted.Title())
	})
}

type memoryPropertyRepository struct {
	properties map[uuid.UUID]domain.Property
	order      []uuid.UUID

	createErr  error
	findAllErr error
	findErr    error
	updateErr  error

	createCalls  int
	findAllCalls int
	findCalls    int
	updateCalls  int
}

var _ ports.PropertyRepository = (*memoryPropertyRepository)(nil)

func newMemoryPropertyRepository() *memoryPropertyRepository {
	return &memoryPropertyRepository{
		properties: make(map[uuid.UUID]domain.Property),
	}
}

func (r *memoryPropertyRepository) Create(_ context.Context, property domain.Property) error {
	r.createCalls++
	if r.createErr != nil {
		return r.createErr
	}

	if _, exists := r.properties[property.ID()]; !exists {
		r.order = append(r.order, property.ID())
	}
	r.properties[property.ID()] = property
	return nil
}

func (r *memoryPropertyRepository) FindAll(_ context.Context) ([]domain.Property, error) {
	r.findAllCalls++
	if r.findAllErr != nil {
		return nil, r.findAllErr
	}

	properties := make([]domain.Property, 0, len(r.order))
	for _, id := range r.order {
		properties = append(properties, r.properties[id])
	}
	return properties, nil
}

func (r *memoryPropertyRepository) Find(_ context.Context, id uuid.UUID) (domain.Property, error) {
	r.findCalls++
	if r.findErr != nil {
		return domain.Property{}, r.findErr
	}

	property, ok := r.properties[id]
	if !ok {
		return domain.Property{}, errPropertyMiss
	}
	return property, nil
}

func (r *memoryPropertyRepository) Update(_ context.Context, property domain.Property) error {
	r.updateCalls++
	if r.updateErr != nil {
		return r.updateErr
	}
	if _, ok := r.properties[property.ID()]; !ok {
		return errPropertyMiss
	}

	r.properties[property.ID()] = property
	return nil
}

func (r *memoryPropertyRepository) mustStore(t *testing.T, properties ...domain.Property) {
	t.Helper()

	for _, property := range properties {
		require.NoError(t, r.Create(context.Background(), property))
	}
	r.createCalls = 0
}

func mustNewPropertyService(t *testing.T, repo ports.PropertyRepository) *app.PropertyService {
	t.Helper()

	service, err := app.NewPropertyService(repo)
	require.NoError(t, err)
	return service
}

func mustNewDomainProperty(t *testing.T, title string) domain.Property {
	t.Helper()

	return mustNewDomainPropertyWithID(t, uuid.New(), title)
}

func mustNewDomainPropertyWithID(t *testing.T, id uuid.UUID, title string) domain.Property {
	t.Helper()

	now := time.Now().UTC().Add(-time.Hour)
	if id == uuid.Nil {
		id = uuid.New()
	}

	property, err := domain.NewProperty(id, title, uuid.New(), now, uuid.New(), now)
	require.NoError(t, err)
	return property
}

func assertPropertyOutput(t *testing.T, property domain.Property, output app.PropertyOutput) {
	t.Helper()

	assert.Equal(t, property.ID(), output.ID)
	assert.Equal(t, property.Title(), output.Title)
	assert.Equal(t, property.CreatedBy(), output.CreatedBy)
	assert.Equal(t, property.CreatedAt(), output.CreatedAt)
	assert.Equal(t, property.UpdatedBy(), output.UpdatedBy)
	assert.Equal(t, property.UpdatedAt(), output.UpdatedAt)
}
