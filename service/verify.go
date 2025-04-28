package service

import (
	"go-web/appError"
	buildconditionsmap "go-web/buildConditionsMap"
	"go-web/common"
	getheader "go-web/getHeader"
	"go-web/jwtutil"
	"go-web/models"
	"go-web/redisutil"
	"go-web/req"
	"go-web/resp"
	"go-web/utils"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func isPassword(db *gorm.DB, req *req.UserReq, input models.UserBasic) (bool, bool, error) {
	if len(req.Password) < 8 || len(req.Password) > 20 {
		return true, true, appError.NewAppError(500, "密码长度必须在8-20位之间!!", nil)
	}
	if len(req.NewPassword) < 8 || len(req.NewPassword) > 20 {
		return true, true, appError.NewAppError(500, "密码长度必须在8-20位之间!!", nil)
	}
	if req.Password != req.NewPassword {
		return true, true, appError.NewAppError(500, "密码不一致请核对后在重新输入!!", nil)
	}
	conditions := buildconditionsmap.ConditionConfig{
		ExactFields: []string{"name", "phone"},
	}
	err := common.QueryOne(db, &models.UserBasic{}, &input, conditions, nil)
	if err != nil {
		return true, true, appError.NewAppError(500, "已有注册，请更换名称和手机号", nil)
	}
	return false, false, nil
}
func isValidata(req *req.UserReq) (bool, bool, error) {
	isValidateEmail := utils.ValidateEmail(req.Email)
	if !isValidateEmail {
		return true, true, appError.NewAppError(500, "邮箱错误", nil)
	}
	isValidatePhone := utils.ValidateEmail(req.Phone)
	if !isValidatePhone {
		return true, true, appError.NewAppError(500, "手机号码错误", nil)
	}
	return false, false, nil
}

func saveIsValidata(db *gorm.DB, req *req.UserReq, c *gin.Context, input models.UserBasic) (error, bool) {
	input.ClentIp = c.ClientIP()                        // 从 gin 上下文获取客户端 IP
	input.ClentPort = c.Request.RemoteAddr              // 获取客户端端口（通过远程地址）
	input.LoginTime = time.Now()                        // 当前时间作为登录时间
	input.HeartbeatTime = time.Now()                    // 当前时间作为心跳时间
	input.LoginOutTime = time.Now().Add(24 * time.Hour) // 假设默认的登出时间是24小时后
	// 从请求头中获取设备信息（User-Agent）
	input.DeviceInfo = getheader.GetDeviceInfoHeader(c)
	// 设置创建时间和更新时间
	input.CreatedAt = time.Now()
	input.UpdatedAt = time.Now()
	b, done, err2 := isValidata(req)
	if done {
		return err2, b
	}
	input.Email = req.Email
	input.Name = req.Name
	input.Phone = req.Phone
	b2, done2, err3 := isPassword(db, req, input)
	if done2 {
		return err3, b2
	}

	input.Password = utils.CalculateMD5(req.Password)
	return nil, false
}

func verifyUserNameAndPhone(name string, phone string, db *gorm.DB, user resp.UserResp) (*resp.UserResp, error, bool) {
	conditions := buildconditionsmap.ConditionConfig{
		ExactFields: []string{},
	}
	if name != "" {
		conditions.ExactFields = append(conditions.ExactFields, "name")
	}
	if phone != "" {
		conditions.ExactFields = append(conditions.ExactFields, "phone")
	}

	err := common.QueryOne(db, &models.UserBasic{}, &user, conditions, nil)
	if err != nil {
		return nil, appError.NewAppError(500, "无此用户", nil), true
	}
	return nil, nil, false
}

func updateNewPassword(db *gorm.DB, user models.UserBasic, userReq req.UserReq) (error, bool) {
	conditions := buildconditionsmap.ConditionConfig{
		ExactFields: []string{"id"},
	}
	err := common.QueryOne(db, &models.UserBasic{}, &user, conditions, nil)
	if err != nil {
		return appError.NewAppError(500, "无此用户", nil), true
	}
	user.Password = utils.CalculateMD5(userReq.Password)
	return nil, false
}

func funcUserLogin(name string, password string, db *gorm.DB, user models.UserBasic, c *gin.Context) (*resp.UserResp, error, bool) {
	conditions := buildconditionsmap.ConditionConfig{
		ExactFields: []string{},
	}
	if name != "" {
		conditions.ExactFields = append(conditions.ExactFields, "name")
	}
	err := common.QueryOne(db, &models.UserBasic{}, &user, conditions, nil)
	if err != nil {
		return nil, appError.NewAppError(500, "无此用户", nil), true
	}

	isPassword := utils.CheckPasswordHash(password, user.Password)
	if !isPassword {
		return nil, appError.NewAppError(500, "密码错误，请核对后在重新登录", nil), true
	}

	updates := map[string]interface{}{
		"clentIp":    c.ClientIP(),
		"clentPort":  c.Request.RemoteAddr,
		"loginTime":  time.Now(),
		"deviceInfo": getheader.GetDeviceInfoHeader(c),
	}
	err = common.UpdateOneByID(db, &models.UserBasic{}, user.ID, updates)
	if err != nil {
		return nil, appError.NewAppError(500, "更新失败", nil), false
	}
	token, err := jwtutil.GenerateToken(user.ID, redisutil.TokenPrefix, time.Hour*6)
	if err != nil {

		return nil, appError.NewAppError(50001, "生成 Token 失败", nil), false
	}
	err = redisutil.SetToken(c, token, user.ID, time.Hour*6)
	if err != nil {
		return nil, appError.NewAppError(50002, "缓存 Token 失败", nil), false
	}
	return nil, nil, false
}

func queryUser(db *gorm.DB, user models.UserBasic, c *gin.Context) error {
	conditions := buildconditionsmap.ConditionConfig{
		ExactFields: []string{"id"},
	}

	err := common.QueryOne(db, &models.UserBasic{}, &user, conditions, nil)
	if err != nil {
		return appError.NewAppError(500, "无此用户", nil)
	}
	updates := map[string]interface{}{
		"loginOutTime": time.Now(),
	}
	err = common.UpdateOneByID(db, &models.UserBasic{}, user.ID, updates)
	if err != nil {
		return appError.NewAppError(500, "更新失败", nil)
	}

	token, _ := utils.GetTokenFromHeader(c)
	redisutil.DeleteRedisKey(c, redisutil.TokenPrefix+token)
	return nil
}
