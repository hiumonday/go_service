package services

import (
	"context"
	"errors"
	"go_service/internal/dto"
	"go_service/internal/models"
	"go_service/internal/repositories"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IAssetService interface {
	CreateFolder(ctx context.Context, req *dto.CreateFolderRequest, ownerID uuid.UUID) (*models.Folder, error)
	GetFolder(ctx context.Context, folderID, userID uuid.UUID) (*models.Folder, error)
	UpdateFolder(ctx context.Context, folderID, userID uuid.UUID, req *dto.UpdateFolderRequest) (*models.Folder, error)
	DeleteFolder(ctx context.Context, folderID, userID uuid.UUID) error

	CreateNote(ctx context.Context, folderID, ownerID uuid.UUID, req *dto.CreateNoteRequest) (*models.Note, error)
	GetNote(ctx context.Context, noteID, userID uuid.UUID) (*models.Note, error)
	UpdateNote(ctx context.Context, noteID, userID uuid.UUID, req *dto.UpdateNoteRequest) (*models.Note, error)
	DeleteNote(ctx context.Context, noteID, userID uuid.UUID) error

	ShareResource(ctx context.Context, resourceID, ownerID uuid.UUID, resourceType string, req *dto.ShareRequest) error
	RevokeShare(ctx context.Context, resourceID, ownerID, targetUserID uuid.UUID, resourceType string) error

	GetTeamAssets(ctx context.Context, teamID, userID uuid.UUID) ([]models.Folder, error)
	GetUserAssets(ctx context.Context, targetUserID, currentUserID uuid.UUID) ([]models.Folder, error)
}

type AssetService struct {
	assetRepo repositories.IAssetRepository
	teamRepo  repositories.ITeamRepository
}

func NewAssetService(assetRepo repositories.IAssetRepository, teamRepo repositories.ITeamRepository) *AssetService {
	return &AssetService{
		assetRepo: assetRepo,
		teamRepo:  teamRepo,
	}
}

// CreateFolder creates a new folder
func (s *AssetService) CreateFolder(ctx context.Context, req *dto.CreateFolderRequest, ownerID uuid.UUID) (*models.Folder, error) {
	teamID, err := uuid.Parse(req.TeamID)
	if err != nil {
		return nil, err
	}

	// Check if user belongs to the team
	isMember, err := s.teamRepo.IsUserInTeam(ctx, teamID, ownerID)
	if err != nil {
		return nil, err
	}

	if !isMember {
		return nil, errors.New("user is not a member of this team")
	}

	folder := &models.Folder{
		ID:      uuid.New(),
		Name:    req.Name,
		OwnerID: ownerID,
		TeamID:  teamID,

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.assetRepo.CreateFolder(ctx, folder); err != nil {
		return nil, err
	}

	return folder, nil
}

// GetFolder retrieves a folder by ID if the user has access to it
func (s *AssetService) GetFolder(ctx context.Context, folderID, userID uuid.UUID) (*models.Folder, error) {
	folder, err := s.assetRepo.GetFolderByID(ctx, folderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("folder not found")
		}
		return nil, err
	}

	// Check if user is owner
	if folder.OwnerID == userID {
		return folder, nil
	}

	// Check if folder is shared with user
	share, err := s.assetRepo.GetShare(ctx, folderID, userID, "folder")
	if err == nil && share != nil {
		return folder, nil
	}

	// Check if user is team manager
	isManager, err := s.teamRepo.IsManager(ctx, folder.TeamID, userID)
	if err != nil {
		return nil, err
	}
	if isManager {
		return folder, nil
	}

	return nil, errors.New("you don't have access to this folder")
}

// UpdateFolder updates a folder if the user has write access
func (s *AssetService) UpdateFolder(ctx context.Context, folderID, userID uuid.UUID, req *dto.UpdateFolderRequest) (*models.Folder, error) {
	folder, err := s.assetRepo.GetFolderByID(ctx, folderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("folder not found")
		}
		return nil, err
	}

	// Check if user is owner
	if folder.OwnerID != userID {
		// Check if folder is shared with write permissions
		share, err := s.assetRepo.GetShare(ctx, folderID, userID, "folder")
		if err != nil || share == nil || share.Permission != "write" {
			return nil, errors.New("you don't have write access to this folder")
		}
	}

	// Update folder fields
	if req.Name != "" {
		folder.Name = req.Name
	}

	folder.UpdatedAt = time.Now()

	if err := s.assetRepo.UpdateFolder(ctx, folder); err != nil {
		return nil, err
	}

	return folder, nil
}

// DeleteFolder deletes a folder if the user is the owner
func (s *AssetService) DeleteFolder(ctx context.Context, folderID, userID uuid.UUID) error {
	folder, err := s.assetRepo.GetFolderByID(ctx, folderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("folder not found")
		}
		return err
	}

	// Only the owner can delete the folder
	if folder.OwnerID != userID {
		isManager, err := s.teamRepo.IsManager(ctx, folder.TeamID, userID)
		if err != nil {
			return err
		}
		if !isManager {
			return errors.New("only the folder owner or team manager can delete this folder")
		}
	}

	return s.assetRepo.DeleteFolder(ctx, folderID)
}

// CreateNote creates a new note inside a folder
func (s *AssetService) CreateNote(ctx context.Context, folderID, ownerID uuid.UUID, req *dto.CreateNoteRequest) (*models.Note, error) {
	// Check if folder exists and user has access to it
	folder, err := s.GetFolder(ctx, folderID, ownerID)
	if err != nil {
		return nil, err
	}

	// Check if user has write access to folder
	if folder.OwnerID != ownerID {
		share, err := s.assetRepo.GetShare(ctx, folderID, ownerID, "folder")
		if err != nil || share == nil || share.Permission != "write" {
			return nil, errors.New("you don't have write access to this folder")
		}
	}

	note := &models.Note{
		ID:        uuid.New(),
		FolderID:  folderID,
		OwnerID:   ownerID,
		Title:     req.Title,
		Content:   req.Content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.assetRepo.CreateNote(ctx, note); err != nil {
		return nil, err
	}

	return note, nil
}

// GetNote retrieves a note by ID if the user has access to it
func (s *AssetService) GetNote(ctx context.Context, noteID, userID uuid.UUID) (*models.Note, error) {
	note, err := s.assetRepo.GetNoteByID(ctx, noteID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("note not found")
		}
		return nil, err
	}

	// Check if user is owner
	if note.OwnerID == userID {
		return note, nil
	}

	// Check if note is shared with user
	share, err := s.assetRepo.GetShare(ctx, noteID, userID, "note")
	if err == nil && share != nil {
		return note, nil
	}

	// Check if parent folder is shared with user
	folderShare, err := s.assetRepo.GetShare(ctx, note.FolderID, userID, "folder")
	if err == nil && folderShare != nil {
		return note, nil
	}

	// Get folder to check team manager access
	folder, err := s.assetRepo.GetFolderByID(ctx, note.FolderID)
	if err != nil {
		return nil, err
	}

	// Check if user is team manager
	isManager, err := s.teamRepo.IsManager(ctx, folder.TeamID, userID)
	if err != nil {
		return nil, err
	}
	if isManager {
		return note, nil
	}

	return nil, errors.New("you don't have access to this note")
}

// UpdateNote updates a note if the user has write access
func (s *AssetService) UpdateNote(ctx context.Context, noteID, userID uuid.UUID, req *dto.UpdateNoteRequest) (*models.Note, error) {
	note, err := s.assetRepo.GetNoteByID(ctx, noteID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("note not found")
		}
		return nil, err
	}

	// Check permissions
	if note.OwnerID != userID {
		// Check direct note share
		noteShare, err := s.assetRepo.GetShare(ctx, noteID, userID, "note")
		if err == nil && noteShare != nil && noteShare.Permission == "write" {
			// User has direct write permission on note
		} else {
			// Check folder share
			folderShare, err := s.assetRepo.GetShare(ctx, note.FolderID, userID, "folder")
			if err != nil || folderShare == nil || folderShare.Permission != "write" {
				return nil, errors.New("you don't have write access to this note")
			}
		}
	}

	// Update note fields
	if req.Title != "" {
		note.Title = req.Title
	}
	if req.Content != "" {
		note.Content = req.Content
	}
	note.UpdatedAt = time.Now()

	if err := s.assetRepo.UpdateNote(ctx, note); err != nil {
		return nil, err
	}

	return note, nil
}

// DeleteNote deletes a note if the user is the owner or has appropriate permissions
func (s *AssetService) DeleteNote(ctx context.Context, noteID, userID uuid.UUID) error {
	note, err := s.assetRepo.GetNoteByID(ctx, noteID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("note not found")
		}
		return err
	}

	// Only the owner or folder owner can delete the note
	if note.OwnerID != userID {
		// Check if user is folder owner
		folder, err := s.assetRepo.GetFolderByID(ctx, note.FolderID)
		if err != nil {
			return err
		}

		if folder.OwnerID != userID {
			// Check if user is team manager
			isManager, err := s.teamRepo.IsManager(ctx, folder.TeamID, userID)
			if err != nil {
				return err
			}
			if !isManager {
				return errors.New("only the note owner, folder owner, or team manager can delete this note")
			}
		}
	}

	return s.assetRepo.DeleteNote(ctx, noteID)
}

// ShareResource shares a resource (folder or note) with another user
func (s *AssetService) ShareResource(ctx context.Context, resourceID, ownerID uuid.UUID, resourceType string, req *dto.ShareRequest) error {
	targetUserID, err := uuid.Parse(req.UserID)
	if err != nil {
		return err
	}

	// Verify resource exists and user is the owner
	var resourceOwnerID uuid.UUID
	var teamID uuid.UUID

	if resourceType == "folder" {
		folder, err := s.assetRepo.GetFolderByID(ctx, resourceID)
		if err != nil {
			return errors.New("folder not found")
		}
		resourceOwnerID = folder.OwnerID
		teamID = folder.TeamID
	} else if resourceType == "note" {
		note, err := s.assetRepo.GetNoteByID(ctx, resourceID)
		if err != nil {
			return errors.New("note not found")
		}
		resourceOwnerID = note.OwnerID

		// Get folder to check team manager access
		folder, err := s.assetRepo.GetFolderByID(ctx, note.FolderID)
		if err != nil {
			return err
		}
		teamID = folder.TeamID
	} else {
		return errors.New("invalid resource type")
	}

	// Check if user is owner or team manager
	if resourceOwnerID != ownerID {
		isManager, err := s.teamRepo.IsManager(ctx, teamID, ownerID)
		if err != nil {
			return err
		}
		if !isManager {
			return errors.New("only resource owner or team manager can share resources")
		}
	}

	// Check if target user is a member of the team
	isMember, err := s.teamRepo.IsUserInTeam(ctx, teamID, targetUserID)
	if err != nil {
		return err
	}
	if !isMember {
		return errors.New("can only share with team members")
	}

	// Create share
	share := &models.Share{
		ResourceID:   resourceID,
		ResourceType: resourceType,
		UserID:       targetUserID,
		Permission:   req.Permission,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	return s.assetRepo.CreateShare(ctx, share)
}

// RevokeShare revokes sharing permissions for a resource
func (s *AssetService) RevokeShare(ctx context.Context, resourceID, ownerID, targetUserID uuid.UUID, resourceType string) error {
	// Verify resource exists and user is the owner
	var resourceOwnerID uuid.UUID
	var teamID uuid.UUID

	if resourceType == "folder" {
		folder, err := s.assetRepo.GetFolderByID(ctx, resourceID)
		if err != nil {
			return errors.New("folder not found")
		}
		resourceOwnerID = folder.OwnerID
		teamID = folder.TeamID
	} else if resourceType == "note" {
		note, err := s.assetRepo.GetNoteByID(ctx, resourceID)
		if err != nil {
			return errors.New("note not found")
		}
		resourceOwnerID = note.OwnerID

		// Get folder to check team manager access
		folder, err := s.assetRepo.GetFolderByID(ctx, note.FolderID)
		if err != nil {
			return err
		}
		teamID = folder.TeamID
	} else {
		return errors.New("invalid resource type")
	}

	// Check if user is owner or team manager
	if resourceOwnerID != ownerID {
		isManager, err := s.teamRepo.IsManager(ctx, teamID, ownerID)
		if err != nil {
			return err
		}
		if !isManager {
			return errors.New("only resource owner or team manager can revoke shares")
		}
	}

	// Delete share
	return s.assetRepo.DeleteShare(ctx, resourceID, targetUserID, resourceType)
}

// GetTeamAssets retrieves all folders belonging to a team if the user is a manager
func (s *AssetService) GetTeamAssets(ctx context.Context, teamID, userID uuid.UUID) ([]models.Folder, error) {
	// Check if user is a team manager
	isManager, err := s.teamRepo.IsManager(ctx, teamID, userID)
	if err != nil {
		return nil, err
	}
	if !isManager {
		return nil, errors.New("only team managers can view all team assets")
	}

	// Retrieve team assets
	return s.assetRepo.GetTeamAssets(ctx, teamID)
}

// GetUserAssets retrieves all folders belonging to a specific user if the requester is a manager
func (s *AssetService) GetUserAssets(ctx context.Context, targetUserID, currentUserID uuid.UUID) ([]models.Folder, error) {
	// Get user's team ID
	teamID, err := s.teamRepo.GetUserTeamID(ctx, targetUserID)
	if err != nil {
		return nil, err
	}

	// Check if requester is a manager of the user's team
	isManager, err := s.teamRepo.IsManager(ctx, teamID, currentUserID)
	if err != nil {
		return nil, err
	}
	if !isManager && targetUserID != currentUserID {
		return nil, errors.New("only team managers can view other users' assets")
	}

	// Retrieve user assets
	return s.assetRepo.GetUserAssets(ctx, targetUserID)
}
