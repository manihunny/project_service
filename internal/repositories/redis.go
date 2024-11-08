package repositories

import (
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/redis/go-redis/v9"
	"gitlab.fast-go.ru/fast-go-team/project/config"
)

func InitRedis(config *config.Config) (*redis.Client, error) {
	dbHost := config.RedisDBHost
	dbPort := config.RedisDBPort
	dbUser := config.RedisDBUser
	dbPassword := config.RedisDBPassword

	// Подключение к базе данных
	db := redis.NewClient(&redis.Options{
		Addr:     dbHost + dbPort,
		Username: dbUser,
		Password: dbPassword,
		DB:       0,
	})

	return db, nil
}
