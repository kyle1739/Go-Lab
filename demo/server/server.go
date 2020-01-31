package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/olivere/elastic"

	"./command"
	"./common"
	"./socket"
)

var elasticclient *elastic.Client

func main() {

	loadEnv()

	timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)

	log.Printf("main timeUnix : [%s] ", timeUnix)

	// common.Elasticclient, _ = elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(os.Getenv("elasticSearchHost")))

	// Getting the ES version number is quite common, so there's a shortcut
	// esversion, err := common.Elasticclient.ElasticsearchVersion(os.Getenv("elasticSearchHost"))
	// if err != nil {
	// 	// Handle error
	// 	panic(err)
	// }
	// log.Printf("Elasticsearch version %s\n", esversion)

	//command.EsInit(elasticclient, context.Background())

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/js/", fileHandler)
	http.HandleFunc("/css/", fileHandler)
	http.HandleFunc("/echo", echoHandler)

	http.ListenAndServe(os.Getenv("imWebSocketPort"), nil)

}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	upgrader := &websocket.Upgrader{
		//如果有 cross domain 的需求，可加入這個，不檢查 cross domain
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	connect, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}

	defer func() {
		log.Println("disconnect !!")

		for _, roominfo := range common.Clients[connect].Room {
			delete(common.Rooms[roominfo.Roomname], connect)

			timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
			roominfo := socket.RoomInfo{Roomname: roominfo.Roomname}
			chatMessage := socket.Chatmessage{Historyuid: "sys", From: common.Clients[connect].User, Stamp: timeUnix, Message: "exit room", Style: "1", Roominfo: roominfo}
			logoutBrocast := socket.Cmd_b_player_exit_room_struct{Base_B: socket.Base_B{Cmd: socket.CMD_B_PLAYER_EXIT_ROOM, Stamp: timeUnix}, Payload: chatMessage}
			logoutBrocastJson, _ := json.Marshal(logoutBrocast)

			common.Broadcast(roominfo.Roomname, 1, logoutBrocastJson)

			if len(common.Rooms[roominfo.Roomname]) == 0 {
				delete(common.Rooms, roominfo.Roomname)
				delete(common.RoomsInfo, roominfo.Roomname)
			}
		}

		delete(common.Clients, connect)

		connect.Close()
	}()

	for {
		err = receivePacketHandle(connect)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func fileHandler(w http.ResponseWriter, r *http.Request) {
	//log.Println("r.URL.Path :" , r.URL.Path)
	http.ServeFile(w, r, "../client"+r.URL.Path)
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	//log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "../client/index.html")
}

func receivePacketHandle(connect *websocket.Conn) error {

	//log.Printf("connect : %+v\n", connect)

	//ReadMessage只能讀一次 猜測是因為讀取指標的問題
	mtype, msg, err := connect.ReadMessage()
	if err != nil {
		log.Println("write:", err)
		return err
	}

	timeUnix := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)

	log.Printf("timeUnix : [%s] ", timeUnix)

	var mapResult map[string]interface{}
	//使用 json.Unmarshal(data []byte, v interface{})进行转换,返回 error 信息
	if err := json.Unmarshal([]byte(msg), &mapResult); err != nil {
		panic(err)
	}

	log.Printf("mapResult : %+v\n", mapResult)

	switch mapResult["cmd"] {
	case socket.CMD_C_TOKEN_CHANGE:
		err = command.Tokenchange(connect, mtype, msg)

		if err != nil {
			return err
		}
		break
	case socket.CMD_C_PLAYER_LOGOUT:

		err = command.Playerloguot(connect, mtype, msg)

		if err != nil {
			return err
		}

		break
	case socket.CMD_C_PLAYER_EXIT_ROOM:

		err = command.Playerexitroom(connect, mtype, msg)

		if err != nil {
			return err
		}

		break
	case socket.CMD_C_PLAYER_ENTER_ROOM:

		err = command.Playerenterroom(connect, mtype, msg)

		if err != nil {
			return err
		}

		break
	case socket.CMD_C_PLAYER_SEND_MSG:

		err = command.Playersendmsg(connect, mtype, msg)

		if err != nil {
			return err
		}

		break
	case socket.CMD_C_GET_MEMBER_LIST:

		err = command.Getmemberlist(connect, mtype, msg)

		if err != nil {
			return err
		}

		break
	case socket.CMD_C_GET_ROOM_LIST:

		err = command.Getroomlist(connect, mtype, msg)

		if err != nil {
			return err
		}

		break
	case socket.CMD_C_GET_CHAT_HISTORY:

		err = command.Getchathistory(connect, mtype, msg)

		if err != nil {
			return err
		}

		break
	default:

	}
	// Write message back to the client message owner with:
	//connect.WriteMessage(mtype , msg.Body)
	// Write message to all except this client with:
	//broadcast(hubName , mtype , msg)

	return nil
}

func loadEnv() {

	config, err := ioutil.ReadFile("config/host.json")
	if err != nil {
		log.Fatal("找不到host.json")
	}

	configHost := make(map[string]string)

	err = json.Unmarshal(config, &configHost)
	if err != nil {
		return
	}

	for k, v := range configHost {
		log.Printf("%s : %s\n", k, v)
		_ = os.Setenv(k, v)
	}

}
