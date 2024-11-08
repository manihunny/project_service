package app

import (
	"log/slog"
	"project-service/internal/controllers"
	"project-service/internal/middleware"
	"project-service/internal/services"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupHandlers(r *gin.Engine, projectService services.ProjectService, log *slog.Logger) {
	projectHandler := &controllers.ProjectHandler{Service: projectService, Log: log}

	// Оповещение docker-compose о том, что контейнер готов к работе
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "OK",
		})
	})

	v1 := r.Group("/project/api/v1").Use(middleware.Auth())
	{
		v1.POST("/", projectHandler.CreateProject)
		v1.PUT("/:id", projectHandler.UpdateProject)
		v1.DELETE("/:id", projectHandler.DeleteProject)
		v1.GET("/:id", projectHandler.GetProjectByID)
		v1.GET("/", projectHandler.GetProjects)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
