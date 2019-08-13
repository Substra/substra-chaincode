package main

import (
	"fmt"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// queryFilter returns all elements of the ledger matching some filters
// For now, ok for everything. Later returns if the requester has permission to see it
func queryFilter(stub shim.ChaincodeStubInterface, args []string) (elements interface{}, err error) {
	inp := inputQueryFilter{}
	err = AssetFromJSON(args, &inp)
	if err != nil {
		return
	}
	// check validity of inputs
	validIndexNames := []string{
		"traintuple~worker~status",
		"testtuple~worker~status",
		"testtuple~tag",
		"traintuple~tag"}
	if !stringInSlice(inp.IndexName, validIndexNames) {
		err = fmt.Errorf("invalid indexName filter query: %s", inp.IndexName)
		return
	}
	indexName := inp.IndexName + "~key"
	attributes := strings.Split(strings.Replace(inp.Attributes, " ", "", -1), ",")
	attributes = append([]string{strings.Split(indexName, "~")[0]}, attributes...)

	filteredKeys, err := getKeysFromComposite(stub, indexName, attributes)
	if err != nil {
		err = fmt.Errorf("issue getting keys from composite key %s - %s", indexName, err.Error())
		return
	}
	// get elements with filtererd keys
	switch indexName {
	case "testtuple~worker~status~key", "testtuple~tag~key":
		elements, err = getOutputTesttuples(stub, filteredKeys)
	case "traintuple~worker~status~key", "traintuple~tag~key":
		elements, err = getOutputTraintuples(stub, filteredKeys)
	}
	return
}
