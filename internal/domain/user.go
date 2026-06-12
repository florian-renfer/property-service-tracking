// Package router defines domain specific entities.
package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type (
	User struct {
		id           uuid.UUID
		passwordHash string
		givenName    string
		familyName   string
		roleID       uuid.UUID
		active       bool
		createdAt    time.Time
		lastLoginAt  time.Time
	}

	UserRepository interface {
		Create(ctx context.Context, user User) error
		FindAll(ctx context.Context) ([]User, error)
		Find(ctx context.Context, id uuid.UUID) (User, error)
	}
)

// ID returns the unique identifier of the user.
func (u User) ID() uuid.UUID {
	return u.id
}

// PasswordHash returns the user's hashed password.
func (u User) PasswordHash() string {
	return u.passwordHash
}

// GivenName returns the user's given name.
func (u User) GivenName() string {
	return u.givenName
}

// FamilyName returns the user's family name.
func (u User) FamilyName() string {
	return u.familyName
}

// RoleID returns the user's role identifier.
func (u User) RoleID() uuid.UUID {
	return u.roleID
}

// Active returns whether the user account is active.
func (u User) Active() bool {
	return u.active
}

// CreatedAt returns the creation timestamp of the user.
func (u User) CreatedAt() time.Time {
	return u.createdAt
}

// LastLoginAt returns the user's last login timestamp.
func (u User) LastLoginAt() time.Time {
	return u.lastLoginAt
}
