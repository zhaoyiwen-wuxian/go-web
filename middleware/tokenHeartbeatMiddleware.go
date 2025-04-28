package middleware

import (
	"go-web/appResponse"
	"go-web/jwtutil"
	"go-web/redisutil"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Token心跳更新中间件
func TokenHeartbeatMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取 Token
		token := c.GetHeader("Authorization")
		if token == "" {
			c.Next()
			return
		}

		// 去掉 Bearer 前缀
		token = strings.TrimPrefix(token, "Bearer ")

		// 解析 Token
		claims, err := jwtutil.ParseToken(token)
		if err != nil {
			c.Next()
			return
		}

		// 刷新 Token 的过期时间（心跳）
		err = jwtutil.UpdateTokenExpiration(c, redisutil.TokenPrefix+token, claims.UserID, time.Hour*6)
		if err != nil {
			appResponse.Error(c, http.StatusUnauthorized, "更新 Token 过期时间失败", err.Error())
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)

		c.Next()
	}
}
