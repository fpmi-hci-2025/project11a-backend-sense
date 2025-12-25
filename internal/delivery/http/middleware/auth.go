package middleware

import (
	"context"
	"net/http"
	"strings"

	"sense-backend/internal/infrastructure/jwt"
)

type contextKey string

const userIDKey contextKey = "user_id"
const usernameKey contextKey = "username"
const roleKey contextKey = "role"

// AuthMiddleware validates JWT token (required)
func AuthMiddleware(tokenSvc *jwt.TokenService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, `{"error":"unauthorized","message":"Требуется аутентификация"}`, http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, `{"error":"unauthorized","message":"Неверный формат токена"}`, http.StatusUnauthorized)
				return
			}

			claims, err := tokenSvc.ValidateToken(parts[1])
			if err != nil {
				http.Error(w, `{"error":"unauthorized","message":"Недействительный токен"}`, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
			ctx = context.WithValue(ctx, usernameKey, claims.Username)
			ctx = context.WithValue(ctx, roleKey, claims.Role)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// OptionalAuthMiddleware extracts user info from JWT token if present, but does not require authentication.
// Use this for endpoints that work for both authenticated and unauthenticated users,
// but need to know who the user is (e.g., to check is_liked status).
func OptionalAuthMiddleware(tokenSvc *jwt.TokenService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				// No token provided - proceed without user context
				next.ServeHTTP(w, r)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				// Invalid format - proceed without user context
				next.ServeHTTP(w, r)
				return
			}

			claims, err := tokenSvc.ValidateToken(parts[1])
			if err != nil {
				// Invalid token - proceed without user context
				next.ServeHTTP(w, r)
				return
			}

			// Valid token - set user context
			ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
			ctx = context.WithValue(ctx, usernameKey, claims.Username)
			ctx = context.WithValue(ctx, roleKey, claims.Role)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserID retrieves user ID from context
func GetUserID(ctx context.Context) string {
	if id, ok := ctx.Value(userIDKey).(string); ok {
		return id
	}
	return ""
}

// GetUsername retrieves username from context
func GetUsername(ctx context.Context) string {
	if username, ok := ctx.Value(usernameKey).(string); ok {
		return username
	}
	return ""
}

// GetRole retrieves role from context
func GetRole(ctx context.Context) string {
	if role, ok := ctx.Value(roleKey).(string); ok {
		return role
	}
	return ""
}

