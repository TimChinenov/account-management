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
	baseUrl := getEnvironmentVariableOrDefault("BASE_URL", "localhost")

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

	userStore := users.NewUserStore(db)
	postStore := posts.NewPostStore(db)

	public := router.Group("/api")
	public.POST("/users", userStore.Create)
	public.GET("/users/:id", userStore.Get)
	public.GET("/users", userStore.Search)
	public.PATCH("/users/:id/points", userStore.UpdatePoints)
	public.POST("/login", userStore.Login)

	protected := router.Group("/api/admin")
	protected.Use(middleware.JwtAuthMiddleware())
	protected.GET("/user", userStore.CurrentUser)
	protected.POST("/posts", postStore.Create)
	protected.GET("/posts/:page/:page_count", postStore.Search)
	protected.POST("/posts/vote", postStore.Vote)

	router.Run(fmt.Sprintf("%s:8080", baseUrl))
}

func getEnvironmentVariableOrDefault(environmentVariable string, defaultValue string) string {
	value := os.Getenv(environmentVariable)

	if value == "" {
		value = defaultValue
	}

	return value
}
