package services

import (
	"context"
	"fmt"
	"go_service/internal/models"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/machinebox/graphql"
)

// IUserService định nghĩa các phương thức cần có
type IUserService interface {
	GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
	GetUsersByIDs(ctx context.Context, userIDs []uuid.UUID) ([]models.User, error)
	UserExists(ctx context.Context, userID uuid.UUID) (bool, error)
	CreateUser(ctx context.Context, username, email, password, role string) (*models.User, error)
}

// UserService implements IUserService
type UserService struct {
	client  *graphql.Client
	baseURL string
}

// NewUserService creates a new user service client
func NewUserService() *UserService {
	baseURL := os.Getenv("USER_SERVICE_URL")

	log.Printf("Initializing UserService with URL: %s", baseURL)
	client := graphql.NewClient(baseURL)

	return &UserService{
		client:  client,
		baseURL: baseURL,
	}
}

// GetUserByID fetches a user by their ID
func (s *UserService) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	userIDStr := userID.String()
	log.Printf("Fetching user with ID: %s", userIDStr)

	// Create GraphQL request
	req := graphql.NewRequest(`
        query GetUser($userId: ID!) {
            user(userId: $userId) {
                userId
                username
                email
                role
            }
        }
    `)

	// Set variables
	req.Var("userId", userIDStr)

	// Add timeout to context if not already set
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Execute request
	var response struct {
		User *models.User `json:"user"`
	}

	if err := s.client.Run(ctx, req, &response); err != nil {
		log.Printf("GraphQL request failed: %v", err)
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	// Check if user was found
	if response.User == nil && response.User.ID == uuid.Nil {
		log.Printf("User not found with ID: %s %s", userIDStr, response)
		return nil, fmt.Errorf("user not found with ID: %s", userIDStr)
	}

	return response.User, nil
}

// GetUsersByIDs fetches multiple users by their IDs
func (s *UserService) GetUsersByIDs(ctx context.Context, userIDs []uuid.UUID) ([]models.User, error) {
	if len(userIDs) == 0 {
		return []models.User{}, nil
	}

	// Convert UUID to string array
	userIDStrs := make([]string, len(userIDs))
	for i, id := range userIDs {
		userIDStrs[i] = id.String()
	}

	log.Printf("Fetching %d users", len(userIDs))

	// Create GraphQL request
	req := graphql.NewRequest(`
        query GetUsers($userIds: [ID!]!) {
            users(userIds: $userIds) {
              	userId
                username
                email
                role
            }
        }
    `)

	// Set variables
	req.Var("userIds", userIDStrs)

	// Add timeout to context if not already set
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second) // Longer timeout for multiple users
	defer cancel()

	// Execute request
	var response struct {
		Users []models.User `json:"users"`
	}

	if err := s.client.Run(ctx, req, &response); err != nil {
		log.Printf("GraphQL request failed: %v", err)
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}

	return response.Users, nil
}

// UserExists checks if a user exists
func (s *UserService) UserExists(ctx context.Context, userID uuid.UUID) (bool, error) {
	user, err := s.GetUserByID(ctx, userID)
	if err != nil {
		// Check if error is "user not found" or a different error
		if fmt.Sprintf("%v", err) == fmt.Sprintf("user not found with ID: %s", userID.String()) {
			return false, nil
		}
		return false, err
	}
	return user != nil, nil
}

// CreateUser creates a new user using GraphQL mutation
func (s *UserService) CreateUser(ctx context.Context, username, email, password, role string) (*models.User, error) {
	// Create GraphQL request
	req := graphql.NewRequest(`
        mutation CreateUser($input: CreateUserInput!) {
            createUser(input: $input) {
                code
                success
                message
                user {
                    userId
                    username
                    email
                    role
                }
            }
        }
    `)

	// Set variables
	req.Var("input", map[string]interface{}{
		"username": username,
		"email":    email,
		"password": password,
		"role":     role,
	})

	// Add timeout context if not already set
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Execute request
	var response struct {
		CreateUser struct {
			Code    int         `json:"code"`
			Success bool        `json:"success"`
			Message string      `json:"message"`
			User    models.User `json:"user"`
		} `json:"createUser"`
	}

	if err := s.client.Run(ctx, req, &response); err != nil {
		log.Printf("GraphQL request failed: %v", err)
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Check if user was created successfully
	if !response.CreateUser.Success {
		log.Printf("Failed to create user: %s", response.CreateUser.Message)
		return nil, fmt.Errorf("user creation failed: %s", response.CreateUser.Message)
	}

	return &response.CreateUser.User, nil
}
