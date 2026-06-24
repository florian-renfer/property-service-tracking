package memory_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/florian-renfer/property-service-tracking/internal/adapters/outbound/memory"
	"github.com/florian-renfer/property-service-tracking/internal/app/ports"
	"github.com/florian-renfer/property-service-tracking/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPropertyRepositoryCreateFindAndList(t *testing.T) {
	ctx := context.Background()
	repo := memory.NewPropertyRepository()
	first := mustNewProperty(t, uuid.New(), "Fronhofallee 40-42")
	second := mustNewProperty(t, uuid.New(), "Bahnhofstrasse 1")

	require.NoError(t, repo.Create(ctx, first))
	require.NoError(t, repo.Create(ctx, second))

	found, err := repo.Find(ctx, first.ID())
	require.NoError(t, err)
	assert.Equal(t, first.ID(), found.ID())
	assert.Equal(t, first.Title(), found.Title())

	properties, err := repo.FindAll(ctx)
	require.NoError(t, err)
	require.Len(t, properties, 2)
	assert.Equal(t, first.ID(), properties[0].ID())
	assert.Equal(t, second.ID(), properties[1].ID())
}

func TestPropertyRepositoryFindAllEmpty(t *testing.T) {
	repo := memory.NewPropertyRepository()

	properties, err := repo.FindAll(context.Background())

	require.NoError(t, err)
	assert.Empty(t, properties)
}

func TestPropertyRepositoryFindMissing(t *testing.T) {
	repo := memory.NewPropertyRepository()

	property, err := repo.Find(context.Background(), uuid.New())

	require.ErrorIs(t, err, ports.ErrPropertyNotFound)
	assert.Empty(t, property)
}

func TestPropertyRepositoryUpdate(t *testing.T) {
	ctx := context.Background()
	repo := memory.NewPropertyRepository()
	property := mustNewProperty(t, uuid.New(), "Fronhofallee 40-42")
	require.NoError(t, repo.Create(ctx, property))
	require.NoError(t, property.Rename("Bahnhofstrasse 1", uuid.New(), time.Now().UTC()))

	require.NoError(t, repo.Update(ctx, property))

	found, err := repo.Find(ctx, property.ID())
	require.NoError(t, err)
	assert.Equal(t, "Bahnhofstrasse 1", found.Title())
}

func TestPropertyRepositoryUpdateMissing(t *testing.T) {
	repo := memory.NewPropertyRepository()
	property := mustNewProperty(t, uuid.New(), "Fronhofallee 40-42")

	err := repo.Update(context.Background(), property)

	require.ErrorIs(t, err, ports.ErrPropertyNotFound)
}

func TestPropertyRepositoryNotFoundSupportsErrorsIs(t *testing.T) {
	err := ports.ErrPropertyNotFound

	assert.True(t, errors.Is(err, ports.ErrPropertyNotFound))
}

func mustNewProperty(t *testing.T, id uuid.UUID, title string) domain.Property {
	t.Helper()

	now := time.Now().UTC().Add(-time.Hour)
	property, err := domain.NewProperty(id, title, uuid.New(), now, uuid.New(), now)
	require.NoError(t, err)
	return property
}
