package redis

import (
	"github.com/TheAmirhosssein/cool-password-manage/config"
	"github.com/redis/go-redis/v9"
)

func Client() *redis.Client {
	conf := config.GetConfig()
	opt, err := redis.ParseURL(conf.Redis.URL)
	if err != nil {
		panic(err)
	}

	return redis.NewClient(opt)
}

func GetClient() *redis.Client {
	return Client()
}
