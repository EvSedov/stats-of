package storage

import (
	"os"
	"stats-of/internal/logger"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func InitDb() *redis.Client {
	// Загрузка переменных из файла .env
	err := godotenv.Load()
	if err != nil {
		logger.Log.Fatal("Ошибка при загрузке файла .env", zap.Error(err))
	}

	// Получение конфигурации из переменных окружения
	addr := os.Getenv("REDIS_ADDR")
	password := os.Getenv("REDIS_PASSWORD")
	dbStr := os.Getenv("REDIS_DB") // Получаем значение как строку

	// Преобразуем значение DB из строки в число
	db, err := strconv.Atoi(dbStr)
	if err != nil {
		logger.Log.Fatal("Ошибка при преобразовании REDIS_DB в число", zap.Error(err))
	}

	// Создаем новый клиент Redis с использованием переменных окружения
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db, // Используем преобразованное значение
	})

	// Выполняем команду PING
	pong, err := rdb.Ping().Result()
	if err != nil {
		logger.Log.Info("Ошибка при подключении к Redis", zap.Error(err))
	}

	// Logging with zap
	logger.Log.Info("Ответ от Redis:", zap.String("response", pong))

	return rdb
}

// Функция для поиска ключей в Redis по шаблону
func FindKeysByPattern(rdb *redis.Client, pattern string) ([]string, error) {

	// Используем Scan для поиска ключей по шаблону без блокировки базы данных
	var cursor uint64
	var keys []string
	for {
		var err error
		var k []string
		k, cursor, err = rdb.Scan(cursor, pattern, 0).Result()
		if err != nil {
			return nil, err // Возвращаем ошибку, если что-то пошло не так
		}
		keys = append(keys, k...)
		if cursor == 0 { // Если курсор равен 0, значит обход всех ключей завершен
			break
		}
	}

	return keys, nil // Возвращаем найденные ключи
}
