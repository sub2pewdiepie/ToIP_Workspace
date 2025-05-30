package main

import (
	"log"
	"net/http"
	"space/auth"
	"space/database"
	"space/repositories"
	"space/routes"
	"space/services"

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

// @host      185.221.155.133:8080
// @BasePath  /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {

	err := database.ConnectDatabase()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	router := gin.Default()

	// Initialize dependencies
	userRepo := repositories.NewUserRepository(database.DB)
	authService := services.NewAuthService(userRepo)
	taskRepo := repositories.NewTaskRepository(database.DB)
	taskService := services.NewTaskService(taskRepo)
	taskHandler := routes.NewTaskHandler(taskService)
	groupRepo := repositories.NewGroupRepository(database.DB)
	groupService := services.NewGroupService(groupRepo, userRepo)
	groupHandler := routes.NewGroupHandler(groupService)
	subjectRepo := repositories.NewSubjectRepository(database.DB)
	subjectService := services.NewSubjectService(subjectRepo)
	subjectHandler := routes.NewSubjectHandler(subjectService)
	groupUserRepo := repositories.NewGroupUserRepository(database.DB)
	groupUserService := services.NewGroupUserService(groupUserRepo)
	groupUserHandler := routes.NewGroupUserHandler(groupUserService)
	academicGroupRepo := repositories.NewAcademicGroupRepository(database.DB)
	academicGroupService := services.NewAcademicGroupService(academicGroupRepo)
	academicGroupHandler := routes.NewAcademicGroupHandler(academicGroupService)
	groupModerRepo := repositories.NewGroupModerRepository(database.DB)
	groupModerService := services.NewGroupModerService(groupModerRepo)
	groupModerHandler := routes.NewGroupModerHandler(groupModerService)

	appRepo := repositories.NewGroupApplicationRepository(database.DB)
	appService := services.NewGroupApplicationService(appRepo, groupRepo, groupModerRepo, userRepo, groupUserRepo)
	appHandler := routes.NewGroupApplicationHandler(appService)

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
		protected.GET("/tasks/:id", taskHandler.GetTask)
		protected.POST("/tasks", taskHandler.CreateTask)
		protected.PATCH("/tasks/:id", taskHandler.UpdateTask)
		// Group endpoints
		protected.GET("/groups/:id", groupHandler.GetGroup)
		protected.POST("/groups", groupHandler.CreateGroup)
		protected.PATCH("/groups/:id", groupHandler.UpdateGroup)
		protected.DELETE("/groups/:id", groupHandler.DeleteGroup)
		// Subject endpoints
		protected.GET("/subjects/:id", subjectHandler.GetSubject)
		protected.POST("/subjects", subjectHandler.CreateSubject)
		protected.PATCH("/subjects/:id", subjectHandler.UpdateSubject)
		protected.DELETE("/subjects/:id", subjectHandler.DeleteSubject)
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
		// GroupModer endpoints
		protected.GET("/group-moders/:group_id/:user_id", groupModerHandler.GetGroupModer)
		protected.POST("/group-moders", groupModerHandler.CreateGroupModer)
		protected.DELETE("/group-moders/:group_id/:user_id", groupModerHandler.DeleteGroupModer)
		// Applications
		protected.PATCH("/groups/applications/:user_id/:group_id", appHandler.ReviewApplication)
	}

	router.Run(":8080")
}
