package service

import (
	"go-web/models"
	"log"
	"strconv"

	"gorm.io/gorm"
)

// 处理创建群组
func HandleGroupCreate(db *gorm.DB, msg models.Message) error {
	groupName := msg.Content
	creatorID, _ := strconv.Atoi(msg.From)
	if _, err := CreateGroup(db, groupName, uint(creatorID)); err != nil {
		log.Println("创建群组失败:", err)
		return err
	}
	return nil
}

// 处理加入群组
func HandleGroupJoin(db *gorm.DB, msg models.Message) error {
	groupID, _ := strconv.Atoi(msg.To)
	userID, _ := strconv.Atoi(msg.From)
	if err := AddUserToGroup(db, uint(userID), uint(groupID)); err != nil {
		log.Println("加入群组失败:", err)
		return err
	}
	return nil
}

// 处理离开群组
func HandleGroupLeave(db *gorm.DB, msg models.Message) error {
	groupID, _ := strconv.Atoi(msg.To)
	userID, _ := strconv.Atoi(msg.From)
	if err := DeleteGroup(db, uint(userID), uint(groupID)); err != nil {
		log.Println("离开群组失败:", err)
		return err
	}
	return nil
}
