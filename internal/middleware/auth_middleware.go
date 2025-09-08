package middleware

import (
	"log"
	"net/http"
	"strings"

	"go_service/internal/auth"
	"go_service/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func AuthMiddleware(db *gorm.DB) gin.HandlerFunc {
	userService := services.NewUserService()
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			log.Printf("Token validation failed: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		log.Printf("Token validated successfully for user ID: %s", claims.UserID.String())

		// Fetch user info from database to get role
		userResponse, err := userService.GetUserByID(c.Request.Context(), claims.UserID)
		if err != nil {
			log.Printf("Failed to fetch user data: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Fail to fetch user data"})
			c.Abort()
			return
		}

		if userResponse == nil || userResponse.ID == uuid.Nil {
			log.Printf("User data is empty for ID: %s", claims.UserID.String())
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("role", userResponse.Role)
		c.Next()
	}
}
