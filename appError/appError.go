package appError

import (
	"fmt"
)

// 定义统一的错误结构体
type AppError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (e *AppError) Error() string {
	return fmt.Sprintf("code=%d, message=%s", e.Code, e.Message)
}

// 构造通用错误
func NewAppError(code int, message string, data interface{}) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

func NewAppErrorFromError(err error) *AppError {
	if err == nil {
		return nil
	}
	return &AppError{
		Code:    500,
		Message: err.Error(),
		Data:    nil,
	}
}

// 快捷错误函数
var (
	ErrModelNil            = NewAppError(1001, "Model 不能为空", nil)
	ErrInvalidParams       = NewAppError(1002, "参数错误", nil)
	ErrDBQueryFailed       = NewAppError(1003, "数据库查询失败", nil)
	ErrInternalServerError = NewAppError(500, "系统内部错误", nil)
	ErrHanderFailed        = NewAppError(507, "获取分页参数错误", nil)
)
