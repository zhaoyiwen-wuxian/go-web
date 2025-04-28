package service

import (
	"fmt"
	"go-web/appError"
	buildconditionsmap "go-web/buildConditionsMap"
	"go-web/common"
	"go-web/models"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateGroup 创建群组
func CreateGroup(db *gorm.DB, name string, ownerID uint) (*models.Group, error) {
	group := models.Group{
		Name:    name,
		OwnerID: ownerID,
	}

	// 使用通用 Create 操作
	if err := common.Create(db, &group); err != nil {
		return nil, appError.NewAppError(500, "创建群组失败", err)
	}

	return &group, nil
}

// DeleteGroup 删除群组
func DeleteGroup(db *gorm.DB, groupID uint, ownerID uint) error {
	// 检查群组是否存在并且是群主
	var group models.Group
	if err := common.QueryOne(db, &models.Group{}, &group, buildconditionsmap.ConditionConfig{
		TableFields: map[string][]string{"group": {"id", "owner_id"}},
		ExactFields: []string{"id", "owner_id"},
	}, nil); err != nil {
		return appError.NewAppError(500, "群组不存在或无权限删除", err)
	}

	// 删除群组成员
	if err := common.Delete(db, &models.GroupMember{}, map[string]interface{}{"group_id": groupID}); err != nil {
		return appError.NewAppError(500, "删除群组成员失败", err)
	}

	// 删除群组
	return common.Delete(db, &models.Group{}, map[string]interface{}{"id": groupID, "owner_id": ownerID})
}

// AddUserToGroup 添加用户到群组
func AddUserToGroup(db *gorm.DB, userID, groupID uint) error {
	// 1. 检查群组是否存在
	var group models.Group
	groupConfig := buildconditionsmap.ConditionConfig{
		TableFields: map[string][]string{"group": {"id"}},
		ExactFields: []string{"id"},
		FieldValues: map[string]interface{}{"id": groupID},
	}
	if err := common.QueryOne(db, &models.Group{}, &group, groupConfig, nil); err != nil {
		return appError.NewAppError(404, "群组不存在", err)
	}

	// 2. 检查用户是否已经是群组成员
	var member models.GroupMember
	memberConfig := buildconditionsmap.ConditionConfig{
		TableFields: map[string][]string{"group_member": {"user_id", "group_id"}},
		ExactFields: []string{"user_id", "group_id"},
		FieldValues: map[string]interface{}{"user_id": userID, "group_id": groupID},
	}
	if err := common.QueryOne(db, &models.GroupMember{}, &member, memberConfig, nil); err == nil {
		return appError.NewAppError(400, "用户已经是群组成员", nil)
	}

	// 3. 将用户添加到群组
	groupMember := models.GroupMember{
		UserID:  userID,
		GroupID: groupID,
	}
	if err := common.Create(db, &groupMember); err != nil {
		return appError.NewAppError(500, "添加用户到群组失败", err)
	}

	return nil
}

// 获取用户加入的群组
func GetUserGroups(db *gorm.DB, userID uint, c *gin.Context) ([]models.Group, error) {
	// 创建查询条件
	conditionConfig := buildconditionsmap.ConditionConfig{
		TableFields: map[string][]string{
			"group_member": {"user_id", "group_id"},
		},
		ExactFields: []string{"group_member.user_id"},
		FieldValues: map[string]interface{}{
			"group_member.user_id": userID,
		},
	}

	// 连接查询群组信息
	var groups []models.Group
	joins := []string{"JOIN group_member ON group_member.group_id = group.id"}

	// 执行查询
	result, err := common.QueryAllWithPagination(db, models.Group{}, &groups, conditionConfig, nil, c, joins)
	if err != nil {
		return nil, err
	}

	return result.Data.([]models.Group), nil
}

// SendNotification 把通知消息写入数据库
func SendNotification(db *gorm.DB, msg models.Message) error {
	return common.Create(db, &msg)
}

// SendGroupDissolveNotification 发送群解散通知给所有群成员
func SendGroupDissolveNotification(db *gorm.DB, msg models.Message) error {
	// msg.To 存的是群ID的字符串
	groupID, err := strconv.ParseUint(msg.To, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid group ID: %v", err)
	}

	members, err := GetGroupMembersUsers(db, uint(groupID))
	if err != nil {
		return fmt.Errorf("获取群成员失败: %v", err)
	}

	for _, m := range members {
		notification := models.Message{
			From:    msg.From,                    // 群主或管理员
			To:      strconv.Itoa(int(m.UserID)), // 群成员的 userID
			Type:    "notification",
			Content: fmt.Sprintf("群组 %d 已解散。", groupID),
		}
		if err := SendNotification(db, notification); err != nil {
			return fmt.Errorf("发送群解散通知失败给用户 %d: %v", m.UserID, err)
		}
	}
	return nil
}

// HandleGroupDissolve 处理群解散：删库并广播通知
func HandleGroupDissolve(db *gorm.DB, msg models.Message) error {
	// 把 msg.To 当作群ID
	groupID, err := strconv.ParseUint(msg.To, 10, 64)
	if err != nil {
		return appError.NewAppError(400, "无效的群组ID", err)
	}

	// 1. 删除群成员
	if err := common.Delete(db, &models.GroupMember{}, map[string]interface{}{"group_id": uint(groupID)}); err != nil {
		log.Printf("删除群成员失败: %v", err)
		return appError.NewAppError(500, "删除群成员失败", err)
	}

	// 2. 删除群组
	if err := common.Delete(db, &models.Group{}, map[string]interface{}{"id": uint(groupID)}); err != nil {
		log.Printf("解散群组失败: %v", err)
		return appError.NewAppError(500, "解散群组失败", err)
	}

	// 3. 发送通知
	if err := SendGroupDissolveNotification(db, msg); err != nil {
		log.Printf("发送群解散通知失败: %v", err)
		return appError.NewAppError(500, "发送群解散通知失败", err)
	}

	log.Printf("群组 %d 已解散", groupID)
	return nil
}
