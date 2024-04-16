package redis

import (
	"context"
	"stats-of/internal/logger"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

type (
	// Storage структура для работы с redis клиентом
	Storage struct {
		Client *redis.Client
	}
)

// NewRedisService функция для создания нового экземпляра DBService
func NewRedisClient(opt *Options) *Storage {
	client := redis.NewClient(&redis.Options{
		Addr:     opt.Addr,
		Password: opt.Password,
		DB:       opt.DB,
	})

	return &Storage{Client: client}
}

func (r *Storage) Ping(ctx context.Context) error {
	result, err := r.Client.Ping().Result() // Использование Ping из библиотеки go-redis
	if err != nil {
		return err
	}
	logger.Log.Info("Redis Ping Response", zap.String("response", result))
	return nil
}

// FindKeysByPattern метод для поиска ключей в Redis по шаблону
func (r *Storage) FindKeysByPattern(pattern string) ([]string, error) {
	var cursor uint64
	var keys []string
	for {
		k, nextCursor, err := r.Client.Scan(cursor, pattern, 0).Result()
		if err != nil {
			return nil, err
		}
		keys = append(keys, k...)
		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	return keys, nil
}

func (r *Storage) FindKeyByGetRequest(key string) (string, error) {
	result, err := r.Client.Get(key).Result()
	if err == redis.Nil {
		// Ключ не найден
		logger.Log.Info("Ключ не найден", zap.String("key", key))
		return "", nil // Возвращаем пустую строку без ошибки, если такое поведение приемлемо
	} else if err != nil {
		// Произошла другая ошибка
		return "", err
	}
	// Возвращаем результат, если ключ найден и ошибок нет
	return result, nil
}
