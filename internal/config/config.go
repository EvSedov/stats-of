package config

import (
	"fmt"
	"os"
	"stats-of/internal/logger"
	"strconv"

	"go.uber.org/zap"
)

const (
	defaultServerPort = "8080"
)

type Config struct {
	ServerPort int
}

func LoadFromEnv() (*Config, error) {
	// Логирование начала загрузки конфигурации
	logger.Log.Info("Loading configuration from environment variables")

	conf := &Config{}
	var err error
	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		// Логирование использования порта по умолчанию при отсутствии переменной окружения
		logger.Log.Info("SERVER_PORT not set, using default", zap.String("defaultServerPort", defaultServerPort))
		serverPort = defaultServerPort
	}

	conf.ServerPort, err = strconv.Atoi(serverPort)
	if err != nil {
		// Логирование ошибки при преобразовании порта из строки в число
		logger.Log.Error("Failed to parse SERVER_PORT as integer", zap.String("serverPort", serverPort), zap.Error(err))
		return nil, fmt.Errorf("failed to parse %s as int: %w", serverPort, err)
	}

	// Логирование успешной загрузки конфигурации
	logger.Log.Info("Configuration loaded successfully", zap.Int("serverPort", conf.ServerPort))
	return conf, nil
}
