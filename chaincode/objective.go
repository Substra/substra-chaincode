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
	"sort"
)

// Set is a method of the receiver Objective. It checks the validity of inputObjective and uses its fields to set the Objective.
// Returns the objectiveKey and the dataManagerKey associated to test dataSample
func (objective *Objective) Set(db *LedgerDB, inp inputObjective) (objectiveKey string, dataManagerKey string, err error) {
	dataManagerKey = inp.TestDataset.DataManagerKey
	if dataManagerKey != "" {
		var testOnly bool
		dataSampleKeys := inp.TestDataset.DataSampleKeys
		testOnly, _, err = checkSameDataManager(db, dataManagerKey, dataSampleKeys)
		if err != nil {
			err = errors.BadRequest(err, "invalid test dataSample")
			return
		} else if !testOnly {
			err = errors.BadRequest("test dataSample are not tagged as testOnly dataSample")
			return
		}
		objective.TestDataset = &Dataset{
			DataManagerKey: dataManagerKey,
			DataSampleKeys: dataSampleKeys,
		}
	} else {
		objective.TestDataset = nil
	}
	objective.AssetType = ObjectiveType
	objective.Name = inp.Name
	objective.DescriptionStorageAddress = inp.DescriptionStorageAddress
	objective.Metrics = &HashDressName{
		Name:           inp.MetricsName,
		Hash:           inp.MetricsHash,
		StorageAddress: inp.MetricsStorageAddress,
	}
	owner, err := GetTxCreator(db.cc)
	if err != nil {
		return
	}
	permissions, err := NewPermissions(db, inp.Permissions)
	if err != nil {
		return
	}
	objective.Owner = owner
	objective.Permissions = permissions
	objectiveKey = inp.DescriptionHash
	return
}

// -------------------------------------------------------------------------------------------
// Smart contract related to objectivess
// -------------------------------------------------------------------------------------------

// registerObjective stores a new objective in the ledger.
// If the key exists, it will override the value with the new one
func registerObjective(db *LedgerDB, args []string) (resp outputKey, err error) {
	// convert input strings args to input struct inputObjective
	inp := inputObjective{}
	err = AssetFromJSON(args, &inp)
	if err != nil {
		return
	}

	// check validity of input args and convert it to Objective
	objective := Objective{}
	objectiveKey, dataManagerKey, err := objective.Set(db, inp)
	if err != nil {
		return
	}
	// submit to ledger
	if err = db.Add(objectiveKey, objective); err != nil {
		return
	}
	// create composite key
	if err = db.CreateIndex("objective~owner~key", []string{"objective", objective.Owner, objectiveKey}); err != nil {
		return
	}
	// add objective to dataManager
	err = addObjectiveDataManager(db, dataManagerKey, objectiveKey)
	return outputKey{Key: objectiveKey}, err
}

// queryObjective returns a objective of the ledger given its key
func queryObjective(db *LedgerDB, args []string) (out outputObjective, err error) {
	inp := inputKey{}
	err = AssetFromJSON(args, &inp)
	if err != nil {
		return
	}
	objective, err := db.GetObjective(inp.Key)
	if err != nil {
		return
	}
	out.Fill(inp.Key, objective)
	return
}

// queryObjectives returns all objectives of the ledger
func queryObjectives(db *LedgerDB, args []string) (outObjectives []outputObjective, err error) {
	outObjectives = []outputObjective{}
	if len(args) != 0 {
		err = errors.BadRequest("incorrect number of arguments, expecting nothing")
		return
	}
	elementsKeys, err := db.GetIndexKeys("objective~owner~key", []string{"objective"})
	if err != nil {
		return
	}
	for _, key := range elementsKeys {
		objective, err := db.GetObjective(key)
		if err != nil {
			return outObjectives, err
		}
		var out outputObjective
		out.Fill(key, objective)
		outObjectives = append(outObjectives, out)
	}
	return
}

// getObjectiveLeaderboard returns for an objective, all its certified testtuples with a done status, ordered by their perf
// It can be an ascending sort or not depending on the ascendingOrder value.
func queryObjectiveLeaderboard(db *LedgerDB, args []string) (outputLeaderboard, error) {
	inp := inputLeaderboard{}
	err := AssetFromJSON(args, &inp)
	if err != nil {
		return outputLeaderboard{}, err
	}

	objective, err := db.GetObjective(inp.ObjectiveKey)
	if err != nil {
		return outputLeaderboard{}, err
	}
	outObjective := outputObjective{}
	outObjective.Fill(inp.ObjectiveKey, objective)
	out := outputLeaderboard{Objective: outObjective, Testtuples: []outputBoardTuple{}}

	testtupleKeys, err := db.GetIndexKeys("testtuple~objective~certified~key", []string{"testtuple", inp.ObjectiveKey, "true"})
	if err != nil {
		return outputLeaderboard{}, err
	}

	for _, testtupleKey := range testtupleKeys {
		var boardTuple outputBoardTuple
		testtuple, err := db.GetTesttuple(testtupleKey)
		if err != nil {
			return outputLeaderboard{}, err
		}
		if testtuple.Status != StatusDone {
			continue
		}
		err = boardTuple.Fill(db, testtuple, testtupleKey)
		if err != nil {
			return outputLeaderboard{}, err
		}
		out.Testtuples = append(out.Testtuples, boardTuple)
	}

	if inp.AscendingOrder {
		sort.Sort(out.Testtuples)
	} else {
		sort.Sort(sort.Reverse(out.Testtuples))
	}
	return out, nil
}

// -------------------------------------------------------------------------------------------
// Utils for objectivess
// -------------------------------------------------------------------------------------------

// addObjectiveDataManager associates a objective to a dataManager, more precisely, it adds the objective key to the dataManager
func addObjectiveDataManager(db *LedgerDB, dataManagerKey string, objectiveKey string) error {
	dataManager, err := db.GetDataManager(dataManagerKey)
	if err != nil {
		return nil
	}
	if dataManager.ObjectiveKey != "" {
		return errors.BadRequest("dataManager is already associated with a objective")
	}
	dataManager.ObjectiveKey = objectiveKey
	return db.Put(dataManagerKey, dataManager)
}
