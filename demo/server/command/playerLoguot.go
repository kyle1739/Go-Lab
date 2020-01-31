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

func Playerloguot(connect *websocket.Conn, mtype int, msg []byte) error {

	timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	log.Printf("timeUnix : %s\n", timeUnix)

	var packetLogout socket.Cmd_c_player_logout_struct

	if err := json.Unmarshal([]byte(msg), &packetLogout); err != nil {
		panic(err)
	}

	logout := socket.Cmd_r_player_logout_struct{Base_R: socket.Base_R{Cmd: socket.CMD_R_PLAYER_LOGOUT, Idem: packetLogout.Idem, Stamp: timeUnix, Result: "ok", Exp: socket.Exception{Code: "", Message: ""}}}
	logoutJson, _ := json.Marshal(logout)

	connect.WriteMessage(mtype, logoutJson)

	userinfo := common.Clients[connect].User

	for _, roominfo := range common.Clients[connect].Room {
		delete(common.Rooms[roominfo.Roomname], connect)
		if len(common.Rooms[roominfo.Roomname]) == 0 {
			delete(common.Rooms, roominfo.Roomname)
			delete(common.RoomsInfo, roominfo.Roomname)
		}
	}

	delete(common.Clients, connect)

	connect.Close()

	roominfo := socket.RoomInfo{Roomname: packetLogout.Payload.Roomname}
	chatMessage := socket.Chatmessage{Historyuid: "sys", From: userinfo, Stamp: timeUnix, Message: "exit room", Style: "1", Roominfo: roominfo}
	logoutBrocast := socket.Cmd_b_player_exit_room_struct{Base_B: socket.Base_B{Cmd: socket.CMD_B_PLAYER_EXIT_ROOM, Stamp: timeUnix}, Payload: chatMessage}
	logoutBrocastJson, _ := json.Marshal(logoutBrocast)

	common.Broadcast(packetLogout.Payload.Roomname, mtype, logoutBrocastJson)

	return nil
}
