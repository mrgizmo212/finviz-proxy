package pkg

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_fetchAllNews(t *testing.T) {
	_, cancel := InitChromeAllocator(context.Background())
	defer cancel()

	html, err := fetchAllNews(context.Background(), false)
	assert.NoError(t, err)
	os.WriteFile("news.html", html, 0644)
}

func Test_parseNewsAndBlogs(t *testing.T) {
	html, err := os.ReadFile("news.html")
	assert.NoError(t, err)
	news, blogs, err := parseNewsAndBlogs(html)
	assert.NoError(t, err)
	assert.NotEmpty(t, news)
	assert.NotEmpty(t, blogs)
}
