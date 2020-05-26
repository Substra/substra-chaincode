package main

import (
	"chaincode/substra"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// main function starts up the chaincode in the container during instantiate
func main() {
	// TODO use the same level as the shim or an env variable
	// logger.SetLevel(shim.LogDebug)
	if err := shim.Start(new(substra.Chaincode)); err != nil {
		fmt.Printf("Error starting SubstraChaincode chaincode: %s", err)
	}
}
