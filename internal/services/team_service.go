package services

import (
	"context"
	"errors"
	"go_service/internal/models"
	"go_service/internal/repositories"
	"go_service/pkg/kafka"
	"go_service/pkg/redisclient"
	"log"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type ITeamService interface {
	CreateTeam(ctx context.Context, teamName string, userIDs []uuid.UUID, creatorID uuid.UUID) (*models.Team, []uuid.UUID, error)
	AddMembersToTeam(ctx context.Context, teamID uuid.UUID, userIDs []uuid.UUID, currentUserID uuid.UUID) (int, int, []uuid.UUID, error)
	GetTeamMembers(ctx context.Context, teamID uuid.UUID) ([]models.User, error)
	RemoveMemberFromTeam(ctx context.Context, teamID uuid.UUID, userID uuid.UUID, currentUserID uuid.UUID) error
	RemoveManagerFromTeam(ctx context.Context, teamID uuid.UUID, managerID uuid.UUID, currentUserID uuid.UUID) error
}

type TeamService struct {
	repo        repositories.ITeamRepository
	userService *UserService
	producer    *kafka.Producer
	redisClient *redisclient.TeamCache
}

func NewTeamService(repo repositories.ITeamRepository, producer *kafka.Producer, redisClient *redisclient.TeamCache) *TeamService {
	return &TeamService{
		repo:        repo,
		userService: NewUserService(),
		producer:    producer,
		redisClient: redisClient,
	}
}

// Creates a new team and adds members
func (s *TeamService) CreateTeam(ctx context.Context, teamName string, userIDs []uuid.UUID, creatorID uuid.UUID) (*models.Team, []uuid.UUID, error) {
	team := &models.Team{TeamName: teamName}
	if err := s.repo.CreateTeam(ctx, team); err != nil {
		return nil, nil, err
	}

	// Add creator as leader
	leaderRoster := models.Roster{TeamID: team.ID, UserID: creatorID, Role: "MAIN_MANAGER"}
	if err := s.repo.AddMemberToTeam(ctx, leaderRoster); err != nil {
		return nil, nil, err
	}

	if len(userIDs) == 0 {
		return team, nil, nil
	}

	_, _, failedMembers, err := s.AddMembersToTeam(ctx, team.ID, userIDs, creatorID)
	if err != nil {
		log.Printf("Error adding members to team: %v", err)
	}

	return team, failedMembers, nil
}

// adds members to a team with validation
func (s *TeamService) AddMembersToTeam(ctx context.Context, teamID uuid.UUID, userIDs []uuid.UUID, currentUserID uuid.UUID) (int, int, []uuid.UUID, error) {
	// Sửa lỗi điều kiện kiểm tra quyền
	currentUserRole, err := s.repo.GetUserRoleInTeam(ctx, teamID, currentUserID)
	if err != nil || (currentUserRole != "MANAGER" && currentUserRole != "MAIN_MANAGER") {
		return 0, len(userIDs), userIDs, errors.New("you are not a manager")
	}

	addedCount := 0
	failedCount := 0
	failedMembers := []uuid.UUID{}
	validRosters := make([]models.Roster, 0, len(userIDs))

	// Lấy danh sách members hiện có để kiểm tra nhanh
	existingMembers, err := s.repo.GetTeamMembers(ctx, teamID)
	if err != nil {
		return 0, len(userIDs), userIDs, err
	}

	existingMap := make(map[uuid.UUID]bool)
	for _, member := range existingMembers {
		existingMap[member.UserID] = true
	}

	for _, userID := range userIDs {
		if existingMap[userID] {
			failedMembers = append(failedMembers, userID)
			failedCount++
			continue
		}

		validRosters = append(validRosters, models.Roster{TeamID: teamID, UserID: userID, Role: "MEMBER"})
	}

	// Nếu không có user hợp lệ để thêm
	if len(validRosters) == 0 {
		return 0, failedCount, failedMembers, nil
	}

	for _, roster := range validRosters {
		if err := s.repo.AddMemberToTeam(ctx, roster); err != nil {
			failedMembers = append(failedMembers, roster.UserID)
			failedCount++
		}
	}

	//gửi event kafka
	if s.producer != nil {
		for _, roster := range validRosters {
			err := s.producer.SendTeamEvent(
				kafka.EventMemberAdded,
				teamID,
				currentUserID,
				roster.UserID,
			)
			if err != nil {
				log.Printf("Failed to send Kafka event for member addition: %v", err)
				// Continue processing even if Kafka event fails
			}
		}
	}

	addedCount = len(validRosters)
	return addedCount, failedCount, failedMembers, nil
}
func (s *TeamService) GetTeamMembers(ctx context.Context, teamID uuid.UUID) ([]models.User, error) {
	// Try to get from Redis cache first
	if s.redisClient != nil {
		cachedMembers, err := s.redisClient.GetMembers(ctx, teamID)
		if err == nil && len(cachedMembers) > 0 {
			// Fetch user details from user service
			users, err := s.userService.GetUsersByIDs(ctx, cachedMembers)
			if err == nil && len(users) > 0 {
				log.Printf("Fetched user details for cached members: %v", users)
				return users, nil
			}
			log.Printf("Failed to fetch user details for cached members: %v", err)
		} else if err != nil && err != redis.Nil {
			log.Printf("Redis error: %v", err)
		}
	}
	// If cache miss or error, fetch from DB
	rosters, err := s.repo.GetTeamMembers(ctx, teamID)
	if err != nil {
		return nil, err
	}
	userIDs := make([]uuid.UUID, len(rosters))
	for i, roster := range rosters {
		userIDs[i] = roster.UserID
	}
	// Update Redis cache
	if s.redisClient != nil && len(userIDs) > 0 {
		if err := s.redisClient.StoreMembers(ctx, teamID, userIDs); err != nil {
			log.Printf("Failed to update Redis cache: %v", err)
		}
	}
	// Fetch user details from user service
	users, err := s.userService.GetUsersByIDs(ctx, userIDs)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *TeamService) RemoveMemberFromTeam(ctx context.Context, teamID uuid.UUID, targetID uuid.UUID, currentUserID uuid.UUID) error {
	// Check if current user is a manager
	currentUserRole, err := s.repo.GetUserRoleInTeam(ctx, teamID, currentUserID)
	if err != nil || (currentUserRole != "MANAGER" && currentUserRole != "MAIN_MANAGER") {
		return errors.New("you are not a manager")
	}
	targetUserRole, err := s.repo.GetUserRoleInTeam(ctx, teamID, targetID)
	if err != nil || targetUserRole == "MAIN_MANAGER" || targetUserRole == "MANAGER" {
		return errors.New("cannot remove a manager or main manager")
	}

	if s.producer != nil {
		err := s.producer.SendTeamEvent(
			kafka.EventMemberRemoved,
			teamID,
			currentUserID,
			targetID,
		)
		if err != nil {
			log.Printf("Failed to send Kafka event for member removal: %v", err)
		}
	}

	return s.repo.RemoveMemberFromTeam(ctx, teamID, targetID)

}

func (s *TeamService) RemoveManagerFromTeam(ctx context.Context, teamID uuid.UUID, managerID uuid.UUID, currentUserID uuid.UUID) error {
	// Check if current user is a MAIN_MANAGER
	if managerID == currentUserID {
		return errors.New("you cannot remove yourself")
	}

	currentUserRole, err := s.repo.GetUserRoleInTeam(ctx, teamID, currentUserID)
	if err != nil || currentUserRole != "MAIN_MANAGER" {
		return errors.New("you are not the main manager")
	}
	targetUserRole, err := s.repo.GetUserRoleInTeam(ctx, teamID, managerID)
	if err != nil || targetUserRole != "MANAGER" {
		return errors.New("target user is not a manager")
	}
	if s.producer != nil {
		err := s.producer.SendTeamEvent(
			kafka.EventManagerRemoved,
			teamID,
			currentUserID,
			managerID,
		)
		if err != nil {
			log.Printf("Failed to send Kafka event for manager removal: %v", err)
		}
	}

	return s.repo.RemoveMemberFromTeam(ctx, teamID, managerID)
}
