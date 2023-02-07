package models

import (
	"context"
	"database/sql"
	"example/account-management/services"
	"net/http"
	"strings"

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

	hashedPassword, err := services.HashPassword(newUser.Password)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid password"})
		return
	}

	query := `INSERT INTO users (username, password, score) VALUES ($1, $2, $3);`
	_, err = factory.Storage.QueryContext(context.Background(), query, newUser.Username, hashedPassword, 0)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	query = `SELECT id, username, score FROM users WHERE username=$1;`
	row := factory.Storage.QueryRowContext(context.Background(), query, newUser.Username)

	var id int
	var username string
	var score int

	err = row.Scan(&id, &username, &score)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var userResponse UserResponse = UserResponse{ID: id, Username: username, Score: score}

	c.IndentedJSON(http.StatusCreated, userResponse)
}

func (factory UserFactory) Get(c *gin.Context) {
	userID := c.Param("id")

	query := `SELECT id, username, score FROM users WHERE id=$1;`
	row := factory.Storage.QueryRowContext(context.Background(), query, userID)

	var userResponse UserResponse

	err := row.Scan(&userResponse.ID, &userResponse.Username, &userResponse.Score)

	if err != nil || userResponse.ID == 0 || userResponse.Username == "" {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "user not found"})
		return
	}

	c.IndentedJSON(http.StatusFound, userResponse)
}

func (factory UserFactory) Search(c *gin.Context) {
	usernameSearch, err := c.GetQuery("username")
	usernameSearch = strings.TrimSpace(usernameSearch)

	if !err || usernameSearch == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "no parameters found"})
		return
	}

	query := `SELECT id, username, score FROM users WHERE username LIKE '%' || $1 || '%' LIMIT 10`
	rows, queryErr := factory.Storage.QueryContext(context.Background(), query, usernameSearch)

	if queryErr != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "failed to search"})
		return
	}

	users := make([]UserResponse, 0)
	for rows.Next() {
		var user UserResponse
		queryErr = rows.Scan(&user.ID, &user.Username, &user.Score)

		if queryErr != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "failed to parse users"})
			return
		}

		users = append(users, user)
	}

	queryErr = rows.Err()
	if queryErr != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "failed to parse users"})
		return
	}

	c.IndentedJSON(http.StatusOK, users)
}

func (factory UserFactory) UpdatePoints(c *gin.Context) {
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "incomplete"})
}

func (factory UserFactory) Login(c *gin.Context) {
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "incomplete"})
}

func (factory UserFactory) Logout(c *gin.Context) {
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "incomplete"})
}
