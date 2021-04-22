package user

import "errors"

type User struct {
	Id       int
	Name     string
	Sex      int8
	Password string
}

func (u *User) CheckPassword(p string) error {
	if p != u.Password {
		return errors.New("password error")
	}
	return nil
}
