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
func (tuple *Aggregatetuple) SetFromInput(db *LedgerDB, inp inputAggregatetuple) error {
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
	tuple.Worker = inp.Worker
	return nil
}

// SetFromParents set the status of the aggregate tuple depending on its "parents",
// i.e. the traintuples from which it received the outModels as inModels.
// Also it's InModelKeys are set.
func (tuple *Aggregatetuple) SetFromParents(db *LedgerDB, inModels []string) error {
	var parentStatuses []string
	inModelKeys := tuple.InModelKeys
	permissions, err := NewPermissions(db, OpenPermissions)
	if err != nil {
		return errors.BadRequest(err, "could not generate open permissions")
	}

	for _, parentTraintupleKey := range inModels {
		parentType, err := db.GetAssetType(parentTraintupleKey)
		if err != nil {
			return fmt.Errorf("could not retrieve traintuple type with key %s - %s", parentTraintupleKey, err.Error())
		}

		parentPermissions := Permissions{}

		// get out-model and permissions from parent
		switch parentType {
		case CompositeTraintupleType:
			tuple, err := db.GetCompositeTraintuple(parentTraintupleKey)
			if err == nil {
				// if the parent is composite, always take the "trunk" out-model
				parentPermissions = tuple.OutTrunkModel.Permissions
				parentStatuses = append(parentStatuses, tuple.Status)
			}
		case TraintupleType:
			tuple, err := db.GetTraintuple(parentTraintupleKey)
			if err == nil {
				parentPermissions = tuple.Permissions
				parentStatuses = append(parentStatuses, tuple.Status)
			}
		case AggregatetupleType:
			tuple, err := db.GetAggregatetuple(parentTraintupleKey)
			if err == nil {
				parentPermissions = tuple.Permissions
				parentStatuses = append(parentStatuses, tuple.Status)
			}
		default:
			return fmt.Errorf("aggregate.SetFromParents: Unsupported parent type %s", parentType)
		}

		if err != nil {
			return fmt.Errorf("could not retrieve traintuple type with key %s - %s", parentTraintupleKey, err.Error())
		}

		inModelKeys = append(inModelKeys, parentTraintupleKey)
		permissions = MergePermissions(permissions, parentPermissions)
	}
	tuple.Status = determineStatusFromInModels(parentStatuses)
	tuple.InModelKeys = inModelKeys
	tuple.Permissions = permissions
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
// Use checkComputePlanAvailability to ensure the compute plan exists and no other tuple is registered with the same worker/rank
func (tuple *Aggregatetuple) AddToComputePlan(db *LedgerDB, inp inputAggregatetuple, traintupleKey string, checkComputePlanAvailability bool) error {
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
	tuple.ComputePlanID = inp.ComputePlanID
	if !checkComputePlanAvailability {
		return nil
	}

	var ttKeys []string
	ttKeys, err = db.GetIndexKeys("computePlan~computeplanid~worker~rank~key", []string{"computePlan", inp.ComputePlanID})
	if err != nil {
		return err
	}
	if len(ttKeys) == 0 {
		return errors.BadRequest("cannot find the ComputePlanID %s", inp.ComputePlanID)
	}

	ttKeys, err = db.GetIndexKeys("computePlan~computeplanid~worker~rank~key", []string{"computePlan", inp.ComputePlanID, tuple.Worker, inp.Rank})
	if err != nil {
		return err
	} else if len(ttKeys) > 0 {
		err = errors.BadRequest("ComputePlanID %s with worker %s rank %d already exists", inp.ComputePlanID, tuple.Worker, tuple.Rank)
		return err
	}

	return nil
}

// Save will put in the legder interface both the aggregate tuple with its key
// and all the associated composite keys
func (tuple *Aggregatetuple) Save(db *LedgerDB, aggregatetupleKey string) error {

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
		if err := db.CreateIndex("computePlan~computeplanid~worker~rank~key", []string{"computePlan", tuple.ComputePlanID, tuple.Worker, strconv.Itoa(tuple.Rank), aggregatetupleKey}); err != nil {
			return err
		}
		if err := db.CreateIndex("computeplan~id", []string{"computeplan", tuple.ComputePlanID}); err != nil {
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
// createAggregatetuple is the wrapper for the substra smartcontract createAggregatetuple
func createAggregatetuple(db *LedgerDB, args []string) (map[string]string, error) {
	inp := inputAggregatetuple{}
	err := AssetFromJSON(args, &inp)
	if err != nil {
		return nil, err
	}

	key, err := createAggregatetupleInternal(db, inp, true)
	if err != nil {
		return nil, err
	}

	return map[string]string{"key": key}, nil
}

// createAggregatetupleInternal adds a Aggregatetuple in the ledger
func createAggregatetupleInternal(db *LedgerDB, inp inputAggregatetuple, checkComputePlanAvailability bool) (string, error) {

	aggregatetuple := Aggregatetuple{}
	err := aggregatetuple.SetFromInput(db, inp)
	if err != nil {
		return "", err
	}
	err = aggregatetuple.SetFromParents(db, inp.InModels)
	if err != nil {
		return "", err
	}

	aggregatetupleKey := aggregatetuple.GetKey()
	// Test if the key (ergo the aggregatetuple) already exists
	tupleExists, err := db.KeyExists(aggregatetupleKey)
	if err != nil {
		return "", err
	}
	if tupleExists {
		return "", errors.Conflict("aggregatetuple already exists").WithKey(aggregatetupleKey)
	}
	err = aggregatetuple.AddToComputePlan(db, inp, aggregatetupleKey, checkComputePlanAvailability)
	if err != nil {
		return "", err
	}
	err = aggregatetuple.Save(db, aggregatetupleKey)
	if err != nil {
		return "", err
	}

	err = db.AddTupleEvent(aggregatetupleKey)
	if err != nil {
		return "", err
	}

	return aggregatetupleKey, nil
}

// logStartAggregate modifies a aggregatetuple by changing its status from todo to doing
func logStartAggregate(db *LedgerDB, args []string) (o outputAggregatetuple, err error) {
	status := StatusDoing
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
	if err = aggregatetuple.commitStatusUpdate(db, inp.Key, status); err != nil {
		return
	}
	o.Fill(db, aggregatetuple, inp.Key)
	return
}

// logFailAggregate modifies a aggregatetuple by changing its status to fail and reports associated logs
func logFailAggregate(db *LedgerDB, args []string) (o outputAggregatetuple, err error) {
	status := StatusFailed
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
	if err = aggregatetuple.commitStatusUpdate(db, inp.Key, status); err != nil {
		return
	}

	o.Fill(db, aggregatetuple, inp.Key)

	// update depending tuples
	err = UpdateTesttupleChildren(db, inp.Key, aggregatetuple.Status)
	if err != nil {
		return
	}

	err = UpdateTraintupleChildren(db, inp.Key, o.Status)
	if err != nil {
		return
	}

	return
}

// logSuccessAggregate modifies an aggregateTupl by changing its status from doing to done
// reports logs and associated performances
func logSuccessAggregate(db *LedgerDB, args []string) (o outputAggregatetuple, err error) {
	status := StatusDone
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
	if err = aggregatetuple.commitStatusUpdate(db, aggregatetupleKey, status); err != nil {
		return
	}

	err = UpdateTraintupleChildren(db, aggregatetupleKey, aggregatetuple.Status)
	if err != nil {
		return
	}

	err = UpdateTesttupleChildren(db, aggregatetupleKey, aggregatetuple.Status)
	if err != nil {
		return
	}

	o.Fill(db, aggregatetuple, inp.Key)
	return
}

// queryAggregatetuple returns info about an aggregate tuple given its key
func queryAggregatetuple(db *LedgerDB, args []string) (outputAggregatetuple outputAggregatetuple, err error) {
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
func queryAggregatetuples(db *LedgerDB, args []string) ([]outputAggregatetuple, error) {
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
func getOutputAggregatetuple(db *LedgerDB, aggregatetupleKey string) (outAggreagateTuple outputAggregatetuple, err error) {
	aggregatetuple, err := db.GetAggregatetuple(aggregatetupleKey)
	if err != nil {
		return
	}
	outAggreagateTuple.Fill(db, aggregatetuple, aggregatetupleKey)
	return
}

// UpdateAggregatetupleChild updates the status of a waiting trainuple, given the new parent tuple status
func UpdateAggregatetupleChild(db *LedgerDB, parentAggregatetupleKey string, childAggregatetupleKey string, aggregatetupleStatus string) (childStatus string, err error) {
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

	err = db.AddTupleEvent(childAggregatetupleKey)
	return
}

func (tuple *Aggregatetuple) isReady(db *LedgerDB, newDoneAggregatetupleKey string) (ready bool, err error) {
	return IsReady(db, tuple.InModelKeys, newDoneAggregatetupleKey)
}

// getOutputAggregatetuples takes as input a list of keys and returns a paylaod containing a list of associated retrieved elements
func getOutputAggregatetuples(db *LedgerDB, aggregatetupleKeys []string) (outAggreagateTuples []outputAggregatetuple, err error) {
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
func (tuple *Aggregatetuple) commitStatusUpdate(db *LedgerDB, aggregatetupleKey string, newStatus string) error {
	if tuple.Status == newStatus {
		return nil
	}

	// do not update if previous status is already Done, Failed, Todo, Doing
	if StatusCanceled == newStatus && tuple.Status != StatusWaiting {
		return nil
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
func (tuple *Aggregatetuple) validateNewStatus(db *LedgerDB, status string) error {
	// check validity of worker and change of status
	if err := checkUpdateTuple(db, tuple.Worker, tuple.Status, status); err != nil {
		return err
	}
	return nil
}
