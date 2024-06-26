package log

import (
	"log"
	"msm/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.SugaredLogger

func init() {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
	level, err := zapcore.ParseLevel(config.CF.LogLevel)
	if err != nil {
		log.Printf("日志等级错误！不存在“%v”日志等级", config.CF.LogLevel)
		level = zap.DebugLevel
	}
	atom := zap.NewAtomicLevelAt(level)
	zap.NewDevelopmentConfig()
	config := zap.Config{
		Level:            atom,
		Development:      true,
		Encoding:         "console",
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{"stdout", "info.log"},
		ErrorOutputPaths: []string{"stderr"},
	}
	log, _ := config.Build()
	Logger = log.Sugar()
}
