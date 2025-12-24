package middlewares

import (
	"net/http"
	"strings"

	"github.com/cfex/microservices-in-go/services/users-service/internal/jwt"
	"github.com/gin-gonic/gin"
)

func RequireAuthentication(c *gin.Context) {
	accessToken := c.GetHeader("Authorization")

	if accessToken != "" && strings.HasPrefix(accessToken, "Bearer ") {
		ac := strings.TrimPrefix(accessToken, "Bearer ")
		accessToken = ac
	} else {
		token, err := c.Cookie("access_token")

		if token != "" && err == nil {
			accessToken = token
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
	}

	claims, err := jwt.Jwt.ParseToken(accessToken)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	c.Set("user_id", claims.Payload.UserID)
	c.Set("role", claims.Payload.Role)
	c.Next()
}
