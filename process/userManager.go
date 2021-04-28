package process

import (
	"net"
)

type UserProcess struct {
	Uid  int
	Conn net.Conn
}

var onlineUsers = make(map[int]*UserProcess)

func GetOnlineUsers() map[int]*UserProcess {
	return onlineUsers
}

func Push(process *UserProcess) {
	onlineUsers[process.Uid] = process
}

func Get(id int) *UserProcess {
	u, ok := onlineUsers[id]
	if ok {
		return u
	}
	return nil
}

func Del(id int) {
	delete(onlineUsers, id)
}

func Disconnect(conn net.Conn) {
	for u, process := range GetOnlineUsers() {
		if conn.RemoteAddr() == process.Conn.RemoteAddr() {
			delete(onlineUsers, u)
			break
		}
	}
}
