package services

import (
	"fmt"
	"log/slog"

	"gitlab.fast-go.ru/fast-go-team/project/internal/dto"
	"gitlab.fast-go.ru/fast-go-team/project/internal/models"
	"gitlab.fast-go.ru/fast-go-team/project/internal/repositories"
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
		projectService.Log.Error(fmt.Sprintf("Failed to create project. Error: %s", err.Error()), slog.Any("project_data", project))
		return err
	}
	if err := projectService.Repo.Create(&project); err != nil {
		projectService.Log.Error(fmt.Sprintf("Failed to create project. Error: %s", err.Error()), slog.Any("project_data", project))
		return err
	}
	projectService.Log.Debug("Project was created", slog.Any("project_data", project))
	return nil
}

func (projectService *ProjectServiceGORM) UpdateProject(id uint, projectDTO *dto.ProjectDTO) error {
	project, err := projectService.Repo.FindByID(id)
	if err != nil {
		projectService.Log.Error(fmt.Sprintf("Failed to update project. Error: %s", err.Error()), slog.Any("project_data", project))
		return err
	}
	// Маппинг данных из DTO в модель
	if err = projectDTO.Map(project); err != nil {
		projectService.Log.Error(fmt.Sprintf("Failed to update project. Error: %s", err.Error()), slog.Any("project_data", project))
		return err
	}
	if err = projectService.Repo.Update(project); err != nil {
		projectService.Log.Error(fmt.Sprintf("Failed to update project. Error: %s", err.Error()), slog.Any("project_data", project))
		return err
	}
	projectService.Log.Debug("Project was updated", slog.Any("project_data", project))
	return nil
}

func (projectService *ProjectServiceGORM) DeleteProject(id uint) error {
	if err := projectService.Repo.Delete(id); err != nil {
		projectService.Log.Error(fmt.Sprintf("Failed to delete project. Error: %s", err.Error()), slog.Any("project_id", id))
		return err
	}
	projectService.Log.Debug("Project was deleted", slog.Any("project_id", id))
	return nil
}

func (projectService *ProjectServiceGORM) GetProjectByID(id uint) (*models.Project, error) {
	project, err := projectService.Repo.FindByID(id)
	if err != nil {
		projectService.Log.Error(fmt.Sprintf("Failed to get project. Error: %s", err.Error()), slog.Any("project_id", id))
	}
	projectService.Log.Debug("Project was received from DB", slog.Uint64("project_id", uint64(id)))
	return project, err
}

func (projectService *ProjectServiceGORM) GetProjects() ([]models.Project, error) {
	projects, err := projectService.Repo.FindAll()
	if err != nil {
		projectService.Log.Error(fmt.Sprintf("Failed to get all projects. Error: %s", err.Error()))
	}
	projectService.Log.Debug("Projects was received from DB")
	return projects, err
}

func (s *ProjectServiceGORM) GetRepo() repositories.ProjectRepository {
	return s.Repo
}
