package controllers

import (
	"log/slog"
	"project-service/internal/services"
)

type ProjectHandler struct {
	Service services.ProjectService
	Log     *slog.Logger
}

func NewProjectHandler(projectService services.ProjectService, log *slog.Logger) *ProjectHandler {
	return &ProjectHandler{Service: projectService, Log: log}
}
