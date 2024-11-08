package config

import (
	"fmt"
	"log/slog"

	"github.com/caarlos0/env"
)

type Config struct {
	ServiceAddress         string `env:"SERVICE_ADDRESS" envDefault:"localhost:80"`
	LogLevel               string `env:"SERVICE_LOG_LEVEL" envDefault:"info"`
	DBHost                 string `env:"DB_HOST" envDefault:"localhost"`
	DBPort                 string `env:"DB_PORT" envDefault:"5432"`
	DBUser                 string `env:"DB_USER" envDefault:"postgres"`
	DBPassword             string `env:"DB_PASSWORD" envDefault:"mysecretpassword"`
	DBName                 string `env:"DB_NAME" envDefault:"mydb"`
	RedisDBHost            string `env:"REDIS_DB_HOST" envDefault:"localhost"`
	RedisDBPort            string `env:"REDIS_DB_PORT" envDefault:"5432"`
	RedisDBUser            string `env:"REDIS_DB_USER" envDefault:"redis"`
	RedisDBPassword        string `env:"REDIS_DB_PASSWORD" envDefault:"mysecretpassword"`
	AuthServiceGRPCAddress string `env:"AUTH_SERVICE_GRPC_ADDRESS" envDefault:"localhost:50051"`
}

func NewAppConfig() *Config {
	var AppConfig Config
	err := env.Parse(&AppConfig)
	if err != nil {
		slog.Error("Failed to process env var", slog.String("error", err.Error()))
		return nil
	}
	fmt.Printf("Loaded config: %+v\n", AppConfig)
	return &AppConfig
}
