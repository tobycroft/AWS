package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"io"
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
			defer func() {
				// 发送 websocket 结束包
				conn.WriteMessage(websocket.CloseMessage,
					websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				// 真正关闭 conn
				conn.Close()
			}()
			// 读取一个包
			mt, d, err := conn.ReadMessage()

			if err != nil {
				fmt.Printf("read fail = %v\n", err)
				return
			}

			fmt.Printf("data:%s\n", d)
			// 写入一个包
			err = conn.WriteMessage(mt, d)
			if err != nil {
				fmt.Printf("write fail = %v\n", err)
				return
			}
		}
	})

	// http echo
	r.GET("/http", func(c *gin.Context) {
		io.Copy(c.Writer, c.Request.Body)
	})

	r.Run(":80")

}
