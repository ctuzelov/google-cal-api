package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/ctuzelov/google-cal-api/cmd/server"
	"github.com/ctuzelov/google-cal-api/internal/handler"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
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

func MustLoadConfig(path string) *oauth2.Config {
	config, err := FetchConfig(path)
	if err != nil {
		panic(err)
	}
	return config
}

func FetchConfig(path string) (*oauth2.Config, error) {
	if path == "" {
		return nil, errors.New("path is empty")
	}

	credentials, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config, err := google.ConfigFromJSON(credentials, calendar.CalendarScope)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func SetupLogger() *slog.Logger {
	log := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	return log
}
