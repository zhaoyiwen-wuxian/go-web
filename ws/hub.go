package ws

import (
	"go-web/models"
	"log"
	"sync"
)

type Hub struct {
	// 存储已注册的客户端
	Clients map[string]*Client
	// 注册客户端的通道
	Register chan *Client
	// 注销客户端的通道
	Unregister chan *Client
	// 广播消息的通道
	Broadcast chan models.Message
	// 保护 Clients 映射的互斥锁
	mu sync.Mutex
}

// NewHub 用于初始化并返回一个新的 Hub 实例
func NewHub() *Hub {
	return &Hub{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan models.Message),
		Clients:    make(map[string]*Client),
	}
}

// Run 启动 Hub，处理客户端连接和广播消息
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.Clients[client.UserID] = client
			h.mu.Unlock()
			log.Printf("客户端 %s 已连接", client.UserID)

		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.Clients[client.UserID]; ok {
				delete(h.Clients, client.UserID)
				close(client.Send)
				log.Printf("客户端 %s 已断开连接", client.UserID)
			}
			h.mu.Unlock()

		case message := <-h.Broadcast:
			h.mu.Lock()
			for _, client := range h.Clients {
				select {
				case client.Send <- []byte(message.Content): // 将 message.Content 转换为 []byte
				default:
					// 如果客户端无法接收消息，则关闭连接
					close(client.Send)
					delete(h.Clients, client.UserID)
				}
			}
			h.mu.Unlock()
		}
	}
}
