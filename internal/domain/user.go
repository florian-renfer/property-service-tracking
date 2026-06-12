// Package router defines domain specific entities.
package domain

import (
	"context"
	"errors"
	"net/mail"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrUserPasswordHashBlank  = errors.New("user password hash cannot be blank")
	ErrUserPasswordHashTooLong = errors.New("user password hash exceeds 64 characters")
	ErrUserGivenNameBlank     = errors.New("user given name cannot be blank")
	ErrUserFamilyNameBlank    = errors.New("user family name cannot be blank")
	ErrUserEmailBlank         = errors.New("user email cannot be blank")
	ErrUserEmailInvalid       = errors.New("user email format is invalid")
	ErrUserEmailNotLower      = errors.New("user email must be lowercase")
	ErrUserEmailTooLong       = errors.New("user email exceeds 255 characters")
	ErrUserGivenNameTooLong   = errors.New("user given name exceeds 255 characters")
	ErrUserFamilyNameTooLong  = errors.New("user family name exceeds 255 characters")
	ErrUserRoleIDNil          = errors.New("user role id cannot be nil")
)

const (
	UserPasswordHashMaxLen = 64
	UserEmailMaxLen        = 255
	UserGivenNameMaxLen    = 255
	UserFamilyNameMaxLen   = 255
)

type (
	User struct {
		id           uuid.UUID
		email        string
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

// NewUser constructs a new User with validation.
func NewUser(id uuid.UUID, email, passwordHash, givenName, familyName string, roleID uuid.UUID) (User, error) {
	email = strings.TrimSpace(email)
	if email == "" {
		return User{}, ErrUserEmailBlank
	}
	if email != strings.ToLower(email) {
		return User{}, ErrUserEmailNotLower
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return User{}, ErrUserEmailInvalid
	}
	if len(email) > UserEmailMaxLen {
		return User{}, ErrUserEmailTooLong
	}
	passwordHash = strings.TrimSpace(passwordHash)
	if passwordHash == "" {
		return User{}, ErrUserPasswordHashBlank
	}
	if len(passwordHash) > UserPasswordHashMaxLen {
		return User{}, ErrUserPasswordHashTooLong
	}
	givenName = strings.TrimSpace(givenName)
	if givenName == "" {
		return User{}, ErrUserGivenNameBlank
	}
	if len(givenName) > UserGivenNameMaxLen {
		return User{}, ErrUserGivenNameTooLong
	}
	familyName = strings.TrimSpace(familyName)
	if familyName == "" {
		return User{}, ErrUserFamilyNameBlank
	}
	if len(familyName) > UserFamilyNameMaxLen {
		return User{}, ErrUserFamilyNameTooLong
	}
	if roleID == uuid.Nil {
		return User{}, ErrUserRoleIDNil
	}

	return User{
		id:           id,
		email:        email,
		passwordHash: passwordHash,
		givenName:    givenName,
		familyName:   familyName,
		roleID:       roleID,
		active:       true,
	}, nil
}

// ID returns the unique identifier of the user.
func (u User) ID() uuid.UUID {
	return u.id
}

// Email returns the user's email address.
func (u User) Email() string {
	return u.email
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

// Activate marks the user account as active.
func (u *User) Activate() {
	u.active = true
}

// Deactivate marks the user account as inactive.
func (u *User) Deactivate() {
	u.active = false
}

// SetRoleID updates the user's role.
func (u *User) SetRoleID(roleID uuid.UUID) error {
	if roleID == uuid.Nil {
		return ErrUserRoleIDNil
	}
	u.roleID = roleID
	return nil
}

// RecordLogin updates the last login timestamp.
func (u *User) RecordLogin(t time.Time) {
	u.lastLoginAt = t
}
