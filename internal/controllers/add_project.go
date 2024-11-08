package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.fast-go.ru/fast-go-team/project/internal/dto"
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
