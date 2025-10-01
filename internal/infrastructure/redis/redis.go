package redis

import (
	"github.com/TheAmirhosssein/cool-password-manage/config"
	"github.com/redis/go-redis/v9"
)

var client *redis.Client

func initClient() *redis.Client {
	conf := config.GetConfig()
	opt, err := redis.ParseURL(conf.Redis.URL)
	if err != nil {
		panic(err)
	}

	return redis.NewClient(opt)
}

func GetClient() *redis.Client {
	if client == nil {
		client = initClient()
	}
	return client
}
