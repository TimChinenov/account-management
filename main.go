package main

import (
	"database/sql"
	"example/account-management/services/middleware"
	"example/account-management/services/users"
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

	router.Use(cors.New(cors.Config{
		AllowOrigins:  []string{origin},
		AllowMethods:  []string{"PATCH", "POST", "GET"},
		AllowHeaders:  []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders: []string{"Content-Type"},
	}))

	public := router.Group("/api")
	public.POST("/users", users.UserFactory{Storage: db}.Create)
	public.GET("/users/:id", users.UserFactory{Storage: db}.Get)
	public.GET("/users", users.UserFactory{Storage: db}.Search)
	public.PATCH("/users/:id/points", users.UserFactory{Storage: db}.UpdatePoints)
	public.POST("/login", users.UserFactory{Storage: db}.Login)

	protected := router.Group("/api/admin")
	protected.Use(middleware.JwtAuthMiddleware())
	protected.GET("/user", users.UserFactory{Storage: db}.CurrentUser)

	router.Run(fmt.Sprintf("%s:8080", baseUrl))
}
