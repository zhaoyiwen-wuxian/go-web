package ws

import (
	"go-web/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func ServeWS(hub *Hub, w http.ResponseWriter, r *http.Request, userID string, c *gin.Context) {
	url := r.Header.Get("X-Client-URL")
	conn, _ := upgrader.Upgrade(w, r, nil)
	client := &Client{Conn: conn, Send: make(chan []byte), UserID: userID, url: url}

	if !client.authenticateRequest(r, c) {
		// 如果身份验证失败，关闭连接并返回
		conn.Close()
		log.Println("身份验证失败，关闭连接")
		return
	}

	hub.Register <- client
	go client.readPump(hub, utils.DB, c)
	go client.writePump()
}
