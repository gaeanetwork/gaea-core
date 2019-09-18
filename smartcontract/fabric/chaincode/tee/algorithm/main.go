package main

import (
	"fmt"

	"github.com/gaeanetwork/gaea-core/smartcontract/fabric/chaincode/tee/algorithm/algorithm"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func main() {
	err := shim.Start(new(algorithm.ChaincodeService))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
