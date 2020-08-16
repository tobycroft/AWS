package http

import (
	"github.com/gin-gonic/gin"
	"main.go/config"
)

func Handler(c *gin.Context) {
	key, ok := c.GetPostForm("key")
	if !ok {
		c.JSON(200, map[string]interface{}{
			"code": 400,
			"data": "key",
		})
		c.Abort()
		return
	}
	if key != config.KEY {
		c.JSON(200, map[string]interface{}{
			"code": 403,
			"data": "key_error",
		})
		c.Abort()
		return
	}
	to_user:=
}
