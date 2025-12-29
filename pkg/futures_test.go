package pkg

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_FetchAllFutures(t *testing.T) {
	_, cancel := InitChromeAllocator(context.Background())
	defer cancel()

	futures, err := FetchAllFutures(context.Background(), false)
	assert.NoError(t, err)
	assert.NotNil(t, futures)
	assert.NotEmpty(t, futures)
}
