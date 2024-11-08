package controllers

import (
	"project-service/internal/dto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	// Получаем ID проекта из параметров URL и преобразуем его в int
	id, _ := strconv.Atoi(c.Param("id"))

	var projectDTO dto.ProjectDTO
	if err := c.ShouldBindJSON(&projectDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Service.UpdateProject(uint(id), &projectDTO); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Проект обновлен", "id": id})
}
