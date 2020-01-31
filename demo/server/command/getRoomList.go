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

func Getroomlist(connect *websocket.Conn, mtype int, msg []byte) error {

	timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	log.Printf("timeUnix : %s\n", timeUnix)

	var packetRoomList socket.Cmd_c_get_room_list

	if err := json.Unmarshal([]byte(msg), &packetRoomList); err != nil {
		panic(err)
	}

	roomList := make([]socket.RoomInfo, 0, len(common.RoomsInfo))
	for _, roominfo := range common.RoomsInfo {
		roomList = append(roomList, roominfo)
	}
	// log.Printf("roomList : %+v\n", roomList)

	SendRoomList := socket.Cmd_r_get_room_list_struct{Base_R: socket.Base_R{Cmd: socket.CMD_R_GET_ROOM_LIST, Idem: packetRoomList.Idem, Stamp: timeUnix, Result: "ok", Exp: socket.Exception{Code: "", Message: ""}}, Payload: roomList}
	SendRoomListJson, _ := json.Marshal(SendRoomList)
	connect.WriteMessage(mtype, SendRoomListJson)

	return nil
}
