package db

import (
	"sync"

	"github.com/go-redis/redis"
	"github.com/ribincao/ribin-game-server/config"
	"github.com/ribincao/ribin-game-server/logger"
)

var (
	openRedisDB sync.Once
	redisClient *redis.Client
)

type RedisClientDB struct{}

func initRedisClient() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     config.GlobalConfig.DbConfig.RedisAddr,
		Password: "",
		DB:       0,
	})
	logger.Info("[Engine-Tool] Redis Client Initialized!")
}

func NewRedisDB() *RedisClientDB {
	return &RedisClientDB{}
}

func (rc *RedisClientDB) Test() (string, error) {
	openRedisDB.Do(func() {
		initRedisClient()
	})
	val, err := redisClient.Get("ping").Result()
	return val, err
}
