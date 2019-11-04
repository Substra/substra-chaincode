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

// Set is a method of the receiver CompositeAlgo. It uses inputCompositeAlgo fields to set the CompositeAlgo
// Returns the compositeAlgoKey
func (algo *CompositeAlgo) Set(db LedgerDB, inp inputCompositeAlgo) (algoKey string, err error) {
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

	algo.AssetType = CompositeAlgoType
	algo.Name = inp.Name
	algo.StorageAddress = inp.StorageAddress
	algo.Description = &HashDress{
		Hash:           inp.DescriptionHash,
		StorageAddress: inp.DescriptionStorageAddress,
	}
	algo.Owner = owner
	algo.Permissions = permissions
	return
}

// -------------------------------------------------------------------------------------------
// Smart contracts related to an algo
// -------------------------------------------------------------------------------------------
// registerCompositeAlgo stores a new algo in the ledger.
// If the key exists, it will override the value with the new one
func registerCompositeAlgo(db LedgerDB, args []string) (resp map[string]string, err error) {
	inp := inputCompositeAlgo{}
	err = AssetFromJSON(args, &inp)
	if err != nil {
		return
	}
	// check validity of input args and convert it to CompositeAlgo
	algo := CompositeAlgo{}
	algoKey, err := algo.Set(db, inp)
	if err != nil {
		return
	}
	// submit to ledger
	err = db.Add(algoKey, algo)
	if err != nil {
		return
	}
	// create composite key
	err = db.CreateIndex("compositealgo~owner~key", []string{"compositealgo", algo.Owner, algoKey})
	if err != nil {
		return
	}
	return map[string]string{"key": algoKey}, nil
}

// queryCompositeAlgo returns an algo of the ledger given its key
func queryCompositeAlgo(db LedgerDB, args []string) (out outputCompositeAlgo, err error) {
	inp := inputHash{}
	err = AssetFromJSON(args, &inp)
	if err != nil {
		return
	}
	algo, err := db.GetCompositeAlgo(inp.Key)
	if err != nil {
		return
	}
	out.Fill(inp.Key, algo)
	return
}

// queryAlgos returns all algos of the ledger
func queryCompositeAlgos(db LedgerDB, args []string) (outAlgos []outputCompositeAlgo, err error) {
	outAlgos = []outputCompositeAlgo{}
	if len(args) != 0 {
		err = errors.BadRequest("incorrect number of arguments, expecting nothing")
		return
	}
	elementsKeys, err := db.GetIndexKeys("compositealgo~owner~key", []string{"compositealgo"})
	if err != nil {
		return
	}
	for _, key := range elementsKeys {
		algo, err := db.GetCompositeAlgo(key)
		if err != nil {
			return outAlgos, err
		}
		var out outputCompositeAlgo
		out.Fill(key, algo)
		outAlgos = append(outAlgos, out)
	}
	return
}
