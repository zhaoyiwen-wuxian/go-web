package jwtutil

import (
	"context"
	"go-web/redisutil"
	"time"
)

// 更新 Token 的过期时间
func UpdateTokenExpiration(ctx context.Context, token string, userID uint, expiration time.Duration) error {
	// 刷新 Token 的过期时间
	err := redisutil.SetToken(ctx, redisutil.TokenPrefix+token, userID, expiration)
	if err != nil {
		return err
	}

	return nil
}
