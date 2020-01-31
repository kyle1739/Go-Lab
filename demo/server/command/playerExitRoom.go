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

func Playerexitroom(connect *websocket.Conn, mtype int, msg []byte) error {

	timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	log.Printf("timeUnix : %s\n", timeUnix)

	var packetExitRoom socket.Cmd_c_player_exit_room_struct

	if err := json.Unmarshal([]byte(msg), &packetExitRoom); err != nil {
		panic(err)
	}

	if common.Checkinroom(packetExitRoom.Payload.Roomname, common.Clients[connect]) {

		roominfo := socket.RoomInfo{Roomname: packetExitRoom.Payload.Roomname}
		chatMessage := socket.Chatmessage{Historyuid: "sys", From: common.Clients[connect].User, Stamp: timeUnix, Message: "exit room", Style: "1", Roominfo: roominfo}
		exitroomBrocast := socket.Cmd_b_player_exit_room_struct{Base_B: socket.Base_B{Cmd: socket.CMD_B_PLAYER_EXIT_ROOM, Stamp: timeUnix}, Payload: chatMessage}
		exitroomBrocastJson, err := json.Marshal(exitroomBrocast)
		if err != nil {
			log.Println("json err:", err)
			return err
		}

		common.Broadcast(packetExitRoom.Payload.Roomname, mtype, exitroomBrocastJson)

		delete(common.Clients[connect].Room, packetExitRoom.Payload.Roomname)
		delete(common.Rooms[packetExitRoom.Payload.Roomname], connect)
		exitroom := socket.Cmd_r_player_exit_room_struct{Base_R: socket.Base_R{Cmd: socket.CMD_R_PLAYER_EXIT_ROOM, Idem: packetExitRoom.Idem, Stamp: timeUnix, Result: "ok", Exp: socket.Exception{Code: "", Message: ""}}}
		exitroomJson, _ := json.Marshal(exitroom)
		connect.WriteMessage(mtype, exitroomJson)

		if len(common.Rooms[packetExitRoom.Payload.Roomname]) == 0 {
			delete(common.Rooms, packetExitRoom.Payload.Roomname)
			delete(common.RoomsInfo, packetExitRoom.Payload.Roomname)
		}

		log.Printf("roomsInfo : %+v\n", common.RoomsInfo)

		log.Printf("rooms : %+v\n", common.Rooms)

		log.Printf("clients : %+v\n", common.Clients)

	} else {
		exitroom := socket.Cmd_r_player_exit_room_struct{Base_R: socket.Base_R{Cmd: socket.CMD_R_PLAYER_EXIT_ROOM, Idem: packetExitRoom.Idem, Stamp: timeUnix, Result: "err", Exp: socket.Exception{Code: "101", Message: "NOT_IN_ROOM"}}}
		exitroomJson, _ := json.Marshal(exitroom)
		connect.WriteMessage(mtype, exitroomJson)

		return nil
	}

	return nil
}
