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
		c.HTML(http.StatusOK, "index.html", gin.H{
			"Error": "URL is required. Please enter a valid URL.",
		})
		return
	}

	results, err := analyzer.AnalyzeURL(url)

	// Handle 400, 502, 503 and 504 Error
	if err != nil {

		var errorMessage string
		if strings.Contains(err.Error(), "HTTP status code: 400") {
			errorMessage = "Bad Request (400): Invalid input URL."
		} else if strings.Contains(err.Error(), "HTTP status code: 502") {
			errorMessage = "Bad Gateway (502): Upstream server error."
		} else if strings.Contains(err.Error(), "HTTP status code: 503") {
			errorMessage = "Service Unavailable (503): The server is overloaded or down."
		} else if strings.Contains(err.Error(), "HTTP status code: 504") {
			errorMessage = "Gateway Timeout (504): The server took too long to respond."
		} else {
			errorMessage = "An unexpected error occurred: " + err.Error()
		}

		c.HTML(http.StatusOK, "index.html", gin.H{
			"Error": errorMessage,
		})

		return
	}

	// Pass analysis results to the template
	c.HTML(http.StatusOK, "index.html", gin.H{
		"URL":           results.URL,
		"Title":         results.Title,
		"HTMLVersion":   results.HTMLVersion,
		"Headings":      results.Headings,
		"InternalLinks": results.InternalLinks,
		"ExternalLinks": results.ExternalLinks,
		"BrokenLinks":   results.BrokenLinks,
		"HasLoginForm":  results.HasLoginForm,
	})
}
