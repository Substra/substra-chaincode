package main

import (
	"fmt"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// query returns an element of the ledger given its key
// For now, ok for everything. Later returns if the requester has permission to see it
func query(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	expectedArgs := [1]string{"element key"}
	if nbArgs := len(expectedArgs); nbArgs != len(args) {
		return nil, fmt.Errorf("incorrect arguments, expecting %d args: %s", nbArgs, strings.Join(expectedArgs[:], ", "))
	}

	key := args[0]
	if len(key) != 64 {
		return nil, fmt.Errorf("invalid key")
	}
	return getElementBytes(stub, key)
}

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

// queryAll returns all elements of the ledger given its type
// It works for challenges, dataset, algo, and testtuple, not for traintuples (see queryTraintuples)
// For now, ok for everything. Later returns if the requester has permission to see it
func queryAll(stub shim.ChaincodeStubInterface, args []string, elementType string) ([]byte, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("incorrect number of arguments, expecting nothing")
	}
	var indexName string
	switch elementType {
	case "challenge":
		indexName = "challenge~owner~key"
	case "dataset":
		indexName = "dataset~owner~key"
	case "algo":
		indexName = "algo~challenge~key"
	case "testtuple":
		indexName = "testtuple~traintuple~certified~key"
	default:
		return nil, fmt.Errorf("no element type %s", elementType)
	}
	elementsKeys, err := getKeysFromComposite(stub, indexName, []string{elementType})
	if err != nil {
		return nil, fmt.Errorf("issue getting keys from composite key %s - %s", indexName, err.Error())
	}
	return getElementsPayload(stub, elementsKeys)
}
