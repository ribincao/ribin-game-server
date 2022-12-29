package test

import (
	"fmt"
	"testing"
	"time"

	config "github.com/ribincao/ribin-game-server/config"
	logger "github.com/ribincao/ribin-game-server/logger"
	utils "github.com/ribincao/ribin-game-server/utils"
)

func f() {
	var infoMap map[int]string
	infoMap[1] = "a"
}

func TestRecover(t *testing.T) {
	config.ParseConf("../conf.yaml", config.GlobalConfig)
	logger.InitLogger(config.GlobalConfig.LogConfig)
	utils.GoWithRecover(f)

	time.Sleep(2 * time.Second)
	fmt.Println("Exit.")
}
