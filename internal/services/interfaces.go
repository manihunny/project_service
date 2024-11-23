package services

import (
	"gitlab.fast-go.ru/fast-go-team/project/internal/dto"
	"gitlab.fast-go.ru/fast-go-team/project/internal/models"
	"gitlab.fast-go.ru/fast-go-team/project/internal/repositories"
)

type ProjectService interface {
	CreateProject(projectDTO *dto.ProjectDTO) error
	UpdateProject(id uint, projectDTO *dto.ProjectDTO) error
	DeleteProject(id uint) error
	GetProjectByID(id uint) (*models.Project, error)
	GetProjects() ([]models.Project, error)
	GetRepo() repositories.ProjectRepository
}
