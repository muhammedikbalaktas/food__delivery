package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{

	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}
var connections = make(map[string]*websocket.Conn)

func HandleConnections(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		fmt.Println("invalid parameters")
		return
	}
	role, err := parseToken(token)
	if err != nil {
		fmt.Println("invalid token")
		fmt.Println(err)
		return
	}
	fmt.Println("api connected with role of ", role)

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("connection refused")
		return
	}
	connections[role] = conn
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("error on reading message")
			return
		}
		pushMessages(messageType, string(message))
	}
}
func pushMessages(messageType int, message string) {
	for _, conn := range connections {
		err := conn.WriteMessage(messageType, []byte(message))
		if err != nil {
			fmt.Println(err)
			fmt.Println("error on sending message ", message)
		}
	}

}
