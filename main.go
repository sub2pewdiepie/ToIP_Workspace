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
		subjects := protected.Group("/subjects")
		{
			subjects.GET("/:id", subjectHandler.GetSubject)
			subjects.POST("", subjectHandler.CreateSubject)
			subjects.PATCH("/:id", subjectHandler.UpdateSubject)
			subjects.DELETE("/:id", subjectHandler.DeleteSubject)
			subjects.GET("/my-groups", taskHandler.GetUserSubjects)
		}

		// GroupUser endpoints
		groupuser := protected.Group("/group-users")
		{
			groupuser.GET("/:group_id/:user_id", groupUserHandler.GetGroupUser)
			groupuser.POST("", groupUserHandler.CreateGroupUser)
			groupuser.PATCH("/:group_id/:user_id", groupUserHandler.UpdateGroupUser)
			groupuser.DELETE("/:group_id/:user_id", groupUserHandler.DeleteGroupUser)
		}
		// AcademicGroup endpoints
		academicgroups := protected.Group("/academic-groups")
		{
			academicgroups.GET("/:id", academicGroupHandler.GetAcademicGroup)
			academicgroups.POST("", academicGroupHandler.CreateAcademicGroup)
			academicgroups.PATCH("/:id", academicGroupHandler.UpdateAcademicGroup)
			academicgroups.DELETE("/:id", academicGroupHandler.DeleteAcademicGroup)
			academicgroups.GET("", academicGroupHandler.GetAllAcademicGroups)
		}

		// GroupModer endpoints
		groumoders := protected.Group("/group-moders")
		{
			groumoders.GET("/:group_id/:user_id", groupModerHandler.GetGroupModer)
			groumoders.POST("", groupModerHandler.CreateGroupModer)
			groumoders.DELETE("/:group_id/:user_id", groupModerHandler.DeleteGroupModer)
		}
		// Applications
		applications := protected.Group("/groups/applications")
		{
			applications.POST("", appHandler.CreateApplication)
			applications.GET("/pending", appHandler.GetPendingApplications)
			applications.PATCH("/review/:id", appHandler.ReviewApplication)
		}
	}

	router.Run(":8080")
}
