package main

import (
	"context"
	"os"
	"os/signal"
	"stats-of/internal/app"
	"stats-of/internal/config"
	"stats-of/internal/logger"
	"stats-of/internal/storage"
	"syscall"

	"go.uber.org/zap"
)

func main() {
	logger.InitLogger()
	logger.Log.Info("reading config...")
	// ---------------------------------------------------
	redisService := storage.NewRedisService()
	keys, err := redisService.FindKeysByPattern("*pattern*")
	if err != nil {
		logger.Log.Fatal("Ошибка при поиске ключей", zap.Error(err))
	}
	for _, key := range keys {
		logger.Log.Info("Найден ключ", zap.String("key", key))
	}

	key := "chat:5481:user:121190002:type:CHANNEL"

	// Получение значения по ключу
	value, err := storage.NewRedisService().FindKeyByGetRequest(key)
	if err != nil {
		logger.Log.Fatal("Ошибка при получении значения из Redis", zap.Error(err))
	}

	// Вывод полученного значения
	logger.Log.Info("Полученное значение", zap.String("value", value))
	// ---------------------------------------------------

	config, err := config.LoadFromEnv()
	if err != nil {
		logger.Log.Info("failed to read config")
		os.Exit(1)
	}

	serverCtx, serverStopCtx := context.WithCancel(context.Background())
	app, err := app.New(config)
	if err != nil {
		logger.Log.Info("failed to read config")
		os.Exit(1)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		app.GracefulStop(serverCtx, sig, serverStopCtx)
	}()

	err = app.Run()
	if err != nil {
		logger.Log.Info("failed to read config")
		os.Exit(1)
	}

	<-serverCtx.Done()
}
