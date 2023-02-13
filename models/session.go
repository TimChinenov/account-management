package models

import (
	"crypto/sha256"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func GetSession(c *gin.Context) {
	session := sessions.Default(c)

	value := session.Get("user")
	if value == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "There is no user in this session"})
		return
	}

	c.String(200, value.(string))
}

func SetSession(c *gin.Context, username string) {
	session := sessions.Default(c)

	hasher := sha256.New()
	hasher.Write([]byte(username))
	token := hasher.Sum(nil)

	session.Set(username, token)
	err := session.Save()

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "There is no user in this session"})
		return
	}
	c.String(200, "session set")
}
