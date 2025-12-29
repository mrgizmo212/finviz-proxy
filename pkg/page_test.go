package pkg

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_fetchFinvizPage(t *testing.T) {
	// init chrome allocator
	_, cancel := InitChromeAllocator(context.Background())
	defer cancel()

	html, err := fetchFinvizPage(context.Background(), "", false)
	assert.NoError(t, err)
	os.WriteFile("screener.html", html, 0644)
}
