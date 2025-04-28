package utils

import (
	"context"
	"go-web/redisutil"
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client
var redisOnce sync.Once
var redisWrapper redisutil.IRedisClient

func InitRedis() {
	redisOnce.Do(func() {
		// 获取配置
		conf := Conf.Redis

		// 初始化 Redis 客户端
		RedisClient = redis.NewClient(&redis.Options{
			Addr:         conf.URL,
			Password:     conf.Password,
			DB:           conf.DB,
			PoolSize:     conf.Pool.MaxActive,
			MinIdleConns: conf.Pool.MaxIdle,
			IdleTimeout:  conf.Pool.IdleTimeout,
		})

		// 测试 Redis 连接
		if _, err := RedisClient.Ping(context.Background()).Result(); err != nil {
			log.Fatalf("❌ Redis连接失败: %v", err)
		}
		log.Println("✅ Redis 初始化完成")

		// 将 Redis 客户端包装成接口，方便后续使用
		redisWrapper = redisutil.NewRedisWrapper(RedisClient)
		redisutil.InitRedisManager(redisWrapper)
	})
}

// JSON 存储
func SetJSON(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return redisWrapper.SetJSON(ctx, key, value, expiration)
}

func GetJSON(ctx context.Context, key string, dest interface{}) error {
	return redisWrapper.GetJSON(ctx, key, dest)
}
