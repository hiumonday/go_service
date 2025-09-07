package repositories

import (
	"context"
	"go_service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ITeamRepository interface {
	GetTeamByID(ctx context.Context, teamID uuid.UUID) (*models.Team, error)
	GetUserRoleInTeam(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) (string, error)
	CreateTeam(ctx context.Context, team *models.Team) error
	AddMemberToTeam(ctx context.Context, roster models.Roster) error
	IsUserInTeam(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) (bool, error)
	GetTeamMembers(ctx context.Context, teamID uuid.UUID) ([]models.Roster, error)
	RemoveMemberFromTeam(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) error
}

type TeamRepository struct {
	db *gorm.DB
}

func NewTeamRepository(db *gorm.DB) *TeamRepository {
	return &TeamRepository{db: db}
}

func (r *TeamRepository) GetTeamByID(ctx context.Context, teamID uuid.UUID) (*models.Team, error) {
	var team models.Team
	if err := r.db.WithContext(ctx).First(&team, "id = ?", teamID).Error; err != nil {
		return nil, err
	}
	return &team, nil
}

func (r *TeamRepository) GetUserRoleInTeam(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) (string, error) {
	var member models.Roster
	err := r.db.WithContext(ctx).
		Select("role").
		Where("team_id = ? AND user_id = ?", teamID, userID).
		First(&member).Error

	if err != nil {
		return "", err
	}
	return member.Role, nil
}

// inserts a new team into the database
func (r *TeamRepository) CreateTeam(ctx context.Context, team *models.Team) error {
	return r.db.WithContext(ctx).Create(team).Error
}

// inserts a roster entry for a team member
func (r *TeamRepository) AddMemberToTeam(ctx context.Context, roster models.Roster) error {
	return r.db.WithContext(ctx).Create(&roster).Error
}

func (r *TeamRepository) IsUserInTeam(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Roster{}).
		Where("team_id = ? AND user_id = ?", teamID, userID).
		Count(&count).Error

	return count > 0, err
}

func (r *TeamRepository) GetTeamMembers(ctx context.Context, teamID uuid.UUID) ([]models.Roster, error) {
	var members []models.Roster
	err := r.db.WithContext(ctx).
		Select("user_id").
		Where("team_id = ?", teamID).
		Find(&members).Error

	return members, err
}

func (r *TeamRepository) RemoveMemberFromTeam(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("team_id = ? AND user_id = ?", teamID, userID).
		Delete(&models.Roster{}).Error
}
