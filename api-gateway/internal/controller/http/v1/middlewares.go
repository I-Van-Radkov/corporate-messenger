package v1

import (
	"context"
	"strings"

	"github.com/I-Van-Radkov/corporate-messenger/api-gateway/internal/clients/identity"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(authClient identity.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
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

		token := parts[1]

		ok, err := authClient.IntrospectToken(context.Background(), token)
		if err != nil {
			c.AbortWithStatusJSON(503, gin.H{
				"error": "Indentity service unavailable",
			})
			return
		}

		if !ok {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token"})
			return
		}

		claims, err := parseJWTUnverified(token)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token format"})
			return
		}

		userId, _ := claims["user_id"].(string)
		userRole, _ := claims["user_role"].(string)

		if userId == "" || userRole == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token payload"})
			return
		}

		c.Request.Header.Set("X-User-ID", userId)
		c.Request.Header.Set("X-User-Role", userRole)

		c.Set("user_id", userId)
		c.Set("user_role", userRole)

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
