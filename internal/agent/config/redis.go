package config

import "github.com/redis/go-redis/v9"

var RedisClientQ = redis.NewClient(&redis.Options{
	Addr: Conf.redis_host + ":" + Conf.redis_port,
	DB:   0,
})

var RedisClientA = redis.NewClient(&redis.Options{
	Addr: Conf.redis_host + ":" + Conf.redis_port,
	DB:   1,
})
