package v1

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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

func ExtractUserIdForWs() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Для WebSocket сначала проверяем query параметры
		token := c.Query("token")

		// Если нет в query, пробуем из заголовков (для HTTP запросов)
		if token == "" {
			authHeader := c.GetHeader("Authorization")
			if authHeader == "" {
				c.AbortWithStatusJSON(401, gin.H{"error": "Authorization required"})
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				c.AbortWithStatusJSON(401, gin.H{"error": "Invalid auth format"})
				return
			}
			token = parts[1]
		}

		if token == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Token required"})
			return
		}

		claims, err := parseJWTUnverified(token)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token format"})
			return
		}

		userId, _ := claims["user_id"].(string)
		if userId == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token payload"})
			return
		}

		c.Set("user_id", userId)
		log.Println(userId)
		c.Next()
	}
}

func parseJWTUnverified(tokenString string) (jwt.MapClaims, error) {
	parser := jwt.NewParser(jwt.WithoutClaimsValidation())

	token, _, err := parser.ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return nil, err
	}

	return token.Claims.(jwt.MapClaims), nil
}
