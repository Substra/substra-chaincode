package main

import (
	"fmt"

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
	algo.Permissions = inp.Permissions
	return
}

// -------------------------------------------------------------------------------------------
// Smart contracts related to an algo
// -------------------------------------------------------------------------------------------
// registerAlgo stores a new algo in the ledger.
// If the key exists, it will override the value with the new one
func registerAlgo(stub shim.ChaincodeStubInterface, args []string) (resp map[string]string, err error) {
	inp := inputAlgo{}
	err = AssetFromJSON(args[0], &inp)
	if err != nil {
		return
	}
	// check validity of input args and convert it to Algo
	algo := Algo{}
	algoKey, err := algo.Set(stub, inp)
	if err != nil {
		return
	}
	// check data is not already in ledgert
	if elementBytes, _ := stub.GetState(algoKey); elementBytes != nil {
		err = fmt.Errorf("algo with this hash already exists")
		return
	}
	// submit to ledger
	algoBytes, _ := json.Marshal(algo)
	err = stub.PutState(algoKey, algoBytes)
	if err != nil {
		err = fmt.Errorf("failed to add to ledger algo with key %s with error %s", algoKey, err.Error())
		return
	}
	// create composite key
	err = createCompositeKey(stub, "algo~owner~key", []string{"algo", algo.Owner, algoKey})
	if err != nil {
		return
	}
	return map[string]string{"key": algoKey}, nil
}

// queryAlgo returns an algo of the ledger given its key
func queryAlgo(stub shim.ChaincodeStubInterface, args []string) (out outputAlgo, err error) {
	inp := inputHashe{}
	err = AssetFromJSON(args[0], &inp)
	if err != nil {
		return
	}
	var algo Algo
	if err = getElementStruct(stub, inp.Key, &algo); err != nil {
		return
	}
	out.Fill(inp.Key, algo)
	return
}

// queryAlgos returns all algos of the ledger
func queryAlgos(stub shim.ChaincodeStubInterface, args []string) (outAlgos []outputAlgo, err error) {
	if len(args) != 0 {
		err = fmt.Errorf("incorrect number of arguments, expecting nothing")
		return
	}
	var indexName = "algo~owner~key"
	elementsKeys, err := getKeysFromComposite(stub, indexName, []string{"algo"})
	if err != nil {
		err = fmt.Errorf("issue getting keys from composite key %s - %s", indexName, err.Error())
		return
	}
	for _, key := range elementsKeys {
		var algo Algo
		if err = getElementStruct(stub, key, &algo); err != nil {
			return
		}
		var out outputAlgo
		out.Fill(key, algo)
		outAlgos = append(outAlgos, out)
	}
	return
}
