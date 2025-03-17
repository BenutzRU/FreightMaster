// backend/middleware.go
package backend

import (
	"FreightMaster/backend/database"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("Middleware triggered for path:", c.Request.URL.Path)
		fmt.Println("Request headers:", c.Request.Header)

		username := c.GetHeader("X-Username")
		password := c.GetHeader("X-Password")

		if username == "" || password == "" {
			fmt.Println("No username or password provided")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Требуются имя пользователя и пароль"})
			c.Abort()
			return
		}

		var user database.User
		if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
			fmt.Println("User not found:", username, "Error:", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не найден"})
			c.Abort()
			return
		}

		fmt.Println("Found user:", user.Username, "Hashed Password:", user.Password)
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
			fmt.Println("Invalid password for user:", username, "Error:", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный пароль"})
			c.Abort()
			return
		}

		fmt.Println("User authenticated, ID:", user.ID, "Role:", user.Role)
		c.Set("userID", user.ID)
		c.Set("userRole", user.Role)
		c.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("userRole")
		if !exists || userRole != "admin" {
			fmt.Println("Forbidden: Role check failed, exists:", exists, "role:", userRole)
			c.JSON(http.StatusForbidden, gin.H{"error": "Требуется доступ администратора"})
			c.Abort()
			return
		}
		fmt.Println("Admin access granted")
		c.Next()
	}
}
