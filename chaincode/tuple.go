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
	"fmt"
	"sort"
	"strconv"
)

// List of the possible tuple's status
const (
	StatusDoing   = "doing"
	StatusTodo    = "todo"
	StatusWaiting = "waiting"
	StatusFailed  = "failed"
	StatusDone    = "done"
)

// ------------------------------------------------
// Smart contracts related to multiple tuple types
// ------------------------------------------------

// createComputePlan is the wrapper for the substra smartcontract CreateComputePlan
func createComputePlan(db LedgerDB, args []string) (resp outputComputePlan, err error) {
	inp := inputComputePlan{}
	err = AssetFromJSON(args, &inp)
	if err != nil {
		return
	}
	traintupleKeysByID := map[string]string{}
	resp.TraintupleKeys = []string{}
	var traintuplesTodo []outputTraintuple
	for i, computeTraintuple := range inp.Traintuples {
		inpTraintuple := inputTraintuple{}
		inpTraintuple.AlgoKey = inp.AlgoKey
		inpTraintuple.ObjectiveKey = inp.ObjectiveKey
		inpTraintuple.DataManagerKey = computeTraintuple.DataManagerKey
		inpTraintuple.DataSampleKeys = computeTraintuple.DataSampleKeys
		inpTraintuple.Tag = computeTraintuple.Tag
		inpTraintuple.Rank = strconv.Itoa(i)

		traintuple := Traintuple{}
		err := traintuple.SetFromInput(db, inpTraintuple)
		if err != nil {
			return resp, err
		}

		// Set the inModels by matching the id to traintuples key previously
		// encontered in this compute plan
		for _, InModelID := range computeTraintuple.InModelsIDs {
			inModelKey, ok := traintupleKeysByID[InModelID]
			if !ok {
				return resp, errors.BadRequest("traintuple ID %s: model ID %s not found, check traintuple list order", computeTraintuple.ID, InModelID)
			}
			traintuple.InModelKeys = append(traintuple.InModelKeys, inModelKey)
		}

		traintupleKey := traintuple.GetKey()

		// Set the ComputePlanID
		if i == 0 {
			traintuple.ComputePlanID = traintupleKey
			resp.ComputePlanID = traintuple.ComputePlanID
		} else {
			traintuple.ComputePlanID = resp.ComputePlanID
		}

		// Set status: if it has parents it's waiting
		// if not it's todo and it has to be included in the event
		if len(computeTraintuple.InModelsIDs) > 0 {
			traintuple.Status = StatusWaiting
		} else {
			traintuple.Status = StatusTodo
			out := outputTraintuple{}
			err = out.Fill(db, traintuple, traintupleKey)
			if err != nil {
				return resp, err
			}
			traintuplesTodo = append(traintuplesTodo, out)
		}

		err = traintuple.Save(db, traintupleKey)
		if err != nil {
			return resp, err
		}
		traintupleKeysByID[computeTraintuple.ID] = traintupleKey
		resp.TraintupleKeys = append(resp.TraintupleKeys, traintupleKey)
	}

	resp.TesttupleKeys = []string{}
	for index, computeTesttuple := range inp.Testtuples {
		traintupleKey, ok := traintupleKeysByID[computeTesttuple.TraintupleID]
		if !ok {
			return resp, errors.BadRequest("testtuple index %s: traintuple ID %s not found", index, computeTesttuple.TraintupleID)
		}
		testtuple := Testtuple{}
		testtuple.Model = &Model{TraintupleKey: traintupleKey}
		testtuple.ObjectiveKey = inp.ObjectiveKey
		testtuple.AlgoKey = inp.AlgoKey

		inputTesttuple := inputTesttuple{}
		inputTesttuple.DataManagerKey = computeTesttuple.DataManagerKey
		inputTesttuple.DataSampleKeys = computeTesttuple.DataSampleKeys
		inputTesttuple.Tag = computeTesttuple.Tag
		err = testtuple.SetFromInput(db, inputTesttuple)
		if err != nil {
			return resp, err
		}
		testtuple.Status = StatusWaiting
		testtupleKey := testtuple.GetKey()
		err = testtuple.Save(db, testtupleKey)
		if err != nil {
			return resp, err
		}
		resp.TesttupleKeys = append(resp.TesttupleKeys, testtupleKey)
	}

	event := TuplesEvent{}
	event.SetTraintuples(traintuplesTodo...)
	err = SendTuplesEvent(db.cc, event)
	if err != nil {
		return resp, err
	}

	return resp, err
}

// queryModelDetails returns info about the testtuple and algo related to a traintuple
func queryModelDetails(db LedgerDB, args []string) (outModelDetails outputModelDetails, err error) {
	inp := inputHash{}
	err = AssetFromJSON(args, &inp)
	if err != nil {
		return
	}

	// get associated traintuple
	outModelDetails.Traintuple, err = getOutputTraintuple(db, inp.Key)
	if err != nil {
		return
	}

	// get certified and non-certified testtuples related to traintuple
	testtupleKeys, err := db.GetIndexKeys("testtuple~traintuple~certified~key", []string{"testtuple", inp.Key})
	if err != nil {
		return
	}
	for _, testtupleKey := range testtupleKeys {
		// get testtuple and serialize it
		var outputTesttuple outputTesttuple
		outputTesttuple, err = getOutputTesttuple(db, testtupleKey)
		if err != nil {
			return
		}

		if outputTesttuple.Certified {
			outModelDetails.Testtuple = outputTesttuple
		} else {
			outModelDetails.NonCertifiedTesttuples = append(outModelDetails.NonCertifiedTesttuples, outputTesttuple)
		}
	}
	return
}

// queryModels returns all traintuples and associated testuples
func queryModels(db LedgerDB, args []string) (outModels []outputModel, err error) {
	outModels = []outputModel{}

	if len(args) != 0 {
		err = errors.BadRequest("incorrect number of arguments, expecting nothing")
		return
	}

	traintupleKeys, err := db.GetIndexKeys("traintuple~algo~key", []string{"traintuple"})
	if err != nil {
		return
	}
	for _, traintupleKey := range traintupleKeys {
		var outputModel outputModel

		// get traintuple
		outputModel.Traintuple, err = getOutputTraintuple(db, traintupleKey)
		if err != nil {
			return
		}

		// get associated testtuple
		var testtupleKeys []string
		testtupleKeys, err = db.GetIndexKeys("testtuple~traintuple~certified~key", []string{"testtuple", traintupleKey, "true"})
		if err != nil {
			return
		}
		if len(testtupleKeys) == 1 {
			// get testtuple and serialize it
			testtupleKey := testtupleKeys[0]
			outputModel.Testtuple, err = getOutputTesttuple(db, testtupleKey)
			if err != nil {
				return
			}
		}
		outModels = append(outModels, outputModel)
	}
	return
}

// ----------------------------------------------------------
// Utils for smartcontracts related to  multiple tuple types
// ----------------------------------------------------------

// checkLog checks the validity of logs
func checkLog(log string) (err error) {
	maxLength := 200
	if length := len(log); length > maxLength {
		err = fmt.Errorf("too long log, is %d and should be %d ", length, maxLength)
	}
	return
}

func validateTupleOwner(db LedgerDB, worker string) error {
	txCreator, err := GetTxCreator(db.cc)
	if err != nil {
		return err
	}
	if txCreator != worker {
		return fmt.Errorf("%s is not allowed to update tuple (%s)", txCreator, worker)
	}
	return nil
}

// check validity of traintuple update: consistent status and agent submitting the transaction
func checkUpdateTuple(db LedgerDB, worker string, oldStatus string, newStatus string) error {
	statusPossibilities := map[string]string{
		StatusWaiting: StatusTodo,
		StatusTodo:    StatusDoing,
		StatusDoing:   StatusDone}
	if statusPossibilities[oldStatus] != newStatus && newStatus != StatusFailed {
		return errors.BadRequest("cannot change status from %s to %s", oldStatus, newStatus)
	}
	return nil
}

// HashForKey to generate key for an asset
func HashForKey(objectType string, hashElements ...string) string {
	toHash := objectType
	sort.Strings(hashElements)
	for _, element := range hashElements {
		toHash += "," + element
	}
	sum := sha256.Sum256([]byte(toHash))
	return hex.EncodeToString(sum[:])
}
