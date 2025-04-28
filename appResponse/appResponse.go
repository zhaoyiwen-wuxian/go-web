package appResponse

import (
	"go-web/appError"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 通用响应结构体
type AppResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// 返回的响应结构体
type AppResponsePage struct {
	Code       int         `json:"code"`
	Msg        string      `json:"msg"`
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	PageSize   int         `json:"pageSize"`
	TotalCount int64       `json:"total"`
}

// 成功响应函数
func SuccessPage(c *gin.Context, result PaginatedResult) {
	renderJSONPage(c, 200, AppResponsePage{
		Code:       200,
		Msg:        "成功",
		Data:       result.Data,       // 返回分页数据
		Page:       result.Page,       // 返回当前页
		PageSize:   result.PageSize,   // 返回每页大小
		TotalCount: result.TotalCount, // 返回总记录数
	})
}

// 成功响应
func Success(c *gin.Context, data interface{}) {
	renderJSON(c, 200, AppResponse{
		Code: 200,
		Msg:  "成功",
		Data: data,
	})
}

// 错误响应
func Error(c *gin.Context, code int, msg string, data interface{}) {
	renderJSON(c, http.StatusBadRequest, AppResponse{
		Code: code,
		Msg:  msg,
		Data: data,
	})
}

// 错误响应 (兼容 appError.AppError)
func ErrorFromApp(c *gin.Context, appErr error) {
	if err, ok := appErr.(*appError.AppError); ok {
		Error(c, err.Code, err.Message, err.Data)
	} else {
		Error(c, 50000, "系统错误", appErr.Error())
	}
}

// 内部封装 JSON 响应
func renderJSON(c *gin.Context, status int, resp AppResponse) {
	c.JSON(status, resp)
}

// 内部封装 JSON 响应
func renderJSONPage(c *gin.Context, status int, resp AppResponsePage) {
	c.JSON(status, resp)
}
