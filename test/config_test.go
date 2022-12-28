package test

import (
	"fmt"
	config "ribin-server/config"
	"testing"
)

func TestConfig(t *testing.T) {
	config.ParseConf("../conf.yaml", config.GlobalConfig)
	fmt.Println("ConfigTest Env:", config.GlobalConfig.ServiceConfig.Env)
}
