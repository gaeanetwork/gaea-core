package main

import (
	"fmt"

	"github.com/gaeanetwork/gaea-core/smartcontract/fabric/chaincode/tee/data/data"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func main() {
	err := shim.Start(new(data.SharedDataService))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
