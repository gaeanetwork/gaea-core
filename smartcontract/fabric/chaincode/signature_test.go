package chaincode

import (
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// key pair for test
const (
	privHexForTests = "308187020100301306072a8648ce3d020106082a8648ce3d030107046d306b02010104202d130ea6dac76fcae718fbd20bf146643aa66fe6e5902975d2c5ed6ab3bcb5e2a144034200048f03f8321b00a4466f4bf4be51c91898cd50d8cc64c6ecf53e73443e348d5925a16f88c8952b78ebac2dc277a2cc54c77b4c3c07830f49629b689edf63086293"
	pubHexForTests  = "048f03f8321b00a4466f4bf4be51c91898cd50d8cc64c6ecf53e73443e348d5925a16f88c8952b78ebac2dc277a2cc54c77b4c3c07830f49629b689edf63086293"
)

var (
	preArgs = []string{"Ciphertext", "Summary", "Description", "Owner"}
)

func Test_CheckArgsContainsHashAndSignatures(t *testing.T) {
	args := make([]string, 6)
	copy(args, preArgs)
	privBytes, _ := hex.DecodeString(privHexForTests)
	hashBytes, data, err := GetArgsHashAndSignatures(privBytes, preArgs)
	assert.NoError(t, err)
	args[4], args[5] = hex.EncodeToString(hashBytes), string(data)

	err = CheckArgsContainsHashAndSignatures(args, pubHexForTests)
	assert.NoError(t, err)

	// Invalid hash
	args1 := make([]string, 6)
	copy(args1, args)
	args1[4] = "1234"
	err1 := CheckArgsContainsHashAndSignatures(args1, pubHexForTests)
	assert.Contains(t, err1.Error(), "Invalid args hash")

	// Invalid signature
	args2 := make([]string, 6)
	copy(args2, args)
	sigs1 := []string{"1234"}
	data1, err1 := json.Marshal(sigs1)
	assert.NoError(t, err1)
	args2[5] = string(data1)
	err2 := CheckArgsContainsHashAndSignatures(args2, pubHexForTests)
	assert.Contains(t, err2.Error(), "Failed to verify ecdsa signature")

	// Invalid pubkey
	err3 := CheckArgsContainsHashAndSignatures(args, "0x1234")
	assert.Contains(t, err3.Error(), "Failed to verify ecdsa signature")

	// Invalid args length
	err4 := CheckArgsContainsHashAndSignatures(args[:1], pubHexForTests)
	assert.Contains(t, err4.Error(), "Invalid args structure")

	args3 := make([]string, 8)
	copy(args3, args)
	err5 := CheckArgsContainsHashAndSignatures(args3, pubHexForTests)
	assert.Contains(t, err5.Error(), "Invalid args hash")
}
