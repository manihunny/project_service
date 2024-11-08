package repositories

import (
	"github.com/jinzhu/gorm"
	"gitlab.fast-go.ru/fast-go-team/project/internal/models"
)

type ProjectRepository interface {
	Create(project *models.Project) error
	Update(project *models.Project) error
	Delete(id uint) error
	FindByID(id uint) (*models.Project, error)
	FindAll() ([]models.Project, error)
	GetDB() *gorm.DB
}
