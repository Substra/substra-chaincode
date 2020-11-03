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
	"reflect"
	"sort"
	"strconv"
)

// -------------------------------------------------------------------------------------------
// Methods on receivers testtuple
// -------------------------------------------------------------------------------------------

// SetFromInput is a method of the receiver Testtuple.
// It uses the inputTesttuple to check and set the testtuple's parameters
// which don't depend on previous testtuples values :
//  - AssetType
//  - Creator
//  - Tag
//  - Dataset
//  - Certified
func (testtuple *Testtuple) SetFromInput(db *LedgerDB, inp inputTesttuple) error {
	creator, err := GetTxCreator(db.cc)
	if err != nil {
		return err
	}
	testtuple.Key = inp.Key
	testtuple.Creator = creator
	testtuple.Tag = inp.Tag
	testtuple.Metadata = inp.Metadata
	testtuple.AssetType = TesttupleType

	// Get test dataset from objective
	objective, err := db.GetObjective(inp.ObjectiveKey)
	if err != nil {
		return errors.BadRequest(err, "could not retrieve objective with key %s", inp.ObjectiveKey)
	}
	testtuple.ObjectiveKey = inp.ObjectiveKey
	var objectiveDataManagerKey string
	var objectiveDataSampleKeys []string
	if objective.TestDataset != nil {
		objectiveDataManagerKey = objective.TestDataset.DataManagerKey
		objectiveDataSampleKeys = objective.TestDataset.DataSampleKeys
	}
	// For now we need to sort it but in fine it should be save sorted
	// TODO
	sort.Strings(objectiveDataSampleKeys)

	var dataManagerKey string
	var dataSampleKeys []string
	switch {
	case len(inp.DataManagerKey) > 0 && len(inp.DataSampleKeys) > 0:
		// non-certified testtuple
		// test dataset are specified by the user
		dataSampleKeys = inp.DataSampleKeys
		_, _, err = checkSameDataManager(db, inp.DataManagerKey, dataSampleKeys)
		if err != nil {
			return err
		}
		dataManagerKey = inp.DataManagerKey
		sort.Strings(dataSampleKeys)
		testtuple.Certified = objectiveDataManagerKey == dataManagerKey && reflect.DeepEqual(objectiveDataSampleKeys, dataSampleKeys)
	case len(inp.DataManagerKey) > 0 || len(inp.DataSampleKeys) > 0:
		return errors.BadRequest("invalid input: dataManagerKey and dataSampleKey should be provided together")
	case objective.TestDataset != nil:
		dataSampleKeys = objectiveDataSampleKeys
		dataManagerKey = objectiveDataManagerKey
		testtuple.Certified = true
	default:
		return errors.BadRequest("can not create a certified testtuple, no data associated with objective %s", testtuple.ObjectiveKey)
	}
	// retrieve dataManager owner
	dataManager, err := db.GetDataManager(dataManagerKey)
	if err != nil {
		return errors.BadRequest(err, "could not retrieve dataManager with key %s", dataManagerKey)
	}
	testtuple.Dataset = &TtDataset{
		Key:            dataManager.Key,
		Worker:         dataManager.Owner,
		DataSampleKeys: dataSampleKeys,
		OpenerHash:     dataManager.Opener.Hash,
	}
	return nil
}

// SetFromTraintuple set the parameters of the testuple depending on traintuple
// it depends on. It sets:
//  - AlgoKey
//  - ObjectiveKey
//  - Model
//  - Status
func (testtuple *Testtuple) SetFromTraintuple(db *LedgerDB, traintupleKey string) error {

	var status, tupleCreator string
	var permissions Permissions

	creator, err := GetTxCreator(db.cc)
	if err != nil {
		return err
	}
	testtuple.TraintupleKey = traintupleKey
	traintupleType, err := db.GetAssetType(traintupleKey)
	if err != nil {
		return errors.BadRequest(err, "key %s is not a valid asset", traintupleKey)
	}
	switch traintupleType {
	case TraintupleType:
		traintuple, err := db.GetTraintuple(traintupleKey)
		if err != nil {
			return errors.BadRequest(err, "could not retrieve traintuple with key %s", traintupleKey)
		}
		permissions = traintuple.Permissions
		tupleCreator = traintuple.Creator
		status = traintuple.Status
		testtuple.AlgoKey = traintuple.AlgoKey
		testtuple.ComputePlanKey = traintuple.ComputePlanKey
		testtuple.Rank = traintuple.Rank
	case CompositeTraintupleType:
		compositeTraintuple, err := db.GetCompositeTraintuple(traintupleKey)
		if err != nil {
			return errors.BadRequest(err, "could not retrieve composite traintuple with key %s", traintupleKey)
		}
		permissions = compositeTraintuple.OutHeadModel.Permissions
		tupleCreator = compositeTraintuple.Creator
		status = compositeTraintuple.Status
		testtuple.AlgoKey = compositeTraintuple.AlgoKey
		testtuple.ComputePlanKey = compositeTraintuple.ComputePlanKey
		testtuple.Rank = compositeTraintuple.Rank
	case AggregatetupleType:
		tuple, err := db.GetAggregatetuple(traintupleKey)
		if err != nil {
			return errors.BadRequest(err, "could not retrieve traintuple with key %s", traintupleKey)
		}
		permissions = tuple.Permissions
		tupleCreator = tuple.Creator
		status = tuple.Status
		testtuple.AlgoKey = tuple.AlgoKey
		testtuple.ComputePlanKey = tuple.ComputePlanKey
		testtuple.Rank = tuple.Rank
	default:
		return errors.BadRequest("key %s is not a valid traintuple", traintupleKey)
	}

	if !permissions.CanProcess(tupleCreator, creator) {
		return errors.Forbidden("not authorized to process traintuple %s", traintupleKey)
	}
	switch status {
	case StatusDone:
		testtuple.Status = StatusTodo
	case StatusFailed, StatusAborted:
		return errors.BadRequest(
			"could not register this testtuple, the traintuple %s has a status %s",
			traintupleKey, status)
	default:
		testtuple.Status = StatusWaiting
	}
	return nil
}

// AddToComputePlan add the testtuple to the compute plan of it's model
func (testtuple *Testtuple) AddToComputePlan(db *LedgerDB, testtupleKey string) error {
	if testtuple.ComputePlanKey == "" {
		return nil
	}
	computePlan, err := db.GetComputePlan(testtuple.ComputePlanKey)
	if err != nil {
		return err
	}
	computePlan.AddTuple(TesttupleType, testtupleKey, testtuple.Status)
	err = computePlan.Save(db, testtuple.ComputePlanKey)
	if err != nil {
		return err
	}
	return nil
}

// Save will put in the legder interface both the testtuple with its key
// and all the associated composite keys
func (testtuple *Testtuple) Save(db *LedgerDB, testtupleKey string) error {
	var err error
	if err = db.Add(testtupleKey, testtuple); err != nil {
		return err
	}

	// create composite keys
	if err = db.CreateIndex("testtuple~objective~certified~key", []string{"testtuple", testtuple.ObjectiveKey, strconv.FormatBool(testtuple.Certified), testtupleKey}); err != nil {
		return err
	}
	if err = db.CreateIndex("testtuple~algo~key", []string{"testtuple", testtuple.AlgoKey, testtupleKey}); err != nil {
		return err
	}
	if err = db.CreateIndex("testtuple~worker~status~key", []string{"testtuple", testtuple.Dataset.Worker, testtuple.Status, testtupleKey}); err != nil {
		return err
	}
	if err = db.CreateIndex("testtuple~traintuple~certified~key", []string{"testtuple", testtuple.TraintupleKey, strconv.FormatBool(testtuple.Certified), testtupleKey}); err != nil {
		return err
	}
	if testtuple.Tag != "" {
		err = db.CreateIndex("testtuple~tag~key", []string{"traintuple", testtuple.Tag, testtupleKey})
		if err != nil {
			return err
		}
	}
	return nil
}

// -------------------------------------
// Smart contracts related to testuples
// -------------------------------------

// createTesttuple adds a Testtuple in the ledger
func createTesttuple(db *LedgerDB, args []string) (outputKey, error) {

	inp := inputTesttuple{}
	err := AssetFromJSON(args, &inp)
	if err != nil {
		return outputKey{}, err
	}
	key, err := createTesttupleInternal(db, inp)
	if err != nil {
		return outputKey{}, err
	}

	return outputKey{Key: key}, nil
}

func createTesttupleInternal(db *LedgerDB, inp inputTesttuple) (string, error) {
	// check validity of input arg and set testtuple
	testtuple := Testtuple{}
	err := testtuple.SetFromTraintuple(db, inp.TraintupleKey)
	if err != nil {
		return "", err
	}
	err = testtuple.SetFromInput(db, inp)
	if err != nil {
		return "", err
	}
	err = testtuple.AddToComputePlan(db, testtuple.Key)
	if err != nil {
		return "", err
	}
	err = testtuple.Save(db, testtuple.Key)
	if err != nil {
		return "", err
	}
	err = db.AddTupleEvent(testtuple.Key)
	if err != nil {
		return "", err
	}

	return testtuple.Key, nil
}

// logStartTest modifies a testtuple by changing its status from todo to doing
func logStartTest(db *LedgerDB, args []string) (o outputTesttuple, err error) {
	status := StatusDoing
	inp := inputKey{}
	err = AssetFromJSON(args, &inp)
	if err != nil {
		return
	}

	// get testtuple, check validity of the update, and update its status
	testtuple, err := db.GetTesttuple(inp.Key)
	if err != nil {
		return
	}

	if err = validateTupleOwner(db, testtuple.Dataset.Worker); err != nil {
		return
	}
	if err = testtuple.commitStatusUpdate(db, inp.Key, status); err != nil {
		return
	}
	err = o.Fill(db, testtuple)
	return
}

// logSuccessTest modifies a testtuple by changing its status to done, reports perf and logs
func logSuccessTest(db *LedgerDB, args []string) (o outputTesttuple, err error) {
	status := StatusDone
	inp := inputLogSuccessTest{}
	err = AssetFromJSON(args, &inp)
	if err != nil {
		return
	}

	testtuple, err := db.GetTesttuple(inp.Key)
	if err != nil {
		return
	}

	testtuple.Dataset.Perf = inp.Perf
	testtuple.Log += inp.Log

	if err = validateTupleOwner(db, testtuple.Dataset.Worker); err != nil {
		return
	}
	if err = testtuple.commitStatusUpdate(db, inp.Key, status); err != nil {
		return
	}
	err = o.Fill(db, testtuple)
	return
}

// logFailTest modifies a testtuple by changing its status to fail and reports associated logs
func logFailTest(db *LedgerDB, args []string) (o outputTesttuple, err error) {
	status := StatusFailed
	inp := inputLogFailTest{}
	err = AssetFromJSON(args, &inp)
	if err != nil {
		return
	}

	// get, update and commit testtuple
	testtuple, err := db.GetTesttuple(inp.Key)
	if err != nil {
		return
	}

	testtuple.Log += inp.Log

	if err = validateTupleOwner(db, testtuple.Dataset.Worker); err != nil {
		return
	}
	if err = testtuple.commitStatusUpdate(db, inp.Key, status); err != nil {
		return
	}
	err = o.Fill(db, testtuple)
	return
}

// queryTesttuple returns a testtuple of the ledger given its key
func queryTesttuple(db *LedgerDB, args []string) (out outputTesttuple, err error) {
	inp := inputKey{}
	err = AssetFromJSON(args, &inp)
	if err != nil {
		return
	}
	testtuple, err := db.GetTesttuple(inp.Key)
	if err != nil {
		return
	}
	if testtuple.AssetType != TesttupleType {
		err = errors.NotFound("no element with key %s", inp.Key)
		return
	}
	err = out.Fill(db, testtuple)
	return
}

// queryTesttuples returns all testtuples of the ledger
func queryTesttuples(db *LedgerDB, args []string) ([]outputTesttuple, error) {
	outTesttuples := []outputTesttuple{}

	if len(args) != 0 {
		err := errors.BadRequest("incorrect number of arguments, expecting nothing")
		return outTesttuples, err
	}
	elementsKeys, err := db.GetIndexKeys("testtuple~traintuple~certified~key", []string{"testtuple"})
	if err != nil {
		return outTesttuples, err
	}
	nb := getLimitedNbSliceElements(elementsKeys)
	for _, key := range elementsKeys[:nb] {
		var out outputTesttuple
		out, err = getOutputTesttuple(db, key)
		if err != nil {
			return outTesttuples, err
		}
		outTesttuples = append(outTesttuples, out)
	}
	return outTesttuples, nil
}

// -----------------------------------------------
// Utils for smartcontracts related to testtuples
// -----------------------------------------------

// getOutputTesttuple takes as input a testtuple key and returns the outputTesttuple
func getOutputTesttuple(db *LedgerDB, testtupleKey string) (outTesttuple outputTesttuple, err error) {
	testtuple, err := db.GetTesttuple(testtupleKey)
	if err != nil {
		return
	}
	err = outTesttuple.Fill(db, testtuple)
	return
}

// getOutputTesttuples takes as input a list of keys and returns a paylaod containing a list of associated retrieved elements
func getOutputTesttuples(db *LedgerDB, testtupleKeys []string) (outTesttuples []outputTesttuple, err error) {
	for _, key := range testtupleKeys {
		var outputTesttuple outputTesttuple
		outputTesttuple, err = getOutputTesttuple(db, key)
		if err != nil {
			return
		}
		outTesttuples = append(outTesttuples, outputTesttuple)
	}
	return
}

// validateNewStatus verifies that the new status is consistent with the tuple current status
func (testtuple *Testtuple) validateNewStatus(db *LedgerDB, status string) error {
	// check validity of worker and change of status
	return checkUpdateTuple(db, testtuple.Dataset.Worker, testtuple.Status, status)
}

// commitStatusUpdate update the testtuple status in the ledger
func (testtuple *Testtuple) commitStatusUpdate(db *LedgerDB, testtupleKey string, newStatus string) error {
	if testtuple.Status == newStatus {
		return nil
	}

	if err := testtuple.validateNewStatus(db, newStatus); err != nil {
		return errors.Internal("update testtuple %s failed: %s", testtupleKey, err.Error())
	}

	// do not update if previous status is already Done, Failed, Todo, Doing
	if StatusAborted == newStatus && testtuple.Status != StatusWaiting {
		return nil
	}

	oldStatus := testtuple.Status
	testtuple.Status = newStatus

	if err := db.Put(testtupleKey, testtuple); err != nil {
		return errors.Internal("failed to update testtuple status to %s with key %s", newStatus, testtupleKey)
	}

	// update associated composite key
	indexName := "testtuple~worker~status~key"
	oldAttributes := []string{"testtuple", testtuple.Dataset.Worker, oldStatus, testtupleKey}
	newAttributes := []string{"testtuple", testtuple.Dataset.Worker, testtuple.Status, testtupleKey}
	if err := db.UpdateIndex(indexName, oldAttributes, newAttributes); err != nil {
		return err
	}
	if err := UpdateComputePlanState(db, testtuple.ComputePlanKey, newStatus, testtupleKey); err != nil {
		return err
	}
	logger.Infof("testtuple %s status updated: %s (from=%s)", testtupleKey, newStatus, oldStatus)
	return nil
}
