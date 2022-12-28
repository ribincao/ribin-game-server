package config

import (
	"fmt"
	"testing"
)

func TestConfig(t *testing.T) {
	ParseConf("../conf.yaml", GlobalConfig)
	fmt.Println("ConfigTest Env:", GlobalConfig.ServiceConfig.Env)
}
