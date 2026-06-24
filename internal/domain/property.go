package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	// ErrPropertyIDNil indicates the property ID is missing.
	ErrPropertyIDNil = errors.New("property id cannot be nil")
	// ErrPropertyTitleBlank indicates the title is empty after trimming.
	ErrPropertyTitleBlank = errors.New("property title cannot be blank")
	// ErrPropertyTitleTooLong indicates the title exceeds max length.
	ErrPropertyTitleTooLong = errors.New("property title exceeds 255 characters")
	// ErrPropertyCreatedByNil indicates the creator user ID is missing.
	ErrPropertyCreatedByNil = errors.New("property created by cannot be nil")
	// ErrPropertyCreatedAtZero indicates the creation timestamp is missing.
	ErrPropertyCreatedAtZero = errors.New("property created at cannot be zero")
	// ErrPropertyUpdatedByNil indicates the updater user ID is missing.
	ErrPropertyUpdatedByNil = errors.New("property updated by cannot be nil")
	// ErrPropertyUpdatedAtZero indicates the update timestamp is missing.
	ErrPropertyUpdatedAtZero = errors.New("property updated at cannot be zero")
	// ErrPropertyUpdatedAtBeforeCreatedAt indicates the update timestamp predates creation.
	ErrPropertyUpdatedAtBeforeCreatedAt = errors.New("property updated at cannot be before created at")
)

const (
	// PropertyTitleMaxLen is max allowed length for a property title.
	PropertyTitleMaxLen = 255
)

type (
	// Property represents a managed real estate property.
	Property struct {
		id        uuid.UUID
		title     string
		createdBy uuid.UUID
		createdAt time.Time
		updatedBy uuid.UUID
		updatedAt time.Time
	}
)

// NewProperty constructs a property with validation.
func NewProperty(
	id uuid.UUID,
	title string,
	createdBy uuid.UUID,
	createdAt time.Time,
	updatedBy uuid.UUID,
	updatedAt time.Time,
) (Property, error) {
	if id == uuid.Nil {
		return Property{}, ErrPropertyIDNil
	}
	title = strings.TrimSpace(title)
	if title == "" {
		return Property{}, ErrPropertyTitleBlank
	}
	if len(title) > PropertyTitleMaxLen {
		return Property{}, ErrPropertyTitleTooLong
	}
	if createdBy == uuid.Nil {
		return Property{}, ErrPropertyCreatedByNil
	}
	if createdAt.IsZero() {
		return Property{}, ErrPropertyCreatedAtZero
	}
	if updatedBy == uuid.Nil {
		return Property{}, ErrPropertyUpdatedByNil
	}
	if updatedAt.IsZero() {
		return Property{}, ErrPropertyUpdatedAtZero
	}
	if updatedAt.Before(createdAt) {
		return Property{}, ErrPropertyUpdatedAtBeforeCreatedAt
	}

	return Property{
		id:        id,
		title:     title,
		createdBy: createdBy,
		createdAt: createdAt,
		updatedBy: updatedBy,
		updatedAt: updatedAt,
	}, nil
}

// ID returns the unique identifier of the property.
func (p Property) ID() uuid.UUID {
	return p.id
}

// Title returns the property's display title.
func (p Property) Title() string {
	return p.title
}

// CreatedBy returns the user ID that created the property.
func (p Property) CreatedBy() uuid.UUID {
	return p.createdBy
}

// CreatedAt returns the creation timestamp of the property.
func (p Property) CreatedAt() time.Time {
	return p.createdAt
}

// UpdatedBy returns the user ID that last updated the property.
func (p Property) UpdatedBy() uuid.UUID {
	return p.updatedBy
}

// UpdatedAt returns the last update timestamp of the property.
func (p Property) UpdatedAt() time.Time {
	return p.updatedAt
}

// Rename updates the property title and audit metadata.
func (p *Property) Rename(title string, updatedBy uuid.UUID, updatedAt time.Time) error {
	title = strings.TrimSpace(title)
	if title == "" {
		return ErrPropertyTitleBlank
	}
	if len(title) > PropertyTitleMaxLen {
		return ErrPropertyTitleTooLong
	}
	if updatedBy == uuid.Nil {
		return ErrPropertyUpdatedByNil
	}
	if updatedAt.IsZero() {
		return ErrPropertyUpdatedAtZero
	}
	if updatedAt.Before(p.createdAt) {
		return ErrPropertyUpdatedAtBeforeCreatedAt
	}

	p.title = title
	p.updatedBy = updatedBy
	p.updatedAt = updatedAt
	return nil
}
