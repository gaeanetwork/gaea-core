package pipeline

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_FanIn(t *testing.T) {
	// int channel
	c1 := make(chan interface{})
	count := 5
	go func() {
		defer close(c1)
		for index := 0; index < count; index++ {
			c1 <- index
		}
	}()

	// string channel
	c2 := make(chan interface{})
	go func() {
		defer close(c2)
		for index := 0; index < count; index++ {
			c2 <- fmt.Sprintf("%v", index)
		}
	}()

	index1, index2 := 0, 0
	for value := range FanIn(nil, []<-chan interface{}{c1, c2}...) {
		switch value.(type) {
		case int:
			assert.Equal(t, index1, value)
			index1++
		case string:
			assert.Equal(t, fmt.Sprintf("%v", index2), value)
			index2++
		}
	}
}
