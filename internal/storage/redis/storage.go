package redis

import (
	"context"

	"github.com/go-redis/redis"
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

func (r *Storage) Ping(ctx context.Context) (err error) {
	err = r.Ping(ctx)
	return err
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
	// Используем метод Get клиента Redis для получения значения по ключу
	result, err := r.Client.Get(key).Result()
	if err != nil {
		// Если произошла ошибка, возвращаем пустую строку и саму ошибку
		return "", err
	}
	// Возвращаем результат и nil в качестве ошибки, если всё прошло успешно
	return result, nil
}
