package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"main.go/config"
	"main.go/function/http"
	"main.go/function/ws"
)

var upgrader = websocket.Upgrader{}

func main() {
	r := gin.Default()
	// websocket echo
	r.Any("/", func(c *gin.Context) {
		r := c.Request
		w := c.Writer
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Printf("err = %s\n", err)
			http.Handler(c)
			return
		} else {
			ws_handler(conn)
		}
	})
	r.Any("/test", func(c *gin.Context) {
		c.File("html/index.html")
	})
	r.Any("/favicon.ico", func(c *gin.Context) {
		c.Abort()
	})
	r.Run(config.SERVER_LISTEN_ADDR + ":" + config.SERVER_LISTEN_PORT)
}

func ws_handler(conn *websocket.Conn) {
	defer on_close(conn)
	//连入时发送欢迎消息
	go on_connect(conn)
	for {
		mt, d, err := conn.ReadMessage()
		if mt == -1 {
			break
		}
		if err != nil {
			fmt.Println(mt)
			fmt.Printf("read fail = %v\n", err)
			break
		}
		ws.Handler(string(d))
	}
}

func on_connect(conn *websocket.Conn) {
	err := conn.WriteMessage(1, []byte("连入成功"))
	if err != nil {
		fmt.Printf("write fail = %v\n", err)
		return
	}
}

func on_close(conn *websocket.Conn) {
	// 发送 websocket 结束包
	conn.WriteMessage(websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	// 真正关闭 conn
	conn.Close()
}
