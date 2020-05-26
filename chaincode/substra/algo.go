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
)

// Set is a method of the receiver Algo. It uses inputAlgo fields to set the Algo
// Returns the algoKey
func (algo *Algo) Set(db *LedgerDB, inp inputAlgo) (algoKey string, err error) {
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

	algo.AssetType = AlgoType
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
// registerAlgo stores a new algo in the ledger.
// If the key exists, it will override the value with the new one
func registerAlgo(db *LedgerDB, args []string) (resp outputKey, err error) {
	inp := inputAlgo{}
	err = AssetFromJSON(args, &inp)
	if err != nil {
		return
	}
	// check validity of input args and convert it to Algo
	algo := Algo{}
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
	err = db.CreateIndex("algo~owner~key", []string{"algo", algo.Owner, algoKey})
	if err != nil {
		return
	}
	return outputKey{Key: algoKey}, nil
}

// queryAlgo returns an algo of the ledger given its key
func queryAlgo(db *LedgerDB, args []string) (out outputAlgo, err error) {
	inp := inputKey{}
	err = AssetFromJSON(args, &inp)
	if err != nil {
		return
	}
	algo, err := db.GetAlgo(inp.Key)
	if err != nil {
		return
	}
	out.Fill(inp.Key, algo)
	return
}

// queryAlgos returns all algos of the ledger
func queryAlgos(db *LedgerDB, args []string) (outAlgos []outputAlgo, err error) {
	outAlgos = []outputAlgo{}
	if len(args) != 0 {
		err = errors.BadRequest("incorrect number of arguments, expecting nothing")
		return
	}
	elementsKeys, err := db.GetIndexKeys("algo~owner~key", []string{"algo"})
	if err != nil {
		return
	}
	nb := getLimitedNbSliceElements(elementsKeys)
	for _, key := range elementsKeys[:nb] {
		algo, err := db.GetAlgo(key)
		if err != nil {
			return outAlgos, err
		}
		var out outputAlgo
		out.Fill(key, algo)
		outAlgos = append(outAlgos, out)
	}
	return
}
