package common

import (
	"math/rand"
	"os"
	"strconv"
)

const (
	// StandardIDSize Is the standard hash id size. The data id calculated by hash is 32 bits,
	// and the id is a hex string, so the id size is 32 * 8 / 4 = 64.
	StandardIDSize = 64
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

// ConvertArrayStringToByte convert the array sting to []byte
func ConvertArrayStringToByte(arrayStr []string) [][]byte {
	arrayByte := make([][]byte, len(arrayStr))
	for i, str := range arrayStr {
		arrayByte[i] = []byte(str)
	}
	return arrayByte
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
