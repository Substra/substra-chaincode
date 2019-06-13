package main

import (
	"chaincode/errors"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"reflect"
	"strings"

	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
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

// sliceEqual tells whether a and b contain the same elements.
func sliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for _, v := range a {
		if !stringInSlice(v, b) {
			return false
		}
	}
	return true
}

// getFieldName returns a slice containing field names of a struc
func getFieldNames(v interface{}) (fieldNames []string) {
	e := reflect.ValueOf(v).Elem()
	eType := e.Type()
	for i := 0; i < e.NumField(); i++ {
		varName := eType.Field(i).Name
		fieldNames = append(fieldNames, varName)
	}
	return
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

// stringToInputStruct fills fields of a input struct (such as defined in ledger.go) with elements stored in a slice of string
func stringToInputStruct(args []string, v interface{}) {
	fieldNames := getFieldNames(v)
	e := reflect.ValueOf(v).Elem()
	for i, fn := range fieldNames {
		f := e.FieldByName(fn)
		f.SetString(args[i])
	}
}

// getTxCreator returns the sha256 of the creator of the transaction
func getTxCreator(stub shim.ChaincodeStubInterface) (string, error) {
	// get the agent submitting the transaction
	bCreator, err := stub.GetCreator()
	if err != nil {
		return "", err
	}
	// get pem certificate only. This might be slightly dirty, but this is to avoid installing external packages
	// change it once github.com/hyperledger/fabric/core/chaincode/lib/cid is in fabric chaincode docker
	cert_prefix := "-----BEGIN CERTIFICATE-----"
	cert_suffix := "-----END CERTIFICATE-----\n"
	var creator string
	if sCreator := strings.Split(string(bCreator), cert_prefix); len(sCreator) > 1 {
		creator = strings.Split(sCreator[1], cert_suffix)[0]
	} else {
		creator = "test"
	}
	creator = cert_prefix + creator + cert_suffix
	tt := sha256.Sum256([]byte(creator))
	return hex.EncodeToString(tt[:]), nil
}

// bytesToStruct converts bytes to one a the struct corresponding to elements stored in the ledger
func bytesToStruct(elementBytes []byte, element interface{}) error {
	return json.Unmarshal(elementBytes, &element)
}

// getElementBytes checks if an element is stored in the ledger given its key, and returns associated bytes
func getElementBytes(stub shim.ChaincodeStubInterface, elementKey string) ([]byte, error) {
	elementBytes, err := stub.GetState(elementKey)
	if err != nil {
		return nil, err
	} else if elementBytes == nil {
		return nil, fmt.Errorf("no element with key %s", elementKey)
	}
	return elementBytes, nil
}

// getElementStruct fills an element struct given its key
func getElementStruct(stub shim.ChaincodeStubInterface, elementKey string, element interface{}) error {
	elementBytes, err := getElementBytes(stub, elementKey)
	if err != nil {
		return errors.NotFound(err)
	}
	return bytesToStruct(elementBytes, element)
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

// checkExist checks if keys in a slice correspond to existing elements in the ledger
// returns the slice of already existing elements
func checkExist(stub shim.ChaincodeStubInterface, keys []string) (existingKeys []string) {
	for _, key := range keys {
		if elementBytes, _ := stub.GetState(key); elementBytes != nil {
			existingKeys = append(existingKeys, key)
		}
	}
	return
}

// createCompositeKey creates a composite key given an indexName and attributes
// (combination of attributes to form a key)
func createCompositeKey(stub shim.ChaincodeStubInterface, indexName string, attributes []string) error {
	compositeKey, err := stub.CreateCompositeKey(indexName, attributes)
	if err != nil {
		return err
	}
	value := []byte{0x00}
	if err = stub.PutState(compositeKey, value); err != nil {
		return fmt.Errorf("failed to add composite key with index %s to the ledger", indexName)
	}
	return nil
}

// getKeysFromComposite returns element keys associated with a composite key specified by its indexName and attributes
func getKeysFromComposite(stub shim.ChaincodeStubInterface, indexName string, attributes []string) ([]string, error) {
	elementKeys := make([]string, 0)
	compositeIterator, err := stub.GetStateByPartialCompositeKey(indexName, attributes)
	if err != nil {
		return elementKeys, err
	}
	defer compositeIterator.Close()
	for i := 0; compositeIterator.HasNext(); i++ {
		compositeKey, err := compositeIterator.Next()
		if err != nil {
			return elementKeys, err
		}
		_, compositeKeyParts, err := stub.SplitCompositeKey(compositeKey.Key)
		if err != nil {
			return elementKeys, err
		}
		elementKeys = append(elementKeys, compositeKeyParts[len(compositeKeyParts)-1])
	}
	return elementKeys, nil
}

// updateCompositeKey modifies composite keys
func updateCompositeKey(stub shim.ChaincodeStubInterface, indexName string, oldAttributes []string, newAttributes []string) error {
	oldCompositeKey, err := stub.CreateCompositeKey(indexName, oldAttributes)
	if err != nil {
		return err
	}
	if element, _ := stub.GetState(oldCompositeKey); element == nil {
		return fmt.Errorf("old composite key does not exist - %s", oldCompositeKey)
	}
	if err = stub.DelState(oldCompositeKey); err != nil {
		return err
	}
	newCompositeKey, err := stub.CreateCompositeKey(indexName, newAttributes)
	if err != nil {
		return err
	}
	value := []byte{0x00}
	return stub.PutState(newCompositeKey, value)
}

// getElementsPayload takes as input a list of keys and returns a paylaod containing a list of associated retrieved elements
func getElementsPayload(stub shim.ChaincodeStubInterface, elementsKeys []string) (elements []map[string]interface{}, err error) {

	for _, key := range elementsKeys {
		var element map[string]interface{}
		if err = getElementStruct(stub, key, &element); err != nil {
			return
		}
		element["key"] = key
		elements = append(elements, element)
	}
	return
}

// AssetFromJSON unmarshal a stringify json into the passed interface
// TODO: Validate the interface here if possible
func AssetFromJSON(args string, asset interface{}) error {
	err := json.Unmarshal([]byte(args), &asset)
	if err != nil {
		return errors.BadRequest(err, "Problem when reading json arg: %s, error is:", args)
	}
	return nil
}
