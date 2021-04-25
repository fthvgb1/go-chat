package rdm

import "github.com/go-redis/redis/v8"

var RDM *redis.Client

func GetRdm() *redis.Client {
	if RDM == nil {
		RDM = redis.NewClient(&redis.Options{
			Addr: "127.0.0.1:6379",
			DB:   0,
		})
	}
	return RDM
}
