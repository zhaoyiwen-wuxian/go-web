package ws

import (
	"encoding/json"
	"go-web/models"
	"go-web/redisutil"
	"go-web/utils"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

type Client struct {
	Conn   *websocket.Conn
	Send   chan []byte
	UserID string
	url    string
}

func (c *Client) writePump() {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	var totalMessagesSent int
	var startTime time.Time

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			err := c.Conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Println("write error:", err)
				return
			}
			totalMessagesSent++
			if totalMessagesSent == 1 {
				startTime = time.Now()
			}
		case <-ticker.C:
			// 每秒记录吞吐量和延迟
			duration := time.Since(startTime)
			log.Printf("发送消息数: %d, 消息发送延迟: %v\n", totalMessagesSent, duration)
			totalMessagesSent = 0
			startTime = time.Now()
		}
	}
}
func (c *Client) reconnect() {
	retryInterval := 1 * time.Second
	for {

		conn, _, err := websocket.DefaultDialer.Dial(c.url, nil)
		if err != nil {
			log.Println("WebSocket 连接失败，重试中...")

			time.Sleep(retryInterval)
			retryInterval = retryInterval * 2
			continue
		}

		c.Conn = conn
		log.Println("WebSocket 重新连接成功")
		break
	}
}

func (c *Client) readPump(hub *Hub, db *gorm.DB, con *gin.Context) {
	defer func() {
		hub.Unregister <- c
		c.reconnect()
		c.Conn.Close()
	}()
	c.Conn.SetReadLimit(512)
	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		// 读取消息
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}
		decryptedMessage, err := utils.DecryptMessage(message, []byte(redisutil.MessageKey))
		if err != nil {
			log.Println("消息解密错误:", err)
			continue
		}
		var msg models.Message
		if err := json.Unmarshal(decryptedMessage, &msg); err != nil {
			log.Println("消息解析错误:", err)
			continue
		}

		// 设置消息发送者和时间戳
		msg.From = c.UserID
		msg.Timestamp = time.Now().Unix()

		// 根据消息类型调用 service 包中的函数
		if switchMsgType(msg, db, c, con) {
			return
		}

		// 广播消息（如果需要的话）
		hub.Broadcast <- msg
	}
	c.handleMessage()
}

// 在 WebSocket 握手时进行身份验证
func (c *Client) authenticateRequest(r *http.Request, ctx *gin.Context) bool {
	// 获取请求头中的 Authorization 信息
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		log.Println("身份验证失败: 未提供认证信息")
		return false
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")
	valid := c.Authenticate(token, ctx)
	return valid
}
