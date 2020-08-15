package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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
			return
		} else {
			defer on_close(conn)
			// 读取一个包
			//mt, d, err := conn.ReadMessage()
			//
			//if err != nil {
			//	fmt.Printf("read fail = %v\n", err)
			//	return
			//}
			//
			//fmt.Printf("data:%s\n", d)
			// 写入一个包
			//err = conn.WriteMessage(mt, d)
			//if err != nil {
			//	fmt.Printf("write fail = %v\n", err)
			//	return
			//}
			on_connect(conn)
		}
	})

	r.Run(":80")
}

func on_connect(conn *websocket.Conn) {
	err := conn.WriteMessage(200, []byte("连入成功"))
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
