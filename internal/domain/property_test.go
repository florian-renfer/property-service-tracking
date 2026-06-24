package domain_test

import (
	"strings"
	"testing"
	"time"

	"github.com/florian-renfer/property-service-tracking/internal/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type propertyArgs struct {
	id        uuid.UUID
	title     string
	createdBy uuid.UUID
	createdAt time.Time
	updatedBy uuid.UUID
	updatedAt time.Time
}

var (
	propertyID    = uuid.MustParse("daac0cb9-c0db-42a4-8f77-fdadbd9c4bc3")
	creatorID     = uuid.MustParse("9f6faa73-51aa-4840-89fe-ea62fd1fc5da")
	updaterID     = uuid.MustParse("673b2070-5c63-4f48-b656-fc097cdf7699")
	createdAt     = time.Date(2026, 6, 24, 8, 0, 0, 0, time.UTC)
	updatedAt     = time.Date(2026, 6, 24, 9, 0, 0, 0, time.UTC)
	validProperty = propertyArgs{
		id:        propertyID,
		title:     "Fronhofallee 40-42",
		createdBy: creatorID,
		createdAt: createdAt,
		updatedBy: updaterID,
		updatedAt: updatedAt,
	}
)

func TestNewProperty(t *testing.T) {
	tests := []struct {
		name    string
		args    propertyArgs
		wantErr error
	}{
		{
			name: "valid property",
			args: validProperty,
		},
		{
			name: "title trimmed",
			args: withPropertyArgs(func(args *propertyArgs) {
				args.title = "  Fronhofallee 40-42  "
			}),
		},
		{
			name: "nil id",
			args: withPropertyArgs(func(args *propertyArgs) {
				args.id = uuid.Nil
			}),
			wantErr: domain.ErrPropertyIDNil,
		},
		{
			name: "blank title",
			args: withPropertyArgs(func(args *propertyArgs) {
				args.title = ""
			}),
			wantErr: domain.ErrPropertyTitleBlank,
		},
		{
			name: "whitespace title",
			args: withPropertyArgs(func(args *propertyArgs) {
				args.title = "   "
			}),
			wantErr: domain.ErrPropertyTitleBlank,
		},
		{
			name: "title too long",
			args: withPropertyArgs(func(args *propertyArgs) {
				args.title = strings.Repeat("a", domain.PropertyTitleMaxLen+1)
			}),
			wantErr: domain.ErrPropertyTitleTooLong,
		},
		{
			name: "nil created by",
			args: withPropertyArgs(func(args *propertyArgs) {
				args.createdBy = uuid.Nil
			}),
			wantErr: domain.ErrPropertyCreatedByNil,
		},
		{
			name: "zero created at",
			args: withPropertyArgs(func(args *propertyArgs) {
				args.createdAt = time.Time{}
			}),
			wantErr: domain.ErrPropertyCreatedAtZero,
		},
		{
			name: "nil updated by",
			args: withPropertyArgs(func(args *propertyArgs) {
				args.updatedBy = uuid.Nil
			}),
			wantErr: domain.ErrPropertyUpdatedByNil,
		},
		{
			name: "zero updated at",
			args: withPropertyArgs(func(args *propertyArgs) {
				args.updatedAt = time.Time{}
			}),
			wantErr: domain.ErrPropertyUpdatedAtZero,
		},
		{
			name: "updated at before created at",
			args: withPropertyArgs(func(args *propertyArgs) {
				args.updatedAt = createdAt.Add(-time.Second)
			}),
			wantErr: domain.ErrPropertyUpdatedAtBeforeCreatedAt,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			property, err := domain.NewProperty(
				tt.args.id,
				tt.args.title,
				tt.args.createdBy,
				tt.args.createdAt,
				tt.args.updatedBy,
				tt.args.updatedAt,
			)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.args.id, property.ID())
			assert.Equal(t, strings.TrimSpace(tt.args.title), property.Title())
			assert.Equal(t, tt.args.createdBy, property.CreatedBy())
			assert.Equal(t, tt.args.createdAt, property.CreatedAt())
			assert.Equal(t, tt.args.updatedBy, property.UpdatedBy())
			assert.Equal(t, tt.args.updatedAt, property.UpdatedAt())
		})
	}
}

func TestPropertyRename(t *testing.T) {
	property := mustNewProperty(t)
	newUpdatedAt := updatedAt.Add(time.Hour)

	err := property.Rename("  Bahnhofstrasse 1  ", creatorID, newUpdatedAt)
	require.NoError(t, err)
	assert.Equal(t, "Bahnhofstrasse 1", property.Title())
	assert.Equal(t, creatorID, property.UpdatedBy())
	assert.Equal(t, newUpdatedAt, property.UpdatedAt())
	assert.Equal(t, validProperty.createdBy, property.CreatedBy())
	assert.Equal(t, validProperty.createdAt, property.CreatedAt())
}

func TestPropertyRenameValidationDoesNotMutate(t *testing.T) {
	tests := []struct {
		name      string
		title     string
		updatedBy uuid.UUID
		updatedAt time.Time
		wantErr   error
	}{
		{
			name:      "blank title",
			title:     "   ",
			updatedBy: updaterID,
			updatedAt: updatedAt.Add(time.Hour),
			wantErr:   domain.ErrPropertyTitleBlank,
		},
		{
			name:      "title too long",
			title:     strings.Repeat("a", domain.PropertyTitleMaxLen+1),
			updatedBy: updaterID,
			updatedAt: updatedAt.Add(time.Hour),
			wantErr:   domain.ErrPropertyTitleTooLong,
		},
		{
			name:      "nil updated by",
			title:     "Bahnhofstrasse 1",
			updatedBy: uuid.Nil,
			updatedAt: updatedAt.Add(time.Hour),
			wantErr:   domain.ErrPropertyUpdatedByNil,
		},
		{
			name:      "zero updated at",
			title:     "Bahnhofstrasse 1",
			updatedBy: updaterID,
			updatedAt: time.Time{},
			wantErr:   domain.ErrPropertyUpdatedAtZero,
		},
		{
			name:      "updated before created",
			title:     "Bahnhofstrasse 1",
			updatedBy: updaterID,
			updatedAt: createdAt.Add(-time.Second),
			wantErr:   domain.ErrPropertyUpdatedAtBeforeCreatedAt,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			property := mustNewProperty(t)

			err := property.Rename(tt.title, tt.updatedBy, tt.updatedAt)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, validProperty.title, property.Title())
			assert.Equal(t, validProperty.updatedBy, property.UpdatedBy())
			assert.Equal(t, validProperty.updatedAt, property.UpdatedAt())
		})
	}
}

func withPropertyArgs(change func(*propertyArgs)) propertyArgs {
	args := validProperty
	change(&args)
	return args
}

func mustNewProperty(t *testing.T) domain.Property {
	t.Helper()

	property, err := domain.NewProperty(
		validProperty.id,
		validProperty.title,
		validProperty.createdBy,
		validProperty.createdAt,
		validProperty.updatedBy,
		validProperty.updatedAt,
	)
	require.NoError(t, err)
	return property
}
