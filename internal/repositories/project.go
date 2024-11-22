package repositories

import (
	"context"
	"encoding/json"
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
	DB      *gorm.DB
	RedisDB *redis.Client
	Log     *slog.Logger
}

func NewProjectRepoWithRedis(db *gorm.DB, redisDB *redis.Client, logger *slog.Logger) ProjectRepository {
	return &ProjectRepoRedis{DB: db, RedisDB: redisDB, Log: logger}
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

func (repo *ProjectRepoRedis) Create(project *models.Project) error {
	// Create record in DB
	if err := repo.DB.Create(project).Error; err != nil {
		repo.Log.Error(fmt.Sprintf("Failed to create project in DB. Error: %s", err.Error()), slog.Any("project_data", project))
		return err
	}
	repo.Log.Debug("Project was created in DB", slog.Any("project_data", project))

	// Create record in Redis
	if err := setWithMarshal(repo.RedisDB, "project_"+strconv.FormatUint(uint64(project.ID), 10), project); err != nil {
		repo.Log.Warn("Couldn't create record in Redis", slog.String("error", err.Error()))
	}
	return nil
}

func (repo *ProjectRepoRedis) Update(project *models.Project) error {
	// Update record in DB
	if err := repo.DB.Save(project).Error; err != nil {
		repo.Log.Error(fmt.Sprintf("Failed to update project in DB. Error: %s", err.Error()), slog.Any("project_data", project))
		return err
	}
	repo.Log.Debug("Project was updated in DB", slog.Any("project_data", project))

	// Update record in Redis
	if err := setWithMarshal(repo.RedisDB, "project_"+strconv.FormatUint(uint64(project.ID), 10), project); err != nil {
		repo.Log.Warn("Couldn't update record in Redis", slog.String("error", err.Error()))
	}
	return nil
}

func (repo *ProjectRepoRedis) Delete(id uint) error {
	// Delete record from DB
	if err := repo.DB.Delete(&models.Project{}, id).Error; err != nil {
		repo.Log.Error(fmt.Sprintf("Failed to delete project from DB. Error: %s", err.Error()), slog.Uint64("project_id", uint64(id)))
		return err
	}
	repo.Log.Debug("Project was deleted from DB", slog.Uint64("project_id", uint64(id)))

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
		if err = repo.DB.First(&project, id).Error; err != nil {
			repo.Log.Error(fmt.Sprintf("Failed to get project. Error: %s", err.Error()), slog.Uint64("project_id", uint64(id)))
			return &project, err
		}
		repo.Log.Debug("Project was received from DB", slog.Uint64("project_id", uint64(id)))

		err = setWithMarshal(repo.RedisDB, "project_"+strconv.FormatUint(uint64(project.ID), 10), project)
		if err != nil {
			repo.Log.Warn("Couldn't create record in Redis", slog.String("error", err.Error()))
		}
	}
	return &project, err
}

func (repo *ProjectRepoRedis) FindAll() ([]models.Project, error) {
	var projects []models.Project
	if err := repo.DB.Find(&projects).Error; err != nil {
		repo.Log.Error(fmt.Sprintf("Failed to get projects. Error: %s", err.Error()))
	}
	repo.Log.Debug("Projects was received from DB")
	return projects, nil
}

// GetDB дает доступ к полю DB
func (repo *ProjectRepoRedis) GetDB() *gorm.DB {
	return repo.DB
}
