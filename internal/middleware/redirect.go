package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Redirect redirects a request to given url
func Redirect(url string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, url)
	}
}
