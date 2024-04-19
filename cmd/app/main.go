package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"stats-of/internal/app"
	"stats-of/internal/config"
	"stats-of/internal/logger"
	"syscall"

	"go.uber.org/zap"
)

func RunApp() error {
	// Загрузка конфигурации из переменных окружения
	appConfig, err := config.LoadFromEnv()
	if err != nil {
		logger.Log.Error("Failed to load configuration from environment", zap.Error(err))
		return fmt.Errorf("failed to load configuration from environment: %w", err)
	}
	logger.Log.Info("Configuration loaded successfully")

	// Инициализация контекста сервера
	serverCtx, serverStopCtx := context.WithCancel(context.Background())
	defer serverStopCtx() // Убедитесь, что функция отмены вызывается для предотвращения утечек контекста

	// Создание нового экземпляра приложения
	myApp, err := app.New(appConfig)
	if err != nil {
		logger.Log.Error("Failed to initialize the application", zap.Error(err))
		return fmt.Errorf("failed to initialize the application: %w", err)
	}
	logger.Log.Info("Application initialized successfully")

	// Подготовка к обработке системных сигналов
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		logger.Log.Info("Waiting for system signals")
		s := <-sig
		logger.Log.Info("Received system signal", zap.String("signal", s.String()))
		myApp.GracefulStop(serverCtx, sig, serverStopCtx)
	}()

	// Запуск основного цикла приложения
	logger.Log.Info("Running the application...")
	if err := myApp.Run(); err != nil {
		logger.Log.Error("Failed to run the application", zap.Error(err))
		return fmt.Errorf("failed to run the application: %w", err)
	}

	logger.Log.Info("Application is now running. Waiting for shutdown signal...")
	<-serverCtx.Done()
	logger.Log.Info("Application shutdown completed successfully")
	return nil
}

func main() {
	logger.InitLogger()
	logger.Log.Info("Starting application and reading configuration...")

	if err := RunApp(); err != nil {
		logger.Log.Fatal("Error occurred", zap.Error(err))
	}

	// if err := storagetestsutils.HandleCsvToDb(); err != nil {
	// 	logger.Log.Fatal("Ошибка при обработке данных CSV", zap.Error(err))
	// }
}
