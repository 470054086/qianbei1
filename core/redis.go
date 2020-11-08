package core

import (
	"github.com/go-redis/redis/v8"
)

var g_redis *redis.Client

func init() {
	config := Config().Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password, // no password set
		DB:       config.Db,       // use default DB
	})
	g_redis = rdb
}

func Redis() *redis.Client {
	return g_redis
}
