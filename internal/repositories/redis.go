package repositories

import (
	"fmt"
	"log"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/redis/go-redis/v9"
	"gitlab.fast-go.ru/fast-go-team/project/config"
)

func InitRedis(config *config.Config) (*redis.Client, error) {
	dbHost := config.RedisHost
	dbPort := config.RedisPort
	dbUser := "default"
	dbPassword := config.RedisPassword

	// Формирование строки подключения
	dbURI := fmt.Sprintf("host=%s port=%s username=%s password=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword)

	log.Printf("Connecting to Redis with URI: %s", dbURI)
	
	// Подключение к базе данных
	db := redis.NewClient(&redis.Options{
		Addr:     dbHost + ":" + dbPort,
		Username: dbUser,
		Password: dbPassword,
		DB:       0,
	})

	return db, nil
}
