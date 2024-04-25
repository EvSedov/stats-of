package main

import (
	"stats-of/internal/app"
	"stats-of/internal/logger"

	"go.uber.org/zap"
)

func main() {
	logger.InitLogger()
	logger.Log.Info("Starting application and reading configuration...")

	if err := app.RunApp(); err != nil {
		logger.Log.Fatal("Error occurred", zap.Error(err))
	}
}
