package redisutil

import (
	"context"
	"time"
)

type IRedisClient interface {
	SetJSON(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	GetJSON(ctx context.Context, key string, dest interface{}) error
	SetToken(ctx context.Context, token string, userID uint, expiration time.Duration) error
	GetToken(ctx context.Context, token string) (string, error)
	AcquireLock(ctx context.Context, key string, value string, expiration time.Duration) (bool, error)
	ReleaseLock(ctx context.Context, key string, value string) (bool, error)
	Delete(ctx context.Context, key string) error
}
