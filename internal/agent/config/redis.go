package config

import (
	"context"
	"github.com/redis/go-redis/v9"
)

var RedisClientQ = redis.NewClient(&redis.Options{})
var RedisClientA = redis.NewClient(&redis.Options{})

func init() {

	var RedisURL = Conf.Redis_host + ":" + Conf.Redis_port

	RedisClientQ = redis.NewClient(&redis.Options{
		Addr: RedisURL,
		DB:   0,
	})

	RedisClientA = redis.NewClient(&redis.Options{
		Addr: RedisURL,
		DB:   1,
	})

	_, err := RedisClientQ.Ping(context.Background()).Result()
	if err != nil {
		panic("Ошибка при подключении к Redis: " + err.Error())
	}

	_, err = RedisClientA.Ping(context.Background()).Result()
	if err != nil {
		panic("Ошибка при подключении к Redis: " + err.Error())
	}
}
