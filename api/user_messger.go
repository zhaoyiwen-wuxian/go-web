package api

import (
	"go-web/appError"
	"go-web/appResponse"
	"go-web/models"
	"go-web/service"
	"go-web/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Produce json
// @Param otherUserID path int true "其他用户ID"
// @Success 200 {object} appResponse.PaginatedResult
// @Router /api/v1/messages/personal/{otherUserID} [post]
func GetPersonalMessagesHandler() gin.HandlerFunc {
	return utils.HandlePaginatedRequest[models.Message, appResponse.PaginatedResult](
		func(c *gin.Context, req *models.Message) error {
			// 绑定请求参数
			userID, err := utils.SomeHandler(c)
			if err != nil {
				return appError.ErrInvalidParams
			}
			*&req.To = string(userID)
			return c.ShouldBindQuery(req)
		},
		func(c *gin.Context, req *models.Message) (appResponse.PaginatedResult, *appError.AppError) {

			otherUserID, _ := strconv.Atoi(c.Param("otherUserID"))

			// 调用服务层获取个人消息
			userID, _ := strconv.Atoi(req.To)
			paginatedResult, err := service.GetPersonalMessages(utils.DB, uint(userID), uint(otherUserID), c)
			if err != nil {
				return appResponse.PaginatedResult{}, appError.NewAppErrorFromError(err)
			}
			return paginatedResult, nil
		},
	)
}

// @Produce json
// @Param groupID path int true "群组ID"
// @Success 200 {object} appResponse.PaginatedResult
// @Router /api/v1/messages/group/{groupID} [post]
func GetGroupMessagesHandler() gin.HandlerFunc {
	return utils.HandlePaginatedRequest[models.Message, appResponse.PaginatedResult](
		func(c *gin.Context, req *models.Message) error {
			// 绑定请求参数
			return c.ShouldBindQuery(req)
		},
		func(c *gin.Context, req *models.Message) (appResponse.PaginatedResult, *appError.AppError) {
			// 从 URL 获取 groupID
			groupID, _ := strconv.Atoi(c.Param("groupID"))

			paginatedResult, err := service.GetGroupMessages(utils.DB, uint(groupID), c)
			if err != nil {
				return appResponse.PaginatedResult{}, appError.NewAppErrorFromError(err)
			}
			return paginatedResult, nil
		},
	)
}

// @Produce json
// @Param groupID path int true "群组ID"
// @Success 200 {array} models.UserBasic
// @Router /api/v1/messages/group/{groupID}/members [post]
func GetGroupMembersHandler() gin.HandlerFunc {
	return utils.HandleRequest(
		func(c *gin.Context, req *int) error {
			// 从 URL 获取 groupID
			groupID, err := strconv.Atoi(c.Param("groupID"))
			if err != nil {
				return appError.ErrInvalidParams
			}
			*req = groupID
			return nil
		},
		func(c *gin.Context, groupID *int) (any, *appError.AppError) {

			// 调用服务层获取群组成员
			members, err := service.GetGroupMembers(utils.DB, uint(*groupID), c)
			if err != nil {
				return nil, appError.NewAppErrorFromError(err)
			}
			return members, nil
		},
	)
}

// @Produce json
// @Param userID path int true "用户ID"
// @Success 200 {array} models.Group
// @Router /api/v1/messages/users/groups [post]
func GetUserGroupsHandler() gin.HandlerFunc {
	return utils.HandleRequest(
		func(c *gin.Context, req *int) error {
			// 获取 userID
			userID, err := utils.SomeHandler(c)
			if err != nil {
				return appError.ErrInvalidParams
			}
			*req = userID
			return nil
		},
		func(c *gin.Context, userID *int) (any, *appError.AppError) {
			// 调用服务层获取用户加入的群组
			groups, err := service.GetUserGroups(utils.DB, uint(*userID), c)
			if err != nil {
				return nil, appError.NewAppErrorFromError(err)
			}
			return groups, nil
		},
	)
}

// @Produce json
// @Param userID path int true "用户ID"
// @Param friendID path int true "好友ID"
// @Success 200 {object} models.UserBasic
// @Router /api/v1/messages/users/friends/{friendID} [post]
func GetFriendInfoHandler() gin.HandlerFunc {
	return utils.HandleRequest(
		func(c *gin.Context, req *struct{ UserID, FriendID int }) error {
			// 获取 userID
			userID, err := utils.SomeHandler(c)
			if err != nil {
				return appError.ErrInvalidParams
			}
			friendID, err1 := strconv.Atoi(c.Param("friendID"))
			if err1 != nil {
				return appError.ErrInvalidParams
			}

			req.UserID = userID
			req.FriendID = friendID
			return nil
		},
		func(c *gin.Context, req *struct{ UserID, FriendID int }) (any, *appError.AppError) {
			// 调用服务层获取好友信息
			friend, err := service.GetFriendInfo(utils.DB, uint(req.UserID), uint(req.FriendID))
			if err != nil {
				return nil, appError.NewAppErrorFromError(err)
			}
			return friend, nil
		},
	)
}

// @Produce json
// @Param userID path int true "用户ID"
// @Success 200 {array} models.UserBasic
// @Router /api/v1/messages/users/friends [post]
func GetFriendListHandler() gin.HandlerFunc {
	return utils.HandleRequest(
		func(c *gin.Context, req *int) error {
			// 获取 userID
			userID, err := utils.SomeHandler(c)
			if err != nil {
				return appError.ErrInvalidParams
			}
			*req = userID
			return nil
		},
		func(c *gin.Context, userID *int) (any, *appError.AppError) {
			// 调用服务层获取好友列表
			friends, err := service.GetFriendList(utils.DB, uint(*userID), c)
			if err != nil {
				return nil, appError.NewAppErrorFromError(err)
			}
			return friends, nil
		},
	)
}
