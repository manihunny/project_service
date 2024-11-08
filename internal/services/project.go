package services

import (
	"log/slog"
	"project-service/internal/dto"
	"project-service/internal/models"
	"project-service/internal/repositories"
)

type ProjectServiceGORM struct {
	Repo repositories.ProjectRepository
	Log  *slog.Logger
}

func NewProjectServiceGORM(repo repositories.ProjectRepository, logger *slog.Logger) ProjectService {
	return &ProjectServiceGORM{Repo: repo, Log: logger}
}

func (projectService *ProjectServiceGORM) CreateProject(projectDTO *dto.ProjectDTO) error {
	var project models.Project

	// Маппинг данных из DTO в модель
	if err := projectDTO.Map(&project); err != nil {
		return err
	}
	return projectService.Repo.Create(&project)
}

func (projectService *ProjectServiceGORM) UpdateProject(id uint, projectDTO *dto.ProjectDTO) error {
	project, err := projectService.Repo.FindByID(id)
	if err != nil {
		return err
	}
	// Маппинг данных из DTO в модель
	if err := projectDTO.Map(project); err != nil {
		return err
	}
	return projectService.Repo.Update(project)
}

func (projectService *ProjectServiceGORM) DeleteProject(id uint) error {
	return projectService.Repo.Delete(id)
}

func (projectService *ProjectServiceGORM) GetProjectByID(id uint) (*models.Project, error) {
	return projectService.Repo.FindByID(id)
}

func (projectService *ProjectServiceGORM) GetProjects() ([]models.Project, error) {
	return projectService.Repo.FindAll()
}

func (s *ProjectServiceGORM) GetRepo() repositories.ProjectRepository {
	return s.Repo
}
