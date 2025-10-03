package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/V4T54L/goship/pkg/goship/utils"
	"github.com/golang-jwt/jwt/v5"
)

// TokenManager defines the interface for JWT operations, including token
// generation, parsing, and providing an authentication middleware.
type TokenManager interface {
	// Generate creates a new JWT token with the given claims.
	Generate(claims map[string]interface{}) (string, error)
	// Parse validates and parses a token string, returning the claims if valid.
	Parse(tokenString string) (map[string]interface{}, error)
	// AuthMiddleware returns an HTTP middleware to protect routes.
	AuthMiddleware() func(http.Handler) http.Handler
}

// jwtManagerHS is an implementation of TokenManager using the HS256 algorithm.
type jwtManagerHS struct {
	secret []byte
	ttl    time.Duration
}

// NewJwtManagerHS creates a new TokenManager using the HS256 signing method.
// It returns an error if the secret is empty.
func NewJwtManagerHS(secret string, ttl time.Duration) (TokenManager, error) {
	if secret == "" {
		return nil, fmt.Errorf("JWT secret cannot be empty")
	}
	return &jwtManagerHS{
		secret: []byte(secret),
		ttl:    ttl,
	}, nil
}

// Generate creates a new JWT string signed with the HS256 method.
// It automatically adds "iat" (issued at) and "exp" (expiration) claims.
func (m *jwtManagerHS) Generate(claims map[string]interface{}) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	tokenClaims := token.Claims.(jwt.MapClaims)

	for key, value := range claims {
		tokenClaims[key] = value
	}

	now := time.Now()
	tokenClaims["iat"] = now.Unix()
	tokenClaims["exp"] = now.Add(m.ttl).Unix()

	return token.SignedString(m.secret)
}

// Parse validates a token string's signature and expiration. If valid, it
// returns the claims contained within the token.
func (m *jwtManagerHS) Parse(tokenString string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.secret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// AuthMiddleware returns an HTTP middleware that enforces JWT authentication.
// It expects a Bearer token in the "Authorization" header. If the token is
// valid, it injects the claims into the request context.
func (m *jwtManagerHS) AuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				utils.WriteJSONError(w, http.StatusUnauthorized, "Authorization header is required")
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				utils.WriteJSONError(w, http.StatusUnauthorized, "Invalid Authorization header format")
				return
			}

			tokenString := parts[1]
			claims, err := m.Parse(tokenString)
			if err != nil {
				utils.WriteJSONError(w, http.StatusUnauthorized, "Invalid token")
				return
			}

			ctx := ContextWithClaims(r.Context(), claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
