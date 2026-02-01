// Package middleware provides HTTP middleware for auth and request handling.
package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/urinvitedto-my/backend/internal/config"
	"github.com/urinvitedto-my/backend/internal/handlers"
	"github.com/urinvitedto-my/backend/internal/utils"
)

// auth errors for cleaner error handling
var (
	errMissingAuth   = &authError{"unauthorized", "Missing authorization header"}
	errInvalidFormat = &authError{"unauthorized", "Invalid authorization format"}
	errInvalidToken  = &authError{"unauthorized", "Invalid token"}
	errInvalidClaims = &authError{"unauthorized", "Invalid token claims"}
	errNoEmail       = &authError{"unauthorized", "No email in token"}
	errForbidden     = &authError{"forbidden", "Admin access required"}
)

type authError struct {
	code string
	msg  string
}

// Error implements the error interface.
func (e *authError) Error() string {
	return e.msg
}

// Middleware manages auth and request middleware.
type Middleware struct {
	cfg      *config.Config
	handlers *handlers.Handlers
}

// New creates a new Middleware instance.
func New(cfg *config.Config, h *handlers.Handlers) *Middleware {
	return &Middleware{
		cfg:      cfg,
		handlers: h,
	}
}

// validateJWT parses and validates a JWT token, returning claims if valid.
func (m *Middleware) validateJWT(r *http.Request) (jwt.MapClaims, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, errMissingAuth
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		return nil, errInvalidFormat
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(m.cfg.SupabaseJWTSecret), nil
	})
	if err != nil || !token.Valid {
		return nil, errInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errInvalidClaims
	}

	return claims, nil
}

// extractUserEmail validates JWT and extracts email from claims.
// Returns email and true on success, or writes error response and returns false.
func (m *Middleware) extractUserEmail(
	w http.ResponseWriter,
	r *http.Request,
) (string, bool) {
	claims, err := m.validateJWT(r)
	if err != nil {
		if authErr, ok := err.(*authError); ok {
			utils.WriteError(w, http.StatusUnauthorized, authErr.code, authErr.msg)
		} else {
			utils.WriteError(
				w,
				http.StatusUnauthorized,
				"unauthorized",
				"Authentication failed",
			)
		}
		return "", false
	}

	email, ok := claims["email"].(string)
	if !ok || email == "" {
		utils.WriteError(w, http.StatusUnauthorized, errNoEmail.code, errNoEmail.msg)
		return "", false
	}

	return email, true
}

// Auth verifies JWT and extracts user email into context.
func (m *Middleware) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		email, ok := m.extractUserEmail(w, r)
		if !ok {
			return
		}

		ctx := context.WithValue(r.Context(), handlers.UserEmailKey, email)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireAdmin checks if the authenticated user is an admin.
// Must be used after Auth middleware.
func (m *Middleware) RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		email, ok := r.Context().Value(handlers.UserEmailKey).(string)
		if !ok || email == "" {
			utils.WriteError(
				w,
				http.StatusUnauthorized,
				"unauthorized",
				"Not authenticated",
			)
			return
		}

		isAdmin, err := m.handlers.IsAdmin(r.Context(), email)
		if err != nil || !isAdmin {
			slog.Warn("Non-admin access attempt", "email", email)
			utils.WriteError(
				w,
				http.StatusForbidden,
				errForbidden.code,
				errForbidden.msg,
			)
			return
		}

		next.ServeHTTP(w, r)
	})
}
