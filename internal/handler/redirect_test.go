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
	url := "http://localhost/somepath"

	r := gin.New()
	r.GET("/", handler.Redirect(url))
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)

	r.ServeHTTP(rr, req)

	header := rr.Header().Get("Location")

	assert.Equal(t, http.StatusMovedPermanently, rr.Code)
	assert.Equal(t, url, header)
}
