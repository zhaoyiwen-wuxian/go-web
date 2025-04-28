package middleware

import (
	"fmt"

	"go-web/appError"
	"go-web/appResponse"

	"github.com/gin-gonic/gin"
)

func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("panic recovered: %v\n", r)

				// 返回统一格式错误响应
				appResponse.ErrorFromApp(c, appError.ErrInternalServerError)
				c.Abort()
			}
		}()
		c.Next()
	}
}
