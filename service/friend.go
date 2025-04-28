package service

import (
	"go-web/models"
	"log"
	"strconv"

	"gorm.io/gorm"
)

// 处理好友请求
func HandleFriendRequest(db *gorm.DB, msg models.Message) error {
	fromID, _ := strconv.Atoi(msg.From)
	toID, _ := strconv.Atoi(msg.To)
	if err := AddFriend(db, uint(fromID), uint(toID)); err != nil {
		log.Println("添加好友失败:", err)
		return err
	}
	return nil
}

// 处理接受好友请求
func HandleFriendAccept(db *gorm.DB, msg models.Message) error {
	fromID, _ := strconv.Atoi(msg.From)
	toID, _ := strconv.Atoi(msg.To)
	if err := AcceptFriend(db, uint(fromID), uint(toID)); err != nil {
		log.Println("接受好友请求失败:", err)
		return err
	}
	return nil
}

// 处理拒绝好友请求
func HandleFriendReject(db *gorm.DB, msg models.Message) error {
	fromID, _ := strconv.Atoi(msg.From)
	toID, _ := strconv.Atoi(msg.To)
	if err := RejectFriend(db, uint(fromID), uint(toID)); err != nil {
		log.Println("拒绝好友请求失败:", err)
		return err
	}
	return nil
}

// 处理拉黑好友
func HandleFriendBlock(db *gorm.DB, msg models.Message) error {
	userID, _ := strconv.Atoi(msg.From)
	friendID, _ := strconv.Atoi(msg.To)
	if err := BlockFriend(db, uint(userID), uint(friendID)); err != nil {
		log.Println("拉黑好友失败:", err)
		return err
	}
	return nil
}
