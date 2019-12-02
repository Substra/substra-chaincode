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
// Methods on receivers Aggregatetuple
// -------------------------------------------------------------------------------------------

// SetFromInput is a method of the receiver Aggregatetuple.
// It uses the inputAggregatetuple to check and set the aggregate tuple's parameters
// which don't depend on previous traintuples values :
//  - AssetType
//  - Creator & permissions
//  - Tag
//  - AlgoKey & ObjectiveKey
func (tuple *Aggregatetuple) SetFromInput(db LedgerDB, inp inputAggregatetuple) error {
	creator, err := GetTxCreator(db.cc)
	if err != nil {
		return err
	}
	tuple.AssetType = AggregatetupleType
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
	// tuple.Permissions = MergePermissions(dataManager.Permissions, algo.Permissions)

	tuple.Worker = inp.Worker
	return nil
}

// SetFromParents set the status of the aggregate tuple depending on its "parents",
// i.e. the traintuples from which it received the outModels as inModels.
// Also it's InModelKeys are set.
func (tuple *Aggregatetuple) SetFromParents(db LedgerDB, inModels []string) error {
	status := StatusTodo

	for _, parentTraintupleKey := range inModels {
		hashDress, err := db.GetOutModelHashDress(parentTraintupleKey, TrunkType, []AssetType{TraintupleType, CompositeTraintupleType, AggregatetupleType})
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
func (tuple *Aggregatetuple) GetKey() string {
	hashKeys := []string{tuple.Creator, tuple.AlgoKey}
	hashKeys = append(hashKeys, tuple.InModelKeys...)
	return HashForKey("aggregate-traintuple", hashKeys...)
}

// AddToComputePlan set the aggregate tuple's parameters that determines if it's part of on ComputePlan and how.
// It uses the inputAggregatetuple values as follow:
//  - If neither ComputePlanID nor rank is set it returns immediately
//  - If rank is 0 and ComputePlanID empty, it's start a new one using this traintuple key
//  - If rank and ComputePlanID are set, it checks if there are coherent with previous ones and set it.
func (tuple *Aggregatetuple) AddToComputePlan(db LedgerDB, inp inputAggregatetuple, traintupleKey string) error {
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
	ttKeys, err = db.GetIndexKeys("aggregatetuple~computeplanid~worker~rank~key", []string{"aggregatetuple", inp.ComputePlanID})
	if err != nil {
		return err
	}
	if len(ttKeys) == 0 {
		return errors.BadRequest("cannot find the ComputePlanID %s", inp.ComputePlanID)
	}
	for _, ttKey := range ttKeys {
		FLTraintuple, err := db.GetAggregatetuple(ttKey)
		if err != nil {
			return err
		}
		if FLTraintuple.AlgoKey != inp.AlgoKey {
			return errors.BadRequest("previous traintuple for ComputePlanID %s does not have the same algo key %s", inp.ComputePlanID, inp.AlgoKey)
		}
	}

	ttKeys, err = db.GetIndexKeys("aggregatetuple~computeplanid~worker~rank~key", []string{"aggregatetuple", inp.ComputePlanID, tuple.Worker, inp.Rank})
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
func (tuple *Aggregatetuple) Save(db LedgerDB, aggregatetupleKey string) error {

	// store in ledger
	if err := db.Add(aggregatetupleKey, tuple); err != nil {
		return err
	}

	// create composite keys
	if err := db.CreateIndex("aggregatetuple~algo~key", []string{"aggregatetuple", tuple.AlgoKey, aggregatetupleKey}); err != nil {
		return err
	}
	if err := db.CreateIndex("aggregatetuple~worker~status~key", []string{"aggregatetuple", tuple.Worker, tuple.Status, aggregatetupleKey}); err != nil {
		return err
	}
	for _, inModelKey := range tuple.InModelKeys {
		if err := db.CreateIndex("aggregatetuple~inModel~key", []string{"aggregatetuple", inModelKey, aggregatetupleKey}); err != nil {
			return err
		}
	}
	if tuple.ComputePlanID != "" {
		if err := db.CreateIndex("aggregatetuple~computeplanid~worker~rank~key", []string{"aggregatetuple", tuple.ComputePlanID, tuple.Worker, strconv.Itoa(tuple.Rank), aggregatetupleKey}); err != nil {
			return err
		}
	}
	if tuple.Tag != "" {
		err := db.CreateIndex("aggregatetuple~tag~key", []string{"aggregatetuple", tuple.Tag, aggregatetupleKey})
		if err != nil {
			return err
		}
	}
	return nil
}

// -------------------------------------------------------------------------------------------
// Smart contracts related to aggregate tuples
// -------------------------------------------------------------------------------------------

// createAggregatetuple adds a Aggregatetuple in the ledger
func createAggregatetuple(db LedgerDB, args []string) (map[string]string, error) {
	inp := inputAggregatetuple{}
	err := AssetFromJSON(args, &inp)
	if err != nil {
		return nil, err
	}

	aggregatetuple := Aggregatetuple{}
	err = aggregatetuple.SetFromInput(db, inp)
	if err != nil {
		return nil, err
	}
	err = aggregatetuple.SetFromParents(db, inp.InModels)
	if err != nil {
		return nil, err
	}

	aggregatetupleKey := aggregatetuple.GetKey()
	// Test if the key (ergo the aggregatetuple) already exists
	tupleExists, err := db.KeyExists(aggregatetupleKey)
	if err != nil {
		return nil, err
	}
	if tupleExists {
		return nil, errors.Conflict("aggregatetuple already exists").WithKey(aggregatetupleKey)
	}
	err = aggregatetuple.AddToComputePlan(db, inp, aggregatetupleKey)
	if err != nil {
		return nil, err
	}
	err = aggregatetuple.Save(db, aggregatetupleKey)
	if err != nil {
		return nil, err
	}
	out := outputAggregatetuple{}
	err = out.Fill(db, aggregatetuple, aggregatetupleKey)
	if err != nil {
		return nil, err
	}

	event := TuplesEvent{}
	event.SetAggregatetuples(out)
	err = SendTuplesEvent(db.cc, event)
	if err != nil {
		return nil, err
	}

	return map[string]string{"key": aggregatetupleKey}, nil
}

// logStartAggregateTrain modifies a aggregatetuple by changing its status from todo to doing
func logStartAggregateTrain(db LedgerDB, args []string) (outputAggregatetuple outputAggregatetuple, err error) {
	inp := inputHash{}
	err = AssetFromJSON(args, &inp)
	if err != nil {
		return
	}

	// get aggregatetuple, check validity of the update
	aggregatetuple, err := db.GetAggregatetuple(inp.Key)
	if err != nil {
		return
	}
	if err = validateTupleOwner(db, aggregatetuple.Worker); err != nil {
		return
	}
	if err = aggregatetuple.commitStatusUpdate(db, inp.Key, StatusDoing); err != nil {
		return
	}
	outputAggregatetuple.Fill(db, aggregatetuple, inp.Key)
	return
}

// logFailAggregateTrain modifies a aggregatetuple by changing its status to fail and reports associated logs
func logFailAggregateTrain(db LedgerDB, args []string) (outputAggregatetuple outputAggregatetuple, err error) {
	inp := inputLogFailTrain{}
	err = AssetFromJSON(args, &inp)
	if err != nil {
		return
	}

	// get, update and commit aggregatetuple
	aggregatetuple, err := db.GetAggregatetuple(inp.Key)
	if err != nil {
		return
	}
	aggregatetuple.Log += inp.Log

	if err = validateTupleOwner(db, aggregatetuple.Worker); err != nil {
		return
	}
	if err = aggregatetuple.commitStatusUpdate(db, inp.Key, StatusFailed); err != nil {
		return
	}

	outputAggregatetuple.Fill(db, aggregatetuple, inp.Key)

	// update depending tuples
	event := TuplesEvent{}
	err = UpdateTesttupleChildren(db, inp.Key, aggregatetuple.Status, &event)
	if err != nil {
		return
	}

	err = UpdateTraintupleChildren(db, inp.Key, outputAggregatetuple.Status, &event)
	if err != nil {
		return
	}

	err = SendTuplesEvent(db.cc, event)
	if err != nil {
		return
	}

	return
}

// logSuccessAggregateTrain modifies an aggregateTupl by changing its status from doing to done
// reports logs and associated performances
func logSuccessAggregateTrain(db LedgerDB, args []string) (outputAggregatetuple outputAggregatetuple, err error) {
	inp := inputLogSuccessTrain{}
	err = AssetFromJSON(args, &inp)
	if err != nil {
		return
	}
	aggregatetupleKey := inp.Key

	// get, update and commit aggregatetuple
	aggregatetuple, err := db.GetAggregatetuple(aggregatetupleKey)
	if err != nil {
		return
	}
	aggregatetuple.OutModel = &HashDress{
		Hash:           inp.OutModel.Hash,
		StorageAddress: inp.OutModel.StorageAddress}
	aggregatetuple.Log += inp.Log

	if err = validateTupleOwner(db, aggregatetuple.Worker); err != nil {
		return
	}
	if err = aggregatetuple.commitStatusUpdate(db, aggregatetupleKey, StatusDone); err != nil {
		return
	}

	event := TuplesEvent{}
	err = UpdateTraintupleChildren(db, aggregatetupleKey, aggregatetuple.Status, &event)
	if err != nil {
		return
	}

	err = UpdateTesttupleChildren(db, aggregatetupleKey, aggregatetuple.Status, &event)
	if err != nil {
		return
	}

	outputAggregatetuple.Fill(db, aggregatetuple, inp.Key)
	err = SendTuplesEvent(db.cc, event)
	if err != nil {
		return
	}

	return
}

// queryAggregatetuple returns info about an aggregate tuple given its key
func queryAggregatetuple(db LedgerDB, args []string) (outputAggregatetuple outputAggregatetuple, err error) {
	inp := inputHash{}
	err = AssetFromJSON(args, &inp)
	if err != nil {
		return
	}
	aggregatetuple, err := db.GetAggregatetuple(inp.Key)
	if err != nil {
		return
	}
	if aggregatetuple.AssetType != AggregatetupleType {
		err = errors.NotFound("no element with key %s", inp.Key)
		return
	}
	outputAggregatetuple.Fill(db, aggregatetuple, inp.Key)
	return
}

// queryAggregatetuples returns all aggregate tuples
func queryAggregatetuples(db LedgerDB, args []string) ([]outputAggregatetuple, error) {
	outputAggregatetuples := []outputAggregatetuple{}

	if len(args) != 0 {
		err := errors.BadRequest("incorrect number of arguments, expecting nothing")
		return outputAggregatetuples, err
	}
	elementsKeys, err := db.GetIndexKeys("aggregatetuple~algo~key", []string{"aggregatetuple"})
	if err != nil {
		return outputAggregatetuples, err
	}
	for _, key := range elementsKeys {
		outputAggregatetuple, err := getOutputAggregatetuple(db, key)
		if err != nil {
			return outputAggregatetuples, err
		}
		outputAggregatetuples = append(outputAggregatetuples, outputAggregatetuple)
	}
	return outputAggregatetuples, nil
}

// ----------------------------------------------------------
// Utils for smartcontracts related to aggregate tuples
// ----------------------------------------------------------

// getOutputAggregatetuple takes as input a aggregatetuple key and returns the outputAggregatetuple
func getOutputAggregatetuple(db LedgerDB, aggregatetupleKey string) (outAggreagateTuple outputAggregatetuple, err error) {
	aggregatetuple, err := db.GetAggregatetuple(aggregatetupleKey)
	if err != nil {
		return
	}
	outAggreagateTuple.Fill(db, aggregatetuple, aggregatetupleKey)
	return
}

// UpdateAggregatetupleChild updates the status of a waiting trainuple, given the new parent tuple status
func UpdateAggregatetupleChild(db LedgerDB, parentAggregatetupleKey string, childAggregatetupleKey string, aggregatetupleStatus string, event *TuplesEvent) (childStatus string, err error) {
	// get and update aggregatetuple
	childAggregatetuple, err := db.GetAggregatetuple(childAggregatetupleKey)
	if err != nil {
		return
	}

	childStatus = childAggregatetuple.Status

	// get traintuple new status
	var newStatus string
	if aggregatetupleStatus == StatusFailed {
		newStatus = StatusFailed
	} else if aggregatetupleStatus == StatusDone {
		ready, _err := childAggregatetuple.isReady(db, parentAggregatetupleKey)
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
	if err = childAggregatetuple.commitStatusUpdate(db, childAggregatetupleKey, newStatus); err != nil {
		return
	}

	// update return value after status update
	childStatus = childAggregatetuple.Status

	if newStatus == StatusTodo {
		out := outputAggregatetuple{}
		err = out.Fill(db, childAggregatetuple, childAggregatetupleKey)
		if err != nil {
			return
		}
		event.AddAggregatetuple(out)
	}

	return
}

func (tuple *Aggregatetuple) isReady(db LedgerDB, newDoneAggregatetupleKey string) (ready bool, err error) {
	return IsReady(db, tuple.InModelKeys, newDoneAggregatetupleKey)
}

// getOutputAggregatetuples takes as input a list of keys and returns a paylaod containing a list of associated retrieved elements
func getOutputAggregatetuples(db LedgerDB, aggregatetupleKeys []string) (outAggreagateTuples []outputAggregatetuple, err error) {
	for _, key := range aggregatetupleKeys {
		var outputAggregatetuple outputAggregatetuple
		outputAggregatetuple, err = getOutputAggregatetuple(db, key)
		if err != nil {
			return
		}
		outAggreagateTuples = append(outAggreagateTuples, outputAggregatetuple)
	}
	return
}

// commitStatusUpdate update the aggregatetuple status in the ledger
func (tuple *Aggregatetuple) commitStatusUpdate(db LedgerDB, aggregatetupleKey string, newStatus string) error {
	if tuple.Status == newStatus {
		return fmt.Errorf("cannot update aggregatetuple %s - status already %s", aggregatetupleKey, newStatus)
	}

	if err := tuple.validateNewStatus(db, newStatus); err != nil {
		return fmt.Errorf("update aggregatetuple %s failed: %s", aggregatetupleKey, err.Error())
	}

	oldStatus := tuple.Status
	tuple.Status = newStatus
	if err := db.Put(aggregatetupleKey, tuple); err != nil {
		return fmt.Errorf("failed to update aggregatetuple %s - %s", aggregatetupleKey, err.Error())
	}

	// update associated composite keys
	indexName := "aggregatetuple~worker~status~key"
	oldAttributes := []string{"aggregatetuple", tuple.Worker, oldStatus, aggregatetupleKey}
	newAttributes := []string{"aggregatetuple", tuple.Worker, tuple.Status, aggregatetupleKey}
	if err := db.UpdateIndex(indexName, oldAttributes, newAttributes); err != nil {
		return err
	}
	logger.Infof("aggregatetuple %s status updated: %s (from=%s)", aggregatetupleKey, newStatus, oldStatus)
	return nil
}

// validateNewStatus verifies that the new status is consistent with the tuple current status
func (tuple *Aggregatetuple) validateNewStatus(db LedgerDB, status string) error {
	// check validity of worker and change of status
	if err := checkUpdateTuple(db, tuple.Worker, tuple.Status, status); err != nil {
		return err
	}
	return nil
}
