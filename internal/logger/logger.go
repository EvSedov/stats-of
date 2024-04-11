package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

func InitLogger() {
	logDir := "internal/logger/logs"
	logPath := logDir + "/log.txt"

	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err := os.MkdirAll(logDir, 0755)
		if err != nil {
			panic(err)
		}
	}

	config := zap.NewProductionConfig()
	config.OutputPaths = []string{
		logPath,
		"stdout",
	}
	config.ErrorOutputPaths = []string{
		logPath,
		"stderr",
	}
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	var err error
	Log, err = config.Build()
	if err != nil {
		panic(err)
	}

	Log.Info("Logger initialized")
}
