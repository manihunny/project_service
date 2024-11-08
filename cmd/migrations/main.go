package main

import (
	"log"

	"gitlab.fast-go.ru/fast-go-team/project/config"
	"gitlab.fast-go.ru/fast-go-team/project/internal/models"
	"gitlab.fast-go.ru/fast-go-team/project/internal/repositories"
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
