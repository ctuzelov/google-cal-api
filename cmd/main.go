package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/ctuzelov/google-cal-api/cmd/server"
	"github.com/ctuzelov/google-cal-api/internal/handler"
)

const (
	path = "config/credentials.json"
)

func main() {
	logger := SetupLogger()
	config := MustLoadConfig(path)

	client_handler := handler.New(config, logger)

	srv := new(server.Server)
	err := srv.Run("8080", client_handler.InitRoutes())

	if err != nil {
		logger.Error("error while running http server", err)
		return
	}

	logger.Info("server started")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop
	logger.Info("Stopping application", slog.String("signal", sign.String()))

	srv.MustShutdown(context.Background())
	logger.Info("Application stopped")
}
