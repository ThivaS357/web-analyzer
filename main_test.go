package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAnalyzeHandlerWithRealLogic(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	// Serve HTML templates
	router.LoadHTMLGlob("templates/*")
	router.POST("/analyze", analyzeHandler)

	// Define a test case with a valid URL
	body := strings.NewReader("url=https://google.com")
	req, _ := http.NewRequest(http.MethodPost, "/analyze", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp := httptest.NewRecorder()

	// Call the handler
	router.ServeHTTP(resp, req)

	// Check the response
	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestAnalyzeHandlerWithInvalidURL(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.POST("/analyze", analyzeHandler)

	// Define a test case with an invalid URL
	body := strings.NewReader("url=invalid-url")
	req, _ := http.NewRequest(http.MethodPost, "/analyze", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp := httptest.NewRecorder()

	// Call the handler
	router.ServeHTTP(resp, req)

	// Check the response
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
}
