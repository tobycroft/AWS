package http

import (
	"github.com/gin-gonic/gin"
	"main.go/config"
	"main.go/tuuz/Jsong"
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
	uids, ok := c.GetPostForm("uids")
	if !ok {
		c.JSON(200, map[string]interface{}{
			"code": 400,
			"data": "需要uids来发送数据",
		})
		c.Abort()
		return
	}
	to_users, err := Jsong.JArray(uids)
	if err != nil {
		c.JSON(200, map[string]interface{}{
			"code": 400,
			"data": "to_users_err",
		})
		c.Abort()
		return
	}
	dest, ok := c.GetPostForm("dest")
	if !ok {
		c.JSON(200, map[string]interface{}{
			"code": 400,
			"data": "需要dest参数，dest为chat_id或者gid",
		})
		c.Abort()
		return
	}
	data, ok := c.GetPostForm("data")
	if !ok {
		c.JSON(200, map[string]interface{}{
			"code": 400,
			"data": "data",
		})
		c.Abort()
		return
	}
	Type, ok := c.GetPostForm("type")
	if !ok {
		c.JSON(200, map[string]interface{}{
			"code": 400,
			"data": "需要type类型,chat/message_list/system",
		})
		c.Abort()
		return
	}
	json, jerr := Jsong.JObject(data)
	if jerr != nil {
		c.JSON(200, map[string]interface{}{
			"code": 400,
			"data": "json_err",
		})
		c.Abort()
		return
	}
	json_handler(json, to_users, dest, Type)
}

func json_handler(json map[string]interface{}, to_users []interface{}, dest string, Type string) {
	uids := []interface{}{}
	uidf := []interface{}{}
	data := map[string]interface{}{
		"code": 0,
		"data": json,
		"type": Type,
	}
	switch Type {
	case "system":
		for _, uid := range to_users {

		}
		break

	default:
		break
	}
}
