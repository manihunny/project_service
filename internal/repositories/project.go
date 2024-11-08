package repositories

import (
	"context"
	"log/slog"
	"project-service/internal/models"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/redis/go-redis/v9"
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

func (repo *ProjectRepoRedis) Create(project *models.Project) error {
	// TODO: проверка существования элемента по ключу в Redis, если есть достаём оттуда, иначе из БД и кэшируем в Redis
	err := repo.DB.Create(project).Error
	if err != nil {
		return err
	}
	ctx := context.Background()
	err = repo.RedisDB.Set(ctx, "project_"+strconv.FormatUint(uint64(project.ID), 10), project, 0).Err()
	if err != nil {
		return err
	}
	return err
}

func (repo *ProjectRepoRedis) Update(project *models.Project) error {
	err := repo.DB.Save(project).Error
	if err != nil {
		return err
	}
	ctx := context.Background()
	err = repo.RedisDB.Set(ctx, "project_"+strconv.FormatUint(uint64(project.ID), 10), project, 0).Err()
	if err != nil {
		return err
	}
	return err
}

func (repo *ProjectRepoRedis) Delete(id uint) error {
	err := repo.DB.Delete(&models.Project{}, id).Error
	if err != nil {
		return err
	}
	ctx := context.Background()
	err = repo.RedisDB.Del(ctx, "project_"+strconv.FormatUint(uint64(id), 10)).Err()
	if err != nil {
		return err
	}
	return err
}

func (repo *ProjectRepoRedis) FindByID(id uint) (*models.Project, error) {
	var project models.Project
	ctx := context.Background()
	err := repo.RedisDB.Get(ctx, "project_" + strconv.FormatUint(uint64(id), 10)).Scan(&project)
	if err != nil {
		err = repo.DB.First(&project, id).Error
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
