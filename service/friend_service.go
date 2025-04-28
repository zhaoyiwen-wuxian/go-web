package service

import (
	"go-web/appError"
	buildconditionsmap "go-web/buildConditionsMap"
	"go-web/common"
	"go-web/models"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AddFriend(db *gorm.DB, from, to uint) error {
	// 检查用户是否尝试加自己为好友
	if from == to {
		return appError.NewAppError(500, "不能给自己发送好友请求", nil)
	}

	// 构建查询配置
	config := buildconditionsmap.ConditionConfig{
		TableFields: map[string][]string{
			"friendship": {"user_id", "friend_id"},
		},
		ExactFields: []string{"user_id", "friend_id"}, // 精确匹配 user_id 和 friend_id
	}

	var existingFriendship models.Friendship
	// 使用 QueryAll 查询是否已存在好友关系或请求
	err := common.QueryAll(db, &models.Friendship{}, &existingFriendship, config, nil, nil)
	if err == nil && existingFriendship.ID != 0 {
		// 如果查询成功且已有好友关系或请求存在
		return appError.NewAppError(500, "好友关系已存在或请求已存在", nil)
	}
	// 创建新的好友请求
	friendship := models.Friendship{
		UserID:   from,
		FriendID: to,
		Status:   "apply", // 请求状态为 "apply"
	}

	// 调用 Create 创建新的好友请求
	if err := common.Create(db, &friendship); err != nil {
		// 如果创建失败，返回错误
		return appError.NewAppError(500, "创建好友请求失败", err)
	}

	// 成功创建好友请求
	return nil
}

// AcceptFriend 接受好友请求
func AcceptFriend(db *gorm.DB, from, to uint) error {
	var friendship models.Friendship
	// 查询是否有待处理的请求
	if err := common.QueryOne(db, &models.Friendship{}, &friendship, buildconditionsmap.ConditionConfig{
		TableFields: map[string][]string{"friendship": {"user_id", "friend_id"}},
		ExactFields: []string{"user_id", "friend_id"},
	}, nil); err != nil {
		return appError.NewAppError(500, "未找到请求", nil)
	}

	// 更新请求状态为已接受
	friendship.Status = "accepted"
	return common.Update(db, &models.Friendship{}, map[string]interface{}{"user_id": from, "friend_id": to}, friendship)
}

func BlockFriend(db *gorm.DB, userID, friendID uint) error {
	return common.Update(db, &models.Friendship{}, map[string]interface{}{"user_id": userID, "friend_id": friendID}, map[string]interface{}{"status": "blocked"})
}

// RejectFriend 拒绝好友请求
func RejectFriend(db *gorm.DB, from, to uint) error {
	var friendship models.Friendship
	if err := common.QueryOne(db, &models.Friendship{}, &friendship, buildconditionsmap.ConditionConfig{
		TableFields: map[string][]string{"friendship": {"user_id", "friend_id"}},
		ExactFields: []string{"user_id", "friend_id"},
	}, nil); err != nil {
		return appError.NewAppError(500, "未找到请求", nil)
	}

	// 更新请求状态为已拒绝
	friendship.Status = "rejected"
	return common.Update(db, &models.Friendship{}, map[string]interface{}{"user_id": from, "friend_id": to}, friendship)
}

func GetFriendList(db *gorm.DB, userID uint, c *gin.Context) ([]models.UserBasic, error) {
	// 创建查询条件
	conditionConfig := buildconditionsmap.ConditionConfig{
		TableFields: map[string][]string{
			"friendship": {"user_id", "friend_id", "status"},
		},
		ExactFields: []string{"friendship.user_id"},
		FieldValues: map[string]interface{}{
			"friendship.user_id": userID,
			"friendship.status":  "accepted", // 只查询已接受的好友
		},
	}

	// 连接查询好友的基本信息
	var friends []models.UserBasic
	joins := []string{"JOIN friendship ON friendship.friend_id = user_basic.id"}

	// 执行查询
	result, err := common.QueryAllWithPagination(db, models.UserBasic{}, &friends, conditionConfig, nil, c, joins)
	if err != nil {
		return nil, err
	}

	return result.Data.([]models.UserBasic), nil
}

// 获取特定好友信息
func GetFriendInfo(db *gorm.DB, userID uint, friendID uint) (*models.UserBasic, error) {
	// 创建查询条件
	conditionConfig := buildconditionsmap.ConditionConfig{
		TableFields: map[string][]string{
			"user_basic": {"id", "name", "phone", "email"},
		},
		ExactFields: []string{"id"},
		FieldValues: map[string]interface{}{
			"id": friendID,
		},
	}

	var friend models.UserBasic
	// 查询单条数据
	err := common.QueryOne(db, models.UserBasic{}, &friend, conditionConfig, nil)
	if err != nil {
		return nil, err
	}

	return &friend, nil
}

// 发送好友添加通知给双方
func SendFriendAddNotification(db *gorm.DB, msg models.Message) error {
	// 获取发送者和接收者的ID
	senderID := msg.From
	receiverID := msg.To

	// 创建发送者的通知
	senderNotification := models.Message{
		From:    senderID,
		To:      receiverID,
		Type:    "notification",
		Content: "你和" + receiverID + "成为了好友！",
	}

	// 创建接收者的通知
	receiverNotification := models.Message{
		From:    receiverID,
		To:      senderID,
		Type:    "notification",
		Content: "你和" + senderID + "成为了好友！",
	}

	// 使用通用的插入方法将通知插入到数据库

	if err := common.Create(db, &senderNotification); err != nil {
		log.Printf("发送好友添加通知失败，发送者通知插入失败: %v", err)
		return err
	}
	if err := common.Create(db, &receiverNotification); err != nil {
		log.Printf("发送好友添加通知失败，接收者通知插入失败: %v", err)
		return err
	}

	log.Printf("好友添加通知已成功发送给 %v 和 %v", senderID, receiverID)
	return nil
}
