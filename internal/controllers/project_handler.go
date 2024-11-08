package controllers

import (
	"log/slog"

	"gitlab.fast-go.ru/fast-go-team/project/internal/services"
)

type ProjectHandler struct {
	Service services.ProjectService
	Log     *slog.Logger
}

func NewProjectHandler(projectService services.ProjectService, log *slog.Logger) *ProjectHandler {
	return &ProjectHandler{Service: projectService, Log: log}
}
