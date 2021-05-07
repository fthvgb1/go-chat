package lib

import (
	"chat/dao"
	"chat/message"
	"chat/process"
	"chat/rdm"
	"chat/user"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"strconv"
	"time"
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

func (server *Server) Login(uid int, pw string) (*user.User, error) {
	for _, u := range server.Users {
		if uid == u.Id && pw == u.Password {
			return &u, nil
		}
	}
	u := rdm.GetRdm().HGetAll(context.Background(), user.GetUserKey(uid))
	if u != nil && u.Val()["password"] == pw {
		t := u.Val()
		id, _ := strconv.Atoi(t["id"])
		sex, _ := strconv.Atoi(t["sex"])
		uu := &user.User{
			Id:   id,
			Name: t["name"],
			Sex:  int8(sex),
			//Password: "",
		}
		return uu, nil
	}
	return nil, errors.New("not found this user")
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

func (server *Server) AddUser(name, pw string, sex int8) error {
	server.Id++
	u := user.User{
		Id: server.Id, Name: name, Sex: sex, Password: pw,
	}
	server.Users = append(server.Users, u)
	err := dao.AddUser(u)
	if err != nil {
		return err
	}
	return nil
}

func (server *Server) Process(conn net.Conn) {
	m := make([]byte, 0)
	read, err := conn.Read(m)
	if err != nil {
		return
	}
	fmt.Println(read)
}

func notice(message2 message.Message) {
	for _, userProcess := range process.GetOnlineUsers() {
		err := process.WriteConn(userProcess.Conn, message2)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (server *Server) do(s *message.LoginS, conn net.Conn) (string, error) {
	login, err := server.Login(s.Uid, s.Pw)
	re := message.Message{
		Type: "login_response",
		Data: nil,
	}

	r := message.Correspond{
		Code:  0,
		Msg:   "login fail",
		Error: "pw fail or login not exist ",
	}
	if err == nil {
		r.Code = 1
		r.Msg = login.Name
		r.Error = ""
		r.User = *login
		notice(message.Message{
			Type: "notice",
			Code: 0,
			Msg:  fmt.Sprintf("%s已上线", login.Name),
			Data: nil,
		})
		process.Push(&process.UserProcess{
			Uid:  login.Id,
			Conn: conn,
		})
	} else {
		r.Error += err.Error()
	}
	re.Data = r

	err = process.WriteConn(conn, re)
	return "login_send", err
}

func (s *Server) processConn() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
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
			case "add_user":
				u := c.Msg.Data.(*user.User)
				err := s.AddUser(u.Name, u.Password, u.Sex)
				r := message.Message{
					Type: "add_user_response",
					Code: 1,
					Msg:  "",
					Data: nil,
				}
				if err != nil {
					fmt.Println(err)
					r.Msg = err.Error()
					r.Code = 0
				}
				err = process.WriteConn(c.Conn, r)
				if err != nil {
					fmt.Println(err)
				}
			case "user_message":
				data := c.Msg.Data.(*message.UserMessage)
				now := time.Now().Format("2006-01-02 15:04:05")
				if data.TotUid > 0 {
					to := process.Get(data.TotUid)
					if to == nil {
						to = process.Get(data.FromUid)
						if to == nil {
							break
						}
						toU := dao.UserInfo(data.TotUid)
						err := process.WriteConn(to.Conn, message.Message{
							Type: "notice",
							Code: 0,
							Msg:  "用户" + toU.Name + "已下线",
							Data: nil,
						})
						if err != nil {
							fmt.Println(err)
						}
						break
					}
					err := process.WriteConn(to.Conn, message.Message{
						Type: "user_message",
						Code: 0,
						Msg:  "",
						Data: message.UserMessage{
							FromUid:      data.FromUid,
							FromUserName: data.FromUserName,
							TotUid:       data.TotUid,
							Msg:          data.Msg,
							Datetime:     now,
						},
					})
					if err != nil {
						fmt.Println(err)
					}

				} else {
					fmt.Println(data.FromUserName, ":", data.Msg)
					for _, userProcess := range process.GetOnlineUsers() {
						err := process.WriteConn(userProcess.Conn, message.Message{
							Type: "all_users",
							Code: 0,
							Msg:  "",
							Data: message.AllUser{
								FromUid:      data.FromUid,
								FromUserName: data.FromUserName,
								Msg:          data.Msg,
								DateTime:     now,
							},
						})
						if err != nil {
							fmt.Println(err)
						}
					}
				}
			case "online_users":
				all := process.GetOnlineUsers()
				arr := make([]message.UserPre, 0)
				for _, userProcess := range all {
					v := rdm.GetRdm().HGet(context.Background(), user.GetUserKey(userProcess.Uid), "name").Val()
					arr = append(arr, message.UserPre{
						Id:   userProcess.Uid,
						Name: v,
					})
				}
				err := process.WriteConn(c.Conn, message.Message{
					Type: "online_users",
					Code: 1,
					Msg:  "",
					Data: message.UsersPres{Data: arr},
				})
				if err != nil {
					fmt.Println(err)
				}
			case "offline":
				id := process.Disconnect(c.Conn)
				u := dao.UserInfo(id)
				for _, userProcess := range process.GetOnlineUsers() {
					err := process.WriteConn(userProcess.Conn, message.Message{
						Type: "notice",
						Code: 0,
						Msg:  "用户" + u.Name + "已下线",
						Data: nil,
					})
					if err != nil {
						fmt.Println(err)
					}
				}
			}
		}
	}
}

func (server *Server) read(conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			id := process.Disconnect(conn)
			u := dao.UserInfo(id)
			for _, userProcess := range process.GetOnlineUsers() {
				err := process.WriteConn(userProcess.Conn, message.Message{
					Type: "notice",
					Code: 0,
					Msg:  "用户" + u.Name + "已掉线",
					Data: nil,
				})
				if err != nil {
					fmt.Println(err)
				}
			}
			fmt.Println(err, "ssssssssssssssss")
		}
	}(conn)

	process.Read(conn, Rchan)

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
