package repositories

import (
	"context"
	"go_service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// IAssetRepository defines the interface for asset-related database operations.
type IAssetRepository interface {
	// Folder methods
	CreateFolder(ctx context.Context, folder *models.Folder) error
	GetFolderByID(ctx context.Context, folderID uuid.UUID) (*models.Folder, error)
	UpdateFolder(ctx context.Context, folder *models.Folder) error
	DeleteFolder(ctx context.Context, folderID uuid.UUID) error

	// Note methods
	CreateNote(ctx context.Context, note *models.Note) error
	GetNoteByID(ctx context.Context, noteID uuid.UUID) (*models.Note, error)
	UpdateNote(ctx context.Context, note *models.Note) error
	DeleteNote(ctx context.Context, noteID uuid.UUID) error

	// Share methods
	CreateShare(ctx context.Context, share *models.Share) error
	GetShare(ctx context.Context, resourceID, userID uuid.UUID, resourceType string) (*models.Share, error)
	DeleteShare(ctx context.Context, resourceID, userID uuid.UUID, resourceType string) error

	// Manager methods
	GetTeamAssets(ctx context.Context, teamID uuid.UUID) ([]models.Folder, error)
	GetUserAssets(ctx context.Context, userID uuid.UUID) ([]models.Folder, error)
}

// AssetRepository implements IAssetRepository.
type AssetRepository struct {
	db *gorm.DB
}

// NewAssetRepository creates a new instance of AssetRepository.
func NewAssetRepository(db *gorm.DB) IAssetRepository {
	return &AssetRepository{db: db}
}

// --- Folder Methods Implementation ---

func (r *AssetRepository) CreateFolder(ctx context.Context, folder *models.Folder) error {
	return r.db.WithContext(ctx).Create(folder).Error
}

func (r *AssetRepository) GetFolderByID(ctx context.Context, folderID uuid.UUID) (*models.Folder, error) {
	var folder models.Folder
	err := r.db.WithContext(ctx).Preload("Notes").First(&folder, "id = ?", folderID).Error
	if err != nil {
		return nil, err
	}
	return &folder, nil
}

func (r *AssetRepository) UpdateFolder(ctx context.Context, folder *models.Folder) error {
	return r.db.WithContext(ctx).Save(folder).Error
}

func (r *AssetRepository) DeleteFolder(ctx context.Context, folderID uuid.UUID) error {
	// Using transaction to ensure both folder and its notes are deleted.
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("folder_id = ?", folderID).Delete(&models.Note{}).Error; err != nil {
			return err
		}
		if err := tx.Where("id = ?", folderID).Delete(&models.Folder{}).Error; err != nil {
			return err
		}
		return nil
	})
}

// --- Note Methods Implementation ---

func (r *AssetRepository) CreateNote(ctx context.Context, note *models.Note) error {
	return r.db.WithContext(ctx).Create(note).Error
}

func (r *AssetRepository) GetNoteByID(ctx context.Context, noteID uuid.UUID) (*models.Note, error) {
	var note models.Note
	err := r.db.WithContext(ctx).First(&note, "id = ?", noteID).Error
	if err != nil {
		return nil, err
	}
	return &note, nil
}

func (r *AssetRepository) UpdateNote(ctx context.Context, note *models.Note) error {
	return r.db.WithContext(ctx).Save(note).Error
}

func (r *AssetRepository) DeleteNote(ctx context.Context, noteID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Note{}, "id = ?", noteID).Error
}

// --- Share Methods Implementation ---

func (r *AssetRepository) CreateShare(ctx context.Context, share *models.Share) error {
	return r.db.WithContext(ctx).Create(share).Error
}

func (r *AssetRepository) GetShare(ctx context.Context, resourceID, userID uuid.UUID, resourceType string) (*models.Share, error) {
	var share models.Share
	err := r.db.WithContext(ctx).Where("resource_id = ? AND user_id = ? AND resource_type = ?", resourceID, userID, resourceType).First(&share).Error
	if err != nil {
		return nil, err
	}
	return &share, nil
}

func (r *AssetRepository) DeleteShare(ctx context.Context, resourceID, userID uuid.UUID, resourceType string) error {
	return r.db.WithContext(ctx).Where("resource_id = ? AND user_id = ? AND resource_type = ?", resourceID, userID, resourceType).Delete(&models.Share{}).Error
}

// --- Manager Methods Implementation ---

func (r *AssetRepository) GetTeamAssets(ctx context.Context, teamID uuid.UUID) ([]models.Folder, error) {
	var folders []models.Folder
	// This query can be complex, for now, we get all folders owned by the team.
	// A more advanced version would include shared assets.
	err := r.db.WithContext(ctx).Preload("Notes").Where("team_id = ?", teamID).Find(&folders).Error
	return folders, err
}

func (r *AssetRepository) GetUserAssets(ctx context.Context, userID uuid.UUID) ([]models.Folder, error) {
	var folders []models.Folder
	// This is a simplified query. A complete implementation would require a more complex query
	// involving joins with the shares table.
	err := r.db.WithContext(ctx).Preload("Notes").Where("owner_id = ?", userID).Find(&folders).Error
	return folders, err
}
