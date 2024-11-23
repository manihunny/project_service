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
	RedisHost              string `env:"REDIS_HOST" envDefault:"localhost"`
	RedisPort              string `env:"REDIS_PORT" envDefault:"5432"`
	RedisPassword          string `env:"REDIS_PASSWORD" envDefault:"mysecretpassword"`
	RedisEnabled           string `env:"REDIS_ENABLED" envDefault:"true"`
	AuthServiceGRPCAddress string `env:"AUTH_SERVICE_GRPC_ADDRESS" envDefault:"localhost:50051"`
	AuthEnabled            string `env:"AUTH_ENABLED" envDefault:"true"`
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
