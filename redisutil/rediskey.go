package redisutil

import "fmt"

// Token 缓存的 key 前缀
const TokenPrefix = "token:"

const HeartbeatKey = "user:%d:heartbeat"

const MessageKey = "hdasdhasdbasvydwyuwydha"

// 分布式锁 key 前缀
const LockPrefix = "lock:"

// 获取完整的 Token 缓存 key
func TokenKey(token string) string {
	return fmt.Sprintf("%s%s", TokenPrefix, token)
}

// 获取分布式锁的 key
func LockKey(name string) string {
	return fmt.Sprintf("%s%s", LockPrefix, name)
}
