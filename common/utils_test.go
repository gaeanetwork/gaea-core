package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ContainsStringArray(t *testing.T) {
	src := []string{"1", "2", "3"}

	dest := []string{"1", "2", "3"}
	str, ok := ContainsStringArray(src, dest)
	assert.Equal(t, "", str)
	assert.True(t, ok)

	dest1 := []string{"1", "2"}
	str1, ok1 := ContainsStringArray(src, dest1)
	assert.Equal(t, "", str1)
	assert.True(t, ok1)

	dest2 := []string{"1", "4"}
	str2, ok2 := ContainsStringArray(src, dest2)
	assert.Equal(t, "4", str2)
	assert.False(t, ok2)

	dest3 := []string{"1", "2", "3", "4"}
	str3, ok3 := ContainsStringArray(src, dest3)
	assert.Equal(t, "4", str3)
	assert.False(t, ok3)
}
