package ws

import (
	"log"
	"sync"
	"time"
)

var (
	messageCount int64
	totalLatency time.Duration
	mu           sync.Mutex
)

// 每秒钟记录吞吐量
func LogThroughput() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		<-ticker.C
		mu.Lock()
		log.Printf("吞吐量: %d 消息/秒", messageCount)
		log.Printf("平均延迟: %v", totalLatency/time.Duration(messageCount))
		messageCount = 0
		totalLatency = 0
		mu.Unlock()
	}
}

// 处理消息时记录延迟
func (c *Client) handleMessage() {
	start := time.Now()
	// 处理消息...
	elapsed := time.Since(start)

	mu.Lock()
	messageCount++
	totalLatency += elapsed
	mu.Unlock()
}
