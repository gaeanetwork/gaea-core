package algorithm

import (
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/stretchr/testify/assert"
)

func Test_Init(t *testing.T) {
	stub := shim.NewMockStub("tee", new(ChaincodeService))

	args := [][]byte{}
	response := stub.MockInit("1", args)
	assert.Equal(t, shim.OK, int(response.Status))
	assert.Empty(t, response.Message)

	args1 := [][]byte{[]byte("asdfasdfasdf")}
	response1 := stub.MockInit("2", args1)
	assert.Equal(t, shim.OK, int(response1.Status))
	assert.Empty(t, response1.Message)
}
