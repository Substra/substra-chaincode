package main

import (
	"fmt"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// queryFilter returns all elements of the ledger matching some filters
// For now, ok for everything. Later returns if the requester has permission to see it
func queryFilter(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	expectedArgs := [2]string{"indexName", "attributes"}
	if nbArgs := len(expectedArgs); nbArgs != len(args) {
		return nil, fmt.Errorf("incorrect arguments, expecting %d args: %s", nbArgs, strings.Join(expectedArgs[:], ", "))
	}
	// check validity of inputs
	indexName := args[0]
	validIndexNames := []string{"traintuple~worker~status", "testtuple~worker~status"}
	if !stringInSlice(indexName, validIndexNames) {
		return nil, fmt.Errorf("invalid indexName filter query: %s", indexName)
	}
	indexName = indexName + "~key"
	attributes := strings.Split(strings.Replace(args[1], " ", "", -1), ",")
	attributes = append([]string{strings.Split(indexName, "~")[0]}, attributes...)

	filteredKeys, err := getKeysFromComposite(stub, indexName, attributes)
	if err != nil {
		return nil, fmt.Errorf("issue getting keys from composite key %s - %s", indexName, err.Error())
	}
	// get elements with filtererd keys
	var payload []byte
	switch indexName {
	case "testtuple~worker~status~key":
		payload, err = getElementsPayload(stub, filteredKeys)
	case "traintuple~worker~status~key":
		payload, err = getTraintuplesPayload(stub, filteredKeys)
	}
	return payload, err
}
