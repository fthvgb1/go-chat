package user

import "errors"

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Sex      int8   `json:"sex"`
	Password string `json:"password"`
}

func (u *User) CheckPassword(p string) error {
	if p != u.Password {
		return errors.New("password error")
	}
	return nil
}
