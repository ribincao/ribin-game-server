package db

import (
	"sync"

	"github.com/go-redis/redis"
	"github.com/ribincao/ribin-game-server/config"
	"github.com/ribincao/ribin-game-server/logger"
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
	logger.Info("[Engine-Tool] Redis Cluster Client Initialized!")
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
