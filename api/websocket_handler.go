package api

import (
	"go-web/ws"

	"github.com/gin-gonic/gin"
)

func WebSocketHandler(hub *ws.Hub, c *gin.Context) {

	userID := c.Param("userID")

	ws.ServeWS(hub, c.Writer, c.Request, userID, c)
}
