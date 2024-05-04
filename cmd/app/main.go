package main

import (
	"stats-of/internal/logger"
	"stats-of/internal/storagetestsutils"

	"github.com/go-redis/redis/v8"
)

type CsvDbManager struct {
	FilePath    string
	RedisClient *redis.Client
}

func main() {
	logger.InitLogger()
	logger.Log.Info("Starting application and reading configuration...")

	// Launching the application - can be commented out to just add data to the database
	// if err := app.RunApp(); err != nil {
	// 	logger.Log.Fatal("Error occurred", zap.Error(err))
	// }
	// ------------------------------------------------

	// Initializing database connection - block for adding data from CSV to the database
	// if err := storagetestsutils.HandleCsvToDb(); err != nil {
	// 	logger.Log.Fatal("Error processing CSV data", zap.Error(err))
	// }
	// ------------------------------------------------

	storagetestsutils.RunChatSearch(10)
}
