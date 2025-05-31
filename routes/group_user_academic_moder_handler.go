package routes

import (
	"net/http"
	"space/models"
	"space/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GroupUserHandler handles group-user-related HTTP requests
type GroupUserHandler struct {
	service *services.GroupUserService
}

// AcademicGroupHandler handles academic-group-related HTTP requests
type AcademicGroupHandler struct {
	service *services.AcademicGroupService
}

// GroupModerHandler handles group-moderator-related HTTP requests
type GroupModerHandler struct {
	service *services.GroupModerService
}

func NewGroupUserHandler(service *services.GroupUserService) *GroupUserHandler {
	return &GroupUserHandler{service}
}

func NewAcademicGroupHandler(service *services.AcademicGroupService) *AcademicGroupHandler {
	return &AcademicGroupHandler{service}
}

func NewGroupModerHandler(service *services.GroupModerService) *GroupModerHandler {
	return &GroupModerHandler{service}
}

// GetGroupUser godoc
// @Summary Get a group-user relationship
// @Description Retrieves a group-user relationship with preloaded Group and User data
// @Tags group-users
// @Accept json
// @Produce json
// @Param group_id path int true "Group ID"
// @Param user_id path int true "User ID"
// @Param Authorization header string true "Bearer JWT"
// @Success 200 {object} models.GroupUser
// @Failure 400 {object} map[string]string "Invalid group or user ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Group-user not found"
// @Router /api/group-users/{group_id}/{user_id} [get]
func (h *GroupUserHandler) GetGroupUser(c *gin.Context) {
	groupID, err := strconv.Atoi(c.Param("group_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	groupUser, err := h.service.GetGroupUser(int32(groupID), int32(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group-user not found"})
		return
	}
	c.JSON(http.StatusOK, groupUser)
}

// CreateGroupUser godoc
// @Summary Add a user to a group
// @Description Creates a group-user relationship with a role
// @Tags group-users
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT"
// @Param group_user body models.GroupUser true "GroupUser data"
// @Success 201 {object} models.GroupUser
// @Failure 400 {object} map[string]string "Invalid request body or missing required fields"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /api/group-users [post]
func (h *GroupUserHandler) CreateGroupUser(c *gin.Context) {
	var groupUser models.GroupUser
	if err := c.ShouldBindJSON(&groupUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	if err := h.service.CreateGroupUser(&groupUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, groupUser)
}

// UpdateGroupUser godoc
// @Summary Update a group-user role
// @Description Updates the role of a group-user relationship
// @Tags group-users
// @Accept json
// @Produce json
// @Param group_id path int true "Group ID"
// @Param user_id path int true "User ID"
// @Param Authorization header string true "Bearer JWT"
// @Param group_user body models.GroupUser true "GroupUser data (role only)"
// @Success 200 {object} models.GroupUser
// @Failure 400 {object} map[string]string "Invalid group ID, user ID, or request body"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Group-user not found"
// @Router /api/group-users/{group_id}/{user_id} [patch]
func (h *GroupUserHandler) UpdateGroupUser(c *gin.Context) {
	groupID, err := strconv.Atoi(c.Param("group_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	var input models.GroupUser
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	groupUser, err := h.service.GetGroupUser(int32(groupID), int32(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group-user not found"})
		return
	}
	if input.Role != "" {
		groupUser.Role = input.Role
	}
	if err := h.service.UpdateGroupUser(groupUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, groupUser)
}

// DeleteGroupUser godoc
// @Summary Remove a user from a group
// @Description Deletes a group-user relationship
// @Tags group-users
// @Accept json
// @Produce json
// @Param group_id path int true "Group ID"
// @Param user_id path int true "User ID"
// @Param Authorization header string true "Bearer JWT"
// @Success 200 {object} map[string]string "Group-user deleted"
// @Failure 400 {object} map[string]string "Invalid group or user ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Group-user not found"
// @Router /api/group-users/{group_id}/{user_id} [delete]
func (h *GroupUserHandler) DeleteGroupUser(c *gin.Context) {
	groupID, err := strconv.Atoi(c.Param("group_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	if err := h.service.DeleteGroupUser(int32(groupID), int32(userID)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group-user not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Group-user deleted"})
}

// GetAcademicGroup godoc
// @Summary Get an academic group by ID
// @Description Retrieves an academic group
// @Tags academic-groups
// @Accept json
// @Produce json
// @Param id path int true "Academic Group ID"
// @Param Authorization header string true "Bearer JWT"
// @Success 200 {object} models.AcademicGroup
// @Failure 400 {object} map[string]string "Invalid academic group ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Academic group not found"
// @Router /api/academic-groups/{id} [get]
func (h *AcademicGroupHandler) GetAcademicGroup(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid academic group ID"})
		return
	}
	academicGroup, err := h.service.GetAcademicGroupByID(int32(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Academic group not found"})
		return
	}
	c.JSON(http.StatusOK, academicGroup)
}

// CreateAcademicGroup godoc
// @Summary Create a new academic group
// @Description Creates an academic group with the provided name
// @Tags academic-groups
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT"
// @Param academic_group body models.AcademicGroup true "AcademicGroup data"
// @Success 201 {object} models.AcademicGroup
// @Failure 400 {object} map[string]string "Invalid request body or missing required fields"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /api/academic-groups [post]
func (h *AcademicGroupHandler) CreateAcademicGroup(c *gin.Context) {
	var academicGroup models.AcademicGroup
	if err := c.ShouldBindJSON(&academicGroup); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	if err := h.service.CreateAcademicGroup(&academicGroup); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, academicGroup)
}

// UpdateAcademicGroup godoc
// @Summary Update an academic group
// @Description Updates an academic group's name
// @Tags academic-groups
// @Accept json
// @Produce json
// @Param id path int true "Academic Group ID"
// @Param Authorization header string true "Bearer JWT"
// @Param academic_group body models.AcademicGroup true "AcademicGroup data (name only)"
// @Success 200 {object} models.AcademicGroup
// @Failure 400 {object} map[string]string "Invalid academic group ID or request body"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Academic group not found"
// @Router /api/academic-groups/{id} [patch]
func (h *AcademicGroupHandler) UpdateAcademicGroup(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid academic group ID"})
		return
	}
	var input models.AcademicGroup
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	academicGroup, err := h.service.GetAcademicGroupByID(int32(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Academic group not found"})
		return
	}
	if input.Name != "" {
		academicGroup.Name = input.Name
	}
	if err := h.service.UpdateAcademicGroup(academicGroup); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, academicGroup)
}

// DeleteAcademicGroup godoc
// @Summary Delete an academic group
// @Description Deletes an academic group by ID
// @Tags academic-groups
// @Accept json
// @Produce json
// @Param id path int true "Academic Group ID"
// @Param Authorization header string true "Bearer JWT"
// @Success 200 {object} map[string]string "Academic group deleted"
// @Failure 400 {object} map[string]string "Invalid academic group ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Academic group not found"
// @Router /api/academic-groups/{id} [delete]
func (h *AcademicGroupHandler) DeleteAcademicGroup(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid academic group ID"})
		return
	}
	if err := h.service.DeleteAcademicGroup(int32(id)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Academic group not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Academic group deleted"})
}

// GetGroupModer godoc
// @Summary Get a group-moderator relationship
// @Description Retrieves a group-moderator relationship
// @Tags group-moders
// @Accept json
// @Produce json
// @Param group_id path int true "Group ID"
// @Param user_id path int true "User ID"
// @Param Authorization header string true "Bearer JWT"
// @Success 200 {object} models.GroupModer
// @Failure 400 {object} map[string]string "Invalid group or user ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Group-moderator not found"
// @Router /api/group-moders/{group_id}/{user_id} [get]
func (h *GroupModerHandler) GetGroupModer(c *gin.Context) {
	groupID, err := strconv.Atoi(c.Param("group_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	groupModer, err := h.service.GetGroupModer(int32(groupID), int32(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group-moderator not found"})
		return
	}
	c.JSON(http.StatusOK, groupModer)
}

// CreateGroupModer godoc
// @Summary Add a moderator to a group
// @Description Creates a group-moderator relationship
// @Tags group-moders
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT"
// @Param group_moder body models.GroupModer true "GroupModer data"
// @Success 201 {object} models.GroupModer
// @Failure 400 {object} map[string]string "Invalid request body or moderator already exists"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /api/group-moders [post]
func (h *GroupModerHandler) CreateGroupModer(c *gin.Context) {
	var groupModer models.GroupModer
	if err := c.ShouldBindJSON(&groupModer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	if err := h.service.CreateGroupModer(&groupModer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, groupModer)
}

// DeleteGroupModer godoc
// @Summary Remove a moderator from a group
// @Description Deletes a group-moderator relationship
// @Tags group-moders
// @Accept json
// @Produce json
// @Param group_id path int true "Group ID"
// @Param user_id path int true "User ID"
// @Param Authorization header string true "Bearer JWT"
// @Success 200 {object} map[string]string "Group-moderator deleted"
// @Failure 400 {object} map[string]string "Invalid group or user ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Group-moderator not found"
// @Router /api/group-moders/{group_id}/{user_id} [delete]
func (h *GroupModerHandler) DeleteGroupModer(c *gin.Context) {
	groupID, err := strconv.Atoi(c.Param("group_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	if err := h.service.DeleteGroupModer(int32(groupID), int32(userID)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group-moderator not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Group-moderator deleted"})
}
