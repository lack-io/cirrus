package log

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/lack-io/cirrus/config"
)

var logger *zap.Logger
var sugar *zap.SugaredLogger

func Init(cfg *config.Logger) error {
	var ws zapcore.WriteSyncer
	if cfg != nil {
		lumberJackLogger := &lumberjack.Logger{
			Filename:   cfg.Filename,
			MaxSize:    cfg.MaxAge,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
			LocalTime:  cfg.LocalTime,
			Compress:   cfg.Compress,
		}
		ws = zapcore.AddSync(lumberJackLogger)
	} else {
		ws = zapcore.AddSync(os.Stdout)
	}

	encoder := getEncoder()
	core := zapcore.NewCore(encoder, ws, zapcore.DebugLevel)
	logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	sugar = logger.Sugar()
	return nil
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func Debug(s string) {
	sugar.Debugw(s)
}

func Debugf(format string, v ...interface{}) {
	sugar.Debugf(format, v...)
}

func Info(s string) {
	sugar.Infow(s)
}

func Infof(format string, v ...interface{}) {
	sugar.Infof(format, v...)
}

func Warn(s string) {
	sugar.Warnw(s)
}

func Warnf(format string, v ...interface{}) {
	sugar.Warnf(format, v...)
}

func Error(s string) {
	sugar.Errorw(s)
}

func Errorf(format string, v ...interface{}) {
	sugar.Errorf(format, v...)
}

func Fatal(s string) {
	sugar.Fatalw(s)
}

func Fatalf(format string, v ...interface{}) {
	sugar.Fatalf(format, v...)
}

func Sync() error {
	return sugar.Sync()
}
