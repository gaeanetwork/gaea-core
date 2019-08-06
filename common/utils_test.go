package common

import (
	"os"
	"path/filepath"
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

func Test_FileOrFolderExists(t *testing.T) {
	tmpdir := filepath.Join(os.TempDir(), GetRandomString())

	// Not exists
	exists := FileOrFolderExists(tmpdir)
	assert.False(t, exists)

	// Create
	err := os.MkdirAll(tmpdir, 0755)
	assert.NoError(t, err)
	exists = FileOrFolderExists(tmpdir)
	assert.True(t, exists)

	// Delete
	err = os.RemoveAll(tmpdir)
	assert.NoError(t, err)
	exists = FileOrFolderExists(tmpdir)
	assert.False(t, exists)
}

func Benchmark_GetRandomString(b *testing.B) {
	existsMap := make(map[string]struct{})
	for index := 0; index < b.N; index++ {
		s := GetRandomString()
		if _, exists := existsMap[s]; exists {
			b.FailNow()
		}

		existsMap[s] = struct{}{}
	}
}
