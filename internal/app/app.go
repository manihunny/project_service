package app

import (
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/redis/go-redis/v9"
	"gitlab.fast-go.ru/fast-go-team/project/config"
	"gitlab.fast-go.ru/fast-go-team/project/internal/repositories"
	"gitlab.fast-go.ru/fast-go-team/project/internal/services"
)

type App struct {
	Log           *slog.Logger
	Router        *gin.Engine
	Database      *gorm.DB
	RedisDatabase *redis.Client

	ProjectRepo repositories.ProjectRepository

	ProjectService services.ProjectService

	Config *config.Config
}

func NewApp(log *slog.Logger, config *config.Config) *App {
	return &App{
		Log:    log,
		Config: config,
	}
}

func (a *App) Initialize() {
	a.setupRouter()
	a.setupDatabase()
	a.setupRepositories()
	a.setupServices()
	a.setupHandlersAndRoutes()
}
func (a *App) Run() {
	if err := a.Router.Run(a.Config.ServiceAddress); err != nil {
		a.Log.Error("Failed to run project service", slog.String("error", err.Error()))
	}
}

func (a *App) setupRouter() {
	r := gin.Default()
	r.Use(gin.Recovery())
	a.Router = r
}

func (a *App) setupDatabase() {
	db, err := repositories.InitPostgres(a.Config)
	if err != nil {
		a.Log.Error("Failed to connect to database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	a.Database = db

	redisDB, err := repositories.InitRedis(a.Config)
	if err != nil {
		a.Log.Error("Failed to connect to Redis database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	a.RedisDatabase = redisDB
}

func (a *App) setupRepositories() {
	repoLogger := a.Log.With(slog.String("service", "project"), slog.String("module", "repository"))
	if a.Config.RedisEnabled == "true" {
		a.ProjectRepo = repositories.NewProjectRepoWithRedis(
			&repositories.ProjectRepoPostgres{DB: a.Database, Log: repoLogger},
			a.RedisDatabase,
			repoLogger,
		)
	} else {
		a.ProjectRepo = repositories.NewProjectRepoPostgres(
			a.Database,
			repoLogger,
		)
	}
	
}

func (a *App) setupServices() {
	a.ProjectService = services.NewProjectServiceGORM(a.ProjectRepo, a.Log.With(slog.String("service", "project"), slog.String("module", "service")))
}

func (a *App) setupHandlersAndRoutes() {
	SetupHandlers(a.Router, a.ProjectService, a.Log.With(slog.String("service", "project"), slog.String("module", "transport")))
}
