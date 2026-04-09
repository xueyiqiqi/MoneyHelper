package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

func InitRedis(addr string) {
	RDB = redis.NewClient(&redis.Options{
		Addr: addr,
	})
}

func BlacklistToken(ctx context.Context, token string, expiration time.Duration) error {
	return RDB.Set(ctx, "blacklist:"+token, "revoked", expiration).Err()
}

func IsTokenBlacklisted(ctx context.Context, token string) (bool, error) {
	val, err := RDB.Get(ctx, "blacklist:"+token).Result()
	if errors.Is(err, redis.Nil) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return val == "revoked", nil
}

// InvalidateUserTokens sets a timestamp to invalidate all tokens issued before now
func InvalidateUserTokens(ctx context.Context, userID uint) error {
	return RDB.Set(ctx, fmt.Sprintf("invalidated:user:%d", userID), time.Now().Unix(), 15*time.Minute).Err()
}

// GetUserInvalidationTime returns the timestamp when user tokens were invalidated
func GetUserInvalidationTime(ctx context.Context, userID uint) (int64, error) {
	val, err := RDB.Get(ctx, fmt.Sprintf("invalidated:user:%d", userID)).Int64()
	if errors.Is(err, redis.Nil) {
		return 0, nil
	}
	return val, err
}
