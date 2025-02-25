package log

import (
	"log"

	"github.com/lzh-1625/go_process_manager/config"

	"github.com/timandy/routine"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *logWithAdditional

type logWithAdditional struct {
	*zap.SugaredLogger
	threadLocal routine.ThreadLocal[[]any]
}

func (l *logWithAdditional) Infow(msg string, keysAndValues ...interface{}) {
	keysAndValues = append(keysAndValues, l.threadLocal.Get()...)
	l.SugaredLogger.WithOptions(zap.AddCallerSkip(1)).Infow(msg, keysAndValues...)
}

func (l *logWithAdditional) Debugw(msg string, keysAndValues ...interface{}) {
	keysAndValues = append(keysAndValues, l.threadLocal.Get()...)
	l.SugaredLogger.WithOptions(zap.AddCallerSkip(1)).Debugw(msg, keysAndValues...)
}

func (l *logWithAdditional) Errorw(msg string, keysAndValues ...interface{}) {
	keysAndValues = append(keysAndValues, l.threadLocal.Get()...)
	l.SugaredLogger.WithOptions(zap.AddCallerSkip(1)).Errorw(msg, keysAndValues...)
}

func (l *logWithAdditional) Warnw(msg string, keysAndValues ...interface{}) {
	keysAndValues = append(keysAndValues, l.threadLocal.Get()...)
	l.SugaredLogger.WithOptions(zap.AddCallerSkip(1)).Warnw(msg, keysAndValues...)
}

func (l *logWithAdditional) AddAdditionalInfo(k, v any) {
	l.threadLocal.Set(append(l.threadLocal.Get(), k, v))
}

func (l *logWithAdditional) DeleteAdditionalInfo(layer int) {
	if layer < 0 {
		l.threadLocal.Set([]any{})
		return
	}
	oldKv := l.threadLocal.Get()
	if len(oldKv) < layer*2 {
		l.threadLocal.Set([]any{})
		return
	}
	l.threadLocal.Set(oldKv[:len(oldKv)-2*layer])
}

func InitLog() {

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
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	level, err := zapcore.ParseLevel(config.CF.LogLevel)
	if err != nil {
		log.Printf("日志等级错误！不存在“%v”日志等级", config.CF.LogLevel)
		level = zap.DebugLevel
	}
	atom := zap.NewAtomicLevelAt(level)
	zap.NewDevelopmentConfig()
	var outputPaths []string = []string{"info.log"}
	if !config.CF.Tui { // 不使用tui则打印日志到stdout
		outputPaths = append(outputPaths, "stdout")
	}
	config := zap.Config{
		Level:            atom,
		Development:      true,
		Encoding:         "console",
		EncoderConfig:    encoderConfig,
		OutputPaths:      outputPaths,
		ErrorOutputPaths: []string{"stderr"},
	}
	log, _ := config.Build()
	Logger = &logWithAdditional{
		SugaredLogger: log.Sugar(),
		threadLocal:   routine.NewThreadLocal[[]any](),
	}
}
