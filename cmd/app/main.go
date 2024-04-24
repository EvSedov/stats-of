package main

import (
	"stats-of/internal/logger"
	"stats-of/internal/storagetestsutils"

	"go.uber.org/zap"
)

func main() {
	logger.InitLogger()
	logger.Log.Info("Starting application and reading configuration...")

	// if err := app.RunApp(); err != nil {
	// 	logger.Log.Fatal("Error occurred", zap.Error(err))
	// }

	if err := storagetestsutils.HandleCsvToDb(); err != nil {
		logger.Log.Fatal("Error processing CSV data", zap.Error(err))
	}

	// Инициализация подключения к Redis
	redisClient := storagetestsutils.InitDb()
	if redisClient == nil {
		logger.Log.Fatal("Failed to initialize Redis client", zap.String("reason", "client is nil"))
	}

	// Создание экземпляра CsvDbManager
	manager := storagetestsutils.NewCsvDbManager("", redisClient) // Путь к файлу не используется в данном контексте

	// Вызов AddUsersData с желаемым количеством пользователей
	userCount := 50000 // Примерное количество пользователей для теста
	if err := manager.AddUsersData(userCount); err != nil {
		logger.Log.Fatal("Error adding user data", zap.Error(err))
	}

	logger.Log.Info("Data successfully added for users")
}
