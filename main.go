package main

import (
	"net/http"
	"space/database"
	"space/routes"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title            Сваггер документация api
// @version         1.0
// @description
// @termsOfService  http://swagger.io/terms/

// @contact.name   Иван Васютин, Павел Пронин, Давит Саакови
// @contact.email  vasyutin.i.a@edu.mirea.ru, saakovi.d.@edu.mirea.ru, pronin.p.v@edu.mirea.ru

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      178.208.64.200:8080
// @BasePath  /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {

	database.ConnectDatabase()
	router := gin.Default()
	// Public routes
	router.POST("/login", routes.LoginHandler) // Login route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/hello", func(c *gin.Context) {

		c.String(http.StatusOK, "Hello, World!")

	})

	router.Run(":8080")
}
