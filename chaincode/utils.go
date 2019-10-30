// Copyright 2018 Owkin, inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"chaincode/errors"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"

	"crypto/ecdsa"
	"crypto/rsa"

	"crypto/x509"

	"encoding/pem"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/msp"
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
	bCreator, err := stub.GetCreator()

	if err != nil {
		return "", err
	}

	// get pertain protobuf object
	sID := &msp.SerializedIdentity{}
	err = proto.Unmarshal(bCreator, sID)
	if err != nil {
		return "", err
	}

	// get the cert
	block, _ := pem.Decode(sID.IdBytes)
	if block == nil {
		return "", fmt.Errorf("unable to decode block %s", sID.IdBytes)
	}

	// load it as a pem cert
	var cert *x509.Certificate
	cert, _ = x509.ParseCertificate(block.Bytes)

	creator := "test"
	switch pub := cert.PublicKey.(type) {
	case *rsa.PublicKey:
		rsaPublicKey := cert.PublicKey.(*rsa.PublicKey)
		modulus := rsaPublicKey.N
		hashedModulus := sha256.Sum256(modulus.Bytes())
		creator = hex.EncodeToString(hashedModulus[:])
	case *ecdsa.PublicKey:
		point := fmt.Sprintf("%v:%v", pub.X, pub.Y)
		hashedPoint := sha256.Sum256([]byte(point))
		creator = hex.EncodeToString(hashedPoint[:])
	default:
		return "", fmt.Errorf("can't determine public key algorithm")
	}

	return creator, nil
}
