package models

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Cookies(c *gin.Context) {
	cookie, err := c.Cookie("user")

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": cookie})
}
