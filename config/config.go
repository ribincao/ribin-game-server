package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

var GlobalConfig = &Config{}

type Config struct {
	ServiceConfig *ServiceConfig `yaml:"serviceConfig"`
	LogConfig     *LogConfig     `yaml:"logConfig"`
}

type ServiceConfig struct {
	Env string `yaml:"env"`
}

type LogConfig struct {
	LogPath     string `yaml:"logPath"`
	LogLevel    string `yaml:"logLevel"`
	LogMaxAge   int    `yaml:"logMaxAge"`
	LogMaxSize  int    `yaml:"logMaxSize"`
	LogMode     string `yaml:"logMode"`
	BackupCount int    `yaml:"backupCount"`
}

func ParseConf(path string, conf *Config) *Config {
	y, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(y, conf)
	if err != nil {
		panic(err)
	}
	initEnv(conf)
	return conf
}

func initEnv(conf *Config) {
	os.Setenv("SERVICE_ENV", conf.ServiceConfig.Env)
}
