package pkg

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_fetchParams(t *testing.T) {
	_, cancel := InitChromeAllocator(context.Background())
	defer cancel()
	params, err := FetchParams(context.Background(), false)
	assert.NoError(t, err)
	j, err := json.MarshalIndent(params, "", "  ")
	assert.NoError(t, err)
	println(string(j))
}
