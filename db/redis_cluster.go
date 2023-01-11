package db

import (
	"sync"

	"github.com/go-redis/redis"
	"github.com/ribincao/ribin-game-server/config"
)

var (
	openRedisClusterDB sync.Once
	redisClusterClient *redis.ClusterClient
)

type RedisClusterClientDB struct{}

func initRedisCluster() {
	opts := &redis.ClusterOptions{
		Addrs: []string{config.GlobalConfig.DbConfig.RedisAddr},
	}
	redisClusterClient = redis.NewClusterClient(opts)
}

func NewRedisClusterDB() *RedisClusterClientDB {
	return &RedisClusterClientDB{}
}

func (rc *RedisClusterClientDB) Test() (string, error) {
	openRedisClusterDB.Do(func() {
		initRedisCluster()
	})
	cmd := redisClusterClient.Get("ping")
	return cmd.Val(), cmd.Err()
}
