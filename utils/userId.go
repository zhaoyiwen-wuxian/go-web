package utils

import (
	"errors"
	"go-web/appError"
	"go-web/redisutil"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func SomeHandler(c *gin.Context) (int, *appError.AppError) {
	token, err := GetTokenFromHeader(c)
	if err != nil {

		return 0, appError.NewAppError(401, "token 不存在", nil)
	}

	userID, err := redisutil.GetToken(c, token)
	if err != nil {

		return 0, appError.NewAppError(401, "类型断言失败", nil)
	}

	id, _ := strconv.ParseUint(userID, 10, 64)
	return int(id), nil

}

func GetTokenFromHeader(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", errors.New("请求头中缺少 Authorization")
	}

	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer "), nil
	}

	return authHeader, nil
}
