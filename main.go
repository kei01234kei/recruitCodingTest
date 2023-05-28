package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
)

type User struct {
	UserId   string `json:"user_id"`
	Password string `json:"password"`
}

var userStore = make(map[string]User)

func main() {
	r := gin.Default()

	r.POST("/signup", func(c *gin.Context) {
		var newUser User
		if err := c.ShouldBindJSON(&newUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Account creation failed", "cause": "invalid request body"})
			return
		}

		if newUser.UserId == "" || newUser.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Account creation failed", "cause": "required user_id and password"})
			return
		}

		if len(newUser.UserId) < 6 || len(newUser.UserId) > 20 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Account creation failed", "cause": "user_id should be between 6 and 20 characters"})
			return
		}

		if len(newUser.Password) < 8 || len(newUser.Password) > 20 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Account creation failed", "cause": "password should be between 8 and 20 characters"})
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
