package auth

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type fakeVerifier struct {
	token verifiedToken
	err   error
}

func (v fakeVerifier) Verify(ctx context.Context, rawToken string) (verifiedToken, error) {
	return v.token, v.err
}

func TestMiddlewareAuthenticate(t *testing.T) {
	t.Parallel()

	middleware := &Middleware{
		verifier: fakeVerifier{
			token: verifiedToken{
				subject: "user-123",
				claims: func(v any) error {
					claims := v.(*keycloakClaims)
					claims.Email = "user@example.com"
					claims.RealmAccess.Roles = []string{"GLOBAL_ADMIN", "PROPERTY_MANAGER"}
					return nil
				},
			},
		},
	}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := UserFromContext(r.Context())
		if !ok {
			t.Fatalf("expected user in context")
		}
		if user.Subject != "user-123" {
			t.Fatalf("subject = %q, want %q", user.Subject, "user-123")
		}
		if user.Email != "user@example.com" {
			t.Fatalf("email = %q, want %q", user.Email, "user@example.com")
		}
		if len(user.Roles) != 2 {
			t.Fatalf("roles = %v, want 2 roles", user.Roles)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer token")
	resp := httptest.NewRecorder()

	middleware.Authenticate(next).ServeHTTP(resp, req)

	if resp.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", resp.Code, http.StatusNoContent)
	}
}

func TestMiddlewareAuthenticateRejectsMissingBearer(t *testing.T) {
	t.Parallel()

	middleware := &Middleware{
		verifier: fakeVerifier{
			err: errors.New("should not be called"),
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp := httptest.NewRecorder()

	middleware.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("next handler must not run")
	})).ServeHTTP(resp, req)

	if resp.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", resp.Code, http.StatusUnauthorized)
	}
}

func TestContextWithUser(t *testing.T) {
	t.Parallel()

	want := AuthenticatedUser{Subject: "user-123", Email: "user@example.com", Roles: []string{"ADMIN"}}
	got, ok := UserFromContext(ContextWithUser(context.Background(), want))
	if !ok {
		t.Fatalf("expected user in context")
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("user = %#v, want %#v", got, want)
	}
}
