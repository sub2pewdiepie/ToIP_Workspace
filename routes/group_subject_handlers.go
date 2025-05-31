package routes

import (
	"net/http"
	"space/models"
	"space/models/dto"
	"space/services"
	"strconv"

	"github.com/gin-gonic/gin"
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
// @Description Retrieves a group with preloaded AcademicGroup and Admin data
// @Tags groups
// @Accept json
// @Produce json
// @Param id path int true "Group ID"
// @Param Authorization header string true "Bearer JWT"
// @Success 200 {object} models.Group
// @Failure 400 {object} map[string]string "Invalid group ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Group not found"
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
// @Description Retrieves a paginated list of groups with name, admin name, and academic group
// @Tags groups
// @Accept json
// @Produce json
// @Param page query int false "Page number (default 1)"
// @Param page_size query int false "Items per page (default 10)"
// @Param Authorization header string true "Bearer JWT"
// @Success 200 {object} dto.GetGroupsResponse
// @Failure 400 {object} map[string]string "Invalid page or page_size"
// @Failure 401 {object} map[string]string "Unauthorized"
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

// CreateGroup godoc
// @Summary Create a new group
// @Description Creates a group with the provided details, setting the authenticated user as admin
// @Tags groups
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT"
// @Param group body models.Group true "Group data"
// @Success 201 {object} dto.GroupDTO
// @Failure 400 {object} map[string]string "Invalid request body or missing required fields"
// @Failure 401 {object} map[string]string "Unauthorized"
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
// @Description Updates a group's name or academic_group_id, restricted to the group admin
// @Tags groups
// @Accept json
// @Produce json
// @Param id path int true "Group ID"
// @Param Authorization header string true "Bearer JWT"
// @Param group body models.Group true "Group data (partial)"
// @Success 200 {object} models.Group
// @Failure 400 {object} map[string]string "Invalid group ID, request body, or unauthorized"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Group not found"
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
// @Description Deletes a group by ID
// @Tags groups
// @Accept json
// @Produce json
// @Param id path int true "Group ID"
// @Param Authorization header string true "Bearer JWT"
// @Success 200 {object} map[string]string "Group deleted"
// @Failure 400 {object} map[string]string "Invalid group ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Group not found"
// @Router /api/groups/{id} [delete]
func (h *GroupHandler) DeleteGroup(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}
	if err := h.service.DeleteGroup(int32(id)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
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
	if input.Description != "" {
		subject.Description = input.Description
	}
	if input.GroupID != 0 {
		subject.GroupID = input.GroupID
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
