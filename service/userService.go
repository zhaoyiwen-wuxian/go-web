package service

import (
	"fmt"
	"go-web/appError"
	"go-web/appResponse"
	buildconditionsmap "go-web/buildConditionsMap"
	"go-web/common"
	"go-web/models"
	"go-web/redisutil"
	"go-web/req"
	"go-web/resp"
	"go-web/utils"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 查询用户（可传入 join）
func GetUsers(db *gorm.DB, input models.UserBasic, c *gin.Context) (appResponse.PaginatedResult, error) {
	var users []models.UserBasic

	// 配置字段查询规则
	config := buildconditionsmap.ConditionConfig{
		ExactFields: []string{"name", "age"},
		LikeFields:  []string{"deviceInfo"},
		RangeFields: []string{"phone"},
	}
	//多表查询列如这样如下：
	// config := buildconditionsmap.ConditionConfig{
	// 	TableFields: map[string][]string{
	// 		"users": {"name", "age", "created_at"},
	// 	},
	// 	ExactFields: []string{"name"},
	// 	LikeFields:  []string{"email"},
	// 	RangeFields: []string{"created_at"},
	// 	FieldValues: map[string]interface{}{
	// 		"name": "Tom",
	// 		"created_at": map[string]interface{}{
	// 			"gt": "2023-01-01",
	// 			"lt": "2023-12-31",
	// 		},
	// 	},
	// }
	paginatedResult, err := common.QueryAllWithPagination(db, &models.UserBasic{}, &users, config, []common.SortField{
		{Field: "name", Direction: "asc"},
	}, c, nil)

	return paginatedResult, err
}

// 查询单个用户
func GetUserDetail(db *gorm.DB, id uint) (*resp.UserResp, error) {
	var user resp.UserResp
	conditions := buildconditionsmap.ConditionConfig{
		ExactFields: []string{"id"},
	}
	err := common.QueryOne(db, &models.UserBasic{}, &user, conditions, nil)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// 创建用户
func CreateUser(db *gorm.DB, req *req.UserReq, c *gin.Context) error {
	var input models.UserBasic
	err2, done := saveIsValidata(db, req, c, input)
	if done {
		return err2
	}
	return common.Create(db, input)
}

// 更新用户
func UpdateUser(db *gorm.DB, id uint, input models.UserBasic) error {
	return common.Update(db, &models.UserBasic{}, map[string]interface{}{"id": id}, input)
}

// 修改密码
func UpdatePassword(db *gorm.DB, id uint, userReq req.UserReq) error {
	var user models.UserBasic
	err2, done := updateNewPassword(db, user, userReq)
	if done {
		return err2
	}
	return common.Update(db, &models.UserBasic{}, map[string]interface{}{"id": id}, user)
}

// 根据账户和手机号查询用户
func UserByNameAndPhone(name string, phone string, db *gorm.DB) (*resp.UserResp, error) {
	var user resp.UserResp
	userResp, err2, done := verifyUserNameAndPhone(name, phone, db, user)
	if done {
		return userResp, err2
	}
	return &user, nil
}

// 登录
func LoginUser(name string, password string, db *gorm.DB, c *gin.Context) (*resp.UserResp, error) {
	var user models.UserBasic
	userResp, err2, done := funcUserLogin(name, password, db, user, c)
	if done {
		return userResp, err2
	}

	return &resp.UserResp{
		Name:  user.Name,
		ID:    user.ID,
		Phone: user.Phone,
		Email: user.Email,
	}, nil
}

func LoginOut(db *gorm.DB, c *gin.Context) error {
	var user models.UserBasic
	id, err := utils.SomeHandler(c)

	if err != nil {
		return appError.ErrInvalidParams
	}
	user.ID = uint(id)
	return queryUser(db, user, c)
}

func UpdateHeartbeat(db *gorm.DB, userID uint, c *gin.Context) error {
	heartbeatTime := time.Now().Unix()
	err := redisutil.SetJSONRedis(c, fmt.Sprintf(redisutil.HeartbeatKey, userID), heartbeatTime, 1*time.Hour)
	if err != nil {
		return appError.NewAppError(500, "心跳更新失败", nil)
	}
	return common.Update(db, &models.UserBasic{}, map[string]interface{}{"id": userID}, map[string]interface{}{"heartbeat_time": time.Now()})
}
