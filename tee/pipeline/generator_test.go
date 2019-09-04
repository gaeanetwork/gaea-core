package pipeline

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Generator(t *testing.T) {
	// test string
	result := Generator(nil, []interface{}{"1", "1"}...)
	for value := range result {
		v, ok := value.(string)
		assert.True(t, ok)
		assert.Equal(t, v, "1")
	}

	// test int
	result = Generator(nil, []interface{}{1, 1}...)
	for value := range result {
		v, ok := value.(int)
		assert.True(t, ok)
		assert.Equal(t, v, 1)
	}

	// test done
	done := make(chan interface{})
	close(done)
	result = Generator(done, nil)
	count := 0
	for range result {
		count++
	}
	assert.Zero(t, count)

	// test nil
	result = Generator(nil, nil)
	v, ok := <-result
	assert.True(t, ok)
	assert.Nil(t, v)
}
