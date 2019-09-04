package main

import (
	"encoding/json"
	"strings"

	"github.com/gaeanetwork/gaea-core/tee"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// key pair for test
const (
	PrivHexForTest = "0x307702010104204368376222802d1a941f2eb0b7186a2c75f75e368946f923ad37e7c7718c2d7aa00a06082a8648ce3d030107a14403420004e72b30244d2eda1d4b911f8b1fbadafd34017e2d76188924a1d0459a4a6c5c87e122548a9cb9fc93b84af373838af3d0687b81456550a8aae4bec35a9b438e01"
	PubHexForTest  = "0x04e72b30244d2eda1d4b911f8b1fbadafd34017e2d76188924a1d0459a4a6c5c87e122548a9cb9fc93b84af373838af3d0687b81456550a8aae4bec35a9b438e01"
)

// MockTrustedExecutionEnv for tests
type MockTrustedExecutionEnv struct {
}

// Init the system chaincode
func (t *MockTrustedExecutionEnv) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// Invoke the chaincode
func (t *MockTrustedExecutionEnv) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	_, args := stub.GetFunctionAndParameters()
	var sigs []string
	if len(args) > 0 {
		if args[0] == "err" {
			return shim.Error("err")
		} else if args[0] == "jsonerr" {
			return shim.Success([]byte("asdf"))
		} else if strings.Contains(args[0], "sig") {
			sigs = []string{"signature"}
		}
	}

	tee := tee.SharedData{
		Owner:      PubHexForTest,
		Signatures: sigs,
	}
	data, err := json.Marshal(tee)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(data)
}

// MockTEE for tests
type MockTEE struct {
	Result   []byte
	ErrorMsg string
}

// Init the system chaincode
func (t *MockTEE) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// Invoke the chaincode
func (t *MockTEE) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	if t.ErrorMsg != "" {
		return shim.Error(t.ErrorMsg)
	}

	return shim.Success(t.Result)
}
