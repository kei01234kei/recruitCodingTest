package main

import (
	"database/sql"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	UserId   string `json:"user_id"`
	Password string `json:"password"`
}

func main() {
	db, err := sql.Open("sqlite3", "./user.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	db.Exec("CREATE TABLE IF NOT EXISTS users (user_id TEXT, password TEXT)")

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

		var user User
		db.QueryRow("SELECT * FROM users WHERE user_id=?", newUser.UserId).Scan(&user.UserId, &user.Password)
		if user.UserId != "" {
			c.JSON(http.StatusConflict, gin.H{"message": "Account creation failed", "cause": "already same user_id is used"})
			return
		}

		_, err = db.Exec("INSERT INTO users (user_id, password) VALUES (?, ?)", newUser.UserId, newUser.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Account creation failed", "cause": "failed to insert user to database"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Account successfully created", "user": gin.H{"user_id": newUser.UserId, "nickname": newUser.UserId}})
	})

	r.Run()
}
