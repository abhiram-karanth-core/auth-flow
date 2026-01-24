package auth

import (
	redisclient "authflow/internal/redis"
	"time"

	"github.com/redis/go-redis/v9"
)

func RevokeToken(rdb *redis.Client, claims *Claims) error {
	ttl := time.Until(claims.ExpiresAt.Time)

	return rdb.Set(
		redisclient.Ctx,
		"revoked:"+claims.ID,
		"1",
		ttl,
	).Err()
}
