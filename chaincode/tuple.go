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
)

// List of the possible tuple's status
const (
	StatusDoing    = "doing"
	StatusTodo     = "todo"
	StatusWaiting  = "waiting"
	StatusFailed   = "failed"
	StatusDone     = "done"
	StatusCanceled = "canceled"
	// The status aborted is still under discussion so the logic is already
	// implemented but it's value is the same as canceled for now.
	StatusAborted = "canceled"
)

// ------------------------------------------------
// Smart contracts related to multiple tuple types
// ------------------------------------------------

// queryModelDetails returns info about the testtuple and algo related to a traintuple
func queryModelDetails(db *LedgerDB, args []string) (outModelDetails outputModelDetails, err error) {
	inp := inputKey{}
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
	case AggregatetupleType:
		var out outputAggregatetuple
		out, err = getOutputAggregatetuple(db, inp.Key)
		if err != nil {
			return
		}
		outModelDetails.Aggregatetuple = &out
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

type queryModelsBookmarks struct {
	Traintuple          string `json:"traintuple"`
	CompositeTraintuple string `json:"composite_traintuple"`
	Aggregatetuple      string `json:"aggregatetuple"`
}

type inputQueryModelsBookmarks struct {
	Bookmark string `json:"bookmark"`
}

// queryModels returns all traintuples and associated testuples
func queryModels(db *LedgerDB, args []string) (outModels []outputModel, bookmark string, err error) {
	outModels = []outputModel{}
	bookmarks := queryModelsBookmarks{}

	if len(args) > 1 {
		err = errors.BadRequest("incorrect number of arguments, expecting at most one argument")
		return
	}

	if len(args) == 1 && args[0] != "" {
		inp := inputQueryModelsBookmarks{}
		err = json.Unmarshal([]byte(args[0]), &inp)
		if err != nil {
			return
		}
		if inp.Bookmark != "" {
			err = json.Unmarshal([]byte(inp.Bookmark), &bookmarks)
			if err != nil {
				return
			}
		}
	}

	// populate from regular traintuples
	traintupleKeys, traintupleBookmark, err := db.GetIndexKeysWithPagination("traintuple~algo~key", []string{"traintuple"}, OutputPageSize/3, bookmarks.Traintuple)
	bookmarks.Traintuple = traintupleBookmark

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
		outModels = append(outModels, outputModel)
	}

	// populate from composite traintuples
	compositeTraintupleKeys, compositeTraintupleBookmark, err := db.GetIndexKeysWithPagination("compositeTraintuple~algo~key", []string{"compositeTraintuple"}, OutputPageSize/3, bookmarks.CompositeTraintuple)
	bookmarks.CompositeTraintuple = compositeTraintupleBookmark

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
		outModels = append(outModels, outputModel)
	}

	// populate from composite traintuples
	aggregatetupleKeys, aggregatetupleBookmark, err := db.GetIndexKeysWithPagination("aggregatetuple~algo~key", []string{"aggregatetuple"}, OutputPageSize/3, bookmarks.Aggregatetuple)
	bookmarks.Aggregatetuple = aggregatetupleBookmark

	if err != nil {
		return
	}
	for _, aggregatetupleKey := range aggregatetupleKeys {
		var outputModel outputModel
		var out outputAggregatetuple

		out, err = getOutputAggregatetuple(db, aggregatetupleKey)
		if err != nil {
			return
		}
		outputModel.Aggregatetuple = &out
		outModels = append(outModels, outputModel)
	}

	bookmarkBytes, err := json.Marshal(bookmarks)
	bookmark = string(bookmarkBytes)
	return
}

func queryModelPermissions(db *LedgerDB, args []string) (outputPermissions, error) {
	var out outputPermissions
	inp := inputKey{}
	err := AssetFromJSON(args, &inp)
	if err != nil {
		return out, err
	}
	modelKey := inp.Key
	keys, err := db.GetIndexKeys("tuple~modelKey~key", []string{"tuple", modelKey})
	if err != nil {
		return out, err
	}
	if len(keys) == 0 {
		return out, errors.NotFound("Could not find a model for key %s", modelKey)
	}
	tupleKey := keys[0]
	tupleType, err := db.GetAssetType(tupleKey)
	if err != nil {
		return out, errors.Internal(err, "queryModelPermissions: could not retrieve model type with tupleKey %s", tupleKey)
	}

	modelDetails, err := getModelDetails(db, modelKey, tupleKey, tupleType)
	if err != nil {
		return out, err
	}
	out.Fill(modelDetails.Permissions)

	return out, nil
}

// ----------------------------------------------------------
// Utils for smartcontracts related to  multiple tuple types
// ----------------------------------------------------------

func validateTupleOwner(db *LedgerDB, worker string) error {
	txCreator, err := GetTxCreator(db.cc)
	if err != nil {
		return err
	}
	if txCreator != worker {
		return errors.Forbidden("%s is not allowed to update tuple (%s)", txCreator, worker)
	}
	return nil
}

// check validity of traintuple update: consistent status and agent submitting the transaction
func checkUpdateTuple(db *LedgerDB, worker string, oldStatus string, newStatus string) error {
	if StatusAborted == newStatus {
		return nil
	}

	statusPossibilities := map[string]string{
		StatusWaiting: StatusTodo,
		StatusTodo:    StatusDoing,
		StatusDoing:   StatusDone}
	if statusPossibilities[oldStatus] != newStatus && newStatus != StatusFailed {
		return errors.BadRequest("cannot change status from %s to %s", oldStatus, newStatus)
	}
	return nil
}

func determineStatusFromInModels(statuses []string) string {
	if stringInSlice(StatusFailed, statuses) {
		return StatusFailed
	}

	if stringInSlice(StatusAborted, statuses) {
		return StatusAborted
	}

	for _, s := range statuses {
		if s != StatusDone {
			return StatusWaiting
		}
	}
	return StatusTodo
}

func determineTupleStatus(db *LedgerDB, tupleStatus, computePlanKey string) (string, error) {
	if tupleStatus != StatusWaiting || computePlanKey == "" {
		return tupleStatus, nil
	}
	computePlan, err := db.GetComputePlan(computePlanKey)
	if err != nil {
		return "", err
	}
	if stringInSlice(computePlan.State.Status, []string{StatusFailed, StatusCanceled}) {
		return StatusAborted, nil
	}
	return tupleStatus, nil
}

func createModelIndex(db *LedgerDB, modelKey, tupleKey string) error {
	return db.CreateIndex("tuple~modelKey~key", []string{"tuple", modelKey, tupleKey})
}

type modelDetails struct {
	Permissions Permissions
	Owner       string
}

func getModelDetails(db *LedgerDB, modelKey string, tupleKey string, tupleType AssetType) (modelDetails, error) {
	// By default model is public processable
	model := modelDetails{}

	if tupleType == TraintupleType {
		tuple, err := db.GetTraintuple(tupleKey)
		if err != nil {
			return model, errors.Internal(err, "getModelDetails: cannot get traintuple")
		}
		model.Permissions = tuple.Permissions
		model.Owner = tuple.Dataset.Worker
	}

	if tupleType == AggregatetupleType {
		tuple, err := db.GetAggregatetuple(tupleKey)
		if err != nil {
			return model, errors.Internal(err, "getModelDetails: cannot get aggregatetuple")
		}
		model.Permissions = tuple.Permissions
		model.Owner = tuple.Worker
	}

	if tupleType == CompositeTraintupleType {
		tuple, err := db.GetCompositeTraintuple(tupleKey)
		if err != nil {
			return model, errors.Internal(err, "getModelDetails: cannot get composite traintuple")
		}
		// if `modelKey` refers to the head out-model, return the head out-model permissions
		// if `modelKey` refers to the trunk out-model, default to "public processable")
		if tuple.OutHeadModel.OutModel.Key == modelKey {
			model.Permissions = tuple.OutHeadModel.Permissions
		} else {
			model.Permissions = tuple.OutTrunkModel.Permissions
		}
		model.Owner = tuple.Dataset.Worker
	}
	return model, nil
}
