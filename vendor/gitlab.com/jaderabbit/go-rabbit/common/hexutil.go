package common

import (
	"encoding/hex"
	"fmt"
)

// BytesToHex encodes b as a hex string with 0x prefix.
func BytesToHex(b []byte) string {
	enc := make([]byte, len(b)*2+2)
	copy(enc, "0x")
	hex.Encode(enc[2:], b)
	return string(enc)
}

// HexToBytes decodes a hex string with 0x prefix.
func HexToBytes(input string) ([]byte, error) {
	if len(input) == 0 {
		return nil, fmt.Errorf("Error hex string is empty")
	}

	// MissingPrefix
	if !Has0xPrefix(input) {
		return nil, fmt.Errorf("Error hex string does not have a 0x / 0X prefix")
	}

	return hex.DecodeString(input[2:])
}

// Has0xPrefix returns true if input start with 0x, otherwise false
func Has0xPrefix(input string) bool {
	return len(input) >= 2 && input[0] == '0' && (input[1] == 'x' || input[1] == 'X')
}
