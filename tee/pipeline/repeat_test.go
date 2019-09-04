package pipeline

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_RepeatFn(t *testing.T) {
	getRand := func() interface{} { return rand.Int63n(150) }
	resultStream := RepeatFn(nil, getRand, 10)
	count := 0
	for result := range resultStream {
		r, ok := result.(int64)
		assert.True(t, ok)
		assert.True(t, r < 150)

		count++
	}

	assert.Equal(t, count, 10)
}
