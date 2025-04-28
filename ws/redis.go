package ws

import (
	"context"
	"encoding/json"
	"go-web/models"
	"go-web/redisutil"
	"log"
)

// 发布消息到 Redis
func PublishMessage(msg models.Message) {
	// 将消息序列化为 JSON 格式
	data, err := json.Marshal(msg)
	if err != nil {
		log.Println("Error marshaling message:", err)
		return
	}

	// 使用 redisutil 包的 SetJSON 方法存储到 Redis
	err = redisutil.SetJSONRedis(context.Background(), "ws_channel", data, 0) // 0 表示永久有效
	if err != nil {
		log.Println("Error setting JSON in Redis:", err)
	}
}

// 订阅 Redis 消息并广播
func Subscribe(hub *Hub) {
	// 获取订阅的消息
	go func() {
		for {
			// 从 Redis 中读取消息
			var msg []byte
			err := redisutil.GetJSONRedis(context.Background(), "ws_channel", &msg)
			if err != nil {
				log.Println("Error getting JSON from Redis:", err)
				continue
			}

			// 将获取到的消息解码为 Message 对象
			var m models.Message
			if err := json.Unmarshal(msg, &m); err != nil {
				log.Println("Error unmarshaling message:", err)
				continue
			}

			// 广播消息到 WebSocket 客户端
			hub.Broadcast <- m
		}
	}()
}
