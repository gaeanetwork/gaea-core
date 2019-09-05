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
	PrivHexForTest = "308187020100301306072a8648ce3d020106082a8648ce3d030107046d306b02010104202d130ea6dac76fcae718fbd20bf146643aa66fe6e5902975d2c5ed6ab3bcb5e2a144034200048f03f8321b00a4466f4bf4be51c91898cd50d8cc64c6ecf53e73443e348d5925a16f88c8952b78ebac2dc277a2cc54c77b4c3c07830f49629b689edf63086293"
	PubHexForTest  = "048f03f8321b00a4466f4bf4be51c91898cd50d8cc64c6ecf53e73443e348d5925a16f88c8952b78ebac2dc277a2cc54c77b4c3c07830f49629b689edf63086293"
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
