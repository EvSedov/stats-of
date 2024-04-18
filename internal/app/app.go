package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"stats-of/internal/config"
	"stats-of/internal/entities"
	"stats-of/internal/healthz"
	"stats-of/internal/logger"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

var appInfo = &entities.AppInfo{
	Name:         "stats-of",
	BuildVersion: "0.0.1",
	BuildTime:    "Wed, 10 Apr 2024 22:25:51",
	GitTag:       "no git tag",
	GitHash:      "no git hash",
}

type App struct {
	server *http.Server
}

func New(config *config.Config) (*App, error) {
	// Логирование начала создания нового экземпляра приложения
	logger.Log.Info("Initializing new application instance", zap.Int("ServerPort", config.ServerPort))

	const (
		defaultHTTPServerWriteTimeout = time.Second * 15
		defaultHTTPServerReadTimeout  = time.Second * 15
	)

	app := new(App)
	mux := http.NewServeMux()

	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/healthz", healthz.MakeHandler(appInfo))

	app.server = &http.Server{
		Handler:      mux,
		Addr:         ":" + strconv.Itoa(config.ServerPort),
		WriteTimeout: defaultHTTPServerWriteTimeout,
		ReadTimeout:  defaultHTTPServerReadTimeout,
	}

	// Логирование завершения инициализации сервера
	logger.Log.Info("HTTP server configured", zap.String("address", app.server.Addr),
		zap.Duration("writeTimeout", app.server.WriteTimeout),
		zap.Duration("readTimeout", app.server.ReadTimeout))

	return app, nil
}

func (a *App) Run() error {
	// Логирование попытки запуска сервера
	logger.Log.Info("Starting HTTP server", zap.String("address", a.server.Addr))

	err := a.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		// Логирование ошибки, если сервер не был закрыт нормально
		logger.Log.Error("HTTP server stopped with error", zap.Error(err))
		return fmt.Errorf("server was stopped with error: %w", err)
	}

	// Логирование нормального закрытия сервера
	logger.Log.Info("HTTP server stopped gracefully")
	return nil
}

func (a *App) stop(ctx context.Context) error {
	err := a.server.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("server was shutdown with error: %w", err)
	}

	return nil
}

func (a *App) GracefulStop(serverCtx context.Context, sig <-chan os.Signal, serverStopCtx context.CancelFunc) {
	<-sig
	var timeOut = 30 * time.Second
	shutdownCtx, shutdownStopCtx := context.WithTimeout(serverCtx, timeOut)

	go func() {
		<-shutdownCtx.Done()
		if shutdownCtx.Err() == context.DeadlineExceeded {
			os.Exit(1)
		}
	}()

	err := a.stop(shutdownCtx)
	if err != nil {
		os.Exit(1)
	}

	serverStopCtx()
	shutdownStopCtx()
}
