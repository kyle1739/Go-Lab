package command

import (
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/gorilla/websocket"

	"../common"
	"../socket"
)

func Playerenterroom(connect *websocket.Conn, mtype int, msg []byte) error {

	timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	log.Printf("timeUnix : %s\n", timeUnix)

	var packetEnterRoom socket.Cmd_c_player_enter_room_struct

	if err := json.Unmarshal([]byte(msg), &packetEnterRoom); err != nil {
		panic(err)
		return err
	}

	if common.Checkinroom(packetEnterRoom.Payload.Roomname, common.Clients[connect]) {
		enterroom := socket.Cmd_r_player_enter_room_struct{Base_R: socket.Base_R{Cmd: socket.CMD_R_PLAYER_ENTER_ROOM, Idem: packetEnterRoom.Idem, Stamp: timeUnix, Result: "err", Exp: socket.Exception{Code: "101", Message: "IN_ROOM"}}}
		enterroomJson, _ := json.Marshal(enterroom)
		connect.WriteMessage(mtype, enterroomJson)
		return nil
	}
	client := common.Clients[connect]
	client.Room[packetEnterRoom.Payload.Roomname] = socket.RoomInfo{Roomname: packetEnterRoom.Payload.Roomname}
	if common.Rooms[packetEnterRoom.Payload.Roomname] == nil {
		var adduser = make(map[*websocket.Conn]socket.Client)
		adduser[connect] = client
		common.Rooms[packetEnterRoom.Payload.Roomname] = adduser
		common.RoomsInfo[packetEnterRoom.Payload.Roomname] = socket.RoomInfo{Roomname: packetEnterRoom.Payload.Roomname}
	} else {
		common.Rooms[packetEnterRoom.Payload.Roomname][connect] = client
	}

	enterroom := socket.Cmd_r_player_enter_room_struct{Base_R: socket.Base_R{Cmd: socket.CMD_R_PLAYER_ENTER_ROOM, Idem: packetEnterRoom.Idem, Stamp: timeUnix, Result: "ok", Exp: socket.Exception{Code: "", Message: ""}}}
	enterroomJson, _ := json.Marshal(enterroom)
	connect.WriteMessage(mtype, enterroomJson)

	roominfo := socket.RoomInfo{Roomname: packetEnterRoom.Payload.Roomname}
	chatMessage := socket.Chatmessage{Historyuid: "sys", From: common.Clients[connect].User, Stamp: timeUnix, Message: "enter room", Style: "1", Roominfo: roominfo}
	enterroomBrocast := socket.Cmd_b_player_enter_room_struct{Base_B: socket.Base_B{Cmd: socket.CMD_B_PLAYER_ENTER_ROOM, Stamp: timeUnix}, Payload: chatMessage}
	enterroomBrocastJson, _ := json.Marshal(enterroomBrocast)

	common.Broadcast(packetEnterRoom.Payload.Roomname, mtype, enterroomBrocastJson)

	return nil
}
