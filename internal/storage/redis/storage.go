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

	// Логирование при создании клиента
	logger.Log.Info("Creating new Redis client", zap.String("address", opt.Addr), zap.Int("db", opt.DB))

	// Проверка соединения с Redis
	_, err := client.Ping().Result()
	if err != nil {
		logger.Log.Error("Failed to connect to Redis", zap.Error(err))
	} else {
		logger.Log.Info("Connected to Redis successfully")
	}

	return &Storage{Client: client}
}

func (r *Storage) Ping(ctx context.Context) error {
	// Логирование перед отправкой запроса
	logger.Log.Info("Sending ping to Redis")

	// Отправка ping и получение результата
	result, err := r.Client.Ping().Result()

	// Логирование ошибки, если она произошла
	if err != nil {
		logger.Log.Error("Failed to ping Redis", zap.Error(err))
		return err
	}

	// Логирование успешного получения ответа
	logger.Log.Info("Redis Ping Response", zap.String("response", result))
	return nil
}

// FindKeysByPattern метод для поиска ключей в Redis по шаблону
func (r *Storage) FindKeysByPattern(pattern string) ([]string, error) {
	// Логирование начала операции
	logger.Log.Info("Starting key search by pattern", zap.String("pattern", pattern))

	var cursor uint64
	var keys []string
	for {
		// Выполнение команды SCAN для поиска ключей
		k, nextCursor, err := r.Client.Scan(cursor, pattern, 0).Result()
		if err != nil {
			logger.Log.Error("Failed to scan keys", zap.Error(err))
			return nil, err
		}

		keys = append(keys, k...)
		cursor = nextCursor

		// Логирование промежуточных результатов
		logger.Log.Info("Batch of keys fetched", zap.Strings("keys", k), zap.Uint64("nextCursor", nextCursor))

		if cursor == 0 {
			break
		}
	}

	// Логирование успешного завершения операции
	logger.Log.Info("Key search completed", zap.Int("totalKeys", len(keys)))
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
