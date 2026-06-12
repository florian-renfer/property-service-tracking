// Package domain defines core business entities and contracts.
package domain

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	// ErrRoleLabelBlank indicates the role label is empty after trimming.
	ErrRoleLabelBlank = errors.New("role label cannot be blank")
	// ErrRoleLabelNotUpper indicates the role label is not uppercase.
	ErrRoleLabelNotUpper = errors.New("role label must be uppercase")
	// ErrRoleLabelTooLong indicates the role label exceeds max length.
	ErrRoleLabelTooLong = errors.New("role label exceeds 100 characters")
	// ErrRoleDescriptionTooLong indicates the role description exceeds max length.
	ErrRoleDescriptionTooLong = errors.New("role description exceeds 255 characters")
)

const (
	// RoleLabelMaxLen is max allowed length for role label.
	RoleLabelMaxLen = 100
	// RoleDescriptionMaxLen is max allowed length for role description.
	RoleDescriptionMaxLen = 255
)

type (
	// Role represents an authorization role in the system.
	Role struct {
		id          uuid.UUID
		label       string
		description string
		createdAt   time.Time
		updatedAt   time.Time
	}

	// RoleRepository defines persistence operations for roles.
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
