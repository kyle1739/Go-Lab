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

func Playersendmsg(connect *websocket.Conn, mtype int, msg []byte) error {

	timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	log.Printf("timeUnix : %s\n", timeUnix)

	var packetSendMsg socket.Cmd_c_player_send_msg_struct

	if err := json.Unmarshal([]byte(msg), &packetSendMsg); err != nil {
		panic(err)
		return err
	}

	if !common.Checkinroom(packetSendMsg.Payload.Roominfo.Roomname, common.Clients[connect]) {
		SendMsg := socket.Cmd_r_player_send_msg_struct{Base_R: socket.Base_R{Cmd: socket.CMD_R_PLAYER_SEND_MSG, Idem: packetSendMsg.Idem, Stamp: timeUnix, Result: "err", Exp: socket.Exception{Code: "101", Message: "NOT_IN_ROOM"}}}
		SendMsgJson, _ := json.Marshal(SendMsg)
		connect.WriteMessage(mtype, SendMsgJson)
		return nil
	}
	log.Printf("packetSendMsg : %+v\n", packetSendMsg)

	SendMsg := socket.Cmd_r_player_send_msg_struct{Base_R: socket.Base_R{Cmd: socket.CMD_R_PLAYER_SEND_MSG, Idem: packetSendMsg.Idem, Stamp: timeUnix, Result: "ok", Exp: socket.Exception{Code: "", Message: ""}}}
	SendMsgJson, _ := json.Marshal(SendMsg)

	connect.WriteMessage(mtype, SendMsgJson)
	historyUid := common.GetId().String()
	roominfo := socket.RoomInfo{Roomname: packetSendMsg.Payload.Roominfo.Roomname}
	chatMessage := socket.Chatmessage{Historyuid: historyUid, From: common.Clients[connect].User, Stamp: timeUnix, Message: packetSendMsg.Payload.Message, Style: packetSendMsg.Payload.Style, Roominfo: roominfo}
	sendMsgBrocast := socket.Cmd_b_player_speak_struct{Base_B: socket.Base_B{Cmd: socket.CMD_B_PLAYER_SPEAK, Stamp: timeUnix}, Payload: chatMessage}
	sendMsgBrocastJson, _ := json.Marshal(sendMsgBrocast)

	common.Broadcast(packetSendMsg.Payload.Roominfo.Roomname, mtype, sendMsgBrocastJson)

	return nil
}
