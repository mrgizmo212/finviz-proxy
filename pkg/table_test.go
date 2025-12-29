package pkg

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseTable(t *testing.T) {
	_, cancel := InitChromeAllocator(context.Background())
	defer cancel()

	page, err := fetchFinvizPage(context.Background(), "", false)
	assert.NoError(t, err)
	table, err := parseTable(page)
	assert.NoError(t, err)
	j, err := json.MarshalIndent(table, "", "  ")
	assert.NoError(t, err)
	println(string(j))
}