package common

import (
	"github.com/gorilla/websocket"
	"github.com/olivere/elastic"

	"../socket"
)

var Elasticclient *elastic.Client

const Maxchathistory int = 10

const Oncechathistorylong int64 = 1000 * 60 * 60 * 24 * 30

var Rooms = make(map[string]map[*websocket.Conn]socket.Client)
var RoomsInfo = make(map[string]socket.RoomInfo)
var Clients = make(map[*websocket.Conn]socket.Client)
