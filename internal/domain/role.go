// Package router defines domain specific entities.
package domain

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrRoleLabelBlank         = errors.New("role label cannot be blank")
	ErrRoleLabelNotUpper      = errors.New("role label must be uppercase")
	ErrRoleLabelTooLong       = errors.New("role label exceeds 100 characters")
	ErrRoleDescriptionTooLong = errors.New("role description exceeds 255 characters")
)

const (
	RoleLabelMaxLen       = 100
	RoleDescriptionMaxLen = 255
)

type (
	Role struct {
		id          uuid.UUID
		label       string
		description string
		createdAt   time.Time
		updatedAt   time.Time
	}

	RoleRepository interface {
		Create(ctx context.Context, role Role) error
		FindAll(ctx context.Context) ([]Role, error)
		Find(ctx context.Context, id uuid.UUID) (Role, error)
	}
)

// NewRole constructs a new Role with validation.
func NewRole(id uuid.UUID, label, description string) (Role, error) {
	label = strings.TrimSpace(label)
	if label == "" {
		return Role{}, ErrRoleLabelBlank
	}
	if label != strings.ToUpper(label) {
		return Role{}, ErrRoleLabelNotUpper
	}
	if len(label) > RoleLabelMaxLen {
		return Role{}, ErrRoleLabelTooLong
	}
	description = strings.TrimSpace(description)
	if len(description) > RoleDescriptionMaxLen {
		return Role{}, ErrRoleDescriptionTooLong
	}

	return Role{
		id:          id,
		label:       label,
		description: description,
	}, nil
}

// ID returns the unique identifier of the role.
func (r Role) ID() uuid.UUID {
	return r.id
}

// Label returns the role's label.
func (r Role) Label() string {
	return r.label
}

// Description returns the role's description.
func (r Role) Description() string {
	return r.description
}

// CreatedAt returns the creation timestamp of the role.
func (r Role) CreatedAt() time.Time {
	return r.createdAt
}

// UpdatedAt returns the last update timestamp of the role.
func (r Role) UpdatedAt() time.Time {
	return r.updatedAt
}

// SetDescription updates the role's description.
func (r *Role) SetDescription(description string) {
	r.description = description
}
