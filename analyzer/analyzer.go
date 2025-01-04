package analyzer

import (
	"fmt"
	"net/http"
	"strings"

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

	// Perform HTML parsing
	parseHTML(doc, result, url)

	// Detect HTML Version after parsing
	result.HTMLVersion = detectHTMLVersion(doc)

	return result, nil
}

// Recursive HTML parser
func parseHTML(node *html.Node, result *AnalysisResult, baseURL string) {
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
					if strings.HasPrefix(attr.Val, "http") {
						result.ExternalLinks++
						if isLinkBroken(href) {
							result.BrokenLinks++
						}
					} else {
						result.InternalLinks++
						absoluteURL := resolveURL(baseURL, href)
						if isLinkBroken(absoluteURL) {
							result.BrokenLinks++
						}
					}
				}
			}
		case "form":
			for _, attr := range node.Attr {
				if attr.Key == "action" && strings.Contains(attr.Val, "login") {
					result.HasLoginForm = true
				}
			}
		}
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		parseHTML(c, result, baseURL)
	}
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

// Check if a link is broken
func isLinkBroken(link string) bool {
	resp, err := http.Get(link)
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
