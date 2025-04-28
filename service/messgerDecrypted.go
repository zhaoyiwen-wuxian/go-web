package service

import (
	"go-web/models"
	"go-web/utils"
	"log"
)

// DecryptMessage 批量解密消息内容，解密失败的消息将保留原内容
func DecryptMessage(messages []models.Message, key []byte) ([]models.Message, error) {
	for i, msg := range messages {
		decrypted, err := utils.DecryptMessage([]byte(msg.Content), key)
		if err != nil {
			log.Printf("解密消息 ID %d 时失败: %v", msg.ID, err)
			continue // 解密失败跳过
		}
		messages[i].Content = string(decrypted)
	}
	return messages, nil
}
