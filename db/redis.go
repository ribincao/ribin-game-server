package db

import (
	"sync"

	"github.com/go-redis/redis"
)

var (
	openRedisDB sync.Once
	redisClient *redis.Client
)

type RedisClientDB struct{}

func initRedis() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func NewRedisDB() *RedisClientDB {
	return &RedisClientDB{}
}

func (rd *RedisClientDB) Test() (string, error) {
	openRedisDB.Do(func() {
		initRedis()
	})
	val, err := redisClient.Get("ping").Result()
	return val, err
}
