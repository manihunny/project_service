package main

import (
	"log/slog"
	"os"

	"gitlab.fast-go.ru/fast-go-team/project/config"
	"gitlab.fast-go.ru/fast-go-team/project/internal/app"
)

func main() {
	//загрузка конфигурации и установка уровня логирования
	appConfig := config.NewAppConfig()

	//logger
	logger := configureLogger(appConfig.LogLevel)

	//запуск сервиса
	application := app.NewApp(logger, appConfig)
	application.Initialize()
	application.Run()
}

func configureLogger(logLevel string) *slog.Logger {
	level := new(slog.LevelVar)

	switch logLevel {
	case "debug":
		level.Set(slog.LevelDebug)
	case "info":
		level.Set(slog.LevelInfo)
	case "warn":
		level.Set(slog.LevelWarn)
	case "error":
		level.Set(slog.LevelError)
	default:
		level.Set(slog.LevelError)
	}

	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))
}
