package routes

import (
	"net/http"
	"space/models"
	"space/models/dto"
	"space/services"
	"space/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// GroupHandler handles group-related HTTP requests
type GroupHandler struct {
	service *services.GroupService
}

// SubjectHandler handles subject-related HTTP requests
type SubjectHandler struct {
	service *services.SubjectService
}

func NewGroupHandler(service *services.GroupService) *GroupHandler {
	return &GroupHandler{service}
}

func NewSubjectHandler(service *services.SubjectService) *SubjectHandler {
	return &SubjectHandler{service}
}

// GetGroup godoc
// @Summary Get a group by ID
// @Description Retrieves a group by ID with preloaded AcademicGroup and Admin data. Accessible to any authenticated user.
// @Tags groups
// @Accept json
// @Produce json
// @Param id path int true "Group ID" example(1)
// @Param Authorization header string true "Bearer JWT"
// @Success 200 {object} dto.GroupDTO
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/groups/{id} [get]
func (h *GroupHandler) GetGroup(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}
	group, err := h.service.GetGroupByID(int32(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}
	c.JSON(http.StatusOK, group)
}

// GetAllGroups godoc
// @Summary Get all groups with pagination
// @Description Retrieves a paginated list of groups with name, admin username, and academic group. Accessible to any authenticated user.
// @Tags groups
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1) example(1)
// @Param page_size query int false "Items per page" default(10) example(10)
// @Param Authorization header string true "Bearer JWT"
// @Success 200 {object} dto.GetGroupsResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/groups [get]
func (h *GroupHandler) GetAllGroups(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
		return
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page size"})
		return
	}

	groups, total, err := h.service.GetAllGroups(int32(page), int32(pageSize))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response := dto.GetGroupsResponse{
		Groups: groups,
		Pagination: dto.PaginationMeta{
			Page:     page,
			PageSize: pageSize,
			Total:    total,
			Pages:    (total + int64(pageSize) - 1) / int64(pageSize),
		},
	}
	c.JSON(http.StatusOK, response)
}

// GetAvailableGroups godoc
// @Summary Get groups available to apply to
// @Description Retrieves a paginated list of groups where the authenticated user is not a member, admin, or moderator.
// @Tags groups
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1) example(1)
// @Param page_size query int false "Items per page" default(10) example(10)
// @Param Authorization header string true "Bearer JWT"
// @Success 200 {object} dto.GetGroupsResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/groups/available [get]
func (h *GroupHandler) GetAvailableGroups(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		utils.Logger.WithFields(logrus.Fields{
			"error": err,
			"page":  pageStr,
		}).Error("Invalid page number")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
		return
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		utils.Logger.WithFields(logrus.Fields{
			"error":     err,
			"page_size": pageSizeStr,
		}).Error("Invalid page size")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page size"})
		return
	}

	username, exists := c.Get("username")
	if !exists {
		utils.Logger.Error("Unauthorized: username not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	utils.Logger.WithFields(logrus.Fields{
		"username":  username,
		"page":      page,
		"page_size": pageSize,
	}).Debug("Fetching available groups")

	groups, total, err := h.service.GetAvailableGroups(c, username.(string), page, pageSize)
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error":    err,
			"username": username,
		}).Error("Failed to fetch available groups")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch available groups"})
		return
	}

	response := dto.GetGroupsResponse{
		Groups: groups,
		Pagination: dto.PaginationMeta{
			Page:     page,
			PageSize: pageSize,
			Total:    total,
			Pages:    (total + int64(pageSize) - 1) / int64(pageSize),
		},
	}

	utils.Logger.WithFields(logrus.Fields{
		"username": username,
		"total":    total,
		"page":     page,
	}).Info("Available groups fetched successfully")
	c.JSON(http.StatusOK, response)
}

// CreateGroup godoc
// @Summary Create a new group
// @Description Creates a group with the provided details, setting the authenticated user as admin.
// @Tags groups
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT"
// @Param group body dto.CreateGroupRequest true "Group data"
// @Success 201 {object} dto.GroupDTO
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/groups [post]
func (h *GroupHandler) CreateGroup(c *gin.Context) {
	var group models.Group
	if err := c.ShouldBindJSON(&group); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	if err := h.service.CreateGroup(c, &group); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	groupDTO := dto.GroupDTO{
		ID:              group.ID,
		Name:            group.Name,
		AdminUsername:   group.Admin.Username,
		AcademicGroupID: group.AcademicGroupID,
		AcademicGroup:   group.AcademicGroup.Name,
	}
	c.JSON(http.StatusCreated, groupDTO)
}

// UpdateGroup godoc
// @Summary Update a group
// @Description Updates a group's name or academic group ID, restricted to the group admin.
// @Tags groups
// @Accept json
// @Produce json
// @Param id path int true "Group ID" example(1)
// @Param Authorization header string true "Bearer JWT"
// @Param group body dto.UpdateGroupRequest true "Group data (partial)"
// @Success 200 {object} dto.GroupDTO
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/groups/{id} [patch]
func (h *GroupHandler) UpdateGroup(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}
	var input models.Group
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	group, err := h.service.GetGroupByID(int32(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}
	if input.Name != "" {
		group.Name = input.Name
	}
	if input.AcademicGroupID != 0 {
		group.AcademicGroupID = input.AcademicGroupID
	}
	if err := h.service.UpdateGroup(c, group); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, group)
}

// DeleteGroup godoc
// @Summary Delete a group
// @Description Deletes a group by ID, restricted to the group admin.
// @Tags groups
// @Accept json
// @Produce json
// @Param id path int true "Group ID" example(1)
// @Param Authorization header string true "Bearer JWT"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/groups/{id} [delete]
func (h *GroupHandler) DeleteGroup(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error":    err,
			"group_id": idStr,
		}).Error("Invalid group ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	username, exists := c.Get("username")
	if !exists {
		utils.Logger.Error("Unauthorized: username not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	utils.Logger.WithFields(logrus.Fields{
		"username": username,
		"group_id": id,
	}).Debug("Attempting to delete group")

	err = h.service.DeleteGroup(int32(id), username.(string))
	if err != nil {
		switch err.Error() {
		case "group not found":
			c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		case "user not found":
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		case "only group admin can delete the group":
			c.JSON(http.StatusForbidden, gin.H{"error": "Only group admin can delete the group"})
		default:
			utils.Logger.WithFields(logrus.Fields{
				"error":    err,
				"username": username,
				"group_id": id,
			}).Error("Failed to delete group")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete group"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Group deleted"})
}

// GetSubject godoc
// @Summary Get a subject by ID
// @Description Retrieves a subject with preloaded Group data
// @Tags subjects
// @Accept json
// @Produce json
// @Param id path int true "Subject ID"
// @Param Authorization header string true "Bearer JWT"
// @Success 200 {object} models.Subject
// @Failure 400 {object} map[string]string "Invalid subject ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Subject not found"
// @Router /api/subjects/{id} [get]
func (h *SubjectHandler) GetSubject(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subject ID"})
		return
	}
	subject, err := h.service.GetSubjectByID(int32(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subject not found"})
		return
	}
	c.JSON(http.StatusOK, subject)
}

// CreateSubject godoc
// @Summary Create a new subject
// @Description Creates a subject with the provided details
// @Tags subjects
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT"
// @Param subject body models.Subject true "Subject data"
// @Success 201 {object} models.Subject
// @Failure 400 {object} map[string]string "Invalid request body or missing required fields"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /api/subjects [post]
func (h *SubjectHandler) CreateSubject(c *gin.Context) {
	var subject models.Subject
	if err := c.ShouldBindJSON(&subject); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	if err := h.service.CreateSubject(&subject); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, subject)
}

// UpdateSubject godoc
// @Summary Update a subject
// @Description Updates a subject's name, description, or group_id
// @Tags subjects
// @Accept json
// @Produce json
// @Param id path int true "Subject ID"
// @Param Authorization header string true "Bearer JWT"
// @Param subject body models.Subject true "Subject data (partial)"
// @Success 200 {object} models.Subject
// @Failure 400 {object} map[string]string "Invalid subject ID or request body"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Subject not found"
// @Router /api/subjects/{id} [patch]
func (h *SubjectHandler) UpdateSubject(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subject ID"})
		return
	}
	var input models.Subject
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	subject, err := h.service.GetSubjectByID(int32(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subject not found"})
		return
	}
	if input.Name != "" {
		subject.Name = input.Name
	}

	if input.AcademicGroupID != 0 {
		subject.AcademicGroupID = input.AcademicGroupID
	}
	if err := h.service.UpdateSubject(subject); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, subject)
}

// DeleteSubject godoc
// @Summary Delete a subject
// @Description Deletes a subject by ID
// @Tags subjects
// @Accept json
// @Produce json
// @Param id path int true "Subject ID"
// @Param Authorization header string true "Bearer JWT"
// @Success 200 {object} map[string]string "Subject deleted"
// @Failure 400 {object} map[string]string "Invalid subject ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Subject not found"
// @Router /api/subjects/{id} [delete]
func (h *SubjectHandler) DeleteSubject(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subject ID"})
		return
	}
	if err := h.service.DeleteSubject(int32(id)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subject not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Subject deleted"})
}

// GetGroupModerators godoc
// @Summary Get group moderators and admin
// @Description Retrieves the admin and moderators of a group, accessible to admins or moderators of the group.
// @Tags groups
// @Accept json
// @Produce json
// @Param id path int true "Group ID" example(1)
// @Param Authorization header string true "Bearer JWT"
// @Success 200 {object} dto.ModeratorsResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/groups/{id}/moderators [get]
func (h *GroupHandler) GetGroupModerators(c *gin.Context) {
	groupIDStr := c.Param("id")
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error":    err,
			"group_id": groupIDStr,
		}).Error("Invalid group ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	username, exists := c.Get("username")
	if !exists {
		utils.Logger.Error("Unauthorized: username not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	isAuthorized, err := h.service.IsAdminOrModerator(int32(groupID), username.(string))
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error":    err,
			"username": username,
			"group_id": groupID,
		}).Error("Failed to check authorization")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check authorization"})
		return
	}
	if !isAuthorized {
		utils.Logger.WithFields(logrus.Fields{
			"username": username,
			"group_id": groupID,
		}).Warn("Forbidden: user is not admin or moderator")
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: admin or moderator role required"})
		return
	}

	response, err := h.service.GetGroupModeratorsAndAdmin(int32(groupID))
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error":    err,
			"group_id": groupID,
		}).Error("Failed to fetch moderators and admin")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch moderators and admin"})
		return
	}

	utils.Logger.WithFields(logrus.Fields{
		"username":        username,
		"group_id":        groupID,
		"admin_id":        response.Admin.UserID,
		"moderator_count": len(response.Moderators),
	}).Info("Successfully fetched moderators and admin")
	c.JSON(http.StatusOK, response)
}

// GetGroupUsers godoc
// @Summary Get users in a group
// @Description Retrieves a list of users who are members of the specified group, accessible to group members.
// @Tags groups
// @Accept json
// @Produce json
// @Param id path int true "Group ID" example(1)
// @Param Authorization header string true "Bearer JWT"
// @Success 200 {array} dto.UserDTO
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/groups/{id}/users [get]
func (h *GroupHandler) GetGroupUsers(c *gin.Context) {
	groupIDStr := c.Param("id")
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error":    err,
			"group_id": groupIDStr,
		}).Error("Invalid group ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	username, exists := c.Get("username")
	if !exists {
		utils.Logger.Error("Unauthorized: username not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	isMember, err := h.service.IsGroupMember(int32(groupID), username.(string))
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error":    err,
			"username": username,
			"group_id": groupID,
		}).Error("Failed to check group membership")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check group membership"})
		return
	}
	if !isMember {
		utils.Logger.WithFields(logrus.Fields{
			"username": username,
			"group_id": groupID,
		}).Warn("Forbidden: user is not a group member")
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: group membership required"})
		return
	}

	users, err := h.service.GetGroupUsers(int32(groupID))
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error":    err,
			"group_id": groupID,
		}).Error("Failed to fetch group users")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch group users"})
		return
	}

	utils.Logger.WithFields(logrus.Fields{
		"username": username,
		"group_id": groupID,
		"count":    len(users),
	}).Info("Successfully fetched group users")
	c.JSON(http.StatusOK, users)
}

// GetUserGroups godoc
// @Summary Get user's groups
// @Description Retrieves a paginated list of groups where the authenticated user is a member, moderator, or admin.
// @Tags groups
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT"
// @Param page query int false "Page number" default(1) example(1)
// @Param page_size query int false "Items per page" default(10) example(10)
// @Success 200 {object} dto.GetGroupsResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/groups/my-groups [get]
func (h *GroupHandler) GetUserGroups(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		utils.Logger.Error("Unauthorized: username not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page"})
		return
	}
	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if err != nil || pageSize < 1 || pageSize > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page_size"})
		return
	}

	response, err := h.service.GetUserGroups(username.(string), page, pageSize)
	if err != nil {
		if err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		utils.Logger.WithFields(logrus.Fields{
			"error":     err,
			"username":  username,
			"page":      page,
			"page_size": pageSize,
		}).Error("Failed to fetch user's groups")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch groups"})
		return
	}

	utils.Logger.WithFields(logrus.Fields{
		"username":  username,
		"page":      page,
		"page_size": pageSize,
		"count":     len(response.Groups),
	}).Info("Successfully fetched user's groups")
	c.JSON(http.StatusOK, response)
}
