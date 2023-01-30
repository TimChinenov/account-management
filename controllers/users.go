package controllers

import (
	"example/account-management/models"
	"example/account-management/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

var users = []models.User{}

func CreateUser(c *gin.Context) {
	var newUser models.UserRequest

	if err := c.BindJSON((&newUser)); err != nil {
		return
	}

	for _, user := range users {
		if user.Username == newUser.Username {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "username already exists"})
			return
		}
	}

	hashedPassword, err := services.HashPassword(newUser.Password)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid password"})
		return
	}

	var createdUser = models.User{ID: len(users) + 1, Username: newUser.Username, Password: hashedPassword, Score: 0}

	users = append(users, createdUser)
	c.IndentedJSON(http.StatusCreated, toResponse(createdUser))
}

func LoginUser(c *gin.Context) {
	var loginUser models.UserRequest

	if err := c.BindJSON((&loginUser)); err != nil {
		return
	}

	for _, user := range users {
		if user.Username == loginUser.Username {
			if services.CheckPasswordHash(loginUser.Password, user.Password) {
				c.IndentedJSON(http.StatusOK, toResponse(user))
				return
			}
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "user or password invalid 2"})
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "user or password invalid 1"})
}

func GetUser(c *gin.Context) {
	var user models.UserRequest

	if err := c.BindJSON((&user)); err != nil {
		return
	}

	for _, user := range users {
		if user.Username == user.Username {
			c.IndentedJSON(http.StatusOK, toResponse(user))
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "user with requested username not found"})
}

func toResponse(user models.User) models.UserResponse {
	response := models.UserResponse{ID: user.ID, Username: user.Username, Score: user.Score}
	return response
}
