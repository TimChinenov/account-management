package models

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (factory UserFactory) CurrentUser(c *gin.Context) {
	userId, err := ExtractTokenId(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := factory.GetUserById(int(userId))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data ": user})
}
