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

func Getchathistory(connect *websocket.Conn, mtype int, msg []byte) error {

	timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	log.Printf("timeUnix : %s\n", timeUnix)

	var packetChatHistory socket.Cmd_c_get_chat_history

	if err := json.Unmarshal([]byte(msg), &packetChatHistory); err != nil {
		panic(err)
	}

	log.Printf("packetChatHistory : %+v\n", packetChatHistory)

	if packetChatHistory.Payload.Historyuid == "" {
		packetChatHistory.Payload.Historyuid = common.GetId().String()
	}

	rangeStamp := common.Timetoshift(common.Oncechathistorylong)
	endStamp, _ := strconv.ParseInt(packetChatHistory.Payload.Historyuid, 10, 64)
	startStamp := endStamp - rangeStamp
	historyUid := strconv.FormatInt(startStamp, 10)

	SendChatHistoryList := socket.Cmd_r_get_chat_history{Base_R: socket.Base_R{Cmd: socket.CMD_R_GET_CHAT_HISTORY, Idem: packetChatHistory.Idem, Stamp: timeUnix, Result: "ok", Exp: socket.Exception{Code: "", Message: ""}}}

	SendChatHistoryList.Payload.Historyuid = historyUid
	SendChatHistoryList.Payload.Message = make([]socket.Chatmessage, 0)
	SendChatHistoryListJson, _ := json.Marshal(SendChatHistoryList)
	connect.WriteMessage(mtype, SendChatHistoryListJson)

	return nil
}
