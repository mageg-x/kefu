// middleware/auth.go
package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"kefu-server/utils"
	"kefu-server/utils/logger"
	"kefu-server/utils/response"
)

// AuthMiddleware returns gin.HandlerFunc
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Errorf("authorization token not provided")
			response.ResponseError(c, http.StatusUnauthorized, response.ErrCodeUnauthorized)
			c.Abort()
			return
		}

		// Check if it starts with "Bearer "
		tokenStr := ""
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenStr = authHeader[7:] // Remove "Bearer "
		} else {
			logger.Errorf("invalid token format")
			response.ResponseError(c, http.StatusUnauthorized, response.ErrCodeTokenInvalid)
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(tokenStr)
		if err != nil {
			logger.Errorf("token invalid or expired: %v", err)
			response.ResponseError(c, http.StatusUnauthorized, response.ErrCodeTokenExpired)
			c.Abort()
			return
		}

		// Store user information in context for subsequent handlers
		c.Set("userID", claims.UserID)
		c.Set("userName", claims.UserName)
		c.Set("role", claims.Role)

		c.Next() // Continue to next middleware or handler
	}
}
