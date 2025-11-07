package middleware

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"

	"warehouse-control/internal/domain"
	"warehouse-control/internal/interfaces"
)

func MiddlewareLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] %s %s", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

type contextKey string

const userCtxKey contextKey = "user"

func JWTAuth(tp interfaces.TokenProvide, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			http.Error(w, "missing or invalid token", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(auth, "Bearer ")
		token, err := tp.Parse(tokenStr)
		if err != nil || !token.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "invalid claims", http.StatusUnauthorized)
			return
		}

		uidFloat, ok := claims["user_id"].(float64)
		if !ok {
			http.Error(w, "not claim user_id", http.StatusUnauthorized)
			return
		}

		username, ok := claims["username"].(string)
		if !ok {
			http.Error(w, "not claim username", http.StatusUnauthorized)
			return
		}

		role, ok := claims["role"].(string)
		if !ok {
			http.Error(w, "not claim role", http.StatusUnauthorized)
			return
		}

		user := &domain.User{
			ID:       int64(uidFloat),
			Username: username,
			Role:     domain.Role(role),
		}

		ctx := context.WithValue(r.Context(), userCtxKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UserFromContext(r *http.Request) *domain.User {
	user, _ := r.Context().Value(userCtxKey).(*domain.User)
	return user
}

func RequireAnyRole(roles []domain.Role, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := UserFromContext(r)
		if user == nil {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		for _, role := range roles {
			if user.Role == role {
				next.ServeHTTP(w, r)
				return
			}
		}

		http.Error(w, "forbidden", http.StatusForbidden)
	})
}

func WithDBUser(db *sql.DB) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := r.Context().Value(userCtxKey)

			if u, ok := user.(*domain.User); ok {
				_, _ = db.Exec("SELECT set_config('jwt.user', $1, false)", u.Username)
			}

			next.ServeHTTP(w, r)
		})
	}
}
