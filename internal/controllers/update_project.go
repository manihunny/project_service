package controllers

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.fast-go.ru/fast-go-team/project/internal/dto"
)

func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.Log.Error("Failed to update project. Error: Invalid project ID", slog.String("method", c.Request.Method), slog.Int("code", http.StatusBadRequest), slog.String("url", c.Request.URL.Path), slog.String("client", c.ClientIP()), slog.String("project_id", idStr))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	var projectDTO dto.ProjectDTO
	if err := c.ShouldBindJSON(&projectDTO); err != nil {
		h.Log.Error(fmt.Sprintf("Failed to update project. Error: %s", err.Error()), slog.String("method", c.Request.Method), slog.Int("code", http.StatusBadRequest), slog.String("url", c.Request.URL.Path), slog.String("client", c.ClientIP()))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Service.UpdateProject(uint(id), &projectDTO); err != nil {
		h.Log.Error(fmt.Sprintf("Failed to update project. Error: %s", err.Error()), slog.String("method", c.Request.Method), slog.Int("code", http.StatusInternalServerError), slog.String("url", c.Request.URL.Path), slog.String("client", c.ClientIP()), slog.Uint64("project_id", id), slog.Any("project_data", projectDTO))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.Log.Info("Project was updated", slog.String("method", c.Request.Method), slog.Int("code", http.StatusOK), slog.String("url", c.Request.URL.Path), slog.String("client", c.ClientIP()), slog.Uint64("project_id", id), slog.Any("project_data", projectDTO))

	c.JSON(http.StatusOK, gin.H{"message": "Project deleted successfully", "id": id})
}
