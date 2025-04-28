package service

import (
	"go-web/appError"
	buildconditionsmap "go-web/buildConditionsMap"
	"go-web/common"
	"go-web/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AddMemberToGroup 添加群成员
func AddMemberToGroup(db *gorm.DB, groupID uint, userID uint, role string) error {
	// 检查用户是否已经是成员
	var groupMember models.GroupMember
	err := common.QueryOne(db, &models.GroupMember{}, &groupMember, buildconditionsmap.ConditionConfig{
		TableFields: map[string][]string{"group_member": {"group_id", "user_id"}},
		ExactFields: []string{"group_id", "user_id"},
	}, nil)

	if err == nil && groupMember.ID != 0 {
		return appError.NewAppError(500, "用户已是群成员", nil)
	}

	// 添加用户为群成员
	groupMember = models.GroupMember{
		GroupID: groupID,
		UserID:  userID,
		Role:    role,
	}

	// 使用通用的插入方法
	return common.Create(db, &groupMember)
}

// RemoveMemberFromGroup 移除群成员
func RemoveMemberFromGroup(db *gorm.DB, groupID uint, userID uint) error {
	// 删除群成员
	return common.Delete(db, &models.GroupMember{}, map[string]interface{}{"group_id": groupID, "user_id": userID})
}

// 获取群组成员
func GetGroupMembers(db *gorm.DB, groupID uint, c *gin.Context) ([]models.UserBasic, error) {
	// 创建查询条件
	conditionConfig := buildconditionsmap.ConditionConfig{
		TableFields: map[string][]string{
			"group_member": {"user_id", "group_id"},
		},
		ExactFields: []string{"group_member.group_id"},
		FieldValues: map[string]interface{}{
			"group_member.group_id": groupID,
		},
	}

	// 连接查询用户信息
	var members []models.UserBasic
	joins := []string{"JOIN user_basic ON user_basic.id = group_member.user_id"}

	// 执行查询
	result, err := common.QueryAllWithPagination(db, models.UserBasic{}, &members, conditionConfig, nil, c, joins)
	if err != nil {
		return nil, err
	}

	return result.Data.([]models.UserBasic), nil
}

// GetGroupMembers 读取某个群的所有成员
func GetGroupMembersUsers(db *gorm.DB, groupID uint) ([]models.GroupMember, error) {
	var members []models.GroupMember
	// 不分页地查询所有 group_member 表中 group_id = ?
	config := buildconditionsmap.ConditionConfig{
		TableFields: map[string][]string{"group_member": {"group_id"}},
		ExactFields: []string{"group_id"},
		FieldValues: map[string]interface{}{"group_id": groupID},
	}
	if err := common.QueryAll(db, &models.GroupMember{}, &members, config, nil, nil); err != nil {
		return nil, err
	}
	return members, nil
}
