package main

import (
	"project-service/config"
	"project-service/internal/models"
	"project-service/internal/repositories"
	"log"
)

func main() {
	appConfig := config.NewAppConfig()
	db, err := repositories.InitPostgres(appConfig)
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}
	defer db.Close()
	log.Println("Running migrations...")
	db.AutoMigrate(
		&models.Project{},
	)
	log.Println("Migrations completed")
}
