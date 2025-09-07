package redisclient

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type TeamCache struct {
	client *redis.Client
}

func NewTeamCache(client *redis.Client) *TeamCache {
	return &TeamCache{
		client: client,
	}
}

func (tc *TeamCache) GetTeamMembersKey(teamID uuid.UUID) string {
	return fmt.Sprintf("team:%s:members", teamID.String())
}

// GetMembers retrieves team members from Redis cache
func (tc *TeamCache) GetMembers(ctx context.Context, teamID uuid.UUID) ([]uuid.UUID, error) {
	if tc.client == nil {
		return nil, fmt.Errorf("Redis client not initialized")
	}

	key := tc.GetTeamMembersKey(teamID)

	// Get all members from Redis set
	memberStrings, err := tc.client.SMembers(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			// Cache miss
			return nil, nil
		}
		return nil, err
	}

	// Convert string IDs to UUID
	userIDs := make([]uuid.UUID, 0, len(memberStrings))
	for _, memberStr := range memberStrings {
		memberID, err := uuid.Parse(memberStr)
		if err != nil {
			log.Printf("Invalid UUID in cache: %s", memberStr)
			continue
		}
		userIDs = append(userIDs, memberID)
	}

	// Refresh expiration
	tc.client.Expire(ctx, key, 24*time.Hour)

	return userIDs, nil
}

// StoreMembers stores team members in Redis
func (tc *TeamCache) StoreMembers(ctx context.Context, teamID uuid.UUID, userIDs []uuid.UUID) error {
	if tc.client == nil {
		return fmt.Errorf("Redis client not initialized")
	}

	key := tc.GetTeamMembersKey(teamID)

	// Use pipeline for efficiency
	pipe := tc.client.Pipeline()
	pipe.Del(ctx, key)

	if len(userIDs) > 0 {
		members := make([]interface{}, len(userIDs))
		for i, id := range userIDs {
			members[i] = id.String()
		}
		pipe.SAdd(ctx, key, members...)
		pipe.Expire(ctx, key, 24*time.Hour)
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

// AddMember adds a member to the team members cache
func (tc *TeamCache) AddMember(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) error {
	if tc.client == nil {
		return fmt.Errorf("Redis client not initialized")
	}

	key := tc.GetTeamMembersKey(teamID)

	// Add member to set
	err := tc.client.SAdd(ctx, key, userID.String()).Err()
	if err != nil {
		return err
	}
	return nil
}

// RemoveMember removes a member from the team members cache
func (tc *TeamCache) RemoveMember(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) error {
	if tc.client == nil {
		return fmt.Errorf("redis client not initialized")
	}

	key := tc.GetTeamMembersKey(teamID)

	// Remove member from set
	return tc.client.SRem(ctx, key, userID.String()).Err()
}

// SMembers is a wrapper around the Redis SMembers command
func (tc *TeamCache) SMembers(ctx context.Context, key string) *redis.StringSliceCmd {
	return tc.client.SMembers(ctx, key)
}

// Expire is a wrapper around the Redis Expire command
func (tc *TeamCache) Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	return tc.client.Expire(ctx, key, expiration)
}

// Pipeline returns a Redis pipeline
func (tc *TeamCache) Pipeline() redis.Pipeliner {
	return tc.client.Pipeline()
}
