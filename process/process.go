package process

import (
	"chat/message"
	"encoding/json"
	"fmt"
	"net"
	"reflect"
)

type Ms struct {
	Msg  message.Message
	Conn net.Conn
}

type Transfer struct {
	Conn net.Conn
	Buf  [64 << 10]byte
}

func WriteConn(c net.Conn, message message.Message) error {
	m, _ := json.Marshal(message)
	_, err := c.Write(m)
	if err != nil {
		return err
	}
	return nil
}

func (t *Transfer) ReadConn(c chan Ms) {

	defer func() {
		err := t.Conn.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}()

	for {
		var info = t.Buf
		i, err := t.Conn.Read(info[:])
		if err != nil {
			fmt.Println(err, "errrrrrrrrrrrrr")
			return
		}
		msg := message.Message{}
		err = json.Unmarshal(info[:i], &msg)
		if err != nil {
			fmt.Println(err)
		}
		tt, ok := message.MsgType[msg.Type]
		var mm = msg
		if ok {
			mm.Data = reflect.New(tt).Interface()
			err = json.Unmarshal(info[:i], &mm)
			if err != nil {
				fmt.Println(err)
			}
		}
		w := Ms{
			Msg:  mm,
			Conn: t.Conn,
		}
		c <- w
	}
}

func Read(conn net.Conn, c chan Ms) {
	t := Transfer{
		Conn: conn,
		Buf:  [65536]byte{},
	}
	t.ReadConn(c)
}
