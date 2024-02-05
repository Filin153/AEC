package config

import "github.com/redis/go-redis/v9"

var RedisClient = redis.NewClient(&redis.Options{
	Addr: Conf.redis_host + ":" + Conf.redis_port,
	DB:   0,
})
