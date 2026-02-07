package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORS middleware handles cross-origin requests
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Allow requests from all origins
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		// Allowed request headers
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Encryption-Enabled, X-Device-ID")
		// Allowed request methods
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		// Allowed exposed response headers
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		// Allow sending credentials
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle OPTIONS requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
