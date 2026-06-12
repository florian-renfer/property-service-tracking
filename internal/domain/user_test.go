package domain

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewUser(t *testing.T) {
	validID := uuid.New()
	validRoleID := uuid.New()

	tests := []struct {
		name       string
		id         uuid.UUID
		email      string
		passHash   string
		givenName  string
		familyName string
		roleID     uuid.UUID
		wantErr    error
	}{
		{
			name:       "valid user",
			id:         validID,
			email:      "john.doe@example.com",
			passHash:   "$2a$12$xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
			givenName:  "John",
			familyName: "Doe",
			roleID:     validRoleID,
			wantErr:    nil,
		},
		{
			name:       "email blank",
			id:         validID,
			email:      "",
			passHash:   "$2a$12$xxxx",
			givenName:  "John",
			familyName: "Doe",
			roleID:     validRoleID,
			wantErr:    ErrUserEmailBlank,
		},
		{
			name:       "email whitespace only",
			id:         validID,
			email:      "   ",
			passHash:   "$2a$12$xxxx",
			givenName:  "John",
			familyName: "Doe",
			roleID:     validRoleID,
			wantErr:    ErrUserEmailBlank,
		},
		{
			name:       "email not lowercase",
			id:         validID,
			email:      "John.Doe@Example.com",
			passHash:   "$2a$12$xxxx",
			givenName:  "John",
			familyName: "Doe",
			roleID:     validRoleID,
			wantErr:    ErrUserEmailNotLower,
		},
		{
			name:       "email invalid format",
			id:         validID,
			email:      "invalid-email",
			passHash:   "$2a$12$xxxx",
			givenName:  "John",
			familyName: "Doe",
			roleID:     validRoleID,
			wantErr:    ErrUserEmailInvalid,
		},
		{
			name:       "email too long",
			id:         validID,
			email:      "a@" + strings.Repeat("a", UserEmailMaxLen),
			passHash:   "$2a$12$xxxx",
			givenName:  "John",
			familyName: "Doe",
			roleID:     validRoleID,
			wantErr:    ErrUserEmailTooLong,
		},
		{
			name:       "password hash blank",
			id:         validID,
			email:      "john@example.com",
			passHash:   "",
			givenName:  "John",
			familyName: "Doe",
			roleID:     validRoleID,
			wantErr:    ErrUserPasswordHashBlank,
		},
		{
			name:       "password hash too long",
			id:         validID,
			email:      "john@example.com",
			passHash:   "x" + strings.Repeat("x", UserPasswordHashMaxLen),
			givenName:  "John",
			familyName: "Doe",
			roleID:     validRoleID,
			wantErr:    ErrUserPasswordHashTooLong,
		},
		{
			name:       "given name blank",
			id:         validID,
			email:      "john@example.com",
			passHash:   "$2a$12$xxxx",
			givenName:  "",
			familyName: "Doe",
			roleID:     validRoleID,
			wantErr:    ErrUserGivenNameBlank,
		},
		{
			name:       "given name too long",
			id:         validID,
			email:      "john@example.com",
			passHash:   "$2a$12$xxxx",
			givenName:  "g" + strings.Repeat("g", UserGivenNameMaxLen),
			familyName: "Doe",
			roleID:     validRoleID,
			wantErr:    ErrUserGivenNameTooLong,
		},
		{
			name:       "family name blank",
			id:         validID,
			email:      "john@example.com",
			passHash:   "$2a$12$xxxx",
			givenName:  "John",
			familyName: "",
			roleID:     validRoleID,
			wantErr:    ErrUserFamilyNameBlank,
		},
		{
			name:       "family name too long",
			id:         validID,
			email:      "john@example.com",
			passHash:   "$2a$12$xxxx",
			givenName:  "John",
			familyName: "f" + strings.Repeat("f", UserFamilyNameMaxLen),
			roleID:     validRoleID,
			wantErr:    ErrUserFamilyNameTooLong,
		},
		{
			name:       "role id nil",
			id:         validID,
			email:      "john@example.com",
			passHash:   "$2a$12$xxxx",
			givenName:  "John",
			familyName: "Doe",
			roleID:     uuid.Nil,
			wantErr:    ErrUserRoleIDNil,
		},
		{
			name:       "email trimmed",
			id:         validID,
			email:      "  john@example.com  ",
			passHash:   "$2a$12$xxxx",
			givenName:  "John",
			familyName: "Doe",
			roleID:     validRoleID,
			wantErr:    nil,
		},
		{
			name:       "given name trimmed",
			id:         validID,
			email:      "john@example.com",
			passHash:   "$2a$12$xxxx",
			givenName:  "  John  ",
			familyName: "Doe",
			roleID:     validRoleID,
			wantErr:    nil,
		},
		{
			name:       "family name trimmed",
			id:         validID,
			email:      "john@example.com",
			passHash:   "$2a$12$xxxx",
			givenName:  "John",
			familyName: "  Doe  ",
			roleID:     validRoleID,
			wantErr:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := NewUser(tt.id, tt.email, tt.passHash, tt.givenName, tt.familyName, tt.roleID)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("NewUser() error = %v, want %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("NewUser() unexpected error = %v", err)
				return
			}

			if user.ID() != tt.id {
				t.Errorf("User.ID() = %v, want %v", user.ID(), tt.id)
			}
			if user.Active() != true {
				t.Errorf("User.Active() = %v, want true", user.Active())
			}
		})
	}
}

func TestUserID(t *testing.T) {
	id := uuid.New()
	roleID := uuid.New()
	user, _ := NewUser(id, "john@example.com", "$2a$12$xxxx", "John", "Doe", roleID)

	if user.ID() != id {
		t.Errorf("User.ID() = %v, want %v", user.ID(), id)
	}
}

func TestUserEmail(t *testing.T) {
	email := "john@example.com"
	roleID := uuid.New()
	user, _ := NewUser(uuid.New(), email, "$2a$12$xxxx", "John", "Doe", roleID)

	if user.Email() != email {
		t.Errorf("User.Email() = %q, want %q", user.Email(), email)
	}
}

func TestUserPasswordHash(t *testing.T) {
	hash := "$2a$12$xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
	roleID := uuid.New()
	user, _ := NewUser(uuid.New(), "john@example.com", hash, "John", "Doe", roleID)

	if user.PasswordHash() != hash {
		t.Errorf("User.PasswordHash() = %q, want %q", user.PasswordHash(), hash)
	}
}

func TestUserGivenName(t *testing.T) {
	given := "John"
	roleID := uuid.New()
	user, _ := NewUser(uuid.New(), "john@example.com", "$2a$12$xxxx", given, "Doe", roleID)

	if user.GivenName() != given {
		t.Errorf("User.GivenName() = %q, want %q", user.GivenName(), given)
	}
}

func TestUserFamilyName(t *testing.T) {
	family := "Doe"
	roleID := uuid.New()
	user, _ := NewUser(uuid.New(), "john@example.com", "$2a$12$xxxx", "John", family, roleID)

	if user.FamilyName() != family {
		t.Errorf("User.FamilyName() = %q, want %q", user.FamilyName(), family)
	}
}

func TestUserRoleID(t *testing.T) {
	roleID := uuid.New()
	user, _ := NewUser(uuid.New(), "john@example.com", "$2a$12$xxxx", "John", "Doe", roleID)

	if user.RoleID() != roleID {
		t.Errorf("User.RoleID() = %v, want %v", user.RoleID(), roleID)
	}
}

func TestUserActive(t *testing.T) {
	roleID := uuid.New()
	user, _ := NewUser(uuid.New(), "john@example.com", "$2a$12$xxxx", "John", "Doe", roleID)

	if user.Active() != true {
		t.Errorf("User.Active() = %v, want true", user.Active())
	}
}

func TestUserCreatedAt(t *testing.T) {
	roleID := uuid.New()
	user, _ := NewUser(uuid.New(), "john@example.com", "$2a$12$xxxx", "John", "Doe", roleID)

	// CreatedAt should be zero time if not set during construction
	if !user.CreatedAt().IsZero() && user.CreatedAt().After(time.Now().Add(1*time.Second)) {
		t.Errorf("User.CreatedAt() = %v, want zero or current time", user.CreatedAt())
	}
}

func TestUserLastLoginAt(t *testing.T) {
	roleID := uuid.New()
	user, _ := NewUser(uuid.New(), "john@example.com", "$2a$12$xxxx", "John", "Doe", roleID)

	// LastLoginAt should be zero time initially
	if !user.LastLoginAt().IsZero() {
		t.Errorf("User.LastLoginAt() = %v, want zero time", user.LastLoginAt())
	}
}

func TestUserActivate(t *testing.T) {
	roleID := uuid.New()
	user, _ := NewUser(uuid.New(), "john@example.com", "$2a$12$xxxx", "John", "Doe", roleID)

	user.Deactivate()
	if user.Active() {
		t.Errorf("After Deactivate(), Active() = true, want false")
	}

	user.Activate()
	if !user.Active() {
		t.Errorf("After Activate(), Active() = false, want true")
	}
}

func TestUserDeactivate(t *testing.T) {
	roleID := uuid.New()
	user, _ := NewUser(uuid.New(), "john@example.com", "$2a$12$xxxx", "John", "Doe", roleID)

	if !user.Active() {
		t.Errorf("NewUser() creates inactive user, want active")
	}

	user.Deactivate()
	if user.Active() {
		t.Errorf("After Deactivate(), Active() = true, want false")
	}
}

func TestUserSetRoleID(t *testing.T) {
	oldRoleID := uuid.New()
	newRoleID := uuid.New()
	user, _ := NewUser(uuid.New(), "john@example.com", "$2a$12$xxxx", "John", "Doe", oldRoleID)

	if user.RoleID() != oldRoleID {
		t.Errorf("Initial RoleID() = %v, want %v", user.RoleID(), oldRoleID)
	}

	err := user.SetRoleID(newRoleID)
	if err != nil {
		t.Errorf("SetRoleID() unexpected error = %v", err)
	}

	if user.RoleID() != newRoleID {
		t.Errorf("After SetRoleID(), RoleID() = %v, want %v", user.RoleID(), newRoleID)
	}
}

func TestUserSetRoleIDNil(t *testing.T) {
	roleID := uuid.New()
	user, _ := NewUser(uuid.New(), "john@example.com", "$2a$12$xxxx", "John", "Doe", roleID)

	err := user.SetRoleID(uuid.Nil)
	if !errors.Is(err, ErrUserRoleIDNil) {
		t.Errorf("SetRoleID(uuid.Nil) error = %v, want %v", err, ErrUserRoleIDNil)
	}

	if user.RoleID() != roleID {
		t.Errorf("After failed SetRoleID(), RoleID() changed to %v, want %v", user.RoleID(), roleID)
	}
}

func TestUserRecordLogin(t *testing.T) {
	roleID := uuid.New()
	user, _ := NewUser(uuid.New(), "john@example.com", "$2a$12$xxxx", "John", "Doe", roleID)

	loginTime := time.Date(2026, 6, 12, 10, 30, 0, 0, time.UTC)
	user.RecordLogin(loginTime)

	if user.LastLoginAt() != loginTime {
		t.Errorf("After RecordLogin(), LastLoginAt() = %v, want %v", user.LastLoginAt(), loginTime)
	}
}

func TestUserRecordLoginMultiple(t *testing.T) {
	roleID := uuid.New()
	user, _ := NewUser(uuid.New(), "john@example.com", "$2a$12$xxxx", "John", "Doe", roleID)

	time1 := time.Date(2026, 6, 12, 10, 0, 0, 0, time.UTC)
	time2 := time.Date(2026, 6, 12, 15, 0, 0, 0, time.UTC)

	user.RecordLogin(time1)
	if user.LastLoginAt() != time1 {
		t.Errorf("First RecordLogin(), LastLoginAt() = %v, want %v", user.LastLoginAt(), time1)
	}

	user.RecordLogin(time2)
	if user.LastLoginAt() != time2 {
		t.Errorf("Second RecordLogin(), LastLoginAt() = %v, want %v", user.LastLoginAt(), time2)
	}
}
