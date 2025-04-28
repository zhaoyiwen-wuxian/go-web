package redisutil

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client

// RedisWrapper 实现了 IRedisClient 接口
type RedisWrapper struct {
	Client *redis.Client
}

// NewRedisWrapper 返回 RedisWrapper 的指针实例
func NewRedisWrapper(client *redis.Client) IRedisClient {
	return &RedisWrapper{Client: client}
}

// SetJSON 设置 JSON 数据
func (r *RedisWrapper) SetJSON(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.Client.Set(ctx, key, data, expiration).Err()
}

// GetJSON 获取 JSON 数据
func (r *RedisWrapper) GetJSON(ctx context.Context, key string, dest interface{}) error {
	data, err := r.Client.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(data), dest)
}

// SetToken 设置 Token 数据
func (r *RedisWrapper) SetToken(ctx context.Context, token string, userID uint, expiration time.Duration) error {
	return r.Client.Set(ctx, "token:"+token, userID, expiration).Err()
}

// GetToken 获取 Token 对应的值
func (r *RedisWrapper) GetToken(ctx context.Context, token string) (string, error) {
	return r.Client.Get(ctx, "token:"+token).Result()
}

// AcquireLock 获取分布式锁
func (r *RedisWrapper) AcquireLock(ctx context.Context, key string, value string, expiration time.Duration) (bool, error) {
	return r.Client.SetNX(ctx, "lock:"+key, value, expiration).Result()
}

// ReleaseLock 释放分布式锁
func (r *RedisWrapper) ReleaseLock(ctx context.Context, key string, value string) (bool, error) {
	lua := `
	if redis.call("get", KEYS[1]) == ARGV[1] then
		return redis.call("del", KEYS[1])
	else
		return 0
	end`
	result, err := r.Client.Eval(ctx, lua, []string{"lock:" + key}, value).Result()
	if err != nil {
		return false, err
	}
	return result.(int64) == 1, nil
}

func (r *RedisWrapper) Delete(ctx context.Context, key string) error {
	return r.Client.Del(ctx, key).Err()
}
