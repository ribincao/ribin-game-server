package db

import (
	"sync"

	"github.com/go-redis/redis"
	"github.com/ribincao/ribin-game-server/config"
)

var (
	openRedisDB sync.Once
	redisClient *redis.Client
)

type RedisClientDB struct{}

func initRedis() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     config.GlobalConfig.DbConfig.RedisAddr,
		Password: "",
		DB:       0,
	})
}

func NewRedisDB() *RedisClientDB {
	return &RedisClientDB{}
}

func (rc *RedisClientDB) Test() (string, error) {
	openRedisDB.Do(func() {
		initRedis()
	})
	val, err := redisClient.Get("ping").Result()
	return val, err
}
