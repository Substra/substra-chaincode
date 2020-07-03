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

package substra

import (
	"chaincode/errors"
	"encoding/json"
	"fmt"
	"math/rand"
	"unicode"
	"unicode/utf8"

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
func AssetFromJSON(body string, asset interface{}) error {
	err := json.Unmarshal([]byte(body), &asset)
	if err != nil {
		return errors.BadRequest(err, "problem when reading json arg: %s, error is:", body)
	}
	v := validator.New()
	err = v.Struct(asset)
	if err != nil {
		return errors.BadRequest(err, "inputs validation failed: %s, error is:", body)
	}
	return nil
}

// LowerFirst returns the input string with the first letter lowercased
func LowerFirst(s string) string {
	if s == "" {
		return ""
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToLower(r)) + s[n:]
}

// String returns a string representation for an asset type
func (assetType AssetType) String() string {
	switch assetType {
	case ObjectiveType:
		return "Objective"
	case DataManagerType:
		return "DataManager"
	case DataSampleType:
		return "DataSample"
	case AlgoType:
		return "Algo"
	case CompositeAlgoType:
		return "CompositeAlgo"
	case AggregateAlgoType:
		return "AggregateAlgo"
	case TraintupleType:
		return "Traintuple"
	case CompositeTraintupleType:
		return "CompositeTraintuple"
	case AggregatetupleType:
		return "Aggregatetuple"
	case TesttupleType:
		return "Testtuple"
	case ComputePlanType:
		return "ComputePlan"
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

func initMapOutput(m map[string]string) map[string]string {
	if m == nil {
		return map[string]string{}
	}
	return m
}
