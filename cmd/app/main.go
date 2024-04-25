package main

import (
	"context"
	"os"
	"os/signal"
	"stats-of/internal/app"
	"stats-of/internal/config"
	"stats-of/internal/logger"
	"syscall"
)

func main() {

	serverCtx, serverStopCtx := context.WithCancel(context.Background())
	config := config.InitConfig()
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
		logger.Log.Info("failed to run app")
		os.Exit(1)
	}

	<-serverCtx.Done()
}
