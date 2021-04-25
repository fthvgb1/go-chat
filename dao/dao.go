package dao

import (
	"chat/rdm"
	"chat/user"
	"chat/utils"
	"context"
	"errors"
	"strconv"
)

var Ctx context.Context

func init() {
	Ctx = context.Background()
}

func AddUser(user user.User) error {
	getRdm := rdm.GetRdm()
	if user.Name == "" {
		return errors.New("user name can't be empty")
	}
	if user.Password == "" {
		return errors.New("user password can't be empty")
	}

	if utils.IsContain(user.Sex, []int8{1, 2}) < 0 {
		return errors.New("user name can't be empty")
	}

	r := getRdm.HMSet(Ctx, "go_chat_Users", map[string]interface{}{
		"id":   strconv.Itoa(user.Id),
		"name": user.Name,
		"pw":   user.Password,
		"sex":  strconv.Itoa(int(user.Sex)),
	})
	if e := r.Err(); e != nil {
		return e
	}
	return nil
}
