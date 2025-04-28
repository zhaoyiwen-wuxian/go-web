package middleware

import (
	"go-web/appError"
	"go-web/appResponse"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 从 URL 路径参数中获取分页信息
func GetPaginationParamsFromPath(c *gin.Context) (int, int, error) {
	defaultPage := 1
	defaultPageSize := 10

	pageStr := c.Param("page")
	pageSizeStr := c.Param("pageSize")

	page := defaultPage
	if pageStr != "" {
		var err error
		page, err = strconv.Atoi(pageStr)
		if err != nil || page <= 0 {
			return 0, 0, appError.ErrInternalServerError
		}
	}

	pageSize := defaultPageSize
	if pageSizeStr != "" {
		var err error
		pageSize, err = strconv.Atoi(pageSizeStr)
		if err != nil || pageSize <= 0 {
			return 0, 0, appError.ErrInternalServerError
		}
	}

	return page, pageSize, nil
}

// 中间件：从 URL 路径参数中提取分页参数
func PaginationMiddlewareFromPath() gin.HandlerFunc {
	return func(c *gin.Context) {
		page, pageSize, err := GetPaginationParamsFromPath(c)
		if err != nil {
			appResponse.ErrorFromApp(c, err)
			c.Abort()
			return
		}
		c.Set("page", page)
		c.Set("pageSize", pageSize)
		c.Next()
	}
}
