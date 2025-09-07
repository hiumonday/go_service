package router

import (
	"go_service/internal/handlers"
	"go_service/internal/middleware"
	"go_service/internal/repositories"
	"go_service/internal/services"
	"go_service/pkg/kafka"
	"go_service/pkg/redisclient"

	"github.com/gin-gonic/gin"

	"gorm.io/gorm"
)

func SetupRouter(router *gin.Engine, db *gorm.DB, producer *kafka.Producer, redis_client *redisclient.TeamCache) {
	//Repositories
	teamRepo := repositories.NewTeamRepository(db)

	//Services
	teamService := services.NewTeamService(teamRepo, producer, redis_client)

	// Handlers
	teamHandler := handlers.NewTeamHandler(teamService)
	importHandler := handlers.NewImportHandler()

	//v1 api
	v1 := router.Group("/api/v1")

	protectedRoutes := v1.Group("/")
	protectedRoutes.Use(middleware.AuthMiddleware(db))

	// Set up all routes
	TeamRoutes(protectedRoutes, teamHandler)
	ImportRoutes(protectedRoutes, importHandler)
}
