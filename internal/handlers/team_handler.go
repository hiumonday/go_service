package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"go_service/internal/dto"
	"go_service/internal/services"
	"go_service/pkg/responses"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TeamHandler struct {
	service services.ITeamService
}

func NewTeamHandler(service services.ITeamService) *TeamHandler {
	return &TeamHandler{
		service: service,
	}
}

// POST /teams (only managers can create teams)
func (h *TeamHandler) CreateTeam(c *gin.Context) {
	role, exists := c.Get("role")
	if !exists || strings.ToUpper(role.(string)) != "MANAGER" {
		responses.Error(c, http.StatusForbidden, fmt.Errorf("insufficient permissions"), "Only managers can create teams")
		return
	}

	var req dto.CreateTeamReq
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.Error(c, http.StatusBadRequest, err, "Invalid request format")
		return
	}

	userID, _ := c.Get("user_id")
	creatorID := userID.(uuid.UUID)

	team, failedMembers, err := h.service.CreateTeam(c.Request.Context(), req.TeamName, req.UserIDs, creatorID)
	if err != nil {
		responses.Error(c, http.StatusBadRequest, err, "Failed to create team")
		return
	}

	response := gin.H{
		"success": true,
		"message": "Team created successfully",
		"data": gin.H{
			"team":          team,
			"failedMembers": failedMembers,
		},
	}
	responses.JSON(c, http.StatusCreated, response)
}

// POST /teams/:teamId/members
func (h *TeamHandler) AddMemberToTeam(c *gin.Context) {
	teamIDStr := c.Param("teamId")
	teamID, err := uuid.Parse(teamIDStr)
	if err != nil {
		responses.Error(c, http.StatusBadRequest, err, "Invalid team ID format")
		return
	}

	var req dto.AddMemberReq
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.Error(c, http.StatusBadRequest, err, "Invalid request format")
		return
	}

	userID, _ := c.Get("user_id")
	currentUserID := userID.(uuid.UUID)

	addedCount, failedCount, failedMembers, err := h.service.AddMembersToTeam(c.Request.Context(), teamID, req.UserIDs, currentUserID)
	if err != nil {
		responses.Error(c, http.StatusInternalServerError, err, "Failed to add members to team")
		return
	}

	response := gin.H{
		"success": true,
		"message": "Members added successfully",
		"data": gin.H{
			"addedCount":    addedCount,
			"failedCount":   failedCount,
			"failedMembers": failedMembers,
		},
	}
	responses.JSON(c, http.StatusOK, response)
}

// GET /teams/:teamId/members
func (h *TeamHandler) GetTeamMembers(c *gin.Context) {
	teamIDStr := c.Param("teamId")
	teamID, err := uuid.Parse(teamIDStr)
	if err != nil {
		responses.Error(c, http.StatusBadRequest, err, "Invalid team ID format")
		return
	}

	members, err := h.service.GetTeamMembers(c.Request.Context(), teamID)
	if err != nil {
		responses.Error(c, http.StatusInternalServerError, err, "Failed to retrieve team members")
		return
	}
	response := gin.H{
		"success": true,
		"message": "Team members retrieved successfully",
		"data":    members,
	}
	responses.JSON(c, http.StatusOK, response)
}

// DELETE /teams/:teamId/members/:memberId
func (h *TeamHandler) RemoveMember(c *gin.Context) {
	teamIDStr := c.Param("teamId")
	teamID, err := uuid.Parse(teamIDStr)
	if err != nil {
		responses.Error(c, http.StatusBadRequest, err, "Invalid team ID format")
		return
	}

	memberIDStr := c.Param("memberId")
	memberID, err := uuid.Parse(memberIDStr)
	if err != nil {
		responses.Error(c, http.StatusBadRequest, err, "Invalid member ID format")
		return
	}

	userID, _ := c.Get("user_id")
	currentUserID := userID.(uuid.UUID)

	err = h.service.RemoveMemberFromTeam(c.Request.Context(), teamID, memberID, currentUserID)
	if err != nil {
		responses.Error(c, http.StatusInternalServerError, err, "Failed to remove member from team")
		return
	}

	response := gin.H{
		"success": true,
		"message": "Member removed successfully",
	}
	responses.JSON(c, http.StatusOK, response)
}

// DELETE /teams/:teamId/managers/:managerId (only MAIN_MANAGER can remove a manager)
func (h *TeamHandler) RemoveManagerFromTeam(c *gin.Context) {
	teamIDStr := c.Param("teamId")
	teamID, err := uuid.Parse(teamIDStr)
	if err != nil {
		responses.Error(c, http.StatusBadRequest, err, "Invalid team ID format")
		return
	}

	managerIdStr := c.Param("managerId")
	managerId, err := uuid.Parse(managerIdStr)
	if err != nil {
		responses.Error(c, http.StatusBadRequest, err, "Invalid member ID format")
		return
	}

	userID, _ := c.Get("user_id")
	currentUserID := userID.(uuid.UUID)

	err = h.service.RemoveManagerFromTeam(c.Request.Context(), teamID, managerId, currentUserID)
	if err != nil {
		responses.Error(c, http.StatusInternalServerError, err, "Failed to remove manager from team")
		return
	}

	response := gin.H{
		"success": true,
		"message": "Manager removed successfully",
	}
	responses.JSON(c, http.StatusOK, response)
}
