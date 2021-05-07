package message

import (
	"chat/user"
	"reflect"
)

var MsgType = map[string]reflect.Type{
	"login_send":     reflect.TypeOf(&LoginS{}).Elem(),
	"login_response": reflect.TypeOf(&Correspond{}).Elem(),
	"add_user":       reflect.TypeOf(&user.User{}).Elem(),
	"user_message":   reflect.TypeOf(&UserMessage{}).Elem(),
	"online_users":   reflect.TypeOf(&UsersPres{}).Elem(),
	"all_users":      reflect.TypeOf(&AllUser{}).Elem(),
}

type Message struct {
	//Id int
	Type string      `json:"type"`
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type AllUser struct {
	FromUid      int
	FromUserName string
	Msg          string
	DateTime     string
}

type UsersPres struct {
	Data []UserPre
}

type UserPre struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type UserMessage struct {
	FromUid      int
	FromUserName string
	TotUid       int
	Msg          string
	Datetime     string
}

type LoginS struct {
	Uid  int
	Pw   string
	Name string
}

type Correspond struct {
	Code  int
	Msg   string
	Error string
	User  user.User
}
