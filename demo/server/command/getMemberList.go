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

func Getmemberlist(connect *websocket.Conn, mtype int, msg []byte) error {

	timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	log.Printf("timeUnix : %s\n", timeUnix)

	var packetRoomList socket.Cmd_c_get_member_list

	if err := json.Unmarshal([]byte(msg), &packetRoomList); err != nil {
		panic(err)
	}

	memberList := make([]socket.User, 0, len(common.Rooms[packetRoomList.Payload.Roomname]))
	for _, v := range common.Rooms[packetRoomList.Payload.Roomname] {
		memberList = append(memberList, v.User)
	}

	SendMemberList := socket.Cmd_r_get_member_list_struct{Base_R: socket.Base_R{Cmd: socket.CMD_R_GET_MEMBER_LIST, Idem: packetRoomList.Idem, Stamp: timeUnix, Result: "ok", Exp: socket.Exception{Code: "", Message: ""}}, Payload: memberList}
	SendMemberListJson, _ := json.Marshal(SendMemberList)

	connect.WriteMessage(mtype, SendMemberListJson)

	return nil
}
