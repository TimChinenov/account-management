package main

import (
	"database/sql"
	"example/account-management/models"
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "password"
	dbname   = "postgres"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	router := gin.Default()

	// TODO: delete this in production
	router.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"http://localhost:3000"},
		AllowMethods:  []string{"PATCH", "POST", "GET"},
		AllowHeaders:  []string{"Origin", "Content-Type"},
		ExposeHeaders: []string{"Content-Type"},
	}))

	public := router.Group("/api")
	public.POST("/users/", models.UserFactory{Storage: db}.Create)
	public.GET("/users/:id", models.UserFactory{Storage: db}.Get)
	public.GET("/users", models.UserFactory{Storage: db}.Search)
	public.PATCH("/users/:id/points", models.UserFactory{Storage: db}.UpdatePoints)
	public.POST("/login/", models.UserFactory{Storage: db}.Login)
	public.POST("/logout", models.UserFactory{Storage: db}.Logout)

	protected := router.Group("/api/admin")
	protected.Use(models.JwtAuthMiddleware())
	protected.GET("/user", models.UserFactory{Storage: db}.CurrentUser)

	router.Run("localhost:8080")

}
