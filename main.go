package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
)

type User struct {
	UserId   string `json:"user_id" binding:"required,min=6,max=20"`
	Password string `json:"password" binding:"required,min=8,max=20"`
}

var userStore = make(map[string]User)

func main() {
	r := gin.Default()

	r.POST("/signup", func(c *gin.Context) {
		var newUser User
		if err := c.ShouldBindJSON(&newUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Account creation failed", "cause": "invalid format"})
			return
		}

		matched, err := regexp.MatchString("^[a-zA-Z0-9]+$", newUser.UserId)
		if err != nil || !matched {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Account creation failed", "cause": "invalid user_id format"})
			return
		}

		matched, err = regexp.MatchString("^[\\x21-\\x7E]+$", newUser.Password)
		if err != nil || !matched {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Account creation failed", "cause": "invalid password format"})
			return
		}

		if _, exists := userStore[newUser.UserId]; exists {
			c.JSON(http.StatusConflict, gin.H{"message": "Account creation failed", "cause": "already same user_id is used"})
			return
		}

		userStore[newUser.UserId] = newUser
		c.JSON(http.StatusOK, gin.H{"message": "Account successfully created", "user": gin.H{"user_id": newUser.UserId, "nickname": newUser.UserId}})
	})

	r.Run()
}
