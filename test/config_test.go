package test

import (
	"fmt"
	"testing"

	"github.com/ribincao/ribin-game-server/config"
)

func TestConfig(t *testing.T) {
	config.ParseConf("../conf.yaml", config.GlobalConfig)
	fmt.Println("ConfigTest Env:", config.GlobalConfig.ServiceConfig.Env)
}
