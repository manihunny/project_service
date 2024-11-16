package repositories

import (
	"context"
	"encoding/json"
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
	return repo.DB.Create(project).Error
}

func (repo *ProjectRepoPostgres) Update(project *models.Project) error {
	return repo.DB.Save(project).Error
}

func (repo *ProjectRepoPostgres) Delete(id uint) error {
	return repo.DB.Delete(&models.Project{}, id).Error
}

func (repo *ProjectRepoPostgres) FindByID(id uint) (*models.Project, error) {
	var project models.Project
	err := repo.DB.First(&project, id).Error
	return &project, err
}

func (repo *ProjectRepoPostgres) FindAll() ([]models.Project, error) {
	var projects []models.Project
	err := repo.DB.Find(&projects).Error
	return projects, err
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
	err := repo.DB.Create(project).Error
	if err != nil {
		return err
	}

	// Create record in Redis
	err = setWithMarshal(repo.RedisDB, "project_"+strconv.FormatUint(uint64(project.ID), 10), project)
	if err != nil {
		repo.Log.Warn("Couldn't create record in Redis", slog.String("error", err.Error()))
	}

	return nil
}

func (repo *ProjectRepoRedis) Update(project *models.Project) error {
	// Update record in DB
	err := repo.DB.Save(project).Error
	if err != nil {
		return err
	}

	// Update record in Redis
	err = setWithMarshal(repo.RedisDB, "project_"+strconv.FormatUint(uint64(project.ID), 10), project)
	if err != nil {
		repo.Log.Warn("Couldn't update record in Redis", slog.String("error", err.Error()))
	}
	return nil
}

func (repo *ProjectRepoRedis) Delete(id uint) error {
	err := repo.DB.Delete(&models.Project{}, id).Error
	if err != nil {
		return err
	}
	ctx := context.Background()
	err = repo.RedisDB.Del(ctx, "project_"+strconv.FormatUint(uint64(id), 10)).Err()
	if err != nil {
		repo.Log.Warn("Couldn't delete record in Redis", slog.String("error", err.Error()))
	}
	return nil
}

func (repo *ProjectRepoRedis) FindByID(id uint) (*models.Project, error) {
	var project models.Project

	// Try get record from Redis, otherwise get from DB
	err := getWithUnmarshal(repo.RedisDB, "project_"+strconv.FormatUint(uint64(id), 10), &project)
	if err != nil {
		err = repo.DB.First(&project, id).Error
		if err != nil {
			return nil, err
		}
		err = setWithMarshal(repo.RedisDB, "project_"+strconv.FormatUint(uint64(project.ID), 10), project)
		if err != nil {
			repo.Log.Warn("Couldn't create record in Redis", slog.String("error", err.Error()))
		}
	}

	return &project, err
}

func (repo *ProjectRepoRedis) FindAll() ([]models.Project, error) {
	var projects []models.Project
	err := repo.DB.Find(&projects).Error
	return projects, err
}

// GetDB дает доступ к полю DB
func (repo *ProjectRepoRedis) GetDB() *gorm.DB {
	return repo.DB
}
