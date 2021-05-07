package main

import (
	"bufio"
	"chat/message"
	"chat/process"
	"chat/user"
	"fmt"
	"net"
	"os"
)

var menu1 = make(chan int)
var menu2 = make(chan int)
var rc = make(chan process.Ms, 20)
var thisUser user.User

func login(conn net.Conn) error {
	var id int
	var pw string
	for {
		stdin := bufio.NewReader(os.Stdin)
		_, err := fmt.Fscanf(stdin, "%d %s", &id, &pw)
		if err != nil {
			fmt.Println(err, "请重新输入")
		} else {
			break
		}
	}
	var msg = message.LoginS{
		Uid: id, Pw: pw,
	}
	m := message.Message{Type: "login_send", Data: msg}
	err := process.WriteConn(conn, m)
	if err != nil {
		return err
	}
	return nil
}

func sendMessage(conn net.Conn) {
	id := 0
	for {
		fmt.Print("请输入用户id(0为所有):")
		_, err := fmt.Fscanf(bufio.NewReader(os.Stdin), "%d", &id)
		if err != nil {
			fmt.Println(err)
		} else {
			break
		}
	}
	for {
		fmt.Print("请输入内容: ")
		var i string
		in := bufio.NewReader(os.Stdin)
		_, err := fmt.Fscanf(in, "%s", &i)
		if err != nil {
			fmt.Println(err)
			return
		}
		if i == "exit" {
			break
		}
		uid := id
		err = process.WriteConn(conn, message.Message{
			Type: "user_message",
			Code: 1,
			Msg:  "",
			Data: message.UserMessage{
				TotUid:       uid,
				FromUid:      thisUser.Id,
				FromUserName: thisUser.Name,
				Msg:          i,
			},
		})
		if err != nil {
			fmt.Println(err)
			//return
		}
	}
}

func showMenu(name string, ms process.Ms) {

	for {
		fmt.Printf("-----------------------------欢迎%s登录---------------------------\n", name)
		fmt.Printf("\t\t\t 1.显示在线用户列表\n")
		fmt.Printf("\t\t\t 2.发送消息\n")
		fmt.Printf("\t\t\t 3.信息列表\n")
		fmt.Printf("\t\t\t 4.返回上一级菜单\n")
		fmt.Printf("\t\t\t 请选择(1-4):")
		var k int
		_, err := fmt.Scanf("%d", &k)
		if err != nil {
			fmt.Println(err)
			continue
		}
		switch k {
		case 1:
			fmt.Println("在线用户列表")
			err := process.WriteConn(ms.Conn, message.Message{
				Type: "online_users",
				Code: 0,
				Msg:  "",
				Data: nil,
			})
			if err != nil {
				fmt.Println(err)
			}
			<-menu2
		case 2:
			sendMessage(ms.Conn)
		case 4:
			err := process.WriteConn(ms.Conn, message.Message{
				Type: "offline",
				Code: 0,
				Msg:  "",
				Data: nil,
			})
			if err != nil {
				fmt.Println(err)
			}
			menu1 <- 1
			return
		}
	}

}

func addUser(conn net.Conn) {
	fmt.Print("请输入姓名、密码、性别(1男2女)空格隔开:")
	var name, pw string
	var sex int8
	_, err := fmt.Fscanf(bufio.NewReader(os.Stdin), "%s %s %d", &name, &pw, &sex)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = process.WriteConn(conn, message.Message{
		Type: "add_user",
		Code: 0,
		Msg:  "",
		Data: user.User{
			Id:       0,
			Name:     name,
			Sex:      sex,
			Password: pw,
		},
	})
	if err != nil {
		fmt.Println(err)
		return
	}
}

func handleMsg() { //处理
	for {
		select {
		case c := <-rc:
			switch c.Msg.Type {
			case "login_response":
				if r, ok := c.Msg.Data.(*message.Correspond); ok {
					if r.Error == "" {
						fmt.Println("登录成功！")
						thisUser = r.User
						go showMenu(r.Msg, c)
					} else {
						fmt.Println("登录失败", r.Error)
						menu1 <- 1
					}
				} else {
					fmt.Println("登录失败")
				}
			case "add_user_response":
				if c.Msg.Code == 1 {
					fmt.Println("添加用户成功")
				} else {
					fmt.Println(c.Msg.Msg)
				}
				menu1 <- 1
			case "user_message":
				m := c.Msg.Data.(*message.UserMessage)
				fmt.Printf("\r%s  %s\n%s\n", m.FromUserName, m.Datetime, m.Msg)
			case "online_users":
				list := c.Msg.Data.(*message.UsersPres)
				fmt.Printf("%s\t%s\n", "id", "昵称")
				for _, pre := range list.Data {
					fmt.Printf("%d\t%s\n", pre.Id, pre.Name)
				}
				menu2 <- 1
			case "notice":
				fmt.Printf("\n系统:%s\n", c.Msg.Msg)
			case "all_users":
				m := c.Msg.Data.(*message.AllUser)
				fmt.Printf("\r%s(%d)  %s\n%s\n", m.FromUserName, m.FromUid, m.DateTime, m.Msg)
			}
		}
	}
}

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8989")

	if err != nil {
		fmt.Println(err)
		return
	}
	go handleMsg()
	go process.Read(conn, rc)
	var i int
	var loop = true

	for {
		fmt.Println("-----------------------------欢迎登录---------------------------")
		fmt.Printf("\t\t\t 1.登录聊天室\n")
		fmt.Printf("\t\t\t 2.注册用户\n")
		fmt.Printf("\t\t\t 3.退出系统\n")
		fmt.Printf("\t\t\t 请选择(1-3):")
		_, err := fmt.Scanf("%d", &i)
		fmt.Println()
		if err != nil {
			return
		}
		switch i {
		case 1:
			fmt.Print("请输入用户id和密码，空格隔开:")
			err = login(conn)
			if err != nil {
				fmt.Println("login fail :", err)
			}
			<-menu1
		case 2:
			addUser(conn)
			<-menu1
		case 3:
			//s.Store()
			loop = false
		default:
			fmt.Println("输入有误，请重新输入")
		}
		if !loop {
			return
		}
	}
}
