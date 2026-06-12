package domain

import (
	"errors"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestNewRole(t *testing.T) {
	tests := []struct {
		name        string
		id          uuid.UUID
		label       string
		description string
		wantErr     error
	}{
		{
			name:        "valid role",
			id:          uuid.New(),
			label:       "ADMIN",
			description: "Administrator role",
			wantErr:     nil,
		},
		{
			name:        "valid role without description",
			id:          uuid.New(),
			label:       "USER",
			description: "",
			wantErr:     nil,
		},
		{
			name:        "label blank",
			id:          uuid.New(),
			label:       "",
			description: "test",
			wantErr:     ErrRoleLabelBlank,
		},
		{
			name:        "label with whitespace only",
			id:          uuid.New(),
			label:       "   ",
			description: "test",
			wantErr:     ErrRoleLabelBlank,
		},
		{
			name:        "label not uppercase",
			id:          uuid.New(),
			label:       "Admin",
			description: "test",
			wantErr:     ErrRoleLabelNotUpper,
		},
		{
			name:        "label lowercase",
			id:          uuid.New(),
			label:       "admin",
			description: "test",
			wantErr:     ErrRoleLabelNotUpper,
		},
		{
			name:        "label exceeds max length",
			id:          uuid.New(),
			label:       "A" + strings.Repeat("A", RoleLabelMaxLen),
			description: "test",
			wantErr:     ErrRoleLabelTooLong,
		},
		{
			name:        "description exceeds max length",
			id:          uuid.New(),
			label:       "ADMIN",
			description: "d" + strings.Repeat("d", RoleDescriptionMaxLen),
			wantErr:     ErrRoleDescriptionTooLong,
		},
		{
			name:        "label trimmed whitespace",
			id:          uuid.New(),
			label:       "  ADMIN  ",
			description: "test",
			wantErr:     nil,
		},
		{
			name:        "description trimmed whitespace",
			id:          uuid.New(),
			label:       "ADMIN",
			description: "  description  ",
			wantErr:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			role, err := NewRole(tt.id, tt.label, tt.description)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("NewRole() error = %v, want %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("NewRole() unexpected error = %v", err)
				return
			}

			if role.ID() != tt.id {
				t.Errorf("Role.ID() = %v, want %v", role.ID(), tt.id)
			}
			// Label and description are trimmed
			if role.Label() != strings.TrimSpace(tt.label) {
				t.Errorf("Role.Label() = %q, want %q", role.Label(), strings.TrimSpace(tt.label))
			}
			if role.Description() != strings.TrimSpace(tt.description) {
				t.Errorf("Role.Description() = %q, want %q", role.Description(), strings.TrimSpace(tt.description))
			}
		})
	}
}

func TestRoleID(t *testing.T) {
	id := uuid.New()
	role, _ := NewRole(id, "ADMIN", "Admin role")

	if role.ID() != id {
		t.Errorf("Role.ID() = %v, want %v", role.ID(), id)
	}
}

func TestRoleLabel(t *testing.T) {
	role, _ := NewRole(uuid.New(), "EDITOR", "")

	if role.Label() != "EDITOR" {
		t.Errorf("Role.Label() = %q, want %q", role.Label(), "EDITOR")
	}
}

func TestRoleDescription(t *testing.T) {
	desc := "This is a description"
	role, _ := NewRole(uuid.New(), "VIEWER", desc)

	if role.Description() != desc {
		t.Errorf("Role.Description() = %q, want %q", role.Description(), desc)
	}
}

func TestRoleSetDescription(t *testing.T) {
	role, _ := NewRole(uuid.New(), "ADMIN", "Old description")
	newDesc := "New description"

	role.SetDescription(newDesc)

	if role.Description() != newDesc {
		t.Errorf("After SetDescription(), Description() = %q, want %q", role.Description(), newDesc)
	}
}

func TestRoleSetDescriptionEmpty(t *testing.T) {
	role, _ := NewRole(uuid.New(), "ADMIN", "Original")

	role.SetDescription("")

	if role.Description() != "" {
		t.Errorf("After SetDescription(\"\"), Description() = %q, want empty", role.Description())
	}
}

func TestRoleSetDescriptionWithWhitespace(t *testing.T) {
	role, _ := NewRole(uuid.New(), "ADMIN", "Original")

	role.SetDescription("  trimmed  ")

	// SetDescription doesn't trim - it uses raw input
	if role.Description() != "  trimmed  " {
		t.Errorf("After SetDescription(\"  trimmed  \"), Description() = %q, want \"  trimmed  \"", role.Description())
	}
}

func TestRoleCreatedAt(t *testing.T) {
	role, _ := NewRole(uuid.New(), "ADMIN", "Admin role")

	// CreatedAt should be zero time if not set during construction
	if !role.CreatedAt().IsZero() {
		t.Errorf("Role.CreatedAt() = %v, want zero time", role.CreatedAt())
	}
}

func TestRoleUpdatedAt(t *testing.T) {
	role, _ := NewRole(uuid.New(), "ADMIN", "Admin role")

	// UpdatedAt should be zero time if not set during construction
	if !role.UpdatedAt().IsZero() {
		t.Errorf("Role.UpdatedAt() = %v, want zero time", role.UpdatedAt())
	}
}
