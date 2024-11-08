package repositories

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"gitlab.fast-go.ru/fast-go-team/project/config"
)

func InitPostgres(config *config.Config) (*gorm.DB, error) {
	dbHost := config.DBHost
	dbPort := config.DBPort
	dbUser := config.DBUser
	dbPassword := config.DBPassword
	dbName := config.DBName

	log.Printf("DB_HOST: %s, DB_PORT: %s, DB_USER: %s, DB_NAME: %s", dbHost, dbPort, dbUser, dbName)

	// Формирование строки подключения
	dbURI := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbName, dbPassword)

	log.Printf("Connecting to DB with URI: %s", dbURI)

	// Подключение к базе данных
	db, err := gorm.Open("postgres", dbURI)
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		return nil, err
	}
	return db, nil
}
