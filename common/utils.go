package common

import (
	"math/rand"
	"os"
	"strconv"
)

// ContainsStringArray return true if the dest is a subset of src, otherwise return false and unmatched string
func ContainsStringArray(src []string, dest []string) (string, bool) {
	stringSet := make(map[string]struct{})
	for _, str := range src {
		stringSet[str] = struct{}{}
	}

	for _, str := range dest {
		if _, ok := stringSet[str]; !ok {
			return str, false
		}
	}

	return "", true
}

// FileOrFolderExists checks if a file or folder exists
func FileOrFolderExists(fileOrFolder string) bool {
	_, err := os.Stat(fileOrFolder)
	return !os.IsNotExist(err)
}

// GetRandomString returns a non-negative pseudo-random 63-bit integer as an int64 string from the default Source.
func GetRandomString() string {
	return strconv.FormatInt(rand.Int63(), 10)
}
