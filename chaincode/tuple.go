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

// queryModelDetails returns info about the testtuple and algo related to a traintuple
func queryModelDetails(db LedgerDB, args []string) (outModelDetails outputModelDetails, err error) {
	inp := inputHash{}
	err = AssetFromJSON(args, &inp)
	if err != nil {
		return
	}

	// get traintuple type
	traintupleType, err := db.GetAssetType(inp.Key)
	if err != nil {
		return
	}
	switch traintupleType {
	case TraintupleType:
		var out outputTraintuple
		out, err = getOutputTraintuple(db, inp.Key)
		if err != nil {
			return
		}
		outModelDetails.Traintuple = &out
	case CompositeTraintupleType:
		var out outputCompositeTraintuple
		out, err = getOutputCompositeTraintuple(db, inp.Key)
		if err != nil {
			return
		}
		outModelDetails.CompositeTraintuple = &out
	case AggregateTupleType:
		var out outputAggregateTuple
		out, err = getOutputAggregateTuple(db, inp.Key)
		if err != nil {
			return
		}
		outModelDetails.AggregateTuple = &out
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

	// populate from regular traintuples
	traintupleKeys, err := db.GetIndexKeys("traintuple~algo~key", []string{"traintuple"})
	if err != nil {
		return
	}
	for _, traintupleKey := range traintupleKeys {
		var outputModel outputModel
		var out outputTraintuple

		out, err = getOutputTraintuple(db, traintupleKey)
		if err != nil {
			return
		}
		outputModel.Traintuple = &out
		outputModel.Testtuple, err = getCertifiedOutputTesttuple(db, traintupleKey)
		outModels = append(outModels, outputModel)
	}

	// populate from composite traintuples
	compositeTraintupleKeys, err := db.GetIndexKeys("compositeTraintuple~algo~key", []string{"compositeTraintuple"})
	if err != nil {
		return
	}
	for _, compositeTraintupleKey := range compositeTraintupleKeys {
		var outputModel outputModel
		var out outputCompositeTraintuple

		out, err = getOutputCompositeTraintuple(db, compositeTraintupleKey)
		if err != nil {
			return
		}
		outputModel.CompositeTraintuple = &out
		outputModel.Testtuple, err = getCertifiedOutputTesttuple(db, compositeTraintupleKey)
		outModels = append(outModels, outputModel)
	}

	return
}

// ----------------------------------------------------------
// Utils for smartcontracts related to  multiple tuple types
// ----------------------------------------------------------

/* Unused
// checkLog checks the validity of logs
func checkLog(log string) (err error) {
	maxLength := 200
	if length := len(log); length > maxLength {
		err = fmt.Errorf("too long log, is %d and should be %d ", length, maxLength)
	}
	return
}
*/

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

func getCertifiedOutputTesttuple(db LedgerDB, traintupleKey string) (outputTesttuple, error) {
	var out outputTesttuple
	// get associated testtuple
	var testtupleKeys []string
	testtupleKeys, err := db.GetIndexKeys("testtuple~traintuple~certified~key", []string{"testtuple", traintupleKey, "true"})
	if err != nil {
		return out, err
	}
	if len(testtupleKeys) == 0 {
		return out, nil
	}
	// get testtuple and serialize it
	testtupleKey := testtupleKeys[0]
	out, err = getOutputTesttuple(db, testtupleKey)
	if err != nil {
		return out, err
	}

	return out, nil
}
