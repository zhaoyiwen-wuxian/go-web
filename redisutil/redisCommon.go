package redisutil

import (
	"context"
	"go-web/appError"
	"go-web/models"
	"strconv"
	"time"
)

var wrapper IRedisClient

// InitRedisManager 在应用初始化时调用一次
func InitRedisManager(client IRedisClient) {
	wrapper = client
}

// SetToken 封装后的 token 缓存函数
func SetToken(ctx context.Context, token string, userID uint, expiration time.Duration) error {
	return wrapper.SetToken(ctx, token, userID, expiration)
}

func GetToken(ctx context.Context, token string) (string, error) {
	return wrapper.GetToken(ctx, token)
}

// AcquireLock 获取分布式锁
func AcquireLockRedis(ctx context.Context, key string, value string, expiration time.Duration) (bool, error) {
	return wrapper.AcquireLock(ctx, key, value, expiration)
}

func ReleaseLockRedis(ctx context.Context, key string, value string) (bool, error) {
	return wrapper.ReleaseLock(ctx, key, value)
}

// SetJSON 设置 JSON 数据
func SetJSONRedis(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return wrapper.SetJSON(ctx, key, value, expiration)
}

// GetJSON 获取 JSON 数据
func GetJSONRedis(ctx context.Context, key string, dest interface{}) error {
	return wrapper.GetJSON(ctx, key, dest)
}

func DeleteRedisKey(ctx context.Context, key string) error {
	return wrapper.Delete(ctx, key)
}

func CacheMessage(ctx context.Context, messageID uint, message models.Message, expiration time.Duration) error {
	// 将消息存储到 Redis，使用 messageID 作为键
	err := wrapper.SetJSON(ctx, "message:"+strconv.Itoa(int(messageID)), message, expiration)
	if err != nil {
		return appError.NewAppError(500, "消息缓存失败", err)
	}
	return nil
}

func GetCachedMessage(ctx context.Context, messageID uint) (models.Message, error) {
	var message models.Message
	err := wrapper.GetJSON(ctx, "message:"+strconv.Itoa(int(messageID)), &message)
	if err != nil {
		return message, appError.NewAppError(500, "从 Redis 获取消息失败", err)
	}
	return message, nil
}

// 用户ID 转字符串（可适配 int/uint 等）
func toStr(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case int:
		return strconv.Itoa(val)
	case uint:
		return strconv.Itoa(int(val))
	default:
		return ""
	}
}
