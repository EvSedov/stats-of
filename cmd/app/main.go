package main

import (
	"stats-of/internal/app"
	"stats-of/internal/logger"
	"stats-of/internal/storagetestsutils"

	"go.uber.org/zap"
)

func main() {
	logger.InitLogger()
	logger.Log.Info("Starting application and reading configuration...")

	// Launching the application - can be commented out to just add data to the database
	if err := app.RunApp(); err != nil {
		logger.Log.Fatal("Error occurred", zap.Error(err))
	}
	// ------------------------------------------------

	// Initializing database connection - block for adding data from CSV to the database
	if err := storagetestsutils.HandleCsvToDb(); err != nil {
		logger.Log.Fatal("Error processing CSV data", zap.Error(err))
	}
	// ------------------------------------------------

	// Initializing connection to Redis - block for adding test data for Redis stress testing
	redisClient := storagetestsutils.InitDb()
	if redisClient == nil {
		logger.Log.Fatal("Failed to initialize Redis client", zap.String("reason", "client is nil"))
	}

	// Creating an instance of CsvDbManager
	manager := storagetestsutils.NewCsvDbManager("", redisClient) // File path is not used in this context

	// Calling AddUsersData with the desired number of users
	userCount := 50000 // Approximate number of users for the test
	if err := manager.AddUsersData(userCount); err != nil {
		logger.Log.Fatal("Error adding user data", zap.Error(err))
	}

	logger.Log.Info("Data successfully added for users")
	// ------------------------------------------------
}
