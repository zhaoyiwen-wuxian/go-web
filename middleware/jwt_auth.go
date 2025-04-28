package middleware

import (
	"net/http"
	"strings"

	"go-web/appResponse"
	"go-web/jwtutil"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware 验证 JWT Token 的中间件
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取 Token
		token := c.GetHeader("Authorization")
		if token == "" {
			appResponse.Error(c, http.StatusUnauthorized, "Token 是必须的", nil)
			c.Abort()
			return
		}

		// 去掉 Bearer 前缀
		token = strings.TrimPrefix(token, "Bearer ")

		// 解析 Token
		claims, err := jwtutil.ParseToken(token)
		if err != nil {
			appResponse.Error(c, http.StatusUnauthorized, "无效的 Token", err.Error())
			c.Abort()
			return
		}

		// 检查 Token 是否在 Redis 中有效
		_, err = jwtutil.GetTokenUserID(c, token)
		if err != nil {
			appResponse.Error(c, http.StatusUnauthorized, "Redis 中未找到 Token，可能已失效", err.Error())
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)

		c.Next()
	}
}
