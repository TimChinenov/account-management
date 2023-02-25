package main

import (
	"database/sql"
	"example/account-management/models"
	"fmt"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	origin := os.Getenv("ORIGIN")
	baseUrl := os.Getenv("BASE_URL")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host,
		port,
		user,
		password,
		dbname)

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
		AllowOrigins:  []string{origin},
		AllowMethods:  []string{"PATCH", "POST", "GET"},
		AllowHeaders:  []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders: []string{"Content-Type"},
	}))

	public := router.Group("/api")
	public.POST("/users", models.UserFactory{Storage: db}.Create)
	public.GET("/users/:id", models.UserFactory{Storage: db}.Get)
	public.GET("/users", models.UserFactory{Storage: db}.Search)
	public.PATCH("/users/:id/points", models.UserFactory{Storage: db}.UpdatePoints)
	public.POST("/login", models.UserFactory{Storage: db}.Login)

	protected := router.Group("/api/admin")
	protected.Use(models.JwtAuthMiddleware())
	protected.GET("/user", models.UserFactory{Storage: db}.CurrentUser)
	// protected.POST("/logout", models.UserFactory{Storage: db}.Logout)

	router.Run(fmt.Sprintf("%s:8080", baseUrl))
}
