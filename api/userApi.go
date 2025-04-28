package api

import (
	"go-web/appError"
	"go-web/appResponse"
	"go-web/models"
	"go-web/req"
	"go-web/service"
	"go-web/utils"

	"github.com/gin-gonic/gin"
)

// @Summary 获取用户列表
// @Tags 用户管理
// @Accept json
// @Produce json
// @Success 200 {array} models.UserBasic
// @Router /api/v1/users [get]
func GetUsersHandler() gin.HandlerFunc {
	return utils.HandlePaginatedRequest[models.UserBasic, appResponse.PaginatedResult](
		func(c *gin.Context, req *models.UserBasic) error {
			return c.ShouldBindQuery(req)
		},
		func(c *gin.Context, req *models.UserBasic) (appResponse.PaginatedResult, *appError.AppError) {
			// 获取分页数据
			paginatedResult, err := service.GetUsers(utils.DB, *req, c)
			if err != nil {
				return appResponse.PaginatedResult{}, appError.NewAppErrorFromError(err)
			}
			return paginatedResult, nil
		},
	)
}

// @Summary 获取用户详情
// @Tags 用户管理
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} models.UserBasic
// @Router /api/v1/users [get]
func GetUserDetailHandler() gin.HandlerFunc {
	return utils.HandleRequest(
		func(c *gin.Context, req *int) error {
			id, err := utils.SomeHandler(c)
			if err != nil {
				return appError.ErrInvalidParams
			}
			*req = id
			return nil
		},
		func(c *gin.Context, id *int) (any, *appError.AppError) {
			user, err := service.GetUserDetail(utils.DB, uint(*id))
			if err != nil {
				return nil, appError.NewAppErrorFromError(err)
			}
			return user, nil
		},
	)
}

// @Summary 创建用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param user body models.UserBasic true "用户信息"
// @Success 201 {object} models.User
// @Router /api/v1/users [post]
func CreateUserHandler() gin.HandlerFunc {
	return utils.HandleRequest(
		func(c *gin.Context, req *req.UserReq) error {
			return c.ShouldBindJSON(req)
		},
		func(c *gin.Context, req *req.UserReq) (any, *appError.AppError) {
			if err := service.CreateUser(utils.DB, req, c); err != nil {
				return nil, appError.NewAppErrorFromError(err)
			}
			return nil, nil
		},
	)
}

// @Summary 修改密码
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param user body req.UserReq true "用户信息"
// @Success 201 {object} nil
// @Router /api/v1/users/updatePassword [post]
func UpdateUserPasswordHandler() gin.HandlerFunc {
	return utils.HandleRequest(
		func(c *gin.Context, req *req.UserReq) error {
			return c.ShouldBindJSON(req)
		},
		func(c *gin.Context, req *req.UserReq) (any, *appError.AppError) {
			id, err := utils.SomeHandler(c)

			if err != nil {
				return nil, appError.ErrInvalidParams
			}
			req.ID = uint(id)

			if err := service.UpdatePassword(utils.DB, req.ID, *req); err != nil {
				return nil, appError.NewAppErrorFromError(err)
			}
			return nil, nil
		},
	)
}

// @Summary 根据账户和手机号查询用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param user body req.UserReq true "用户信息"
// @Success 201 {object} models.UserResp
// @Router /api/v1/users/userByNameAndPhone [post]
func UserByNameAndPhone() gin.HandlerFunc {
	return utils.HandleRequest(
		func(c *gin.Context, req *req.UserReq) error {
			return c.ShouldBindJSON(req)
		},
		func(c *gin.Context, req *req.UserReq) (any, *appError.AppError) {
			resp, err := service.UserByNameAndPhone(req.Name, req.Phone, utils.DB)
			if err != nil {
				return nil, appError.NewAppErrorFromError(err)
			}
			return resp, nil
		},
	)
}

// @Summary 登录
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param user body req.UserReq true "用户信息"
// @Success 201 {object} models.UserResp
// @Router /api/v1/users/login [post]
func LoginUser() gin.HandlerFunc {
	return utils.HandleRequest(
		func(c *gin.Context, req *req.UserReq) error {
			return c.ShouldBindJSON(req)
		},
		func(c *gin.Context, req *req.UserReq) (any, *appError.AppError) {
			resp, err := service.LoginUser(req.Name, req.Password, utils.DB, c)
			if err != nil {
				return nil, appError.NewAppErrorFromError(err)
			}
			return resp, nil
		},
	)
}

// @Summary 退出
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param user body req.UserReq true "用户信息"
// @Success 201 {object}
// @Router /api/v1/users/loginOut [post]
func LoginOut() gin.HandlerFunc {
	return utils.HandleRequest(
		func(c *gin.Context, req *req.UserReq) error {
			return c.ShouldBindJSON(req)
		},
		func(c *gin.Context, req *req.UserReq) (any, *appError.AppError) {
			if err := service.LoginOut(utils.DB, c); err != nil {
				return nil, appError.NewAppErrorFromError(err)
			}
			return nil, nil
		},
	)
}
