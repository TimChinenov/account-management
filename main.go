package main

import (
	"example/account-management/controllers"

	"github.com/gin-gonic/gin"
)

func main() {
	router := setupRouter()

	router.Run("localhost:8080")
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.POST("/login", controllers.LoginUser)
	router.POST("/users", controllers.CreateUser)

	return router
}
