package app

import (
	"log/slog"
	"os"
	"time"

	"github.com/gin-contrib/cors"
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
		a.Log.Error("Failed to run auth service", slog.String("error", err.Error()))
	}
}

func (a *App) setupRouter() {
	r := gin.Default()
	r.Use(gin.Recovery())
	//TODO наверное вынести и это в конфиг
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://fast-go.ru", "https://fast-go.ru"},
		AllowMethods:     []string{"POST"},
		AllowHeaders:     []string{"Authorization"},
		ExposeHeaders:    []string{"Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

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
	a.ProjectRepo = repositories.NewProjectRepoWithRedis(
		a.Database,
		a.RedisDatabase,
		a.Log.With(slog.String("repository", "project")))
}

func (a *App) setupServices() {
	a.ProjectService = services.NewProjectServiceGORM(a.ProjectRepo, a.Log.With(slog.String("service", "project")))
}

func (a *App) setupHandlersAndRoutes() {
	SetupHandlers(a.Router, a.ProjectService, a.Log)
}
