package main

import (
	"database/sql"
	"example/account-management/models"
	"fmt"

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

	router.POST("/users", models.UserFactory{Storage: db}.Create)
	router.GET("/users/:id", models.UserFactory{Storage: db}.Get)
	router.GET("/users", models.UserFactory{Storage: db}.Search)

	router.PATCH("/users/:id/points", models.UserFactory{Storage: db}.UpdatePoints)
	router.POST("/users/login", models.UserFactory{Storage: db}.Login)
	router.POST("/users/logout", models.UserFactory{Storage: db}.Logout)

	router.Run("localhost:8080")
}
