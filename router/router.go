package router

import (
	"example/account-management/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.POST("/login", controllers.LoginUser)
	router.POST("/users", controllers.CreateUser)

	return router
}
