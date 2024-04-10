package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"stats-of/internal/app"
	"stats-of/internal/config"
)

func main() {
	fmt.Println("reading config...")
	config, err := config.LoadFromEnv()
	if err != nil {
		fmt.Println("failed to read config")
		os.Exit(1)
	}

	serverCtx, serverStopCtx := context.WithCancel(context.Background())
	app, err := app.New(config)
	if err != nil {
		fmt.Println("failed to read config")
		os.Exit(1)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		app.GracefulStop(serverCtx, sig, serverStopCtx)
	}()

	err = app.Run()
	if err != nil {
		fmt.Println("failed to read config")
		os.Exit(1)
	}

	<-serverCtx.Done()
}
