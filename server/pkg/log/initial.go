package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

var Logger *zap.Logger
var sugaredLogger *zap.SugaredLogger

func InitLogger() {
	appLog := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    500, // megabytes
		MaxBackups: 7,
		MaxAge:     7, // days
	})
	errLog := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "logs/err.log",
		MaxSize:    500, // megabytes
		MaxBackups: 7,
		MaxAge:     7, // days
	})
	stdout := zapcore.AddSync(os.Stdout)
	config := newConfig()

	core := zapcore.NewTee(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(*config),
			zapcore.NewMultiWriteSyncer(appLog, stdout),
			zap.InfoLevel,
		),
		zapcore.NewCore(
			zapcore.NewJSONEncoder(*config),
			zapcore.NewMultiWriteSyncer(errLog),
			zap.ErrorLevel,
		),
	)
	caller := zap.AddCaller()
	Logger = zap.New(core,caller,zap.AddCallerSkip(1))
	sugaredLogger = Logger.Sugar()
}

func newConfig() *zapcore.EncoderConfig {
	return &zapcore.EncoderConfig{
		MessageKey:     "message",
		TimeKey:        "time",
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		LevelKey:       "level",
		CallerKey:      "caller",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		FunctionKey:    zapcore.OmitKey,
		StacktraceKey:  "trace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeDuration: zapcore.StringDurationEncoder,
	}
}