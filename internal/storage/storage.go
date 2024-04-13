package storage

import (
	"os"
	"strconv"

	"stats-of/internal/logger"

	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

// RedisService структура для работы с Redis
type RedisService struct {
	Client *redis.Client
}

// NewRedisService функция для создания нового экземпляра RedisService
func NewRedisService() *RedisService {
	// Загрузка переменных окружения из файла .env
	err := godotenv.Load()
	if err != nil {
		logger.Log.Fatal("Ошибка при загрузке файла .env", zap.Error(err))
	}

	addr := os.Getenv("REDIS_ADDR")
	password := os.Getenv("REDIS_PASSWORD")
	dbStr := os.Getenv("REDIS_DB")

	db, err := strconv.Atoi(dbStr)
	if err != nil {
		logger.Log.Fatal("Ошибка при преобразовании REDIS_DB в число", zap.Error(err))
	}

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	pong, err := client.Ping().Result()
	if err != nil {
		logger.Log.Fatal("Ошибка при подключении к Redis", zap.Error(err))
	} else {
		logger.Log.Info("Успешное подключение к Redis", zap.String("pong", pong))
	}

	return &RedisService{Client: client}
}

// FindKeysByPattern метод для поиска ключей в Redis по шаблону
func (rs *RedisService) FindKeysByPattern(pattern string) ([]string, error) {
	var cursor uint64
	var keys []string
	for {
		k, nextCursor, err := rs.Client.Scan(cursor, pattern, 0).Result()
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

func (rs *RedisService) FindKeyByGetRequest(key string) (string, error) {
	// Используем метод Get клиента Redis для получения значения по ключу
	result, err := rs.Client.Get(key).Result()
	if err != nil {
		// Если произошла ошибка, возвращаем пустую строку и саму ошибку
		return "", err
	}
	// Возвращаем результат и nil в качестве ошибки, если всё прошло успешно
	return result, nil
}
