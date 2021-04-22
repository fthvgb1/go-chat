package lib

import (
	"chat/message"
	"chat/process"
	"chat/user"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
)

type Server struct {
	Id       int
	Users    []user.User
	Messages message.Message
}

const SavePath = "./data.json"

var Rchan = make(chan process.Ms, 100000)

func NewServer() *Server {
	data, err := ioutil.ReadFile(SavePath)
	if err == nil {
		s := &Server{}
		x := json.Unmarshal(data, s)
		if x != nil {
			return &Server{}
		}
		return s
	}
	return &Server{}
}

func (server *Server) Login(uid int, pw string) error {
	for _, u := range server.Users {
		if uid == u.Id && pw == u.Password {
			return nil
		}
	}
	return errors.New("not found this u")
}

func (server *Server) Store() {
	data, err := json.Marshal(server)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = ioutil.WriteFile(SavePath, data, 0755)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (server *Server) AddUser(name, pw string, sex int8) {
	server.Id++
	u := user.User{
		server.Id, name, sex, pw,
	}
	server.Users = append(server.Users, u)
}

func (server *Server) Process(conn net.Conn) {
	m := make([]byte, 0)
	read, err := conn.Read(m)
	if err != nil {
		return
	}
	fmt.Println(read)
}

func (server *Server) do(s *message.LoginS, conn net.Conn) (string, error) {
	err := server.Login(s.Uid, s.Pw)
	re := message.Message{
		Type: "login_response",
		Data: nil,
	}
	re.Data = message.Correspond{
		Code:  1,
		Msg:   "success",
		Error: "",
	}
	if err != nil {
		re.Data = message.Correspond{
			Code:  0,
			Msg:   "login fail",
			Error: "pw fail or user not exist",
		}
	}

	err = process.WriteConn(conn, re)
	return "login_send", err
}

func (s *Server) processConn() {
	for {
		select {
		case c := <-Rchan:
			switch c.Msg.Type {
			case "login_send":
				l := c.Msg.Data.(*message.LoginS)
				_, err := s.do(l, c.Conn)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}
}

func (server *Server) read(conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println(err, "ssssssssssssssss")
		}
	}(conn)

	process.Read(conn, Rchan)

	/*for {
	var info = make([]byte, 65536)
	i, err := conn.Read(info)
	if err != nil {
		fmt.Println(err, "errrrrrrrrrrrrr")
		return
	}
	msg := message.Message{}
	err = json.Unmarshal(info[:i], &msg)
	if err != nil {
		fmt.Println(err)
	}
	t := message.MsgType[msg.Type]
	mm := message.Message{}
	mm.Data = reflect.New(t).Interface()
	err = json.Unmarshal(info[:i], &mm)
	if err != nil {
		fmt.Println(err)
	}
	switch msg.Type {
	case "login_send":
		if d,ok:=mm.Data.(*message.LoginS);ok{
			_, err := server.do(d,conn)
			if err != nil {
				fmt.Println(err)
				return
			}
		}

	}

	/*switch reflect.TypeOf(mm.Data){
	case reflect.TypeOf(&LoginS{}):
		x :=mm.Data.(*LoginS)
		e := s.Login(x.Uid,x.Pw)
		re := Message{
			Type: "login_response",
			Data: nil,
		}
		if e != nil {
			re.Data=Correspond{
				Code:  0,
				Msg:   "login fail",
				Error: "pw fail or user not exist",
			}
		}else{
			re.Data=Correspond{
				Code:  1,
				Msg:   "success",
				Error: "",
			}
		}
		m, _ :=json.Marshal(re)
		_, err := conn.Write(m)
		if err != nil {
			fmt.Println(err)
		}
	}*/

}

func (server *Server) Start() {
	listen, err := net.Listen("tcp", "0.0.0.0:8989")
	if err != nil {
		fmt.Println(err)
		return
	}
	go server.processConn()
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go server.read(conn)
	}
}
