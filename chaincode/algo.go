package main

import (
	"fmt"
	"strings"

	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"gopkg.in/go-playground/validator.v9"
)

// Set is a method of the receiver Algo. It checks the validity of inputAlgo and uses its fields to set the Algo
// Returns the algoKey
func (algo *Algo) Set(stub shim.ChaincodeStubInterface, inp inputAlgo) (algoKey string, err error) {
	// checking validity of submitted fields
	validate := validator.New()
	if err = validate.Struct(inp); err != nil {
		err = fmt.Errorf("invalid algo inputs %s", err.Error())
		return
	}
	// check associated challenge exists
	_, err = getElementBytes(stub, inp.ChallengeKey)
	if err != nil {
		return
	}
	algoKey = inp.Hash
	// find associated owner
	owner, err := getTxCreator(stub)
	if err != nil {
		return
	}
	// set algo
	algo.Name = inp.Name
	algo.StorageAddress = inp.StorageAddress
	algo.Description = &HashDress{
		Hash:           inp.DescriptionHash,
		StorageAddress: inp.DescriptionStorageAddress,
	}
	algo.Owner = owner
	algo.ChallengeKey = inp.ChallengeKey
	algo.Permissions = inp.Permissions
	return
}

// -------------------------------------------------------------------------------------------
// Smart contracts related to an algo
// -------------------------------------------------------------------------------------------
// registerAlgo stores a new algo in the ledger.
// If the key exists, it will override the value with the new one
func registerAlgo(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	expectedArgs := getFieldNames(&inputAlgo{})
	if nbArgs := len(expectedArgs); nbArgs != len(args) {
		return nil, fmt.Errorf("incorrect arguments, expecting %d args: %s", nbArgs, strings.Join(expectedArgs, ", "))
	}

	// convert input strings args to input struct inputAlgo
	inp := inputAlgo{}
	stringToInputStruct(args, &inp)
	// check validity of input args and convert it to Algo
	algo := Algo{}
	algoKey, err := algo.Set(stub, inp)
	if err != nil {
		return nil, err
	}
	// check data is not already in ledgert
	if elementBytes, _ := stub.GetState(algoKey); elementBytes != nil {
		return nil, fmt.Errorf("algo with this hash already exists")
	}
	// submit to ledger
	algoBytes, _ := json.Marshal(algo)
	err = stub.PutState(algoKey, algoBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to add to ledger algo with key %s with error %s", algoKey, err.Error())
	}
	// create composite key
	err = createCompositeKey(stub, "algo~challenge~key", []string{"algo", algo.ChallengeKey, algoKey})
	if err != nil {
		return nil, err
	}
	return []byte(algoKey), nil
}

// queryAlgo returns an algo of the ledger given its key
func queryAlgo(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 || len(args[0]) != 64 {
		return nil, fmt.Errorf("incorrect arguments, expecting key, received: %s", args[0])
	}
	key := args[0]
	var algo Algo
	if err := getElementStruct(stub, key, &algo); err != nil {
		return nil, err
	}
	var out outputAlgo
	out.Fill(key, algo)
	return json.Marshal(out)
}

// queryAlgos returns all algos of the ledger
func queryAlgos(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("incorrect number of arguments, expecting nothing")
	}
	var indexName = "algo~challenge~key"
	elementsKeys, err := getKeysFromComposite(stub, indexName, []string{"algo"})
	if err != nil {
		return nil, fmt.Errorf("issue getting keys from composite key %s - %s", indexName, err.Error())
	}
	var outAlgos []outputAlgo
	for _, key := range elementsKeys {
		var algo Algo
		if err := getElementStruct(stub, key, &algo); err != nil {
			return nil, err
		}
		var out outputAlgo
		out.Fill(key, algo)
		outAlgos = append(outAlgos, out)
	}
	return json.Marshal(outAlgos)
}
