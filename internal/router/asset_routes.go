package router

import (
	"go_service/internal/handlers"
	"go_service/internal/repositories"
	"go_service/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AssetRoutes(router *gin.RouterGroup, db *gorm.DB) {
	assetRepo := repositories.NewAssetRepository(db)
	teamRepo := repositories.NewTeamRepository(db)
	assetService := services.NewAssetService(assetRepo, teamRepo)
	assetHandler := handlers.NewAssetHandler(assetService)

	assetRouter := router.Group("/assets")
	{
		// Folder routes
		folders := assetRouter.Group("/folders")
		{
			folders.POST("", assetHandler.CreateFolder)
			folders.GET("/:folderId", assetHandler.GetFolder)
			folders.PUT("/:folderId", assetHandler.UpdateFolder)
			folders.DELETE("/:folderId", assetHandler.DeleteFolder)
		}

		// Note routes
		notes := assetRouter.Group("/notes")
		{
			notes.GET("/:noteId", assetHandler.GetNote)
			notes.PUT("/:noteId", assetHandler.UpdateNote)
			notes.DELETE("/:noteId", assetHandler.DeleteNote)
		}

		// Nested note routes
		folderNotes := folders.Group("/:folderId/notes")
		{
			folderNotes.POST("", assetHandler.CreateNote)
		}

		// Sharing routes
		folderShares := folders.Group("/:folderId/shares")
		{
			folderShares.POST("", assetHandler.ShareFolder)
			folderShares.DELETE("/:userId", assetHandler.RevokeFolderShare)
		}

		noteShares := notes.Group("/:noteId/shares")
		{
			noteShares.POST("", assetHandler.ShareNote)
			noteShares.DELETE("/:userId", assetHandler.RevokeNoteShare)
		}
	}

	managerRouter := router.Group("/manager")
	{
		managerRouter.GET("/teams/:teamId/assets", assetHandler.GetTeamAssets)
		managerRouter.GET("/users/:userId/assets", assetHandler.GetUserAssets)
	}
}
