package system

import (
	"fmt"
	"math/big"
	"net"
)

// ReverseArray reverse array
func ReverseArray(arrayStr []string) []string {
	for i, j := 0, len(arrayStr)-1; i < j; i, j = i+1, j-1 {
		arrayStr[i], arrayStr[j] = arrayStr[j], arrayStr[i]
	}
	return arrayStr
}

// FormatError format the error of number of arguments
func FormatError(operation string, actual, expected int) string {
	return fmt.Sprintf("failed to %s, Incorrect number of arguments(%d). Expecting %d", operation, actual, expected)
}

// InetNtoA Convert uint to net.IP
func InetNtoA(ip int64) string {
	return fmt.Sprintf("%d.%d.%d.%d",
		byte(ip>>24), byte(ip>>16), byte(ip>>8), byte(ip))
}

// InetAtoN Convert net.IP to int64
func InetAtoN(ip string) int64 {
	ret := big.NewInt(0)
	ret.SetBytes(net.ParseIP(ip).To4())
	return ret.Int64()
}
