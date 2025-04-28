package utils

import (
	"go-web/appError"
	"go-web/appResponse"

	"github.com/gin-gonic/gin"
)

// 泛型处理请求封装
func HandleRequest[Req any](
	bindFunc func(c *gin.Context, req *Req) error,
	handlerFunc func(c *gin.Context, req *Req) (any, *appError.AppError),
) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req Req
		if err := bindFunc(c, &req); err != nil {
			appResponse.ErrorFromApp(c, appError.NewAppError(400, err.Error(), nil))
			return
		}
		result, err := handlerFunc(c, &req)
		if err != nil {
			appResponse.ErrorFromApp(c, err)
			return
		}
		appResponse.Success(c, result)
	}
}

// HandleRequestPage 用于处理分页请求的通用函数
func HandlePaginatedRequest[Req any, T any](
	bindFunc func(c *gin.Context, req *Req) error,
	handlerFunc func(c *gin.Context, req *Req) (appResponse.PaginatedResult, *appError.AppError),
) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req Req
		// 绑定查询参数到结构体
		if err := bindFunc(c, &req); err != nil {
			appResponse.ErrorFromApp(c, appError.NewAppError(400, err.Error(), nil))
			return
		}
		// 调用 handler 函数获取分页数据
		res, err := handlerFunc(c, &req)
		if err != nil {
			appResponse.ErrorFromApp(c, err)
			return
		}
		// 返回分页结果
		appResponse.SuccessPage(c, res)
	}
}
