package utils

import (
	"time"
)

// RedisConfig 存储 Redis 配置信息
type RedisConfig struct {
	URL      string `yaml:"url"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
	Pool     struct {
		MaxIdle     int           `yaml:"max_idle"`
		MaxActive   int           `yaml:"max_active"`
		IdleTimeout time.Duration `yaml:"idle_timeout"`
	} `yaml:"pool"`
}
