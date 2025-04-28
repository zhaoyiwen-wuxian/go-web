package router

import (
	"go-web/api"
	"go-web/middleware"
	"go-web/service"
	"go-web/ws"

	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.TokenHeartbeatMiddleware())
	hub := ws.NewHub()
	go hub.Run()
	ws.Subscribe(hub)
	go ws.LogThroughput()
	r.Use(middleware.RecoveryMiddleware())
	r.GET("/index", service.GetIndex)
	v1 := r.Group("/api/v1")
	{
		v1.POST("/uploadMedia", api.UploadMedia())
		user := v1.Group("/users")
		{

			user.POST("", api.CreateUserHandler())
			user.POST("/login", api.LoginUser())
		}
		// 需要认证的路由
		auth := v1.Group("/")
		auth.Use(middleware.JWTAuthMiddleware())

		{
			// Authenticated Routes
			auth.POST("/updatePassword", api.UpdateUserPasswordHandler())
			auth.POST("/UserByNameAndPhone", api.UserByNameAndPhone())
			auth.GET("/users", api.GetUsersHandler())
			auth.GET("/users", api.GetUserDetailHandler())
			auth.POST("/user/out", api.LoginOut())

			auth.GET("/ws/:userID", func(c *gin.Context) {
				api.WebSocketHandler(hub, c)
			})
			messages := auth.Group("/messages")
			{
				// 获取个人聊天记录（改为 POST）
				messages.POST("/personal/:userID/:otherUserID", api.GetPersonalMessagesHandler())

				// 获取群组消息（改为 POST）
				messages.POST("/group/:groupID", api.GetGroupMessagesHandler())

				// 获取群组成员（改为 POST）
				messages.POST("/group/:groupID/members", api.GetGroupMembersHandler())

				// 获取用户加入的群组（改为 POST）
				messages.POST("/users/:userID/groups", api.GetUserGroupsHandler())

				// 获取好友信息（改为 POST）
				messages.POST("/users/:userID/friends/:friendID", api.GetFriendInfoHandler())

				// 获取好友列表（改为 POST）
				messages.POST("/users/:userID/friends", api.GetFriendListHandler())
			}

		}
		return r
	}
}
