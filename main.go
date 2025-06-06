package main

import (
	"log"
	"net/http"
	"space/auth"
	"space/database"
	"space/repositories"
	"space/routes"
	"space/services"
	"space/utils"

	_ "space/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           СВАГА
// @version         1.0
// @description
// @termsOfService  http://swagger.io/terms/

// @contact.name   Иван Васютин, Павел Пронин, Давит Саакови
// @contact.email  vasyutin.i.a@edu.mirea.ru, saakovi.d.@edu.mirea.ru, pronin.p.v@edu.mirea.ru

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      4edu.su
// @BasePath  /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	utils.Init()
	utils.Logger.Info("Starting application")
	err := database.ConnectDatabase()
	if err != nil {
		utils.Logger.WithField("error", err).Fatal("Failed to connect to database")
		// log.Fatalf("failed to connect to database: %v", err)
	}
	router := gin.Default()

	// Initialize dependencies
	userRepo := repositories.NewUserRepository(database.DB)
	authService := services.NewAuthService(userRepo)

	academicGroupRepo := repositories.NewAcademicGroupRepository(database.DB)
	academicGroupService := services.NewAcademicGroupService(academicGroupRepo)
	academicGroupHandler := routes.NewAcademicGroupHandler(academicGroupService)

	groupuserRepo := repositories.NewGroupUserRepository(database.DB)
	// groupuserService := services.NewGroupUserService(groupuserRepo)
	groupModerRepo := repositories.NewGroupModerRepository(database.DB)
	groupModerService := services.NewGroupModerService(groupModerRepo)
	groupModerHandler := routes.NewGroupModerHandler(groupModerService)

	groupRepo := repositories.NewGroupRepository(database.DB, userRepo)
	groupService := services.NewGroupService(groupRepo, userRepo, groupuserRepo, groupModerRepo)
	groupHandler := routes.NewGroupHandler(groupService)

	subjectRepo := repositories.NewSubjectRepository(database.DB)
	subjectService := services.NewSubjectService(subjectRepo, groupRepo, userRepo)
	subjectHandler := routes.NewSubjectHandler(subjectService)

	taskRepo := repositories.NewTaskRepository(database.DB)
	taskService := services.NewTaskService(taskRepo)
	taskHandler := routes.NewTaskHandler(taskService, groupService, subjectService)

	groupUserRepo := repositories.NewGroupUserRepository(database.DB)
	groupUserService := services.NewGroupUserService(groupUserRepo)
	groupUserHandler := routes.NewGroupUserHandler(groupUserService)

	appRepo := repositories.NewGroupApplicationRepository(database.DB)
	appService := services.NewGroupApplicationService(appRepo, groupRepo, groupModerRepo, userRepo, groupUserRepo)
	appHandler := routes.NewGroupApplicationHandler(appService)

	// Seed database

	if err := database.SeedAcademicGroups(database.DB, academicGroupRepo); err != nil {
		log.Fatalf("failed to seed academic groups: %v", err)
	}
	if err := database.SeedSubjects(database.DB, subjectRepo, academicGroupRepo); err != nil {
		log.Fatalf("failed to seed academic groups: %v", err)
	}

	// Public routes
	router.POST("/login", routes.LoginHandler(authService))
	router.POST("/register", routes.RegisterHandler(authService))
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, World!")
	})

	// Protected routes
	protected := router.Group("/api")
	protected.Use(auth.AuthMiddleware())
	{
		// Task endpoints
		tasks := protected.Group("/tasks")
		{
			tasks.GET("", taskHandler.GetGroupTasks)
			tasks.GET("/my-groups", taskHandler.GetMyGroupTasks)
			tasks.GET("/:id", taskHandler.GetTask)
			tasks.POST("/", taskHandler.CreateTask)
			tasks.DELETE("/:id", taskHandler.DeleteTask)
			// protected.PATCH("/tasks/:id", taskHandler.UpdateTask)
			tasks.PATCH("/:id/verify", taskHandler.VerifyTask)
		}
		// Group endpoints
		groups := protected.Group("/groups")
		{
			groups.GET("/:id", groupHandler.GetGroup)
			groups.POST("", groupHandler.CreateGroup)
			groups.PATCH("/:id", groupHandler.UpdateGroup)
			groups.DELETE("/:id", groupHandler.DeleteGroup)
			groups.GET("/available", groupHandler.GetAvailableGroups)
			groups.GET("/:id/users", groupHandler.GetGroupUsers)
			groups.GET("/my-groups", groupHandler.GetUserGroups)
			groups.GET("/:id/subjects/:subject_id/tasks", taskHandler.GetTasksBySubject)
			groups.GET("/:id/subjects", taskHandler.GetSubjectsByGroup)
		}

		// Subject endpoints
		protected.GET("/subjects/:id", subjectHandler.GetSubject)
		protected.POST("/subjects", subjectHandler.CreateSubject)
		protected.PATCH("/subjects/:id", subjectHandler.UpdateSubject)
		protected.DELETE("/subjects/:id", subjectHandler.DeleteSubject)
		protected.GET("/subjects/my-groups", taskHandler.GetUserSubjects)

		// GroupUser endpoints
		protected.GET("/group-users/:group_id/:user_id", groupUserHandler.GetGroupUser)
		protected.POST("/group-users", groupUserHandler.CreateGroupUser)
		protected.PATCH("/group-users/:group_id/:user_id", groupUserHandler.UpdateGroupUser)
		protected.DELETE("/group-users/:group_id/:user_id", groupUserHandler.DeleteGroupUser)
		// AcademicGroup endpoints
		protected.GET("/academic-groups/:id", academicGroupHandler.GetAcademicGroup)
		protected.POST("/academic-groups", academicGroupHandler.CreateAcademicGroup)
		protected.PATCH("/academic-groups/:id", academicGroupHandler.UpdateAcademicGroup)
		protected.DELETE("/academic-groups/:id", academicGroupHandler.DeleteAcademicGroup)
		protected.GET("/academic-groups", academicGroupHandler.GetAllAcademicGroups)
		// GroupModer endpoints
		protected.GET("/group-moders/:group_id/:user_id", groupModerHandler.GetGroupModer)
		protected.POST("/group-moders", groupModerHandler.CreateGroupModer)
		protected.DELETE("/group-moders/:group_id/:user_id", groupModerHandler.DeleteGroupModer)
		// Applications
		protected.POST("/groups/applications", appHandler.CreateApplication)
		protected.GET("/groups/applications/pending", appHandler.GetPendingApplications)
		protected.PATCH("/groups/applications/:group_id/review", appHandler.ReviewApplication)
	}

	router.Run(":8080")
}
