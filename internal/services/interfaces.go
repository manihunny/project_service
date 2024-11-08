package services

import (
	"project-service/internal/dto"
	"project-service/internal/models"
	"project-service/internal/repositories"
)

type ProjectService interface {
	CreateProject(projectDTO *dto.ProjectDTO) error
	UpdateProject(id uint, projectDTO *dto.ProjectDTO) error
	DeleteProject(id uint) error
	GetProjectByID(id uint) (*models.Project, error)
	GetProjects() ([]models.Project, error)
	GetRepo() repositories.ProjectRepository
}
