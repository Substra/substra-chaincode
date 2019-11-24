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

// SetFromInput is a method of the receiver AggregateTuple.
// It uses the inputAggregateTuple to check and set the aggregate tuple's parameters
// which don't depend on previous traintuples values :
//  - AssetType
//  - Creator & permissions
//  - Tag
//  - AlgoKey & ObjectiveKey
func (tuple *AggregateTuple) SetFromInput(db LedgerDB, inp inputAggregateTuple) error {
	creator, err := GetTxCreator(db.cc)
	if err != nil {
		return err
	}
	tuple.AssetType = AggregateTupleType
	tuple.Creator = creator
	tuple.Tag = inp.Tag
	algo, err := db.GetAggregateAlgo(inp.AlgoKey)
	if err != nil {
		return errors.BadRequest(err, "could not retrieve algo with key %s", inp.AlgoKey)
	}
	if !algo.Permissions.CanProcess(algo.Owner, creator) {
		return errors.Forbidden("not authorized to process algo %s", inp.AlgoKey)
	}
	tuple.AlgoKey = inp.AlgoKey

	// check objective exists
	objective, err := db.GetObjective(inp.ObjectiveKey)
	if err != nil {
		return errors.BadRequest(err, "could not retrieve objective with key %s", inp.ObjectiveKey)
	}
	if !objective.Permissions.CanProcess(objective.Owner, creator) {
		return errors.Forbidden("not authorized to process objective %s", inp.ObjectiveKey)
	}
	tuple.ObjectiveKey = inp.ObjectiveKey

	// TODO (aggregate): uncomment + add test
	// traintuple.Permissions = MergePermissions(dataManager.Permissions, algo.Permissions)

	tuple.Worker = inp.Worker
	return nil
}

// SetFromParents set the status of the aggregate tuple depending on its "parents",
// i.e. the traintuples from which it received the outModels as inModels.
// Also it's InModelKeys are set.
func (tuple *AggregateTuple) SetFromParents(db LedgerDB, inModels []string) error {
	status := StatusTodo

	for _, parentTraintupleKey := range inModels {
		hashDress, err := db.GetOutModelHashDress(parentTraintupleKey, TrunkType, []AssetType{TraintupleType, CompositeTraintupleType, AggregateTupleType})
		if err != nil {
			return err
		}
		if hashDress == nil {
			status = StatusWaiting
		}

		tuple.InModelKeys = append(tuple.InModelKeys, parentTraintupleKey)
	}

	tuple.Status = status
	return nil
}

// GetKey return the key of the aggregate tuple depending on its key parameters.
func (tuple *AggregateTuple) GetKey() string {
	hashKeys := []string{tuple.Creator, tuple.AlgoKey}
	hashKeys = append(hashKeys, tuple.InModelKeys...)
	return HashForKey("aggregate-traintuple", hashKeys...)
}

// AddToComputePlan set the aggregate tuple's parameters that determines if it's part of on ComputePlan and how.
// It uses the inputAggregateTuple values as follow:
//  - If neither ComputePlanID nor rank is set it returns immediately
//  - If rank is 0 and ComputePlanID empty, it's start a new one using this traintuple key
//  - If rank and ComputePlanID are set, it checks if there are coherent with previous ones and set it.
func (tuple *AggregateTuple) AddToComputePlan(db LedgerDB, inp inputAggregateTuple, traintupleKey string) error {
	// check ComputePlanID and Rank and set it when required
	var err error
	if inp.Rank == "" {
		if inp.ComputePlanID != "" {
			return errors.BadRequest("invalid inputs, a ComputePlan should have a rank")
		}
		return nil
	}
	tuple.Rank, err = strconv.Atoi(inp.Rank)
	if err != nil {
		return err
	}
	if inp.ComputePlanID == "" {
		if tuple.Rank != 0 {
			err = errors.BadRequest("invalid inputs, a new ComputePlan should have a rank 0")
			return err
		}
		tuple.ComputePlanID = traintupleKey
		return nil
	}
	var ttKeys []string
	ttKeys, err = db.GetIndexKeys("aggregateTuple~computeplanid~worker~rank~key", []string{"aggregateTuple", inp.ComputePlanID})
	if err != nil {
		return err
	}
	if len(ttKeys) == 0 {
		return errors.BadRequest("cannot find the ComputePlanID %s", inp.ComputePlanID)
	}
	for _, ttKey := range ttKeys {
		FLTraintuple, err := db.GetAggregateTuple(ttKey)
		if err != nil {
			return err
		}
		if FLTraintuple.AlgoKey != inp.AlgoKey {
			return errors.BadRequest("previous traintuple for ComputePlanID %s does not have the same algo key %s", inp.ComputePlanID, inp.AlgoKey)
		}
	}

	ttKeys, err = db.GetIndexKeys("aggregateTuple~computeplanid~worker~rank~key", []string{"aggregateTuple", inp.ComputePlanID, tuple.Worker, inp.Rank})
	if err != nil {
		return err
	} else if len(ttKeys) > 0 {
		err = errors.BadRequest("ComputePlanID %s with worker %s rank %d already exists", inp.ComputePlanID, tuple.Worker, tuple.Rank)
		return err
	}

	tuple.ComputePlanID = inp.ComputePlanID

	return nil
}

// Save will put in the legder interface both the aggregate tuple with its key
// and all the associated composite keys
func (tuple *AggregateTuple) Save(db LedgerDB, traintupleKey string) error {

	// store in ledger
	if err := db.Add(traintupleKey, tuple); err != nil {
		return err
	}

	// create composite keys
	if err := db.CreateIndex("aggregateTuple~algo~key", []string{"aggregateTuple", tuple.AlgoKey, traintupleKey}); err != nil {
		return err
	}
	if err := db.CreateIndex("aggregateTuple~worker~status~key", []string{"aggregateTuple", tuple.Worker, tuple.Status, traintupleKey}); err != nil {
		return err
	}
	for _, inModelKey := range tuple.InModelKeys {
		if err := db.CreateIndex("aggregateTuple~inModel~key", []string{"aggregateTuple", inModelKey, traintupleKey}); err != nil {
			return err
		}
	}
	if tuple.ComputePlanID != "" {
		if err := db.CreateIndex("aggregateTuple~computeplanid~worker~rank~key", []string{"aggregateTuple", tuple.ComputePlanID, tuple.Worker, strconv.Itoa(tuple.Rank), traintupleKey}); err != nil {
			return err
		}
	}
	if tuple.Tag != "" {
		err := db.CreateIndex("aggregateTuple~tag~key", []string{"aggregateTuple", tuple.Tag, traintupleKey})
		if err != nil {
			return err
		}
	}
	return nil
}

// -------------------------------------------------------------------------------------------
// Smart contracts related to aggregate tuples
// -------------------------------------------------------------------------------------------

// createAggregateTuple adds a AggregateTuple in the ledger
func createAggregateTuple(db LedgerDB, args []string) (map[string]string, error) {
	inp := inputAggregateTuple{}
	err := AssetFromJSON(args, &inp)
	if err != nil {
		return nil, err
	}

	traintuple := AggregateTuple{}
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
	out := outputAggregateTuple{}
	err = out.Fill(db, traintuple, traintupleKey)
	if err != nil {
		return nil, err
	}

	event := TuplesEvent{}
	event.SetAggregateTuples(out)
	err = SendTuplesEvent(db.cc, event)
	if err != nil {
		return nil, err
	}

	return map[string]string{"key": traintupleKey}, nil
}

// logStartAggregateTrain modifies a traintuple by changing its status from todo to doing
func logStartAggregateTrain(db LedgerDB, args []string) (outputTraintuple outputAggregateTuple, err error) {
	inp := inputHash{}
	err = AssetFromJSON(args, &inp)
	if err != nil {
		return
	}

	// get traintuple, check validity of the update
	traintuple, err := db.GetAggregateTuple(inp.Key)
	if err != nil {
		return
	}
	if err = validateTupleOwner(db, traintuple.Worker); err != nil {
		return
	}
	if err = traintuple.commitStatusUpdate(db, inp.Key, StatusDoing); err != nil {
		return
	}
	outputTraintuple.Fill(db, traintuple, inp.Key)
	return
}

// logFailAggregateTrain modifies a traintuple by changing its status to fail and reports associated logs
func logFailAggregateTrain(db LedgerDB, args []string) (outputTraintuple outputAggregateTuple, err error) {
	inp := inputLogFailTrain{}
	err = AssetFromJSON(args, &inp)
	if err != nil {
		return
	}

	// get, update and commit traintuple
	traintuple, err := db.GetAggregateTuple(inp.Key)
	if err != nil {
		return
	}
	traintuple.Log += inp.Log

	if err = validateTupleOwner(db, traintuple.Worker); err != nil {
		return
	}
	if err = traintuple.commitStatusUpdate(db, inp.Key, StatusFailed); err != nil {
		return
	}

	outputTraintuple.Fill(db, traintuple, inp.Key)

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

// logSuccessAggregateTrain modifies a traintuple by changing its status from doing to done
// reports logs and associated performances
func logSuccessAggregateTrain(db LedgerDB, args []string) (outputTraintuple outputAggregateTuple, err error) {
	inp := inputLogSuccessTrain{}
	err = AssetFromJSON(args, &inp)
	if err != nil {
		return
	}
	traintupleKey := inp.Key

	// get, update and commit traintuple
	traintuple, err := db.GetAggregateTuple(traintupleKey)
	if err != nil {
		return
	}
	traintuple.OutModel = &HashDress{
		Hash:           inp.OutModel.Hash,
		StorageAddress: inp.OutModel.StorageAddress}
	traintuple.Log += inp.Log

	if err = validateTupleOwner(db, traintuple.Worker); err != nil {
		return
	}
	if err = traintuple.commitStatusUpdate(db, traintupleKey, StatusDone); err != nil {
		return
	}

	event := TuplesEvent{}
	err = UpdateTraintupleChildren(db, traintupleKey, traintuple.Status, &event)
	if err != nil {
		return
	}

	err = UpdateTesttupleChildren(db, traintupleKey, traintuple.Status, &event)
	if err != nil {
		return
	}

	outputTraintuple.Fill(db, traintuple, inp.Key)
	err = SendTuplesEvent(db.cc, event)
	if err != nil {
		return
	}

	return
}

// queryAggregateTuple returns info about an aggregate tuple given its key
func queryAggregateTuple(db LedgerDB, args []string) (outputTraintuple outputAggregateTuple, err error) {
	inp := inputHash{}
	err = AssetFromJSON(args, &inp)
	if err != nil {
		return
	}
	traintuple, err := db.GetAggregateTuple(inp.Key)
	if err != nil {
		return
	}
	if traintuple.AssetType != AggregateTupleType {
		err = errors.NotFound("no element with key %s", inp.Key)
		return
	}
	outputTraintuple.Fill(db, traintuple, inp.Key)
	return
}

// queryAggregateTuples returns all aggregate tuples
func queryAggregateTuples(db LedgerDB, args []string) ([]outputAggregateTuple, error) {
	outTraintuples := []outputAggregateTuple{}

	if len(args) != 0 {
		err := errors.BadRequest("incorrect number of arguments, expecting nothing")
		return outTraintuples, err
	}
	elementsKeys, err := db.GetIndexKeys("aggregateTuple~algo~key", []string{"aggregateTuple"})
	if err != nil {
		return outTraintuples, err
	}
	for _, key := range elementsKeys {
		outputTraintuple, err := getOutputAggregateTuple(db, key)
		if err != nil {
			return outTraintuples, err
		}
		outTraintuples = append(outTraintuples, outputTraintuple)
	}
	return outTraintuples, nil
}

// ----------------------------------------------------------
// Utils for smartcontracts related to aggregate tuples
// ----------------------------------------------------------

// getOutputAggregateTuple takes as input a traintuple key and returns the outputAggregateTuple
func getOutputAggregateTuple(db LedgerDB, traintupleKey string) (outTraintuple outputAggregateTuple, err error) {
	traintuple, err := db.GetAggregateTuple(traintupleKey)
	if err != nil {
		return
	}
	outTraintuple.Fill(db, traintuple, traintupleKey)
	return
}

// UpdateAggregateTupleChild updates the status of a waiting trainuple, given the new parent tuple status
func UpdateAggregateTupleChild(db LedgerDB, parentTraintupleKey string, childTraintupleKey string, traintupleStatus string, event *TuplesEvent) (childStatus string, err error) {
	// get and update traintuple
	childTraintuple, err := db.GetAggregateTuple(childTraintupleKey)
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
		out := outputAggregateTuple{}
		err = out.Fill(db, childTraintuple, childTraintupleKey)
		if err != nil {
			return
		}
		event.AddAggregateTuple(out)
	}

	return
}

func (tuple *AggregateTuple) isReady(db LedgerDB, newDoneTraintupleKey string) (ready bool, err error) {
	return IsReady(db, tuple.InModelKeys, newDoneTraintupleKey)
}

// getOutputAggregateTuples takes as input a list of keys and returns a paylaod containing a list of associated retrieved elements
func getOutputAggregateTuples(db LedgerDB, traintupleKeys []string) (outTraintuples []outputAggregateTuple, err error) {
	for _, key := range traintupleKeys {
		var outputTraintuple outputAggregateTuple
		outputTraintuple, err = getOutputAggregateTuple(db, key)
		if err != nil {
			return
		}
		outTraintuples = append(outTraintuples, outputTraintuple)
	}
	return
}

// commitStatusUpdate update the traintuple status in the ledger
func (tuple *AggregateTuple) commitStatusUpdate(db LedgerDB, traintupleKey string, newStatus string) error {
	if tuple.Status == newStatus {
		return fmt.Errorf("cannot update traintuple %s - status already %s", traintupleKey, newStatus)
	}

	if err := tuple.validateNewStatus(db, newStatus); err != nil {
		return fmt.Errorf("update traintuple %s failed: %s", traintupleKey, err.Error())
	}

	oldStatus := tuple.Status
	tuple.Status = newStatus
	if err := db.Put(traintupleKey, tuple); err != nil {
		return fmt.Errorf("failed to update traintuple %s - %s", traintupleKey, err.Error())
	}

	// update associated composite keys
	indexName := "aggregateTuple~worker~status~key"
	oldAttributes := []string{"aggregateTuple", tuple.Worker, oldStatus, traintupleKey}
	newAttributes := []string{"aggregateTuple", tuple.Worker, tuple.Status, traintupleKey}
	if err := db.UpdateIndex(indexName, oldAttributes, newAttributes); err != nil {
		return err
	}
	logger.Infof("traintuple %s status updated: %s (from=%s)", traintupleKey, newStatus, oldStatus)
	return nil
}

// validateNewStatus verifies that the new status is consistent with the tuple current status
func (tuple *AggregateTuple) validateNewStatus(db LedgerDB, status string) error {
	// check validity of worker and change of status
	if err := checkUpdateTuple(db, tuple.Worker, tuple.Status, status); err != nil {
		return err
	}
	return nil
}
