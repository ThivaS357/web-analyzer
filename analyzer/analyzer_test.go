package analyzer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnalyzeURL_ValidPage(t *testing.T) {
	url := "https://google.com"
	result, err := AnalyzeURL(url)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.Title)
	assert.Greater(t, result.InternalLinks, 0)
}

func TestAnalyzeURL_InvalidURL(t *testing.T) {
	url := "invalid-url"
	result, err := AnalyzeURL(url)
	assert.Error(t, err)
	assert.Nil(t, result)
}
