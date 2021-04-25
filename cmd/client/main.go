package main

import (
	"chat/message"
	"chat/process"
	"chat/user"
	"fmt"
	"net"
)

var cc = make(chan int)
var rc = make(chan process.Ms, 20)

func login(conn net.Conn) error {
	var id int
	var pw string
	_, err := fmt.Scanf("%d %s", &id, &pw)
	if err != nil {
		return err
	}
	var msg = message.LoginS{
		Uid: id, Pw: pw,
	}
	m := message.Message{Type: "login_send", Data: msg}
	err = process.WriteConn(conn, m)
	if err != nil {
		return err
	}
	return nil
}

func showMenu(name string) {

	for {
		fmt.Printf("-----------------------------欢迎%s登录---------------------------\n", name)
		fmt.Printf("\t\t\t 1.显示在线用户列表\n")
		fmt.Printf("\t\t\t 2.发送消息\n")
		fmt.Printf("\t\t\t 3.信息列表\n")
		fmt.Printf("\t\t\t 4.退出\n")
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
		case 2:
			fmt.Println("发送信息")
		case 4:
			return
		}
	}

}

func addUser(conn net.Conn) {
	fmt.Print("请输入姓名、密码、性别(1男2女)空格隔开:")
	var name, pw string
	var sex int8
	_, err := fmt.Scanf("%s %s %d", &name, &pw, &sex)
	if err != nil {
		fmt.Println(err)
		return
	}
	process.WriteConn(conn, message.Message{
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
						showMenu(r.Msg)
					} else {
						fmt.Println("登录失败", r.Error)
					}
				} else {
					fmt.Println("登录失败")
				}
				cc <- 1
			case "add_user_response":
				if c.Msg.Code == 1 {
					fmt.Println("添加用户成功")
				} else {
					fmt.Println(c.Msg.Msg)
				}
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
			<-cc
		case 2:
			addUser(conn)
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
