package common

import (
	"crypto/md5"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var (
	src = rand.NewSource(time.Now().UnixNano())

	errInvalidName = errors.New("invalid name, only numbers, letters, and dash lines are allowed")
)

const letterBytes = "abcdefghijklmnopqrstuvwxyz1234567890"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// RandStringBytesMaskImprSrc Get a random string with a length of n
func RandStringBytesMaskImprSrc(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

// FileOrFolderExists checks if a file or folder exists
func FileOrFolderExists(fileOrFolder string) bool {
	_, err := os.Stat(fileOrFolder)
	return !os.IsNotExist(err)
}

// CreateDateDir create dir, first to create dir, second to assign permissions to dir
func CreateDateDir(basePath string) error {
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		if err := os.MkdirAll(basePath, 0777); err != nil {
			return err
		}

		if err := os.Chmod(basePath, 0777); err != nil {
			return err
		}
	}
	return nil
}

// SaveFile save file
func SaveFile(filePath string, content []byte) error {
	// Create the file directory with appropriate permissions
	// in case it is not present yet.
	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return err
	}

	// Atomic write: create a temporary hidden file first then move it into place.
	f, err := ioutil.TempFile(filepath.Dir(filePath), fmt.Sprint(".", filepath.Base(filePath), ".tmp"))
	if err != nil {
		return err
	}

	if _, err := f.Write(content); err != nil {
		f.Close()
		os.Remove(f.Name())

		return err
	}

	f.Close()

	return os.Rename(f.Name(), filePath)
}

// ValidateString string allowed only a-z, A-Z, 0-9, and -
func ValidateString(str string) error {
	if len(str) == 0 {
		return errors.New("string is empty")
	}

	ok, err := regexp.Match(`^[a-zA-Z0-9\-\_]+$`, []byte(str))
	if ok {
		return nil
	}

	if err != nil {
		return err
	}
	return errInvalidName
}

// ValidateStringAllowedEmpty string allowed empty, a-z, A-Z, 0-9, and -
func ValidateStringAllowedEmpty(str string) error {
	if len(str) == 0 {
		return nil
	}

	ok, err := regexp.Match(`^[a-zA-Z0-9\-\_]+$`, []byte(str))
	if ok {
		return nil
	}

	if err != nil {
		return err
	}
	return errInvalidName
}

// ConvertStringToTime convert string to time, check the str
func ConvertStringToTime(str string) (*time.Time, error) {
	timeLayout := "2006-01-02 15:04:05"

	str = strings.ToLower(strings.Trim(str, " "))
	if len(str) == 0 {
		return nil, errors.New("string is empty")
	}

	convertTime, err := time.ParseInLocation(timeLayout, str, time.Local)
	return &convertTime, err
}

// ConvertArrayStringToByte convert the array sting to []byte
func ConvertArrayStringToByte(arrayStr []string) [][]byte {
	arrayByte := make([][]byte, len(arrayStr))
	for i, str := range arrayStr {
		arrayByte[i] = []byte(str)
	}
	return arrayByte
}

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

// MergeArray merge two arrays and return an array of non-duplicate
func MergeArray(array1 []string, array2 []string) []string {
	mapData := make(map[string]byte)
	for _, str := range array1 {
		if len(str) == 0 {
			continue
		}
		mapData[str] = 0
	}

	for _, str := range array2 {
		if len(str) == 0 {
			continue
		}
		mapData[str] = 0
	}

	distArray := []string{}
	for key := range mapData {
		distArray = append(distArray, key)
	}
	return distArray
}

// RemoveDuplicateArray Remove string from array1 and repeat string in array2
func RemoveDuplicateArray(array1 []string, array2 []string) []string {
	mapData := make(map[string]byte)
	for _, str := range array1 {
		if len(str) == 0 {
			continue
		}
		mapData[str] = 0
	}

	for _, str := range array2 {
		if _, ok := mapData[str]; ok {
			if len(str) == 0 {
				continue
			}
			delete(mapData, str)
		}
	}

	distArray := []string{}
	for key := range mapData {
		distArray = append(distArray, key)
	}
	return distArray
}

// GetRepeateElementArray get a array that the element both exists in array1 and array2
func GetRepeateElementArray(array1 []string, array2 []string) []string {
	mapData := make(map[string]byte)
	for _, str := range array1 {
		if len(str) == 0 {
			continue
		}
		mapData[str] = 0
	}

	distArray := []string{}
	for _, str := range array2 {
		if _, ok := mapData[str]; ok {
			if len(str) == 0 {
				continue
			}
			distArray = append(distArray, str)
		}
	}

	return distArray
}

// ConvertArrayToMap convert a array to map struct
func ConvertArrayToMap(array []string) map[string]byte {
	mapData := make(map[string]byte)
	for _, str := range array {
		mapData[str] = 0
	}

	return mapData
}

// RemoveEmptyElementInArray remove empty element from the array
func RemoveEmptyElementInArray(array []string) []string {
	if len(array) == 0 {
		return array
	}

	k := len(array) - 1
	for i, str := range array {
		if k < i {
			continue
		}

		if len(str) == 0 {
			for ; k > i; k-- {
				if len(array[k]) > 0 {
					array[i], array[k] = array[k], array[i]
					break
				}
			}
		}
	}
	return array[0:k]
}

// ReverseArray reverse array
func ReverseArray(arrayStr []string) []string {
	for i, j := 0, len(arrayStr)-1; i < j; i, j = i+1, j-1 {
		arrayStr[i], arrayStr[j] = arrayStr[j], arrayStr[i]
	}
	return arrayStr
}

// GetUserPassWord get the password of user, first md5, then base64 encode
func GetUserPassWord(password []byte) string {
	cipherText := md5.Sum(password)
	return base64.StdEncoding.EncodeToString(cipherText[:])
}
