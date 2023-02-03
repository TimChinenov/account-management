package models

import (
	"context"
	"database/sql"
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
}

type UserFactory struct {
	Storage
}

func (factory UserFactory) Get(c *gin.Context) {
	c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "username already exists"})
	return
}
