package ws

import (
	"go-web/jwtutil"

	"github.com/gin-gonic/gin"
)

// WebSocket 握手时的身份验证
func (c *Client) Authenticate(token string, ctx *gin.Context) bool {
	userID, err := jwtutil.GetTokenUserID(ctx, token)
	if err != nil {

		return false
	}
	if userID == "" {
		return false
	}

	return true
}
