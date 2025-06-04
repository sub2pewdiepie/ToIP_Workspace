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

type TaskHandler struct {
	taskService  *services.TaskService
	groupService *services.GroupService
}

func NewTaskHandler(taskService *services.TaskService, groupService *services.GroupService) *TaskHandler {
	return &TaskHandler{taskService, groupService}
}

// GetGroupTasks godoc
// @Summary Get tasks in a group
// @Description Get a list of tasks for the specified group
// @Tags tasks
// @Accept json
// @Produce json
// @Param group_id query int true "Group ID"
// @Param Authorization header string true "Bearer JWT"
// @Success 200 {array} dto.TaskDTO
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /api/tasks [get]
func (h *TaskHandler) GetGroupTasks(c *gin.Context) {
	groupIDStr := c.Query("group_id")
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

	isMember, err := h.groupService.IsGroupMember(int32(groupID), username.(string))
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

	tasks, err := h.taskService.GetGroupTasks(int32(groupID))
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error":    err,
			"group_id": groupID,
		}).Error("Failed to fetch group tasks")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch group tasks"})
		return
	}

	utils.Logger.WithFields(logrus.Fields{
		"username": username,
		"group_id": groupID,
		"count":    len(tasks),
	}).Info("Successfully fetched group tasks")
	c.JSON(http.StatusOK, tasks)
}

// CreateTask godoc
// @Summary Create a new task
// @Description Create a new task in the specified group
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body dto.CreateTaskRequest true "Task data"
// @Param Authorization header string true "Bearer JWT"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /api/tasks [post]
func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req dto.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	username, exists := c.Get("username")
	if !exists {
		utils.Logger.Error("Unauthorized: username not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	isMember, err := h.groupService.IsGroupMember(req.GroupID, username.(string))
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error":    err,
			"username": username,
			"group_id": req.GroupID,
		}).Error("Failed to check group membership")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check group membership"})
		return
	}
	if !isMember {
		utils.Logger.WithFields(logrus.Fields{
			"username": username,
			"group_id": req.GroupID,
		}).Warn("Forbidden: user is not a group member")
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: group membership required"})
		return
	}

	user, err := h.groupService.GetUserByUsername(username.(string))
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error":    err,
			"username": username,
		}).Error("Failed to fetch user")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	if err := h.taskService.CreateTask(req.GroupID, user.UserID, req.Title, req.Description); err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error":    err,
			"group_id": req.GroupID,
			"title":    req.Title,
		}).Error("Failed to create task")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	utils.Logger.WithFields(logrus.Fields{
		"username": username,
		"group_id": req.GroupID,
		"title":    req.Title,
	}).Info("Task created successfully")
	c.JSON(http.StatusCreated, gin.H{"message": "Task created successfully"})
}

// VerifyTask godoc
// @Summary Verify a task
// @Description Verify or deny a task's legitimacy
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Param verify body dto.VerificationRequest true "Verification Status"
// @Param Authorization header string true "Bearer JWT"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/tasks/{id}/verify [patch]
func (h *TaskHandler) VerifyTask(c *gin.Context) {
	taskIDStr := c.Param("id")
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error":   err,
			"task_id": taskIDStr,
		}).Error("Invalid task ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var req dto.VerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	username, exists := c.Get("username")
	if !exists {
		utils.Logger.Error("Unauthorized: username not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Fetch task to get group_id
	task, err := h.taskService.GetTaskByID(int32(taskID))
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error":   err,
			"task_id": taskID,
		}).Error("Failed to fetch task")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch task"})
		return
	}

	isAuthorized, err := h.groupService.IsAdminOrModerator(task.GroupID, username.(string))
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error":    err,
			"username": username,
			"group_id": task.GroupID,
		}).Error("Failed to check authorization")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check authorization"})
		return
	}
	if !isAuthorized {
		utils.Logger.WithFields(logrus.Fields{
			"username": username,
			"group_id": task.GroupID,
		}).Warn("Forbidden: admin or moderator role required")
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: admin or moderator role required"})
		return
	}

	if err := h.taskService.VerifyTask(int32(taskID), req.VerificationStatus); err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error":       err,
			"task_id":     taskID,
			"is_verified": req.VerificationStatus,
		}).Error("Failed to verify task")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify task"})
		return
	}

	utils.Logger.WithFields(logrus.Fields{
		"username":    username,
		"task_id":     taskID,
		"is_verified": req.VerificationStatus,
	}).Info("Task verification completed")
	c.JSON(http.StatusOK, gin.H{"message": "Task verification updated"})
}

// GetTask godoc
// @Summary Get a specific task
// @Description Get details of a task by ID, if user is a member of the task's group
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Param Authorization header string true "Bearer JWT"
// @Success 200 {object} dto.TaskDTO
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/tasks/{id} [get]
func (h *TaskHandler) GetTask(c *gin.Context) {
	taskIDStr := c.Param("id")
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error":   err,
			"task_id": taskIDStr,
		}).Error("Invalid task ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	username, exists := c.Get("username")
	if !exists {
		utils.Logger.Error("Unauthorized: username not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	task, err := h.taskService.GetTaskByID(int32(taskID))
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error":   err,
			"task_id": taskID,
		}).Error("Failed to fetch task")
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	isMember, err := h.groupService.IsGroupMember(task.GroupID, username.(string))
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error":    err,
			"username": username,
			"group_id": task.GroupID,
		}).Error("Failed to check group membership")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check group membership"})
		return
	}
	if !isMember {
		utils.Logger.WithFields(logrus.Fields{
			"username": username,
			"group_id": task.GroupID,
		}).Warn("Forbidden: user is not a group member")
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: group membership required"})
		return
	}

	taskDTO := dto.ToTaskDTO(task)
	utils.Logger.WithFields(logrus.Fields{
		"username": username,
		"task_id":  taskID,
		"title":    task.Title,
	}).Info("Successfully fetched task")
	c.JSON(http.StatusOK, taskDTO)
}

// GetMyGroupTasks godoc
// @Summary Get tasks from all user's groups
// @Description Get a list of tasks from all groups the user is a member of
// @Tags tasks
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT"
// @Success 200 {array} dto.TaskDTO
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/tasks/my-groups [get]
func (h *TaskHandler) GetMyGroupTasks(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		utils.Logger.Error("Unauthorized: username not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	groupIDs, err := h.groupService.GetUserGroupIDs(username.(string))
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error":    err,
			"username": username,
		}).Error("Failed to fetch user's groups")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user's groups"})
		return
	}

	if len(groupIDs) == 0 {
		utils.Logger.WithFields(logrus.Fields{
			"username": username,
		}).Info("User is not a member of any groups")
		c.JSON(http.StatusOK, []dto.TaskDTO{})
		return
	}

	tasks, err := h.taskService.GetTasksByGroupIDs(groupIDs)
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error":     err,
			"group_ids": groupIDs,
		}).Error("Failed to fetch tasks for user's groups")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
		return
	}

	utils.Logger.WithFields(logrus.Fields{
		"username":    username,
		"count":       len(tasks),
		"group_count": len(groupIDs),
	}).Info("Successfully fetched tasks from user's groups")
	c.JSON(http.StatusOK, tasks)
}
