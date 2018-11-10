package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/martinohmann/rfoutlet/internal/middleware"
	"github.com/stretchr/testify/assert"
)

func TestCors(t *testing.T) {
	tests := []struct {
		given    []string
		expected string
	}{
		{
			given:    []string{"*"},
			expected: "*",
		},
		{
			given:    []string{"foo.com", "bar.org"},
			expected: "foo.com, bar.org",
		},
	}

	for _, tt := range tests {
		r := gin.New()
		r.Use(middleware.Cors(tt.given...))
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)

		r.ServeHTTP(rr, req)

		header := rr.Header().Get("Access-Control-Allow-Origin")

		assert.Equal(t, tt.expected, header)
	}
}

func TestRedirect(t *testing.T) {
	url := "http://localhost/somepath"

	r := gin.New()
	r.GET("/", middleware.Redirect(url))
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)

	r.ServeHTTP(rr, req)

	header := rr.Header().Get("Location")

	assert.Equal(t, http.StatusMovedPermanently, rr.Code)
	assert.Equal(t, url, header)
}
