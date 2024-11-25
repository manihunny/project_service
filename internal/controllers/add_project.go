package controllers

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.fast-go.ru/fast-go-team/project/internal/dto"
)

func (h *ProjectHandler) CreateProject(c *gin.Context) {
	var projectDTO dto.ProjectDTO
	if err := c.ShouldBindJSON(&projectDTO); err != nil {
		h.Log.Error(fmt.Sprintf("Failed to create project. Error: %s", err.Error()), slog.String("method", c.Request.Method), slog.Int("code", http.StatusBadRequest), slog.String("url", c.Request.URL.Path), slog.String("client", c.ClientIP()))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Service.CreateProject(&projectDTO); err != nil {
		h.Log.Error(fmt.Sprintf("Failed to create project. Error: %s", err.Error()), slog.String("method", c.Request.Method), slog.Int("code", http.StatusInternalServerError), slog.String("url", c.Request.URL.Path), slog.String("client", c.ClientIP()), slog.Any("project_data", projectDTO))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.Log.Info("Project was created", slog.String("method", c.Request.Method), slog.Int("code", http.StatusCreated), slog.String("url", c.Request.URL.Path), slog.String("client", c.ClientIP()), slog.Any("project_data", projectDTO))

	c.JSON(http.StatusCreated, projectDTO)
}
