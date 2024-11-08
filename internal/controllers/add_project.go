package controllers

import (
	"project-service/internal/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *ProjectHandler) CreateProject(c *gin.Context) {
	var projectDTO dto.ProjectDTO
	if err := c.ShouldBindJSON(&projectDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Service.CreateProject(&projectDTO); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, projectDTO)
}
