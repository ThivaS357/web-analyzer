package main

import (
	"net/http"

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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.HTML(http.StatusOK, "index.html", results)
}
