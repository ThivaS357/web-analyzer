package analyzer

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
)

type AnalysisResult struct {
	URL           string
	Title         string
	HTMLVersion   string
	Headings      map[string]int
	InternalLinks int
	ExternalLinks int
	BrokenLinks   int
	HasLoginForm  bool
}

// AnalyzeURL analyzes the webpage and returns details
func AnalyzeURL(url string) (*AnalysisResult, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP status code: %d", resp.StatusCode)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}

	result := &AnalysisResult{
		URL:      url,
		Headings: make(map[string]int),
	}

	// Initialize WaitGroup and Channel

	workerCount := 5 // Limit concurrent goroutines to 10

	semaphore := make(chan struct{}, workerCount)

	var waitGroup sync.WaitGroup

	channel := make(chan int)

	// Perform HTML parsing
	parseHTML(doc, result, url, semaphore, &waitGroup, channel)

	// Wait for all goroutines to complete
	waitGroup.Wait()
	close(channel)

	// The result.BrokenLinks will be updated with the correct count of broken links
	result.BrokenLinks = len(channel)

	// Detect HTML Version after parsing
	result.HTMLVersion = detectHTMLVersion(doc)

	return result, nil
}

// Recursive HTML parser
func parseHTML(node *html.Node, result *AnalysisResult, baseURL string, semaphore chan struct{}, waitGroup *sync.WaitGroup, channel chan int) {
	if node.Type == html.ElementNode {
		switch node.Data {
		case "title":
			if node.FirstChild != nil {
				result.Title = node.FirstChild.Data
			}
		case "h1", "h2", "h3", "h4", "h5", "h6":
			result.Headings[node.Data]++
		case "a":
			for _, attr := range node.Attr {
				if attr.Key == "href" {
					href := attr.Val
					var link string
					if strings.HasPrefix(attr.Val, "http") {
						result.ExternalLinks++
						link = href

					} else {
						result.InternalLinks++
						link = resolveURL(baseURL, href)
					}

					// Use a goroutine with controlled concurrency
					waitGroup.Add(1)
					semaphore <- struct{}{} // Acquire semaphore
					go func(link string) {
						defer waitGroup.Done()
						defer func() { <-semaphore }() // Release semaphore
						if isLinkBroken(link) {
							channel <- 1
						}
					}(link)
				}
			}
		case "form":
			result.HasLoginForm = isLoginFormInline(node)
		}
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		parseHTML(c, result, baseURL, semaphore, waitGroup, channel)
	}
}

// isLoginFormInline detects login forms during traversal
func isLoginFormInline(formNode *html.Node) bool {
	// Check form-level attributes
	for _, attr := range formNode.Attr {
		if attr.Key == "action" || attr.Key == "id" || attr.Key == "class" || attr.Key == "name" {
			if containsLoginKeyword(attr.Val) {
				return true
			}
		}
	}

	// Inline recursive check for child nodes of the current <form>
	return containsLoginKeywordInNode(formNode)
}

// containsLoginKeywordInNode checks attributes and text content for login keywords
func containsLoginKeywordInNode(node *html.Node) bool {
	if node == nil {
		return false
	}

	// Check text content
	if node.Type == html.TextNode && containsLoginKeyword(node.Data) {
		return true
	}

	// Check attributes
	for _, attr := range node.Attr {
		if containsLoginKeyword(attr.Val) {
			return true
		}
	}

	// Recursively check child nodes
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if containsLoginKeywordInNode(c) {
			return true
		}
	}

	return false
}

// containsLoginKeyword checks if a string contains login-related keywords
func containsLoginKeyword(value string) bool {
	keywords := []string{"login", "signin", "sign-in", "sign_in"}
	lowerValue := strings.ToLower(value)
	for _, keyword := range keywords {
		if strings.Contains(lowerValue, keyword) {
			return true
		}
	}
	return false
}

// Detect HTML version from the parsed document
func detectHTMLVersion(node *html.Node) string {
	for n := node; n != nil; n = n.FirstChild {
		if n.Type == html.DoctypeNode {
			doctype := strings.ToLower(n.Data)
			if doctype == "html" {
				return "HTML5"
			}
			if strings.Contains(doctype, "xhtml 1.0") {
				return "XHTML 1.0"
			}
			if strings.Contains(doctype, "html 4.01") {
				return "HTML 4.01"
			}
		}
	}
	return "Unknown"
}

// checkLink concurrently checks if a link is broken
/**func checkLink(link string, channel chan int, wg *sync.WaitGroup, semaphore chan struct{}) {
	defer wg.Done()
	defer func() { <-semaphore }() // Release semaphore
	if isLinkBroken(link) {
		channel <- 1
	}
}**/

// Check if a link is broken
func isLinkBroken(link string) bool {

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(link)
	if err != nil || resp.StatusCode >= 400 {
		return true
	}
	defer resp.Body.Close()
	return false
}

// Resolve relative URLs to absolute
func resolveURL(baseURL, relative string) string {
	if strings.HasPrefix(relative, "http") {
		return relative
	}
	if strings.HasPrefix(relative, "/") {
		return baseURL + relative
	}
	return baseURL + "/" + relative
}
