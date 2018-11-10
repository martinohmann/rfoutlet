package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Healthz responds with status 200 if the application is healthy
func Healthz(c *gin.Context) {
	c.String(http.StatusOK, "ok")
}
