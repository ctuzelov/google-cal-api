package main

import (
	"errors"
	"log/slog"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

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
