package models

import (
	"context"
	"database/sql"
	"example/account-management/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Score    int    `json:"score"`
}

type UserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Score    int    `json:"score"`
}

type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Storage interface {
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
}

type UserFactory struct {
	Storage
}

func (factory UserFactory) Create(c *gin.Context) {
	var newUser UserRequest

	if err := c.BindJSON((&newUser)); err != nil {
		return
	}

	// for _, user := range users {
	// 	if user.Username == newUser.Username {
	// 		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "username already exists"})
	// 		return
	// 	}
	// }

	hashedPassword, err := services.HashPassword(newUser.Password)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid password"})
		return
	}

	query := `INSERT INTO users (username, password, score) VALUES ($1, $2, $3);`
	id := 0
	rows, err := factory.Storage.QueryContext(context.Background(), query, newUser.Username, hashedPassword, 0)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	rows.Scan(&id)

	query = `SELECT id, username, score FROM users WHERE id=$1;`
	row := factory.Storage.QueryRowContext(context.Background(), query, id)

	var username string
	var score int

	err = row.Scan(&id, &username, &score)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	}

	var userResponse UserResponse = UserResponse{ID: id, Username: username, Score: score}

	c.IndentedJSON(http.StatusCreated, userResponse)
}

func (factory UserFactory) Get(c *gin.Context) {
	c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "username already exists"})
	return
}
