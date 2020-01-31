package socket

import "github.com/gorilla/websocket"

type Base_C struct {
	Cmd  string `json:"cmd"`
	Idem string `json:"idem"`
}

type Base_R struct {
	Cmd    string    `json:"cmd"`
	Idem   string    `json:"idem"`
	Stamp  string    `json:"stamp"`
	Result string    `json:"result"`
	Exp    Exception `json:"exp"`
}

type Base_B struct {
	Cmd   string `json:"cmd"`
	Stamp string `json:"stamp"`
}

type Exception struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type User struct {
	Id       string `json:"id"`
	Nickname string `json:"nickName"`
	Icon     string `json:"icon"`
	Role     string `json:"role"`
	Status   string `json:"status"`
}

type Chat struct {
	Message  string   `json:"message"`
	Style    string   `json:"style"`
	Roominfo RoomInfo `json:"roomInfo"`
}

type Chatmessage struct {
	Historyuid string   `json:"historyUid"`
	From       User     `json:"from"`
	Stamp      string   `json:"stamp"`
	Message    string   `json:"message"`
	Style      string   `json:"style"`
	Roominfo   RoomInfo `json:"roomInfo"`
}

type Room struct {
	Users        []User        `json:"users"`
	Chatmessages []Chatmessage `json:"chatMessages"`
}

type LoginInfo struct {
	Roomname string `json:"roomName"`
	Nickname string `json:"nickName"`
}

type RoomInfo struct {
	Roomname string `json:"roomName"`
}

type Chatmessagehistory struct {
	Historyuid string `json:"historyUid"`
	Roomname   string `json:"roomName"`
	Userid     string `json:"userId"`
	Nickname   string `json:"nickName"`
	Icon       string `json:"icon"`
	Stamp      string `json:"stamp"`
	Message    string `json:"message"`
	Style      string `json:"style"`
}

type Client struct {
	Room map[string]RoomInfo
	// The websocket connection.
	Conn *websocket.Conn
	User User
}
