package main

import (
	"database/sql"
	"example/account-management/services/middleware"
	"example/account-management/services/posts"
	"example/account-management/services/users"
	"fmt"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	host := getEnvironmentVariableOrDefault("POSTGRES_HOST", "localhost")
	port := getEnvironmentVariableOrDefault("POSTGRES_PORT", "5432")
	user := getEnvironmentVariableOrDefault("POSTGRES_USER", "postgres")
	password := getEnvironmentVariableOrDefault("POSTGRES_PASSWORD", "password")
	dbname := getEnvironmentVariableOrDefault("POSTGRES_DB", "postgres")
	origin := getEnvironmentVariableOrDefault("ORIGIN", "http://localhost:3000")
	baseUrl := getEnvironmentVariableOrDefault("BASE_URL", "")

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
	protected.POST("/posts", posts.PostFactory{Storage: db}.Create)
	protected.GET("/posts/:page/:page_count", posts.PostFactory{Storage: db}.Search)

	router.Run(fmt.Sprintf("%s:8080", baseUrl))
}

func getEnvironmentVariableOrDefault(environmentVariable string, defaultValue string) string {
	value := os.Getenv(environmentVariable)

	if value == "" {
		value = defaultValue
	}

	return value
}
