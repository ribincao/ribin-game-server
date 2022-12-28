package logger

import (
	"os"
	"time"

	config "ribin-server/config"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

const (
	DefaultLogMaxSize  = 512 // DefaultLogMaxSize is the default log max size
	DefaultLogMaxAge   = 7   // DefaultLogMaxAge is the default log max save time
	DefaultBackupCount = 9   // DefaultBackupCount is the default log max backup count
)

func getLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	case "warn":
		return zapcore.WarnLevel
	case "panic":
		return zapcore.PanicLevel
	}
	return zapcore.DebugLevel
}

func InitLogger(config *config.LogConfig) {
	if config == nil {
		return
	}
	if config.LogMaxSize == 0 {
		config.LogMaxSize = DefaultLogMaxSize
	}
	if config.LogMaxAge == 0 {
		config.LogMaxAge = DefaultLogMaxAge
	}
	if config.BackupCount == 0 {
		config.BackupCount = DefaultBackupCount
	}

	var opts []zapcore.WriteSyncer
	switch config.LogMode {
	case "file":
		opts = append(opts, zapcore.AddSync(&lumberjack.Logger{
			Filename:  config.LogPath,
			MaxSize:   config.LogMaxSize,
			MaxAge:    config.LogMaxAge,
			LocalTime: true,
			Compress:  false,
		}))
	case "console":
		zapcore.AddSync(os.Stdout)
	case "combine":
		opts = append(opts, zapcore.AddSync(&lumberjack.Logger{
			Filename:   config.LogPath,
			MaxSize:    config.LogMaxSize,
			MaxAge:     config.LogMaxAge,
			LocalTime:  true,
			Compress:   false,
			MaxBackups: config.BackupCount,
		}), zapcore.AddSync(os.Stdout))
	default:
		zapcore.AddSync(os.Stdout)
	}

	syncWriter := zapcore.NewMultiWriteSyncer(opts...)

	customTimeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
	}
	customLevelEncoder := func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("[" + level.CapitalString() + "]")
	}
	customCallerEncoder := func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("[" + caller.TrimmedPath() + "]")
	}

	encoderConf := zapcore.EncoderConfig{
		CallerKey:      "caller_line",
		LevelKey:       "level_name",
		MessageKey:     "msg",
		TimeKey:        "time",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeTime:     customTimeEncoder,
		EncodeLevel:    customLevelEncoder,
		EncodeCaller:   customCallerEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}

	encoder := zapcore.NewJSONEncoder(encoderConf)
	core := zapcore.NewCore(encoder, syncWriter, zap.NewAtomicLevelAt(getLevel(config.LogLevel)))
	logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
}

func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	logger.Fatal(msg, fields...)
}
