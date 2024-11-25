package repositories

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"gitlab.fast-go.ru/fast-go-team/project/config"
)

// Init connection to specific DB (with DBName)
func InitPostgres(config *config.Config) (*gorm.DB, error) {
	dbHost := config.DBHost
	dbPort := config.DBPort
	dbUser := config.DBUser
	dbPassword := config.DBPassword
	dbName := config.DBName

	dbURI := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)
	return connectDB(dbURI)
}

// Init connection to database server, not specific DB (without DBName)
func InitPostgresServer(config *config.Config) (*gorm.DB, error) {
	dbHost := config.DBHost
	dbPort := config.DBPort
	dbUser := config.DBUser
	dbPassword := config.DBPassword

	dbURI := fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword)
	return connectDB(dbURI)
}

func connectDB(connString string) (*gorm.DB, error) {
	log.Printf("Connecting to DB Server with URI: %s", connString)
	db, err := gorm.Open("postgres", connString)
	if err != nil {
		log.Printf("Failed to connect to database server: %v", err)
		return nil, err
	}
	return db, nil
}