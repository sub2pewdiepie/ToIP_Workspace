package routes

import (
	"net/http"
	"space/models/dto"
	"space/services"
	"space/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type GroupApplicationHandler struct {
	service *services.GroupApplicationService
}

func NewGroupApplicationHandler(service *services.GroupApplicationService) *GroupApplicationHandler {
	return &GroupApplicationHandler{service}
}

// ReviewStatusRequest represents the request body for reviewing an application
type ReviewStatusRequest struct {
	Username string `json:"username" binding:"required"`
	Status   string `json:"status" binding:"required,oneof=approved rejected"`
}

// ReviewApplication godoc
// @Summary Review a group application
// @Description Approve or reject a group application for a user to join a group. Requires admin or moderator privileges for the group.
// @Tags group_applications
// @Accept json
// @Produce json
// @Param group_id path int true "Group ID"
// @Param body body ReviewStatusRequest true "Review details"
// @Param Authorization header string true "Bearer JWT"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/groups/applications/{group_id}/review [patch]
func (h *GroupApplicationHandler) ReviewApplication(c *gin.Context) {
	groupIDStr := c.Param("group_id")

	utils.Logger.WithFields(logrus.Fields{
		"group_id": groupIDStr,
	}).Debug("Received review application request")

	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error":    err,
			"group_id": groupIDStr,
			"place":    "handler",
		}).Error("Invalid group ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	var req ReviewStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error": err,
			"body":  req,
		}).Error("Invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid username or status"})
		return
	}

	username, exists := c.Get("username")
	if !exists {
		utils.Logger.Error("Unauthorized: username not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	utils.Logger.WithFields(logrus.Fields{
		"reviewer": username,
		"username": req.Username,
		"group_id": groupID,
		"status":   req.Status,
	}).Debug("Reviewing application")

	err = h.service.ReviewApplication(int32(groupID), req.Username, username.(string), req.Status)
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error":    err,
			"reviewer": username,
			"username": req.Username,
			"group_id": groupID,
		}).Error("Failed to review application")
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	utils.Logger.WithFields(logrus.Fields{
		"reviewer": username,
		"username": req.Username,
		"group_id": groupID,
		"status":   req.Status,
	}).Info("Application reviewed successfully")
	c.JSON(http.StatusOK, gin.H{"message": "Application reviewed successfully"})
}

// CreateApplication godoc
// @Summary Apply to a group
// @Description Submit an application to join a group with an optional message.
// @Tags group_applications
// @Accept json
// @Produce json
// @Param body body dto.CreateApplicationRequest true "Application details"
// @Param Authorization header string true "Bearer JWT"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/groups/applications [post]
func (h *GroupApplicationHandler) CreateApplication(c *gin.Context) {
	var req dto.CreateApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error": err,
		}).Error("Invalid input")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// userID, exists := c.Get("userID")
	username, exists := c.Get("username")
	if !exists {
		utils.Logger.Error("Unauthorized: username not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	utils.Logger.WithFields(logrus.Fields{
		"username": username,
		"group_id": req.GroupID,
	}).Debug("Processing group application")

	err := h.service.ApplyToGroup(c, req.GroupID, req.Message)
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error":    err,
			"username": username,
			"group_id": req.GroupID,
		}).Error("Failed to fetch application")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	utils.Logger.WithFields(logrus.Fields{
		// "application_id": application.ApplicationID,
		"username": username,
		"group_id": req.GroupID,
	}).Info("Group application created")
	c.JSON(http.StatusCreated, gin.H{"message": "Application submitted successfully"})
}

// GetPendingApplications godoc
// @Summary Get pending group applications
// @Description Retrieve all pending applications for groups where the user is an admin or moderator.
// @Tags group_applications
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT"
// @Success 200 {array} dto.GroupApplicationDTO
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/groups/applications/pending [get]
func (h *GroupApplicationHandler) GetPendingApplications(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		utils.Logger.Error("Unauthorized: username not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	utils.Logger.WithFields(logrus.Fields{
		"username": username,
	}).Debug("Fetching pending applications")
	applications, err := h.service.GetPendingApplications(c)
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error":    err,
			"username": username,
		}).Error("Failed to fetch pending applications")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	applicationDTOs := make([]dto.GroupApplicationDTO, len(applications))
	for i, app := range applications {
		applicationDTOs[i] = dto.ToGroupApplicationDTO(&app)
	}
	utils.Logger.WithFields(logrus.Fields{
		"username":          username,
		"application_count": len(applications),
	}).Info("Retrieved pending applications")
	c.JSON(http.StatusOK, applicationDTOs)
}
