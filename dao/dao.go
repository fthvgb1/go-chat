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

func AddUser(u user.User) error {
	getRdm := rdm.GetRdm()
	if u.Name == "" {
		return errors.New("user name can't be empty")
	}
	if u.Password == "" {
		return errors.New("user password can't be empty")
	}

	if utils.IsContain(u.Sex, []int8{1, 2}) < 0 {
		return errors.New("user name can't be empty")
	}

	r := getRdm.HMSet(Ctx, user.GetUserKey(u.Id), map[string]interface{}{
		"id":       strconv.Itoa(u.Id),
		"name":     u.Name,
		"password": u.Password,
		"sex":      strconv.Itoa(int(u.Sex)),
	})
	if e := r.Err(); e != nil {
		return e
	}
	return nil
}
