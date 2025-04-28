package service

import (
	"go-web/models"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 处理普通聊天消息
func HandleMessage(db *gorm.DB, msg models.Message) error {
	to, _ := strconv.Atoi(msg.To)
	from, _ := strconv.Atoi(msg.From)
	err := SendMessage(db, uint(from), uint(to), msg.Content, msg.Type, msg.Extra)
	if err != nil {
		log.Println("发送消息失败:", err)
		return err
	}
	return nil
}

// 处理获取消息请求
func HandleGetMessages(db *gorm.DB, con *gin.Context, msg models.Message) ([]models.Message, error) {
	from, _ := strconv.Atoi(msg.From)
	messages, err := GetMessages(db, con, uint(from), con)
	if err != nil {
		log.Println("获取消息失败:", err)
		return nil, err
	}
	return messages.Data.([]models.Message), nil
}
