package main

import (
	"chaincode/errors"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"gopkg.in/go-playground/validator.v9"
)

// stringInSlice check if a string is in a slice
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// inputStructToBytes converts fields of a struct (with string fields only, such as input struct defined in ledger.go) to a [][]byte
func inputStructToBytes(v interface{}) (sb [][]byte, err error) {

	e := reflect.Indirect(reflect.ValueOf(v))
	for i := 0; i < e.NumField(); i++ {
		v := e.Field(i)
		if v.Type().Name() != "string" {
			err = fmt.Errorf("struct should contain only string values")
			return
		}
		varValue := v.String()
		sb = append(sb, []byte(varValue))
	}
	return

}

// bytesToStruct converts bytes to one a the struct corresponding to elements stored in the ledger
func bytesToStruct(elementBytes []byte, element interface{}) error {
	return json.Unmarshal(elementBytes, &element)
}

// checkHashes checks if all elements in a slice are all hashes, returns error if not the case
func checkHashes(hashes []string) (err error) {
	for _, hash := range hashes {
		// check validity of dataSampleHashes
		if len(hash) != 64 {
			err = fmt.Errorf("invalid hash %s", hash)
			return
		}
	}
	return
}

// AssetFromJSON unmarshal a stringify json into the passed interface
func AssetFromJSON(args []string, asset interface{}) error {
	if len(args) != 1 {
		return errors.BadRequest("arguments should only contains 1 json string, received: %s", args)
	}
	arg := args[0]
	err := json.Unmarshal([]byte(arg), &asset)
	if err != nil {
		return errors.BadRequest(err, "problem when reading json arg: %s, error is:", arg)
	}
	v := validator.New()
	err = v.Struct(asset)
	if err != nil {
		return errors.BadRequest(err, "inputs validation failed: %s, error is:", arg)
	}
	return nil
}

// SendTuplesEvent sends an event with updated traintuples and testtuples
// Only one event can be sent per transaction
func SendTuplesEvent(stub shim.ChaincodeStubInterface, event interface{}) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}
	err = stub.SetEvent("tuples-updated", payload)
	if err != nil {
		return err
	}
	return nil
}

// GetTxCreator returns the transaction creator
func GetTxCreator(stub shim.ChaincodeStubInterface) (string, error) {
	// get the agent submitting the transaction
	bCreator, err := stub.GetCreator()
	if err != nil {
		return "", err
	}
	// get pem certificate only. This might be slightly dirty, but this is to avoid installing external packages
	// change it once github.com/hyperledger/fabric/core/chaincode/lib/cid is in fabric chaincode docker
	certPrefix := "-----BEGIN CERTIFICATE-----"
	certSuffix := "-----END CERTIFICATE-----\n"
	var creator string
	if sCreator := strings.Split(string(bCreator), certPrefix); len(sCreator) > 1 {
		creator = strings.Split(sCreator[1], certSuffix)[0]
	} else {
		creator = "test"
	}
	creator = certPrefix + creator + certSuffix
	tt := sha256.Sum256([]byte(creator))
	return hex.EncodeToString(tt[:]), nil
}
