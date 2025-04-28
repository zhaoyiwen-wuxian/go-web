package ws

import (
	"encoding/json"
	"go-web/models"
	"go-web/redisutil"
	"go-web/utils"
	"log"

	"github.com/gorilla/websocket"
)

func (c *Client) SendError(message string) {
	errMsg := models.Message{
		Type:    "error",
		Content: message,
	}
	c.SendMessage(errMsg)
}

func (c *Client) SendMessage(message models.Message) error {
	response, err := json.Marshal(message)
	if err != nil {
		log.Println("消息序列化错误:", err)
		return err
	}
	encryptedMessage, err := utils.EncryptMessage(response, []byte(redisutil.MessageKey))
	if err != nil {
		log.Println("消息加密错误:", err)
		return err
	}
	if err := c.Conn.WriteMessage(websocket.TextMessage, encryptedMessage); err != nil {
		log.Println("发送消息错误:", err)
		return err
	}
	return nil
}
