package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// Cors sets cors headers on response
func Cors(origins ...string) gin.HandlerFunc {
	header := strings.Join(origins, ", ")

	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", header)
		c.Next()
	}
}
