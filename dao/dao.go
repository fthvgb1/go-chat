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

func UserInfo(id int) *user.User {
	rd := rdm.GetRdm()
	ctx := context.Background()
	key := user.GetUserKey(id)
	u := rd.HGetAll(ctx, key).Val()
	return map2user(u)
}

func map2user(hash map[string]string) *user.User {
	id, _ := strconv.Atoi(hash["id"])
	sex, _ := strconv.Atoi(hash["sex"])
	return &user.User{
		Id:       id,
		Name:     hash["name"],
		Sex:      int8(sex),
		Password: hash["password"],
	}
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
