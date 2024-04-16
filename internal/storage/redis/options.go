package redis

import (
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"os"
	"stats-of/internal/logger"
	"strconv"
)

type (
	Options struct {
		Addr     string
		Password string
		DB       int
	}
)

func CreateOptions() (opt *Options, err error) {
	// Загрузка переменных окружения из файла .env
	err = godotenv.Load()
	if err != nil {
		logger.Log.Fatal("Ошибка при загрузке файла .env", zap.Error(err))
		return nil, err
	}

	addr := os.Getenv("REDIS_ADDR")
	password := os.Getenv("REDIS_PASSWORD")
	rdbStr := os.Getenv("REDIS_DB")

	rdb, err := strconv.Atoi(rdbStr)
	if err != nil {
		logger.Log.Fatal("Ошибка при преобразовании REDIS_DB в число", zap.Error(err))
		return nil, err
	}

	return &Options{
		Addr:     addr,
		Password: password,
		DB:       rdb,
	}, nil
}
