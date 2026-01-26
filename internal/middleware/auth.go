package middleware

import (
	"context"
	"net/http"
	"os"

	"authflow/internal/auth"
	"authflow/internal/redis"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

func JWTAuth(rdb *redis.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenStr := extractBearerToken(r)
			if tokenStr == "" {
				http.Error(w, "missing token", http.StatusUnauthorized)
				return
			}

			token, err := jwt.ParseWithClaims(
				tokenStr,
				&auth.Claims{},
				func(t *jwt.Token) (interface{}, error) {
					return []byte(os.Getenv("JWT_SECRET")), nil
				},
			)

			if err != nil || !token.Valid {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			claims := token.Claims.(*auth.Claims)

			// Redis revocation check
			exists, err := rdb.Exists(
				redisclient.Ctx,
				"revoked:"+claims.ID,
			).Result()

			if err != nil || exists == 1 {
				http.Error(w, "token revoked", http.StatusUnauthorized)
				return
			}

			// attach claims to context
			ctx := context.WithValue(r.Context(), "claims", claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
