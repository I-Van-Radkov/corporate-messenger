package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ExtractUserInfoMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetHeader("X-User-ID")
		role := c.GetHeader("X-User-Role")

		if userID == "" || role == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "User info missing",
			})
			return
		}

		c.Set("user_id", userID)
		c.Set("user_role", role)

		c.Next()
	}
}

func RequireRoles(allowedRoles ...ClientRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "User role not found in context",
			})
			return
		}

		for _, allowedRole := range allowedRoles {
			if userRole == allowedRole {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "Unknown role",
		})
	}
}

func RequireAdminOnly() gin.HandlerFunc {
	return RequireRoles(adminRole)
}

func RequireStaff() gin.HandlerFunc {
	return RequireRoles(adminRole, moderatorRole, supportRole)
}
