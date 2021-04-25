package message

import (
	"chat/user"
	"reflect"
)

var MsgType = map[string]reflect.Type{
	"login_send":     reflect.TypeOf(&LoginS{}).Elem(),
	"login_response": reflect.TypeOf(&Correspond{}).Elem(),
	"add_user":       reflect.TypeOf(&user.User{}).Elem(),
}

type Message struct {
	//Id int
	Type string
	Code int
	Msg  string
	Data interface{}
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
}
