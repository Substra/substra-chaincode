package main

import (
	"fmt"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// queryFilter returns all elements of the ledger matching some filters
// For now, ok for everything. Later returns if the requester has permission to see it
func queryFilter(stub shim.ChaincodeStubInterface, args []string) (elements []map[string]interface{}, err error) {
	expectedArgs := [2]string{"indexName", "attributes"}
	if nbArgs := len(expectedArgs); nbArgs != len(args) {
		err = fmt.Errorf("incorrect arguments, expecting %d args: %s", nbArgs, strings.Join(expectedArgs[:], ", "))
		return
	}
	// check validity of inputs
	indexName := args[0]
	validIndexNames := []string{
		"traintuple~worker~status",
		"testtuple~worker~status",
		"testtuple~tag",
		"traintuple~tag"}
	if !stringInSlice(indexName, validIndexNames) {
		err = fmt.Errorf("invalid indexName filter query: %s", indexName)
		return
	}
	indexName = indexName + "~key"
	attributes := strings.Split(strings.Replace(args[1], " ", "", -1), ",")
	attributes = append([]string{strings.Split(indexName, "~")[0]}, attributes...)

	filteredKeys, err := getKeysFromComposite(stub, indexName, attributes)
	if err != nil {
		err = fmt.Errorf("issue getting keys from composite key %s - %s", indexName, err.Error())
		return
	}
	// get elements with filtererd keys
	switch indexName {
	case "testtuple~worker~status~key", "testtuple~tag~key":
		elements, err = getElementsPayload(stub, filteredKeys)
	case "traintuple~worker~status~key", "traintuple~tag~key":
		elements, err = getTraintuplesPayload(stub, filteredKeys)
	}
	return
}
