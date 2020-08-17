package ws

import (
	"fmt"
	"github.com/gorilla/websocket"
	"main.go/config"
	"main.go/tuuz/Calc"
	"main.go/tuuz/Jsong"
	"main.go/tuuz/Net"
)

var User2Conn map[string]*websocket.Conn
var Conn2User map[*websocket.Conn]string
var Room map[string]int

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
	delete(User2Conn, Conn2User[conn])
	delete(Conn2User, conn)
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
		auth_init(conn, data)
		break

	case "join_room", "JOIN_ROOM":
		join_room(conn, data)
		break

	case "exit_room", "EXIT_ROOM":
		exit_room(conn, data)
		break

	case "msg_list", "MSG_LIST":
		msg_list(conn, data)
		break

	case "private_msg", "PRIVATE_MSG":
		private_msg(conn, data)
		break

	case "group_msg":
		group_msg(conn, data)
		break

	case "requst_count":
		requst_count(conn, data)
		break

	case "ping":
		ping(conn, data)
		break

	case "api":
		api(conn, data)
		break

	case "clear_private_unread":
		clear_private_unread(conn, data)
		break

	case "clear_group_unread":
		clear_group_unread(conn, data)
		break

	default:
		break
	}
}

func auth_init(conn *websocket.Conn, data map[string]interface{}) {
	uid := Calc.Any2String(data["uid"])
	token := Calc.Any2String(data["token"])
	if uid == "" || token == "" {
		On_close(conn)
		fmt.Println("uid_not_exists,UID-token不存在")
	}
	ret, err := Net.Post(config.CHAT_URL+config.AuthURL, nil, map[string]interface{}{
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
		rtt, err := Jsong.JObject(ret)
		if err != nil {
			res := map[string]interface{}{
				"code": 400,
				"data": "网络错误请重试",
				"type": data["type"],
			}
			conn.WriteJSON(res)
		} else {
			if rtt["code"] == 0 {
				User2Conn[uid] = conn
				Conn2User[conn] = uid
				Room[uid] = 0
				message := "欢迎" + uid + "连入聊天服务器"
				if config.DEBUG {
					fmt.Println(message)
				}
				res := map[string]interface{}{
					"code": 0,
					"data": "初始化完成",
					"type": data["type"],
				}
				conn.WriteJSON(res)
			} else {
				res := map[string]interface{}{
					"code": -1,
					"data": "未登录",
					"type": data["type"],
				}
				conn.WriteJSON(res)
			}
		}
	}
}

func join_room(conn *websocket.Conn, data map[string]interface{}) {
	if Conn2User[conn] != "" {
		if data["chat_type"] == "private" {
			res := map[string]interface{}{
				"code": 0,
				"data": "已经加入和" + Calc.Any2String(data["id"]),
				"type": data["type"],
			}
			conn.WriteJSON(res)
		} else if data["chat_type"] == "group" {
			res := map[string]interface{}{
				"code": 0,
				"data": "已经加入和" + Calc.Any2String(data["id"]),
				"type": data["type"],
			}
			conn.WriteJSON(res)
		} else {
			res := map[string]interface{}{
				"code": 400,
				"data": "类型不存在",
				"type": data["type"],
			}
			conn.WriteJSON(res)
		}
	} else {
		conn.WriteJSON(map[string]interface{}{"code": -1, "data": "Auth_Fail", "type": data["type"]})
	}
}

func exit_room(conn *websocket.Conn, data map[string]interface{}) {
	if Conn2User[conn] != "" {
		Room[Conn2User[conn]] = 0
		res := map[string]interface{}{
			"code": 0,
			"data": "退出至大厅",
			"type": data["type"],
		}
		conn.WriteJSON(res)
	} else {
		conn.WriteJSON(map[string]interface{}{"code": -1, "data": "Auth_Fail", "type": data["type"]})
	}
}

func msg_list(conn *websocket.Conn, data map[string]interface{}) {
	if Conn2User[conn] != "" {
		ret, err := Net.Post(config.CHAT_URL+config.Msg_list, nil, map[string]interface{}{
			"uid": Conn2User[conn],
			"ip":  conn.RemoteAddr(),
		}, nil, nil)
		if config.DEBUG_REMOTE_RET {
			fmt.Println("DEBUG_REMOTE_RET", ret, err)
		}
		if err != nil {
			res := map[string]interface{}{
				"code": 400,
				"data": "网络错误请重试",
				"type": data["type"],
			}
			conn.WriteJSON(res)
		} else {
			rtt, err := Jsong.JObject(ret)
			if err != nil {
				res := map[string]interface{}{
					"code": 404,
					"data": "消息列表数据不完整",
					"type": data["type"],
				}
				conn.WriteJSON(res)
			} else {
				if rtt["code"] == 0 {
					res := map[string]interface{}{
						"code": 0,
						"data": rtt["data"],
						"type": data["type"],
					}
					conn.WriteJSON(res)
				} else {
					res := map[string]interface{}{
						"code": -1,
						"data": "未登录",
						"type": data["type"],
					}
					conn.WriteJSON(res)
				}
			}
		}
	} else {
		conn.WriteJSON(map[string]interface{}{"code": -1, "data": "Auth_Fail", "type": data["type"]})
	}

}

func private_msg(conn *websocket.Conn, data map[string]interface{}) {
	if Conn2User[conn] != "" {
		ret, err := Net.Post(config.CHAT_URL+config.Private_msg, nil, map[string]interface{}{
			"uid": Conn2User[conn],
			"fid": data["uid"],
			"ip":  conn.RemoteAddr(),
		}, nil, nil)
		if config.DEBUG_REMOTE_RET {
			fmt.Println("DEBUG_REMOTE_RET", ret, err)
		}
		if err != nil {
			res := map[string]interface{}{
				"code": 400,
				"data": "网络错误请重试",
				"type": data["type"],
			}
			conn.WriteJSON(res)
		} else {
			rtt, err := Jsong.JObject(ret)
			if err != nil {
				res := map[string]interface{}{
					"code": 404,
					"data": "消息列表数据不完整",
					"type": data["type"],
				}
				conn.WriteJSON(res)
			} else {
				if rtt["code"] == 0 {
					res := map[string]interface{}{
						"code": 0,
						"data": rtt["data"],
						"type": data["type"],
					}
					conn.WriteJSON(res)
				}
			}
		}
	} else {
		conn.WriteJSON(map[string]interface{}{"code": -1, "data": "Auth_Fail", "type": data["type"]})
	}
}

func group_msg(conn *websocket.Conn, data map[string]interface{}) {
	if Conn2User[conn] != "" {
		ret, err := Net.Post(config.CHAT_URL+config.Group_msg, nil, map[string]interface{}{
			"uid": Conn2User[conn],
			"fid": data["uid"],
			"ip":  conn.RemoteAddr(),
		}, nil, nil)
		if config.DEBUG_REMOTE_RET {
			fmt.Println("DEBUG_REMOTE_RET", ret, err)
		}
		if err != nil {
			res := map[string]interface{}{
				"code": 400,
				"data": "网络错误请重试",
				"type": data["type"],
			}
			conn.WriteJSON(res)
		} else {
			rtt, err := Jsong.JObject(ret)
			if err != nil {
				res := map[string]interface{}{
					"code": 404,
					"data": "消息列表数据不完整",
					"type": data["type"],
				}
				conn.WriteJSON(res)
			} else {
				if rtt["code"] == 0 {
					res := map[string]interface{}{
						"code": 0,
						"data": rtt["data"],
						"type": data["type"],
					}
					conn.WriteJSON(res)
				}
			}
		}
	} else {
		conn.WriteJSON(map[string]interface{}{"code": -1, "data": "Auth_Fail", "type": data["type"]})
	}
}

func requst_count(conn *websocket.Conn, data map[string]interface{}) {
	if Conn2User[conn] != "" {
		ret, err := Net.Post(config.CHAT_URL+config.Request_count, nil, map[string]interface{}{
			"uid": Conn2User[conn],
			"ip":  conn.RemoteAddr(),
		}, nil, nil)
		if config.DEBUG_REMOTE_RET {
			fmt.Println("DEBUG_REMOTE_RET", ret, err)
		}
		if err != nil {
			res := map[string]interface{}{
				"code": 400,
				"data": "网络错误请重试",
				"type": data["type"],
			}
			conn.WriteJSON(res)
		} else {
			rtt, err := Jsong.JObject(ret)
			if err != nil {
				res := map[string]interface{}{
					"code": 404,
					"data": "Req_Count不完整",
					"type": data["type"],
				}
				conn.WriteJSON(res)
			} else {
				if rtt["code"] == 0 {
					res := map[string]interface{}{
						"code": 0,
						"data": rtt["data"],
						"type": data["type"],
					}
					conn.WriteJSON(res)
				}
			}
		}
	} else {
		conn.WriteJSON(map[string]interface{}{"code": -1, "data": "Auth_Fail", "type": data["type"]})
	}
}

func ping(conn *websocket.Conn, data map[string]interface{}) {
	res := map[string]interface{}{
		"code": 0,
		"data": "PONG",
		"type": data["type"],
	}
	conn.WriteJSON(res)
}

func api(conn *websocket.Conn, data map[string]interface{}) {
	if Conn2User[conn] != "" {
		function := Calc.Any2String(data["func"])
		ret, err := Net.Post(config.CHAT_URL+config.ManualAPI+function, nil, map[string]interface{}{
			"uid": Conn2User[conn],
			"ip":  conn.RemoteAddr(),
		}, nil, nil)
		if config.DEBUG_REMOTE_RET {
			fmt.Println("DEBUG_REMOTE_RET", ret, err)
		}
		if err != nil {
			res := map[string]interface{}{
				"code": 400,
				"data": "网络错误请重试",
				"type": data["type"],
			}
			conn.WriteJSON(res)
		} else {
			rtt, err := Jsong.JObject(ret)
			if err != nil {
				res := map[string]interface{}{
					"code": 404,
					"data": "API不完整",
					"type": data["type"],
				}
				conn.WriteJSON(res)
			} else {
				if rtt["code"] == 0 {
					res := map[string]interface{}{
						"code": 0,
						"data": rtt["data"],
						"type": data["type"],
					}
					conn.WriteJSON(res)
				}
			}
		}
	} else {
		conn.WriteJSON(map[string]interface{}{"code": -1, "data": "Auth_Fail", "type": data["type"]})
	}
}

func clear_private_unread(conn *websocket.Conn, data map[string]interface{}) {
	if Conn2User[conn] != "" {
		if data["id"] == nil {
			res := map[string]interface{}{
				"code": 400,
				"data": "id",
				"type": data["type"],
			}
			conn.WriteJSON(res)
			return
		}
		ret, err := Net.Post(config.CHAT_URL+config.Clear_private_unread, nil, map[string]interface{}{
			"uid": Conn2User[conn],
			"fid": data["id"],
			"ip":  conn.RemoteAddr(),
		}, nil, nil)
		if config.DEBUG_REMOTE_RET {
			fmt.Println("DEBUG_REMOTE_RET", ret, err)
		}
		if err != nil {
			res := map[string]interface{}{
				"code": 400,
				"data": "网络错误请重试",
				"type": data["type"],
			}
			conn.WriteJSON(res)
		} else {
			rtt, err := Jsong.JObject(ret)
			if err != nil {
				res := map[string]interface{}{
					"code": 404,
					"data": "API不完整",
					"type": data["type"],
				}
				conn.WriteJSON(res)
			} else {
				if rtt["code"] == 0 {
					res := map[string]interface{}{
						"code": 0,
						"data": rtt["data"],
						"type": data["type"],
					}
					conn.WriteJSON(res)
				}
			}
		}
	} else {
		conn.WriteJSON(map[string]interface{}{"code": -1, "data": "Auth_Fail", "type": data["type"]})
	}
}

func clear_group_unread(conn *websocket.Conn, data map[string]interface{}) {
	if Conn2User[conn] != "" {
		if data["id"] == nil {
			res := map[string]interface{}{
				"code": 400,
				"data": "id",
				"type": data["type"],
			}
			conn.WriteJSON(res)
			return
		}
		ret, err := Net.Post(config.CHAT_URL+config.Clear_group_unread, nil, map[string]interface{}{
			"uid": Conn2User[conn],
			"gid": data["id"],
			"ip":  conn.RemoteAddr(),
		}, nil, nil)
		if config.DEBUG_REMOTE_RET {
			fmt.Println("DEBUG_REMOTE_RET", ret, err)
		}
		if err != nil {
			res := map[string]interface{}{
				"code": 400,
				"data": "网络错误请重试",
				"type": data["type"],
			}
			conn.WriteJSON(res)
		} else {
			rtt, err := Jsong.JObject(ret)
			if err != nil {
				res := map[string]interface{}{
					"code": 404,
					"data": "API不完整",
					"type": data["type"],
				}
				conn.WriteJSON(res)
			} else {
				if rtt["code"] == 0 {
					res := map[string]interface{}{
						"code": 0,
						"data": rtt["data"],
						"type": data["type"],
					}
					conn.WriteJSON(res)
				}
			}
		}
	} else {
		conn.WriteJSON(map[string]interface{}{"code": -1, "data": "Auth_Fail", "type": data["type"]})
	}
}
