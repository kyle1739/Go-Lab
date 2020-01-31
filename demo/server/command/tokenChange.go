package command

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/websocket"

	"../common"
	"../socket"
)

func Tokenchange(connect *websocket.Conn, mtype int, msg []byte) error {

	timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	log.Printf("timeUnix : %s\n", timeUnix)

	var packetToken socket.Cmd_c_token_change_struct

	if err := json.Unmarshal([]byte(msg), &packetToken); err != nil {
		panic(err)
	}

	log.Printf("packetToken : %+v\n", packetToken)

	roomMap := make(map[string]socket.RoomInfo)
	roomMap[packetToken.Payload.Roomname] = socket.RoomInfo{Roomname: packetToken.Payload.Roomname}
	userinfo := socket.User{Id: common.GetId().String(), Nickname: packetToken.Payload.Nickname, Icon: os.Getenv("cdnHost") + "/static/hdp" + strconv.Itoa(rand.Intn(20)) + ".png", Role: "", Status: ""}
	client := socket.Client{Room: roomMap, Conn: connect, User: userinfo}
	// log.Printf("client : %+v\n", client)
	common.Clients[connect] = client
	//log.Printf("clients : %+v\n", clients)
	if common.Rooms[packetToken.Payload.Roomname] == nil {
		var adduser = make(map[*websocket.Conn]socket.Client)
		adduser[connect] = client
		common.Rooms[packetToken.Payload.Roomname] = adduser
		common.RoomsInfo[packetToken.Payload.Roomname] = roomMap[packetToken.Payload.Roomname]
	} else {
		common.Rooms[packetToken.Payload.Roomname][connect] = client
	}
	// log.Printf("rooms : %+v\n", rooms)

	tokenChange := socket.Cmd_r_token_change_struct{Base_R: socket.Base_R{Cmd: socket.CMD_R_TOKEN_CHANGE, Idem: packetToken.Idem, Stamp: timeUnix, Result: "ok", Exp: socket.Exception{Code: "", Message: ""}}}
	tokenChange.Payload = userinfo
	tokenChangeJson, _ := json.Marshal(tokenChange)

	log.Printf("WriteMessage : %+v\n", tokenChangeJson)
	connect.WriteMessage(mtype, tokenChangeJson)

	roominfo := roomMap[packetToken.Payload.Roomname]
	chatMessage := socket.Chatmessage{Historyuid: "sys", From: common.Clients[connect].User, Stamp: timeUnix, Message: "enter room", Style: "1", Roominfo: roominfo}
	logoutBrocast := socket.Cmd_b_player_enter_room_struct{Base_B: socket.Base_B{Cmd: socket.CMD_B_PLAYER_ENTER_ROOM, Stamp: timeUnix}, Payload: chatMessage}
	logoutBrocastJson, err := json.Marshal(logoutBrocast)
	if err != nil {
		log.Println("json err:", err)
	}

	common.Broadcast(packetToken.Payload.Roomname, mtype, logoutBrocastJson)

	return nil
}
