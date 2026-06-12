// Package router defines domain specific entities.
package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type (
	Role struct {
		id          uuid.UUID
		label       string
		description string
		createdAt   time.Time
		updateAt    time.Time
	}

	RoleRepository interface {
		Create(ctx context.Context, role Role) error
		FindAll(ctx context.Context) ([]Role, error)
		Find(ctx context.Context, id uuid.UUID) (Role, error)
	}
)

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
	return r.updateAt
}
