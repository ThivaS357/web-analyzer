package main

import (
	"net/http"
	"strings"

	"web-analyzer/analyzer"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	// Create Gin router
	router := gin.Default()

	// Serve HTML templates
	router.LoadHTMLGlob("templates/*")

	// Routes
	router.GET("/", homeHandler)
	router.POST("/analyze", analyzeHandler)

	// Start server
	router.Run(":8080")
}

// Serve the home page
func homeHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
}

// Handle URL analysis
func analyzeHandler(c *gin.Context) {
	url := c.PostForm("url")
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL is required"})
		return
	}

	results, err := analyzer.AnalyzeURL(url)
	if err != nil {
		handleSpecificErrors(c, err)
		return
	}

	c.HTML(http.StatusOK, "index.html", results)
}

// Handle specific errors
func handleSpecificErrors(c *gin.Context, err error) {
	// Example error messages or custom error handling
	if strings.Contains(err.Error(), "HTTP status code: 400") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request (400): Invalid input URL"})
	} else if strings.Contains(err.Error(), "HTTP status code: 502") {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Bad Gateway (502): Upstream server error"})
	} else if strings.Contains(err.Error(), "HTTP status code: 503") {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Service Unavailable (503): The server is overloaded or down"})
	} else if strings.Contains(err.Error(), "HTTP status code: 504") {
		c.JSON(http.StatusGatewayTimeout, gin.H{"error": "Gateway Timeout (504): The server took too long to respond"})
	} else {
		// General server error
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
