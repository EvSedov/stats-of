package storage

import (
	"context"
	"fmt"
	"stats-of/internal/logger"
	"stats-of/internal/storage/redis"

	"go.uber.org/zap"
)

type (
	StorageType string

	Storage interface {
		// Open() error
		Ping(ctx context.Context) error
		FindKeysByPattern(pattern string) ([]string, error)
		FindKeyByGetRequest(key string) (string, error)
	}
)

const (
	Redis StorageType = "redis"
	// Map   StorageType = "map"
)

var client Storage

// Фабричная функция для создания экземпляра базы данных в соответствии с указанным типом
func NewStorage(storageType StorageType) (Storage, error) {
	logger.Log.Info("Initializing new storage", zap.String("storageType", string(storageType)))

	if storageType == Redis {
		// Создание опций для Redis
		options, err := redis.CreateOptions()
		if err != nil {
			logger.Log.Error("Failed to create Redis options", zap.Error(err))
			return nil, err
		}

		// Создание клиента Redis
		client := redis.NewRedisClient(options)
		logger.Log.Info("Redis client created successfully")
		return client, nil
	}

	// Если тип хранилища не поддерживается или не указан
	logger.Log.Warn("Storage type not supported or not specified", zap.String("storageType", string(storageType)))
	return nil, fmt.Errorf("storage type '%s' is not supported", storageType)
}
