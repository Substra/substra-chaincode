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
	"fmt"
	"strconv"
)

// -------------------------------------------------------------------------------------------
// Methods on receivers traintuple
// -------------------------------------------------------------------------------------------

// SetFromInput is a method of the receiver Traintuple.
// It uses the inputTraintuple to check and set the traintuple's parameters
// which don't depend on previous traintuples values :
//  - AssetType
//  - Creator & permissions
//  - Tag
//  - AlgoKey & ObjectiveKey
//  - Dataset
func (traintuple *Traintuple) SetFromInput(db LedgerDB, inp inputTraintuple) error {

	// TODO later: check permissions
	// find associated creator and check permissions (TODO later)
	creator, err := GetTxCreator(db.cc)
	if err != nil {
		return err
	}
	traintuple.AssetType = TraintupleType
	traintuple.Creator = creator
	traintuple.Tag = inp.Tag
	algo, err := db.GetAlgo(inp.AlgoKey)
	if err != nil {
		return errors.BadRequest(err, "could not retrieve algo with key %s", inp.AlgoKey)
	}
	if !algo.Permissions.CanProcess(algo.Owner, creator) {
		return errors.Forbidden("not authorized to process algo %s", inp.AlgoKey)
	}
	traintuple.AlgoKey = inp.AlgoKey

	// check objective exists
	objective, err := db.GetObjective(inp.ObjectiveKey)
	if err != nil {
		return errors.BadRequest(err, "could not retrieve objective with key %s", inp.ObjectiveKey)
	}
	if !objective.Permissions.CanProcess(objective.Owner, creator) {
		return errors.Forbidden("not authorized to process objective %s", inp.ObjectiveKey)
	}
	traintuple.ObjectiveKey = inp.ObjectiveKey

	// check if DataSampleKeys are from the same dataManager and if they are not test only dataSample
	_, trainOnly, err := checkSameDataManager(db, inp.DataManagerKey, inp.DataSampleKeys)
	if err != nil {
		return err
	}
	if !trainOnly {
		return errors.BadRequest("not possible to create a traintuple with test only data")
	}

	dataManager, err := db.GetDataManager(inp.DataManagerKey)
	if err != nil {
		return errors.BadRequest(err, "could not retrieve dataManager with key %s", inp.DataManagerKey)
	}
	if !dataManager.Permissions.CanProcess(dataManager.Owner, creator) {
		return errors.Forbidden("not authorized to process dataManager %s", inp.DataManagerKey)
	}

	traintuple.Permissions = MergePermissions(dataManager.Permissions, algo.Permissions)

	// fill traintuple.Dataset from dataManager and dataSample
	traintuple.Dataset = &Dataset{
		DataManagerKey: inp.DataManagerKey,
		DataSampleKeys: inp.DataSampleKeys,
	}
	traintuple.Dataset.Worker, err = getDataManagerOwner(db, traintuple.Dataset.DataManagerKey)
	return err
}

// SetFromParents set the status of the traintuple depending on its "parents",
// i.e. the traintuples from which it received the outModels as inModels.
// Also it's InModelKeys are set.
func (traintuple *Traintuple) SetFromParents(db LedgerDB, inModels []string) error {
	status := StatusTodo
	parentTraintupleKeys := inModels
	for _, parentTraintupleKey := range parentTraintupleKeys {
		parentTraintuple, err := db.GetTraintuple(parentTraintupleKey)
		if err != nil {
			err = errors.BadRequest(err, "could not retrieve parent traintuple with key %s %d", parentTraintupleKeys, len(parentTraintupleKeys))
			return err
		}
		// set traintuple to waiting if one of the parent traintuples is not done
		if parentTraintuple.OutModel == nil {
			status = StatusWaiting
		}
		traintuple.InModelKeys = append(traintuple.InModelKeys, parentTraintupleKey)
	}
	traintuple.Status = status
	return nil
}

// GetKey return the key of the traintuple depending on its key parameters.
func (traintuple *Traintuple) GetKey() string {
	hashKeys := []string{traintuple.Creator, traintuple.AlgoKey, traintuple.Dataset.DataManagerKey}
	hashKeys = append(hashKeys, traintuple.Dataset.DataSampleKeys...)
	hashKeys = append(hashKeys, traintuple.InModelKeys...)
	return HashForKey("traintuple", hashKeys...)

}

// AddToComputePlan set the traintuple's parameters that determines if it's part of on ComputePlan and how.
// It uses the inputTraintuple values as follow:
//  - If neither ComputePlanID nor rank is set it returns immediately
//  - If rank is 0 and ComputePlanID empty, it's start a new one using this traintuple key
//  - If rank and ComputePlanID are set, it checks if there are coherent with previous ones and set it.
func (traintuple *Traintuple) AddToComputePlan(db LedgerDB, inp inputTraintuple, traintupleKey string) error {
	// check ComputePlanID and Rank and set it when required
	var err error
	if inp.Rank == "" {
		if inp.ComputePlanID != "" {
			return errors.BadRequest("invalid inputs, a ComputePlan should have a rank")
		}
		return nil
	}
	traintuple.Rank, err = strconv.Atoi(inp.Rank)
	if err != nil {
		return err
	}
	if inp.ComputePlanID == "" {
		if traintuple.Rank != 0 {
			err = errors.BadRequest("invalid inputs, a new ComputePlan should have a rank 0")
			return err
		}
		traintuple.ComputePlanID = traintupleKey
		return nil
	}
	var ttKeys []string
	ttKeys, err = db.GetIndexKeys("traintuple~computeplanid~worker~rank~key", []string{"traintuple", inp.ComputePlanID})
	if err != nil {
		return err
	}
	if len(ttKeys) == 0 {
		return errors.BadRequest("cannot find the ComputePlanID %s", inp.ComputePlanID)
	}
	for _, ttKey := range ttKeys {
		FLTraintuple, err := db.GetTraintuple(ttKey)
		if err != nil {
			return err
		}
		if FLTraintuple.AlgoKey != inp.AlgoKey {
			return errors.BadRequest("previous traintuple for ComputePlanID %s does not have the same algo key %s", inp.ComputePlanID, inp.AlgoKey)
		}
	}

	ttKeys, err = db.GetIndexKeys("traintuple~computeplanid~worker~rank~key", []string{"traintuple", inp.ComputePlanID, traintuple.Dataset.Worker, inp.Rank})
	if err != nil {
		return err
	} else if len(ttKeys) > 0 {
		err = errors.BadRequest("ComputePlanID %s with worker %s rank %d already exists", inp.ComputePlanID, traintuple.Dataset.Worker, traintuple.Rank)
		return err
	}

	traintuple.ComputePlanID = inp.ComputePlanID

	return nil
}

// Save will put in the legder interface both the traintuple with its key
// and all the associated composite keys
func (traintuple *Traintuple) Save(db LedgerDB, traintupleKey string) error {

	// store in ledger
	if err := db.Add(traintupleKey, traintuple); err != nil {
		return err
	}

	// create composite keys
	if err := db.CreateIndex("traintuple~algo~key", []string{"traintuple", traintuple.AlgoKey, traintupleKey}); err != nil {
		return err
	}
	if err := db.CreateIndex("traintuple~worker~status~key", []string{"traintuple", traintuple.Dataset.Worker, traintuple.Status, traintupleKey}); err != nil {
		return err
	}
	for _, inModelKey := range traintuple.InModelKeys {
		if err := db.CreateIndex("traintuple~inModel~key", []string{"traintuple", inModelKey, traintupleKey}); err != nil {
			return err
		}
	}
	if traintuple.ComputePlanID != "" {
		if err := db.CreateIndex("traintuple~computeplanid~worker~rank~key", []string{"traintuple", traintuple.ComputePlanID, traintuple.Dataset.Worker, strconv.Itoa(traintuple.Rank), traintupleKey}); err != nil {
			return err
		}
		if err := db.CreateIndex("computeplan~id", []string{"computeplan", traintuple.ComputePlanID}); err != nil {
			return err
		}
	}
	if traintuple.Tag != "" {
		err := db.CreateIndex("traintuple~tag~key", []string{"traintuple", traintuple.Tag, traintupleKey})
		if err != nil {
			return err
		}
	}
	return nil
}

// -------------------------------------------------------------------------------------------
// Smart contracts related to traintuples
// -------------------------------------------------------------------------------------------

// createTraintuple adds a Traintuple in the ledger
func createTraintuple(db LedgerDB, args []string) (map[string]string, error) {
	inp := inputTraintuple{}
	err := AssetFromJSON(args, &inp)
	if err != nil {
		return nil, err
	}

	traintuple := Traintuple{}
	err = traintuple.SetFromInput(db, inp)
	if err != nil {
		return nil, err
	}
	err = traintuple.SetFromParents(db, inp.InModels)
	if err != nil {
		return nil, err
	}
	traintupleKey := traintuple.GetKey()
	// Test if the key (ergo the traintuple) already exists
	tupleExists, err := db.KeyExists(traintupleKey)
	if err != nil {
		return nil, err
	}
	if tupleExists {
		return nil, errors.Conflict("traintuple already exists").WithKey(traintupleKey)
	}
	err = traintuple.AddToComputePlan(db, inp, traintupleKey)
	if err != nil {
		return nil, err
	}
	err = traintuple.Save(db, traintupleKey)
	if err != nil {
		return nil, err
	}
	out := outputTraintuple{}
	err = out.Fill(db, traintuple, traintupleKey)
	if err != nil {
		return nil, err
	}

	event := TuplesEvent{}
	event.SetTraintuples(out)
	err = SendTuplesEvent(db.cc, event)
	if err != nil {
		return nil, err
	}

	return map[string]string{"key": traintupleKey}, nil
}

// logStartTrain modifies a traintuple by changing its status from todo to doing
func logStartTrain(db LedgerDB, args []string) (outputTraintuple outputTraintuple, err error) {
	inp := inputHash{}
	err = AssetFromJSON(args, &inp)
	if err != nil {
		return
	}

	// get traintuple, check validity of the update
	traintuple, err := db.GetTraintuple(inp.Key)
	if err != nil {
		return
	}
	if err = validateTupleOwner(db, traintuple.Dataset.Worker); err != nil {
		return
	}
	if err = traintuple.commitStatusUpdate(db, inp.Key, StatusDoing); err != nil {
		return
	}
	err = outputTraintuple.Fill(db, traintuple, inp.Key)
	return
}

// logSuccessTrain modifies a traintuple by changing its status from doing to done
// reports logs and associated performances
func logSuccessTrain(db LedgerDB, args []string) (outputTraintuple outputTraintuple, err error) {
	inp := inputLogSuccessTrain{}
	err = AssetFromJSON(args, &inp)
	if err != nil {
		return
	}
	traintupleKey := inp.Key

	// get, update and commit traintuple
	traintuple, err := db.GetTraintuple(traintupleKey)
	if err != nil {
		return
	}
	traintuple.Perf = inp.Perf
	traintuple.OutModel = &HashDress{
		Hash:           inp.OutModel.Hash,
		StorageAddress: inp.OutModel.StorageAddress}
	traintuple.Log += inp.Log

	if err = validateTupleOwner(db, traintuple.Dataset.Worker); err != nil {
		return
	}
	if err = traintuple.commitStatusUpdate(db, traintupleKey, StatusDone); err != nil {
		return
	}

	// update depending tuples
	event := TuplesEvent{}
	err = UpdateTraintupleChildren(db, traintupleKey, traintuple.Status, &event)
	if err != nil {
		return
	}

	err = UpdateTesttupleChildren(db, traintupleKey, traintuple.Status, &event)
	if err != nil {
		return
	}

	err = outputTraintuple.Fill(db, traintuple, inp.Key)
	if err != nil {
		return
	}

	err = SendTuplesEvent(db.cc, event)
	if err != nil {
		return
	}

	return
}

// logFailTrain modifies a traintuple by changing its status to fail and reports associated logs
func logFailTrain(db LedgerDB, args []string) (outputTraintuple outputTraintuple, err error) {
	inp := inputLogFailTrain{}
	err = AssetFromJSON(args, &inp)
	if err != nil {
		return
	}

	// get, update and commit traintuple
	traintuple, err := db.GetTraintuple(inp.Key)
	if err != nil {
		return
	}
	traintuple.Log += inp.Log

	if err = validateTupleOwner(db, traintuple.Dataset.Worker); err != nil {
		return
	}
	if err = traintuple.commitStatusUpdate(db, inp.Key, StatusFailed); err != nil {
		return
	}

	if err = outputTraintuple.Fill(db, traintuple, inp.Key); err != nil {
		return
	}

	// update depending tuples
	event := TuplesEvent{}
	err = UpdateTesttupleChildren(db, inp.Key, traintuple.Status, &event)
	if err != nil {
		return
	}

	err = UpdateTraintupleChildren(db, inp.Key, traintuple.Status, &event)
	if err != nil {
		return
	}

	err = SendTuplesEvent(db.cc, event)
	if err != nil {
		return
	}

	return
}

// queryTraintuple returns info about a traintuple given its key
func queryTraintuple(db LedgerDB, args []string) (outputTraintuple outputTraintuple, err error) {
	inp := inputHash{}
	err = AssetFromJSON(args, &inp)
	if err != nil {
		return
	}
	traintuple, err := db.GetTraintuple(inp.Key)
	if err != nil {
		return
	}
	if traintuple.AssetType != TraintupleType {
		err = errors.NotFound("no element with key %s", inp.Key)
		return
	}
	err = outputTraintuple.Fill(db, traintuple, inp.Key)
	return
}

// queryTraintuples returns all traintuples
func queryTraintuples(db LedgerDB, args []string) ([]outputTraintuple, error) {
	outTraintuples := []outputTraintuple{}

	if len(args) != 0 {
		err := errors.BadRequest("incorrect number of arguments, expecting nothing")
		return outTraintuples, err
	}
	elementsKeys, err := db.GetIndexKeys("traintuple~algo~key", []string{"traintuple"})
	if err != nil {
		return outTraintuples, err
	}
	for _, key := range elementsKeys {
		outputTraintuple, err := getOutputTraintuple(db, key)
		if err != nil {
			return outTraintuples, err
		}
		outTraintuples = append(outTraintuples, outputTraintuple)
	}
	return outTraintuples, nil
}

// -----------------------------------------------
// Utils for smartcontracts related to traintuples
// -----------------------------------------------

// getOutputTraintuple takes as input a traintuple key and returns the outputTraintuple
func getOutputTraintuple(db LedgerDB, traintupleKey string) (outTraintuple outputTraintuple, err error) {
	traintuple, err := db.GetTraintuple(traintupleKey)
	if err != nil {
		return
	}
	err = outTraintuple.Fill(db, traintuple, traintupleKey)
	return
}

// getOutputTraintuples takes as input a list of keys and returns a paylaod containing a list of associated retrieved elements
func getOutputTraintuples(db LedgerDB, traintupleKeys []string) (outTraintuples []outputTraintuple, err error) {
	for _, key := range traintupleKeys {
		var outputTraintuple outputTraintuple
		outputTraintuple, err = getOutputTraintuple(db, key)
		if err != nil {
			return
		}
		outTraintuples = append(outTraintuples, outputTraintuple)
	}
	return
}

// validateNewStatus verifies that the new status is consistent with the tuple current status
func (traintuple *Traintuple) validateNewStatus(db LedgerDB, status string) error {
	// check validity of worker and change of status
	return checkUpdateTuple(db, traintuple.Dataset.Worker, traintuple.Status, status)
}

// UpdateTraintupleChildren updates the status of waiting trainuples  InModels of traintuples once they have been trained (succesfully or failed)
func UpdateTraintupleChildren(db LedgerDB, traintupleKey string, traintupleStatus string, event *TuplesEvent) error {
	// get traintuples having as inModels the input traintuple
	childTraintupleKeys, err := db.GetIndexKeys("traintuple~inModel~key", []string{"traintuple", traintupleKey})
	if err != nil {
		return fmt.Errorf("error while getting associated traintuples to update their inModel")
	}
	childCompositeTraintupleKeys, err := db.GetIndexKeys("compositeTraintuple~inModel~key", []string{"compositeTraintuple", traintupleKey})
	if err != nil {
		return fmt.Errorf("error while getting associated composite traintuples to update their inModel")
	}
	childAggregateTupleKeys, err := db.GetIndexKeys("aggregateTuple~inModel~key", []string{"aggregateTuple", traintupleKey})
	if err != nil {
		return fmt.Errorf("error while getting associated aggregate tuples to update their inModel")
	}

	allChildKeys := append(append(childTraintupleKeys, childCompositeTraintupleKeys...), childAggregateTupleKeys...)

	for _, childTraintupleKey := range allChildKeys {
		childTraintupleType, childTraintupleStatus, err := db.GetGenericTraintuple(childTraintupleKey)
		if err != nil {
			return err
		}

		if childTraintupleStatus == StatusFailed {
			// traintuple is already failed, don't update it
			continue
		}
		if childTraintupleStatus != StatusWaiting {
			return fmt.Errorf("traintuple %s has invalid status : '%s' instead of waiting", childTraintupleKey, childTraintupleStatus)
		}

		// Update the child traintuple and get its new status
		switch childTraintupleType {
		case TraintupleType:
			childTraintupleStatus, err = UpdateTraintupleChild(db, traintupleKey, childTraintupleKey, traintupleStatus, event)
			if err != nil {
				return err
			}
		case CompositeTraintupleType:
			childTraintupleStatus, err = UpdateCompositeTraintupleChild(db, traintupleKey, childTraintupleKey, traintupleStatus, event)
			if err != nil {
				return err
			}
		case AggregateTupleType:
			childTraintupleStatus, err = UpdateAggregateTupleChild(db, traintupleKey, childTraintupleKey, traintupleStatus, event)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("Unknown child traintuple type: %s", childTraintupleType)
		}

		// Recursively call for an update on this child's children
		err = UpdateTesttupleChildren(db, childTraintupleKey, childTraintupleStatus, event)
		if err != nil {
			return err
		}

		err = UpdateTraintupleChildren(db, childTraintupleKey, childTraintupleStatus, event)
		if err != nil {
			return err
		}
	}
	return nil
}

// UpdateTraintupleChild updates the status of a waiting trainuple, given the new parent traintuple status
func UpdateTraintupleChild(db LedgerDB, parentTraintupleKey string, childTraintupleKey string, traintupleStatus string, event *TuplesEvent) (childStatus string, err error) {
	// get and update traintuple
	childTraintuple, err := db.GetTraintuple(childTraintupleKey)
	if err != nil {
		return
	}

	childStatus = childTraintuple.Status

	// get traintuple new status
	var newStatus string
	if traintupleStatus == StatusFailed {
		newStatus = StatusFailed
	} else if traintupleStatus == StatusDone {
		ready, _err := childTraintuple.isReady(db, parentTraintupleKey)
		if _err != nil {
			err = _err
			return
		}
		if ready {
			newStatus = StatusTodo
		}
	}

	// commit new status
	if newStatus == "" {
		return
	}
	if err = childTraintuple.commitStatusUpdate(db, childTraintupleKey, newStatus); err != nil {
		return
	}

	// update return value after status update
	childStatus = childTraintuple.Status

	if newStatus == StatusTodo {
		out := outputTraintuple{}
		err = out.Fill(db, childTraintuple, childTraintupleKey)
		if err != nil {
			return
		}
		event.AddTraintuple(out)
	}

	return
}

func (traintuple *Traintuple) isReady(db LedgerDB, newDoneTraintupleKey string) (ready bool, err error) {
	return IsReady(db, traintuple.InModelKeys, newDoneTraintupleKey)
}

// IsReady checks if inModels of a traintuple have been trained, except the newDoneTraintupleKey (since the transaction is not commited)
func IsReady(db LedgerDB, inModelKeys []string, newDoneTraintupleKey string) (ready bool, err error) {
	for _, key := range inModelKeys {
		// don't check newly done traintuple
		if key == newDoneTraintupleKey {
			continue
		}
		_, status, err := db.GetGenericTraintuple(key)
		if err != nil {
			return false, err
		}
		if status != StatusDone {
			return false, nil
		}
	}
	return true, nil
}

// commitStatusUpdate update the traintuple status in the ledger
func (traintuple *Traintuple) commitStatusUpdate(db LedgerDB, traintupleKey string, newStatus string) error {
	if traintuple.Status == newStatus {
		return fmt.Errorf("cannot update traintuple %s - status already %s", traintupleKey, newStatus)
	}

	if err := traintuple.validateNewStatus(db, newStatus); err != nil {
		return fmt.Errorf("update traintuple %s failed: %s", traintupleKey, err.Error())
	}

	oldStatus := traintuple.Status
	traintuple.Status = newStatus
	if err := db.Put(traintupleKey, traintuple); err != nil {
		return fmt.Errorf("failed to update traintuple %s - %s", traintupleKey, err.Error())
	}

	// update associated composite keys
	indexName := "traintuple~worker~status~key"
	oldAttributes := []string{"traintuple", traintuple.Dataset.Worker, oldStatus, traintupleKey}
	newAttributes := []string{"traintuple", traintuple.Dataset.Worker, traintuple.Status, traintupleKey}
	if err := db.UpdateIndex(indexName, oldAttributes, newAttributes); err != nil {
		return err
	}
	logger.Infof("traintuple %s status updated: %s (from=%s)", traintupleKey, newStatus, oldStatus)
	return nil
}

// UpdateTesttupleChildren update testtuples status associated with a done or failed traintuple
func UpdateTesttupleChildren(db LedgerDB, traintupleKey string, traintupleStatus string, event *TuplesEvent) error {
	var newStatus string
	switch {
	case traintupleStatus == StatusFailed:
		newStatus = StatusFailed
	case traintupleStatus == StatusDone:
		newStatus = StatusTodo
	default:
		return nil
	}

	indexName := "testtuple~traintuple~certified~key"
	// get testtuple associated with this traintuple and updates its status
	testtupleKeys, err := db.GetIndexKeys(indexName, []string{"testtuple", traintupleKey})
	if err != nil {
		return err
	}
	for _, testtupleKey := range testtupleKeys {
		// get and update testtuple
		testtuple, err := db.GetTesttuple(testtupleKey)
		if err != nil {
			return err
		}
		testtuple.TraintupleKey = traintupleKey

		if err := testtuple.commitStatusUpdate(db, testtupleKey, newStatus); err != nil {
			return err
		}

		if newStatus == StatusTodo {
			out := outputTesttuple{}
			err = out.Fill(db, testtupleKey, testtuple)
			if err != nil {
				return err
			}
			event.AddTesttuple(out)
		}
	}
	return nil
}
