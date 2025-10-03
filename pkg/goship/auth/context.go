package auth

import "context"

// contextKey is a private type used for context keys to avoid collisions.
type contextKey string

const (
	ClaimsKey contextKey = "claims"
)

// ContextWithClaims returns a new context with the provided JWT claims.
func ContextWithClaims(ctx context.Context, claims map[string]interface{}) context.Context {
	return context.WithValue(ctx, ClaimsKey, claims)
}

