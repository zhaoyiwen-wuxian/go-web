package ws

import (
	"go-web/models"
	"go-web/service"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func switchMsgType(msg models.Message, db *gorm.DB, c *Client, con *gin.Context) bool {
	switch msg.Type {
	case "friend_request":
		if err := service.HandleFriendRequest(db, msg); err != nil {
			c.SendError("处理好友请求失败")
		}
	case "friend_accept":
		// 处理接受好友请求
		if err := service.HandleFriendAccept(db, msg); err != nil {
			c.SendError("接受好友请求失败")
		} else {
			// 发送好友添加通知给双方
			if err := service.SendFriendAddNotification(db, msg); err != nil {
				c.SendError("发送好友添加通知失败")
			}
		}
	case "friend_reject":
		if err := service.HandleFriendReject(db, msg); err != nil {
			c.SendError("拒绝好友请求失败")
		}
	case "friend_block":
		if err := service.HandleFriendBlock(db, msg); err != nil {
			c.SendError("拉黑好友失败")
		}
	case "group_create":
		if err := service.HandleGroupCreate(db, msg); err != nil {
			c.SendError("创建群组失败")
		}
	case "group_join":
		if err := service.HandleGroupJoin(db, msg); err != nil {
			c.SendError("加入群组失败")
		}
	case "group_leave":
		if err := service.HandleGroupLeave(db, msg); err != nil {
			c.SendError("离开群组失败")
		}
	case "group_dissolve":
		// 群解散时发送通知给群成员
		if err := service.HandleGroupDissolve(db, msg); err != nil {
			c.SendError("群解散失败")
		} else {
			if err := service.SendGroupDissolveNotification(db, msg); err != nil {
				c.SendError("发送群解散通知失败")
			}
		}
	case "message":
		if err := service.HandleMessage(db, msg); err != nil {
			c.SendError("发送消息失败")
		}
	case "get_message":
		messages, err := service.HandleGetMessages(db, con, msg)
		if err != nil {
			c.SendError("获取消息失败")
			break
		}
		// 逐条发送消息
		for _, message := range messages {
			if err := c.SendMessage(message); err != nil {
				c.SendError("处理消息时发生错误")
				return true
			}
		}
	case "heartbeat":
		// 更新心跳
		userID, _ := strconv.Atoi(msg.From)
		if err := service.UpdateHeartbeat(db, uint(userID), con); err != nil {
			c.SendError("更新心跳失败")
		}
	case "emoji", "audio":
		// 通用的聊天消息处理
		if err := service.HandleMessage(db, msg); err != nil {
			c.SendError("发送" + msg.Type + "消息失败")
		}
	case "voice_invite", "video_invite", "call_offer", "call_answer", "call_candidate", "call_leave":
		roomID := msg.To
		Manager.JoinRoom(roomID, c.UserID, c)
		Manager.Broadcast(roomID, c.UserID, msg)

	case "call_accept", "call_reject":
		roomID := msg.To
		Manager.Broadcast(roomID, c.UserID, msg)
	default:
		log.Println("未知消息类型:", msg.Type)
	}
	return false
}
