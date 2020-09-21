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
	"encoding/json"
	"fmt"
	"math/rand"

	"gopkg.in/go-playground/validator.v9"

	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/msp"
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

// typeInSlice check if an AssetType is in a slice
func typeInSlice(a AssetType, list []AssetType) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
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

// GetTxCreator returns the transaction creator
func GetTxCreator(stub shim.ChaincodeStubInterface) (string, error) {
	creator, err := stub.GetCreator()

	if err != nil {
		return "", err
	}

	sID := &msp.SerializedIdentity{}
	err = proto.Unmarshal(creator, sID)
	if err != nil {
		return "", err
	}

	return sID.GetMspid(), nil
}

// String returns a string representation for an asset type
func (assetType AssetType) String() string {
	switch assetType {
	case ObjectiveType:
		return "objective"
	case DataManagerType:
		return "data_manager"
	case DataSampleType:
		return "data_sample"
	case AlgoType:
		return "algo"
	case CompositeAlgoType:
		return "composite_algo"
	case AggregateAlgoType:
		return "aggregate_algo"
	case TraintupleType:
		return "traintuple"
	case CompositeTraintupleType:
		return "composite_traintuple"
	case AggregatetupleType:
		return "aggregatetuple"
	case TesttupleType:
		return "testtuple"
	case ComputePlanType:
		return "compute_plan"
	default:
		return fmt.Sprintf("(unknown asset type: %d)", assetType)
	}
}

var characterRunes = []rune("abcdef0123456789")

// GetRandomHash generate a random string of 64 character
func GetRandomHash() string {
	b := make([]rune, 64)
	for i := range b {
		b[i] = characterRunes[rand.Intn(len(characterRunes))]
	}
	return string(b)
}

// GetNewUUID generates a new UUID
func GetNewUUID() (string, error) {
	a, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return a.String(), nil
}

func initMapOutput(m map[string]string) map[string]string {
	if m == nil {
		return map[string]string{}
	}
	return m
}
