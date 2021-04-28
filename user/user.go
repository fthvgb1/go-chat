package user

import (
	"errors"
	"fmt"
)

const Key = "go_chat_Users:%d"

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Sex      int8   `json:"sex"`
	Password string `json:"password"`
}

func GetUserKey(id int) string {
	return fmt.Sprintf(Key, id)
}

func (u *User) CheckPassword(p string) error {
	if p != u.Password {
		return errors.New("password error")
	}
	return nil
}
