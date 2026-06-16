// Package auth provides HTTP authentication middleware.
package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
)

type contextKey string

const userContextKey contextKey = "authenticatedUser"

// AuthenticatedUser is the authenticated principal exposed to handlers.
type AuthenticatedUser struct {
	Subject string   `json:"subject"`
	Email   string   `json:"email"`
	Roles   []string `json:"roles"`
}

type keycloakClaims struct {
	Email string `json:"email"`
	RealmAccess struct {
		Roles []string `json:"roles"`
	} `json:"realm_access"`
}

type tokenVerifier interface {
	Verify(ctx context.Context, rawToken string) (verifiedToken, error)
}

type verifiedToken struct {
	subject string
	claims  func(v any) error
}

func (t verifiedToken) Subject() string {
	return t.subject
}

func (t verifiedToken) Claims(v any) error {
	return t.claims(v)
}

type oidcVerifier struct {
	verifier *oidc.IDTokenVerifier
}

func (v oidcVerifier) Verify(ctx context.Context, rawToken string) (verifiedToken, error) {
	token, err := v.verifier.Verify(ctx, rawToken)
	if err != nil {
		return verifiedToken{}, err
	}

	return verifiedToken{
		subject: token.Subject,
		claims:  token.Claims,
	}, nil
}

// Middleware validates bearer tokens and stores the authenticated user in context.
type Middleware struct {
	verifier tokenVerifier
}

// NewMiddleware creates a middleware backed by the issuer's JWKS endpoint.
func NewMiddleware(ctx context.Context, issuerURL string, audience string) (*Middleware, error) {
	provider, err := oidc.NewProvider(ctx, issuerURL)
	if err != nil {
		return nil, err
	}

	return &Middleware{
		verifier: oidcVerifier{
			verifier: provider.Verifier(&oidc.Config{
				ClientID: audience,
			}),
		},
	}, nil
}

// Authenticate verifies the bearer token and injects the authenticated user.
func (m *Middleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawToken, err := bearerToken(r)
		if err != nil {
			http.Error(w, "missing or invalid bearer token", http.StatusUnauthorized)
			return
		}

		token, err := m.verifier.Verify(r.Context(), rawToken)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		var claims keycloakClaims
		if err := token.Claims(&claims); err != nil {
			http.Error(w, "invalid token claims", http.StatusUnauthorized)
			return
		}

		user := AuthenticatedUser{
			Subject: token.Subject(),
			Email:   claims.Email,
			Roles:   append([]string(nil), claims.RealmAccess.Roles...),
		}

		next.ServeHTTP(w, r.WithContext(ContextWithUser(r.Context(), user)))
	})
}

// ContextWithUser stores the authenticated user in the context.
func ContextWithUser(ctx context.Context, user AuthenticatedUser) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

// UserFromContext returns the authenticated user from context.
func UserFromContext(ctx context.Context) (AuthenticatedUser, bool) {
	user, ok := ctx.Value(userContextKey).(AuthenticatedUser)
	return user, ok
}

func bearerToken(r *http.Request) (string, error) {
	header := r.Header.Get("Authorization")
	if header == "" {
		return "", errors.New("missing authorization header")
	}

	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || parts[1] == "" {
		return "", errors.New("invalid authorization header")
	}

	return parts[1], nil
}
