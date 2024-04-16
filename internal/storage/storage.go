package storage

import (
	"context"
	"stats-of/internal/storage/redis"
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
	if storageType == Redis {
		options, err := redis.CreateOptions()
		if err != nil {
			return nil, err
		}

		client = redis.NewRedisClient(options)
		return client, nil
	}

	return nil, nil
}
