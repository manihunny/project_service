package app

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gitlab.fast-go.ru/fast-go-team/project/config"
	"gitlab.fast-go.ru/fast-go-team/project/internal/controllers"
	"gitlab.fast-go.ru/fast-go-team/project/internal/middleware"
	"gitlab.fast-go.ru/fast-go-team/project/internal/services"
)

func SetupHandlers(r *gin.Engine, projectService services.ProjectService, log *slog.Logger) {
	projectHandler := &controllers.ProjectHandler{Service: projectService, Log: log}

	// Оповещение docker-compose о том, что контейнер готов к работе
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "OK",
		})
	})
	
	v1 := r.Group("/project/api/v1")
	{
		v1.POST("/", projectHandler.CreateProject)
		v1.PUT("/:id", projectHandler.UpdateProject)
		v1.DELETE("/:id", projectHandler.DeleteProject)
		v1.GET("/:id", projectHandler.GetProjectByID)
		v1.GET("/", projectHandler.GetProjects)
	}

	if config.NewAppConfig().AuthEnabled == "true" {
		v1.Use(middleware.Auth())
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
