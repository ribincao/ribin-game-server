package test

import (
	config "ribin-server/config"
	logger "ribin-server/logger"
	"testing"

	"go.uber.org/zap"
)

func TestLog(t *testing.T) {
	config.ParseConf("../conf.yaml", config.GlobalConfig)
	logger.InitLogger(config.GlobalConfig.LogConfig)
	logger.Debug("DebugTest :", zap.String("Env", config.GlobalConfig.ServiceConfig.Env))
	logger.Info("InfoTest :", zap.String("Env", config.GlobalConfig.ServiceConfig.Env))
	logger.Error("ErrorTest :", zap.String("Env", config.GlobalConfig.ServiceConfig.Env))
	logger.Warn("WarnTest :", zap.String("Env", config.GlobalConfig.ServiceConfig.Env))
	logger.Fatal("FatalTest :", zap.String("Env", config.GlobalConfig.ServiceConfig.Env))
}
