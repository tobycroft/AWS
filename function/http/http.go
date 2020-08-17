package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"main.go/config"
	"main.go/function/ws"
	"main.go/tuuz/Calc"
	"main.go/tuuz/Input"
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
	dest, ok := Input.Post("dest", c, false)
	if !ok {
		return
	}
	data, ok := Input.Post("data", c, false)
	if !ok {
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
			"data": "data需要提交",
		})
		c.Abort()
		return
	}
	json_handler(c, json, to_users, dest, Type)
}

func json_handler(c *gin.Context, json map[string]interface{}, to_users []interface{}, dest string, Type string) {
	fmt.Println("json_handler", json)
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
			conn := ws.User2Conn[Calc.Any2String(uid)]
			if conn != nil {
				uids = append(uids, uid)
				conn.WriteJSON(data)
			}
			uidf = append(uidf, uid)
		}
		break

	case "refresh_list":
		for _, uid := range to_users {
			conn := ws.User2Conn[Calc.Any2String(uid)]
			if conn != nil {
				uids = append(uids, uid)
				conn.WriteJSON(data)
			}
			uidf = append(uidf, uid)
		}
		break

	case "private_chat":
		for _, uid := range to_users {
			if ws.Room[Calc.Any2String(uid)] == dest {
				conn := ws.User2Conn[Calc.Any2String(uid)]
				if conn != nil {
					uids = append(uids, uid)
					conn.WriteJSON(data)
				}
			}
			uidf = append(uidf, uid)
		}
		break

	case "group_chat":
		for _, uid := range to_users {
			if ws.Room[Calc.Any2String(uid)] == dest {
				conn := ws.User2Conn[Calc.Any2String(uid)]
				if conn != nil {
					conn.WriteJSON(data)
				}
			}
			uidf = append(uidf, uid)
		}
		break

	case "request_count":
		for _, uid := range to_users {
			if ws.Room[Calc.Any2String(uid)] == 0 {
				conn := ws.User2Conn[Calc.Any2String(uid)]
				if conn != nil {
					uids = append(uids, uid)
					conn.WriteJSON(data)
				}
			}
			uidf = append(uidf, uid)
		}
		break

	case "push":
		for _, uid := range to_users {
			conn := ws.User2Conn[Calc.Any2String(uid)]
			if conn != nil {
				uids = append(uids, uid)
				conn.WriteJSON(data)
			}
			uidf = append(uidf, uid)
		}
		break

	case "message":
		for _, uid := range to_users {
			conn := ws.User2Conn[Calc.Any2String(uid)]
			if conn != nil {
				uids = append(uids, uid)
				conn.WriteJSON(data)
			}
		}
		break

	case "indoor_message":
		for _, uid := range to_users {
			if ws.Room[Calc.Any2String(uid)] == 0 {
				conn := ws.User2Conn[Calc.Any2String(uid)]
				if conn != nil {
					uids = append(uids, uid)
					conn.WriteJSON(data)
				}
			}
		}
		break

	case "outer_all":
		for _, uid := range to_users {
			if ws.Room[Calc.Any2String(uid)] != 0 {
				conn := ws.User2Conn[Calc.Any2String(uid)]
				if conn != nil {
					uids = append(uids, uid)
					conn.WriteJSON(data)
				}
			}
		}
		break

	case "user_room":
		user_room := map[string]interface{}{}
		for _, uid := range to_users {
			id := Calc.Any2String(uid)
			user_room[id] = ws.Room[id]
		}
		c.JSON(200, map[string]interface{}{
			"code": 0,
			"data": user_room,
		})
		c.Abort()
		return

	default:
		for _, uid := range to_users {
			if ws.Room[Calc.Any2String(uid)] == dest {
				conn := ws.User2Conn[Calc.Any2String(uid)]
				if conn != nil {
					uids = append(uids, uid)
					conn.WriteJSON(data)
				}
			}
			uidf = append(uidf, uid)
		}
		break
	}

	c.JSON(200, map[string]interface{}{
		"code": 0,
		"data": map[string]interface{}{
			"success": uids,
			"fail":    uidf,
		},
	})
	c.Abort()
	return

}
