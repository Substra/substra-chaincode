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
func (algo *AggregateAlgo) Set(db *LedgerDB, inp inputAggregateAlgo) (err error) {
	// find associated owner
	owner, err := GetTxCreator(db.cc)
	if err != nil {
		return
	}

	permissions, err := NewPermissions(db, inp.Permissions)
	if err != nil {
		return
	}

	algo.Key = inp.Key
	algo.AssetType = AggregateAlgoType
	algo.Name = inp.Name
	algo.Checksum = inp.Checksum
	algo.StorageAddress = inp.StorageAddress
	algo.Description = &ChecksumAddress{
		Checksum:       inp.DescriptionChecksum,
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
	err = algo.Set(db, inp)
	if err != nil {
		return
	}
	// submit to ledger
	err = db.Add(inp.Key, algo)
	if err != nil {
		return
	}
	// create aggregate key
	err = db.CreateIndex("aggregateAlgo~owner~key", []string{"aggregateAlgo", algo.Owner, inp.Key})
	if err != nil {
		return
	}
	return outputKey{Key: inp.Key}, nil
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
	out.Fill(algo)
	return
}

// queryAggregateAlgos returns all algos of the ledger
func queryAggregateAlgos(db *LedgerDB, args []string) (outAlgos []outputAggregateAlgo, bookmarks map[string]string, err error) {
	outAlgos = []outputAggregateAlgo{}
	index := "aggregateAlgo~owner~key"
	bookmarks = map[string]string{index: ""}

	if len(args) > 1 {
		err = errors.BadRequest("incorrect number of arguments, expecting at most one argument")
		return
	}

	if len(args) == 1 && args[0] != "" {
		inp := inputBookmarks{}
		err := AssetFromJSON(args, &inp)
		if err != nil {
			return nil, bookmarks, err
		}
		bookmarks = inp.Bookmarks
	}


	elementsKeys, bookmark, err := db.GetIndexKeysWithPagination(index, []string{"aggregateAlgo"}, OutputAssetPaginationHardLimit, bookmarks[index])
	bookmarks[index] = bookmark

	if err != nil {
		return
	}

	for _, key := range elementsKeys {
		algo, err := db.GetAggregateAlgo(key)
		if err != nil {
			return outAlgos, bookmarks, err
		}
		var out outputAggregateAlgo
		out.Fill(algo)
		outAlgos = append(outAlgos, out)
	}
	return
}
