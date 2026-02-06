package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"assetManager/internal/auth"
)

const (
	AuthorizationHeader = "Authorization"
	BearerPrefix        = "Bearer "
	UserIDKey           = "userID"
	UsernameKey         = "username"
)

// AuthMiddleware creates a JWT authentication middleware
func AuthMiddleware(jwtService *auth.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(AuthorizationHeader)
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"Error": "Authorization header required"})
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, BearerPrefix) {
			c.JSON(http.StatusUnauthorized, gin.H{"Error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, BearerPrefix)
		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"Error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Set user info in context
		c.Set(UserIDKey, claims.UserID)
		c.Set(UsernameKey, claims.Username)

		c.Next()
	}
}

// GetUserID retrieves the user ID from the context
func GetUserID(c *gin.Context) int64 {
	userID, exists := c.Get(UserIDKey)
	if !exists {
		return 0
	}
	return userID.(int64)
}

// GetUsername retrieves the username from the context
func GetUsername(c *gin.Context) string {
	username, exists := c.Get(UsernameKey)
	if !exists {
		return ""
	}
	return username.(string)
}
