package main

import (
	"chaincode/errors"
)

// Set is a method of the receiver Algo. It uses inputAlgo fields to set the Algo
// Returns the algoKey
func (algo *Algo) Set(db LedgerDB, inp inputAlgo) (algoKey string, err error) {
	algoKey = inp.Hash
	// find associated owner
	owner, err := GetTxCreator(db.cc)
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
	algo.Permissions = NewPermissions(inp.Permissions, owner)
	return
}

// -------------------------------------------------------------------------------------------
// Smart contracts related to an algo
// -------------------------------------------------------------------------------------------
// registerAlgo stores a new algo in the ledger.
// If the key exists, it will override the value with the new one
func registerAlgo(db LedgerDB, args []string) (resp map[string]string, err error) {
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
	return map[string]string{"key": algoKey}, nil
}

// queryAlgo returns an algo of the ledger given its key
func queryAlgo(db LedgerDB, args []string) (out outputAlgo, err error) {
	inp := inputHash{}
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
func queryAlgos(db LedgerDB, args []string) (outAlgos []outputAlgo, err error) {
	outAlgos = []outputAlgo{}
	if len(args) != 0 {
		err = errors.BadRequest("incorrect number of arguments, expecting nothing")
		return
	}
	elementsKeys, err := db.GetIndexKeys("algo~owner~key", []string{"algo"})
	if err != nil {
		return
	}
	for _, key := range elementsKeys {
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
