package db

import "github.com/ribincao/ribin-game-server/config"

var GlobalDB DB

func InitDB() {
	if config.GlobalConfig.DbConfig.RedisMode == "cluster" {
		GlobalDB = NewRedisClusterDB()
	} else {
		GlobalDB = NewRedisDB()
	}
}
