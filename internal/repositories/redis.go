package repositories

import (
	"context"
	"encoding/json"
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

func setWithMarshal(redisClient *redis.Client, key string, value interface{}) error {
	ctx := context.Background()
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return redisClient.Set(ctx, key, data, 0).Err()
}

func getWithUnmarshal(redisClient *redis.Client, key string, dest interface{}) error {
	ctx := context.Background()
	var data []byte
	err := redisClient.Get(ctx, key).Scan(&data)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &dest)
}
