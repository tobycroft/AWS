package ws

import (
	"fmt"
	"github.com/gorilla/websocket"
	"main.go/config"
	"main.go/tuuz/Calc"
	"main.go/tuuz/Jsong"
	"main.go/tuuz/Net"
)

var User2Ws map[string]*websocket.Conn
var Ws2User map[*websocket.Conn]string

func On_connect(conn *websocket.Conn) {
	//err := conn.WriteMessage(1, []byte("连入成功"))
	str := map[string]interface{}{
		"data": "连入成功",
	}
	err := conn.WriteJSON(str)

	if err != nil {
		fmt.Printf("write fail = %v\n", err)
		return
	}
}

func On_close(conn *websocket.Conn) {
	// 发送 websocket 结束包
	conn.WriteMessage(websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	// 真正关闭 conn
	conn.Close()
}

func Handler(json_str string, conn *websocket.Conn) {
	fmt.Println(json_str)
	json, jerr := Jsong.JObject(json_str)
	if jerr != nil {
		fmt.Println("jsonerr", jerr)
		return
	}
	if config.DEBUG_WS_REQ {
		fmt.Println("DEBUG_WS_REQ", json_str)
	}
	if config.DEBUG_WS_REQ {
		fmt.Println("DEBUG_WS_REQ:type", json["type"])
	}
	data, derr := Jsong.ParseObject(json["data"])
	if derr != nil {
		fmt.Println("jsondataerr", derr)
		return
	}
	switch Calc.Any2String(json["type"]) {
	case "init", "INIT":
		uid := data["uid"]
		token := data["token"]
		if uid == nil || token == nil {
			On_close(conn)
			fmt.Println("uid_not_exists,UID-token不存在")
		}
		ret, err := Net.Post("/api/auth/userauth", nil, map[string]interface{}{
			"uid":   uid,
			"token": token,
			"type":  1,
			"ip":    conn.RemoteAddr(),
		}, nil, nil)
		if config.DEBUG_AUTH {
			fmt.Println("DEBUG_AUTH", ret, err)
		}
		if err != nil {
			res := map[string]interface{}{
				"code": 400,
				"data": "网络错误请重试",
				"type": data["type"],
			}
			conn.WriteJSON(res)
		} else {

		}
		break

	default:
		break
	}
}

func retjson(json map[string]interface{}) string {
	ret, err := Jsong.Encode(json)
	if err != nil {
		return ""
	} else {
		return ret
	}
}
