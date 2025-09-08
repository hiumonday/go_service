package handlers

import (
	"go_service/internal/dto"
	"go_service/internal/services"
	"go_service/pkg/responses"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AssetHandler struct {
	service services.IAssetService
}

func NewAssetHandler(service services.IAssetService) *AssetHandler {
	return &AssetHandler{service: service}
}

// Folder Handlers
func (h *AssetHandler) CreateFolder(c *gin.Context) {
	var req dto.CreateFolderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.Error(c, http.StatusBadRequest, err, "Invalid request format")
		return
	}

	userID, _ := c.Get("user_id")
	folder, err := h.service.CreateFolder(c.Request.Context(), &req, userID.(uuid.UUID))
	if err != nil {
		responses.Error(c, http.StatusInternalServerError, err, "Failed to create folder")
		return
	}
	responses.JSON(c, http.StatusCreated, gin.H{"success": true, "data": folder})
}

func (h *AssetHandler) GetFolder(c *gin.Context) {
	folderID, err := uuid.Parse(c.Param("folderId"))
	if err != nil {
		responses.Error(c, http.StatusBadRequest, err, "Invalid folder ID format")
		return
	}
	userID, _ := c.Get("user_id")
	folder, err := h.service.GetFolder(c.Request.Context(), folderID, userID.(uuid.UUID))
	if err != nil {
		responses.Error(c, http.StatusNotFound, err, "Folder not found or access denied")
		return
	}
	responses.JSON(c, http.StatusOK, gin.H{"success": true, "data": folder})
}

func (h *AssetHandler) UpdateFolder(c *gin.Context) {
	folderID, err := uuid.Parse(c.Param("folderId"))
	if err != nil {
		responses.Error(c, http.StatusBadRequest, err, "Invalid folder ID format")
		return
	}
	var req dto.UpdateFolderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.Error(c, http.StatusBadRequest, err, "Invalid request format")
		return
	}
	userID, _ := c.Get("user_id")
	folder, err := h.service.UpdateFolder(c.Request.Context(), folderID, userID.(uuid.UUID), &req)
	if err != nil {
		responses.Error(c, http.StatusForbidden, err, "Update failed or access denied")
		return
	}
	responses.JSON(c, http.StatusOK, gin.H{"success": true, "data": folder})
}

func (h *AssetHandler) DeleteFolder(c *gin.Context) {
	folderID, err := uuid.Parse(c.Param("folderId"))
	if err != nil {
		responses.Error(c, http.StatusBadRequest, err, "Invalid folder ID format")
		return
	}
	userID, _ := c.Get("user_id")
	err = h.service.DeleteFolder(c.Request.Context(), folderID, userID.(uuid.UUID))
	if err != nil {
		responses.Error(c, http.StatusForbidden, err, "Delete failed or access denied")
		return
	}
	responses.JSON(c, http.StatusOK, gin.H{"success": true, "message": "Folder deleted"})
}

// Note Handlers
func (h *AssetHandler) CreateNote(c *gin.Context) {
	folderID, err := uuid.Parse(c.Param("folderId"))
	if err != nil {
		responses.Error(c, http.StatusBadRequest, err, "Invalid folder ID format")
		return
	}
	var req dto.CreateNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.Error(c, http.StatusBadRequest, err, "Invalid request format")
		return
	}
	userID, _ := c.Get("user_id")
	note, err := h.service.CreateNote(c.Request.Context(), folderID, userID.(uuid.UUID), &req)
	if err != nil {
		responses.Error(c, http.StatusForbidden, err, "Create note failed or access denied")
		return
	}
	responses.JSON(c, http.StatusCreated, gin.H{"success": true, "data": note})
}

func (h *AssetHandler) GetNote(c *gin.Context) {
	noteID, err := uuid.Parse(c.Param("noteId"))
	if err != nil {
		responses.Error(c, http.StatusBadRequest, err, "Invalid note ID format")
		return
	}
	userID, _ := c.Get("user_id")
	note, err := h.service.GetNote(c.Request.Context(), noteID, userID.(uuid.UUID))
	if err != nil {
		responses.Error(c, http.StatusNotFound, err, "Note not found or access denied")
		return
	}
	responses.JSON(c, http.StatusOK, gin.H{"success": true, "data": note})
}

func (h *AssetHandler) UpdateNote(c *gin.Context) {
	noteID, err := uuid.Parse(c.Param("noteId"))
	if err != nil {
		responses.Error(c, http.StatusBadRequest, err, "Invalid note ID format")
		return
	}
	var req dto.UpdateNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.Error(c, http.StatusBadRequest, err, "Invalid request format")
		return
	}
	userID, _ := c.Get("user_id")
	note, err := h.service.UpdateNote(c.Request.Context(), noteID, userID.(uuid.UUID), &req)
	if err != nil {
		responses.Error(c, http.StatusForbidden, err, "Update note failed or access denied")
		return
	}
	responses.JSON(c, http.StatusOK, gin.H{"success": true, "data": note})
}

func (h *AssetHandler) DeleteNote(c *gin.Context) {
	noteID, err := uuid.Parse(c.Param("noteId"))
	if err != nil {
		responses.Error(c, http.StatusBadRequest, err, "Invalid note ID format")
		return
	}
	userID, _ := c.Get("user_id")
	err = h.service.DeleteNote(c.Request.Context(), noteID, userID.(uuid.UUID))
	if err != nil {
		responses.Error(c, http.StatusForbidden, err, "Delete note failed or access denied")
		return
	}
	responses.JSON(c, http.StatusOK, gin.H{"success": true, "message": "Note deleted"})
}

// Sharing Handlers
func (h *AssetHandler) ShareFolder(c *gin.Context) {
	folderID, err := uuid.Parse(c.Param("folderId"))
	if err != nil {
		responses.Error(c, http.StatusBadRequest, err, "Invalid folder ID format")
		return
	}
	var req dto.ShareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.Error(c, http.StatusBadRequest, err, "Invalid request format")
		return
	}
	userID, _ := c.Get("user_id")
	err = h.service.ShareResource(c.Request.Context(), folderID, userID.(uuid.UUID), "folder", &req)
	if err != nil {
		responses.Error(c, http.StatusForbidden, err, "Share folder failed or access denied")
		return
	}
	responses.JSON(c, http.StatusOK, gin.H{"success": true, "message": "Folder shared"})
}

func (h *AssetHandler) RevokeFolderShare(c *gin.Context) {
	folderID, err := uuid.Parse(c.Param("folderId"))
	if err != nil {
		responses.Error(c, http.StatusBadRequest, err, "Invalid folder ID format")
		return
	}
	targetUserID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		responses.Error(c, http.StatusBadRequest, err, "Invalid user ID format")
		return
	}
	userID, _ := c.Get("user_id")
	err = h.service.RevokeShare(c.Request.Context(), folderID, userID.(uuid.UUID), targetUserID, "folder")
	if err != nil {
		responses.Error(c, http.StatusForbidden, err, "Revoke folder share failed or access denied")
		return
	}
	responses.JSON(c, http.StatusOK, gin.H{"success": true, "message": "Folder share revoked"})
}

func (h *AssetHandler) ShareNote(c *gin.Context) {
	noteID, err := uuid.Parse(c.Param("noteId"))
	if err != nil {
		responses.Error(c, http.StatusBadRequest, err, "Invalid note ID format")
		return
	}
	var req dto.ShareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.Error(c, http.StatusBadRequest, err, "Invalid request format")
		return
	}
	userID, _ := c.Get("user_id")
	err = h.service.ShareResource(c.Request.Context(), noteID, userID.(uuid.UUID), "note", &req)
	if err != nil {
		responses.Error(c, http.StatusForbidden, err, "Share note failed or access denied")
		return
	}
	responses.JSON(c, http.StatusOK, gin.H{"success": true, "message": "Note shared"})
}

func (h *AssetHandler) RevokeNoteShare(c *gin.Context) {
	noteID, err := uuid.Parse(c.Param("noteId"))
	if err != nil {
		responses.Error(c, http.StatusBadRequest, err, "Invalid note ID format")
		return
	}
	targetUserID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		responses.Error(c, http.StatusBadRequest, err, "Invalid user ID format")
		return
	}
	userID, _ := c.Get("user_id")
	err = h.service.RevokeShare(c.Request.Context(), noteID, userID.(uuid.UUID), targetUserID, "note")
	if err != nil {
		responses.Error(c, http.StatusForbidden, err, "Revoke note share failed or access denied")
		return
	}
	responses.JSON(c, http.StatusOK, gin.H{"success": true, "message": "Note share revoked"})
}

// Manager-only APIs
func (h *AssetHandler) GetTeamAssets(c *gin.Context) {
	teamID, err := uuid.Parse(c.Param("teamId"))
	if err != nil {
		responses.Error(c, http.StatusBadRequest, err, "Invalid team ID format")
		return
	}
	userID, _ := c.Get("user_id")
	assets, err := h.service.GetTeamAssets(c.Request.Context(), teamID, userID.(uuid.UUID))
	if err != nil {
		responses.Error(c, http.StatusForbidden, err, "Get team assets failed or access denied")
		return
	}
	responses.JSON(c, http.StatusOK, gin.H{"success": true, "data": assets})
}

func (h *AssetHandler) GetUserAssets(c *gin.Context) {
	targetUserID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		responses.Error(c, http.StatusBadRequest, err, "Invalid user ID format")
		return
	}
	userID, _ := c.Get("user_id")
	assets, err := h.service.GetUserAssets(c.Request.Context(), targetUserID, userID.(uuid.UUID))
	if err != nil {
		responses.Error(c, http.StatusForbidden, err, "Get user assets failed or access denied")
		return
	}
	responses.JSON(c, http.StatusOK, gin.H{"success": true, "data": assets})
}
