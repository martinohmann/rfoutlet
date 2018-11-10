package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/martinohmann/rfoutlet/internal/handler"
	"github.com/stretchr/testify/assert"
)

func TestRedirect(t *testing.T) {
	r := gin.New()
	r.GET("/", handler.Healthz)
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "ok", rr.Body.String())
}
