package repositories

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/redis/go-redis/v9"
	"gitlab.fast-go.ru/fast-go-team/project/internal/models"
)

// Структура для работы с Postgres
type ProjectRepoPostgres struct {
	DB  *gorm.DB
	Log *slog.Logger
}

func NewProjectRepoPostgres(db *gorm.DB, logger *slog.Logger) ProjectRepository {
	return &ProjectRepoPostgres{DB: db, Log: logger}
}

func (repo *ProjectRepoPostgres) Create(project *models.Project) error {
	if err := repo.DB.Create(project).Error; err != nil {
		repo.Log.Error(fmt.Sprintf("Failed to create project in DB. Error: %s", err.Error()), slog.Any("project_data", project))
		return err
	}
	repo.Log.Debug("Project was created in DB", slog.Any("project_data", project))
	return nil
}

func (repo *ProjectRepoPostgres) Update(project *models.Project) error {
	if err := repo.DB.Save(project).Error; err != nil {
		repo.Log.Error(fmt.Sprintf("Failed to update project in DB. Error: %s", err.Error()), slog.Any("project_data", project))
		return err
	}
	repo.Log.Debug("Project was updated in DB", slog.Any("project_data", project))
	return nil
}

func (repo *ProjectRepoPostgres) Delete(id uint) error {
	if err := repo.DB.Delete(&models.Project{}, id).Error; err != nil {
		repo.Log.Error(fmt.Sprintf("Failed to delete project from DB. Error: %s", err.Error()), slog.Uint64("project_id", uint64(id)))
		return err
	}
	repo.Log.Debug("Project was deleted from DB", slog.Uint64("project_id", uint64(id)))
	return nil
}

func (repo *ProjectRepoPostgres) FindByID(id uint) (*models.Project, error) {
	var project models.Project
	if err := repo.DB.First(&project, id).Error; err != nil {
		repo.Log.Error(fmt.Sprintf("Failed to get project. Error: %s", err.Error()), slog.Uint64("project_id", uint64(id)))
	}
	repo.Log.Debug("Project was received from DB", slog.Uint64("project_id", uint64(id)))
	return &project, nil
}

func (repo *ProjectRepoPostgres) FindAll() ([]models.Project, error) {
	var projects []models.Project
	if err := repo.DB.Find(&projects).Error; err != nil {
		repo.Log.Error(fmt.Sprintf("Failed to get projects. Error: %s", err.Error()))
	}
	repo.Log.Debug("Projects was received from DB")
	return projects, nil
}

// GetDB дает доступ к полю DB
func (repo *ProjectRepoPostgres) GetDB() *gorm.DB {
	return repo.DB
}

// Структура для работы с Redis (надстройка для Postgres)
type ProjectRepoRedis struct {
	DBRepo  *ProjectRepoPostgres
	RedisDB *redis.Client
	Log     *slog.Logger
}

func NewProjectRepoWithRedis(dbRepo *ProjectRepoPostgres, redisDB *redis.Client, logger *slog.Logger) ProjectRepository {
	return &ProjectRepoRedis{DBRepo: dbRepo, RedisDB: redisDB, Log: logger}
}

func (repo *ProjectRepoRedis) Create(project *models.Project) error {
	// Create record in DB
	err := repo.DBRepo.Create(project)
	if err != nil {
		return err
	}

	// Create record in Redis
	if err := setWithMarshal(repo.RedisDB, "project_"+strconv.FormatUint(uint64(project.ID), 10), project); err != nil {
		repo.Log.Warn("Couldn't create record in Redis", slog.String("error", err.Error()))
	}
	return nil
}

func (repo *ProjectRepoRedis) Update(project *models.Project) error {
	// Update record in DB
	err := repo.DBRepo.Update(project)
	if err != nil {
		return err
	}

	// Update record in Redis
	if err := setWithMarshal(repo.RedisDB, "project_"+strconv.FormatUint(uint64(project.ID), 10), project); err != nil {
		repo.Log.Warn("Couldn't update record in Redis", slog.String("error", err.Error()))
	}
	return nil
}

func (repo *ProjectRepoRedis) Delete(id uint) error {
	// Delete record from DB
	err := repo.DBRepo.Delete(id)
	if err != nil {
		return err
	}

	// Delete record from Redis
	ctx := context.Background()
	if err := repo.RedisDB.Del(ctx, "project_"+strconv.FormatUint(uint64(id), 10)).Err(); err != nil {
		repo.Log.Warn("Couldn't delete record in Redis", slog.String("error", err.Error()))
	}
	return nil
}

func (repo *ProjectRepoRedis) FindByID(id uint) (*models.Project, error) {
	var project models.Project

	// Try get record from Redis, otherwise get from DB and create record in Redis
	err := getWithUnmarshal(repo.RedisDB, "project_"+strconv.FormatUint(uint64(id), 10), &project)
	if err != nil {
		project, err := repo.DBRepo.FindByID(id)
		if err != nil {
			return project, err
		}

		err = setWithMarshal(repo.RedisDB, "project_"+strconv.FormatUint(uint64(project.ID), 10), project)
		if err != nil {
			repo.Log.Warn("Couldn't create record in Redis", slog.String("error", err.Error()))
		}
	}
	return &project, err
}

func (repo *ProjectRepoRedis) FindAll() ([]models.Project, error) {
	return repo.DBRepo.FindAll()
}

// GetDB дает доступ к полю DB
func (repo *ProjectRepoRedis) GetDB() *gorm.DB {
	return repo.DBRepo.DB
}
