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
)

// Set is a method of the receiver Aggregate. It uses inputAggregateAlgo fields to set the AggregateAlgo
// Returns the aggregateAlgoKey
func (algo *AggregateAlgo) Set(db *LedgerDB, inp inputAggregateAlgo) (algoKey string, err error) {
	algoKey = inp.Hash
	// find associated owner
	owner, err := GetTxCreator(db.cc)
	if err != nil {
		return
	}

	permissions, err := NewPermissions(db, inp.Permissions)
	if err != nil {
		return
	}

	algo.AssetType = AggregateAlgoType
	algo.Name = inp.Name
	algo.StorageAddress = inp.StorageAddress
	algo.Description = &HashDress{
		Hash:           inp.DescriptionHash,
		StorageAddress: inp.DescriptionStorageAddress,
	}
	algo.Owner = owner
	algo.Permissions = permissions
	algo.Metadata = inp.Metadata
	return
}

// -------------------------------------------------------------------------------------------
// Smart contracts related to an algo
// -------------------------------------------------------------------------------------------
// registerAggregateAlgo stores a new algo in the ledger.
// If the key exists, it will override the value with the new one
func registerAggregateAlgo(db *LedgerDB, args []string) (resp outputKey, err error) {
	inp := inputAggregateAlgo{}
	err = AssetFromJSON(args, &inp)
	if err != nil {
		return
	}
	// check validity of input args and convert it to CompositeAlgo
	algo := AggregateAlgo{}
	algoKey, err := algo.Set(db, inp)
	if err != nil {
		return
	}
	// submit to ledger
	err = db.Add(algoKey, algo)
	if err != nil {
		return
	}
	// create aggregate key
	err = db.CreateIndex("aggregateAlgo~owner~key", []string{"aggregateAlgo", algo.Owner, algoKey})
	if err != nil {
		return
	}
	return outputKey{Key: algoKey}, nil
}

// queryAggregateAlgo returns an algo of the ledger given its key
func queryAggregateAlgo(db *LedgerDB, args []string) (out outputAggregateAlgo, err error) {
	inp := inputKey{}
	err = AssetFromJSON(args, &inp)
	if err != nil {
		return
	}
	algo, err := db.GetAggregateAlgo(inp.Key)
	if err != nil {
		return
	}
	out.Fill(inp.Key, algo)
	return
}

// queryAggregateAlgos returns all algos of the ledger
func queryAggregateAlgos(db *LedgerDB, args []string) (outAlgos []outputAggregateAlgo, err error) {
	outAlgos = []outputAggregateAlgo{}
	if len(args) != 0 {
		err = errors.BadRequest("incorrect number of arguments, expecting nothing")
		return
	}
	elementsKeys, err := db.GetIndexKeys("aggregateAlgo~owner~key", []string{"aggregateAlgo"})
	if err != nil {
		return
	}
	nb := getLimitedNbSliceElements(elementsKeys)
	for _, key := range elementsKeys[:nb] {
		algo, err := db.GetAggregateAlgo(key)
		if err != nil {
			return outAlgos, err
		}
		var out outputAggregateAlgo
		out.Fill(key, algo)
		outAlgos = append(outAlgos, out)
	}
	return
}
