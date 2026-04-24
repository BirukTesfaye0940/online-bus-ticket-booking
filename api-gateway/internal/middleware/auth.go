package middleware

import (
	"context"
	"net/http"
	"strings"

	pb "github.com/biruk/bus-ticket/api-gateway/internal/proto"
)

type contextKey string

const (
	UserIDKey contextKey = "user_id"
	RoleKey   contextKey = "role"
)

// Auth validates the Bearer token in the Authorization header by calling the Auth
// Service's ValidateToken RPC. On success, user_id and role are injected into context.
func Auth(authClient pb.AuthServiceClient) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := extractBearerToken(r)
			if token == "" {
				http.Error(w, `{"error":"missing authorization header"}`, http.StatusUnauthorized)
				return
			}

			resp, err := authClient.ValidateToken(r.Context(), &pb.ValidateTokenRequest{
				AccessToken: token,
			})
			if err != nil || !resp.Valid {
				http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, resp.UserId)
			ctx = context.WithValue(ctx, RoleKey, resp.Role.String())
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func extractBearerToken(r *http.Request) string {
	header := r.Header.Get("Authorization")
	if header == "" || !strings.HasPrefix(header, "Bearer ") {
		return ""
	}
	return strings.TrimPrefix(header, "Bearer ")
}
