package container

import (
	"github.com/Nonsensersunny/poker_game/conf"
	"github.com/go-redis/redis/v8"
)

var (
	DefaultContainer *Container
)

type Container struct {
	Redis *redis.Client
}

func init() {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     conf.DefaultConfig.Addr(),
		Password: conf.DefaultConfig.Password,
		DB:       conf.DefaultConfig.DB,
	})

	DefaultContainer = &Container{
		Redis: redisClient,
	}
}
