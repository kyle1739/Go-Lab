package common

import "../socket"

func Broadcast(roomname string, mtype int, msg []byte) {

	//log.Printf("roomName : %s\n", roomName)
	targetroom := Rooms[roomname]
	//log.Printf("targetroom : %+v\n", targetroom)

	for clientConnect := range targetroom {
		clientConnect.WriteMessage(mtype, msg)
	}
}

func Checkinroom(roomname string, client socket.Client) bool {
	for _, clientRoom := range client.Room {
		if roomname == clientRoom.Roomname {
			return true
		}
	}
	return false
}
