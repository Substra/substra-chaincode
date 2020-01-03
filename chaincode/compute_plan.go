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
	"sort"
	"strconv"
)

func (inpTraintuple *inputTraintuple) Fill(inpCP inputComputePlanTraintuple, traintupleKeysByID map[string]string) error {
	inpTraintuple.DataManagerKey = inpCP.DataManagerKey
	inpTraintuple.DataSampleKeys = inpCP.DataSampleKeys
	inpTraintuple.AlgoKey = inpCP.AlgoKey
	inpTraintuple.Tag = inpCP.Tag

	// Set the inModels by matching the id to tuples key previously
	// encontered in this compute plan
	for _, InModelID := range inpCP.InModelsIDs {
		inModelKey, ok := traintupleKeysByID[InModelID]
		if !ok {
			return fmt.Errorf("model ID %s not found", InModelID)
		}
		inpTraintuple.InModels = append(inpTraintuple.InModels, inModelKey)
	}

	return nil

}
func (inpAggregatetuple *inputAggregatetuple) Fill(inpCP inputComputePlanAggregatetuple, aggregatetupleKeysByID map[string]string) error {
	inpAggregatetuple.AlgoKey = inpCP.AlgoKey
	inpAggregatetuple.Tag = inpCP.Tag
	inpAggregatetuple.Worker = inpCP.Worker

	// Set the inModels by matching the id to tuples key previously
	// encontered in this compute plan
	for _, InModelID := range inpCP.InModelsIDs {
		inModelKey, ok := aggregatetupleKeysByID[InModelID]
		if !ok {
			return fmt.Errorf("model ID %s not found", InModelID)
		}
		inpAggregatetuple.InModels = append(inpAggregatetuple.InModels, inModelKey)
	}

	return nil

}
func (inpCompositeTraintuple *inputCompositeTraintuple) Fill(inpCP inputComputePlanCompositeTraintuple, traintupleKeysByID map[string]string) error {
	inpCompositeTraintuple.DataManagerKey = inpCP.DataManagerKey
	inpCompositeTraintuple.DataSampleKeys = inpCP.DataSampleKeys
	inpCompositeTraintuple.AlgoKey = inpCP.AlgoKey
	inpCompositeTraintuple.Tag = inpCP.Tag
	inpCompositeTraintuple.OutTrunkModelPermissions = inpCP.OutTrunkModelPermissions

	// Set the inModels by matching the id to traintuples key previously
	// encontered in this compute plan
	if inpCP.InHeadModelID != "" {
		var ok bool
		inpCompositeTraintuple.InHeadModelKey, ok = traintupleKeysByID[inpCP.InHeadModelID]
		if !ok {
			return fmt.Errorf("head model ID %s not found", inpCP.InHeadModelID)
		}
	}
	if inpCP.InTrunkModelID != "" {
		var ok bool
		inpCompositeTraintuple.InTrunkModelKey, ok = traintupleKeysByID[inpCP.InTrunkModelID]
		if !ok {
			return fmt.Errorf("trunk model ID %s not found", inpCP.InTrunkModelID)
		}
	}
	return nil
}

func (inpTesttuple *inputTesttuple) Fill(inpCP inputComputePlanTesttuple, traintupleKeysByID map[string]string) error {
	traintupleKey, ok := traintupleKeysByID[inpCP.TraintupleID]
	if !ok {
		return fmt.Errorf("traintuple ID %s not found", inpCP.TraintupleID)
	}
	inpTesttuple.TraintupleKey = traintupleKey
	inpTesttuple.DataManagerKey = inpCP.DataManagerKey
	inpTesttuple.DataSampleKeys = inpCP.DataSampleKeys
	inpTesttuple.Tag = inpCP.Tag
	inpTesttuple.ObjectiveKey = inpCP.ObjectiveKey

	return nil
}

// createComputePlan is the wrapper for the substra smartcontract CreateComputePlan
func createComputePlan(db *LedgerDB, args []string) (resp outputComputePlan, err error) {
	inp := inputComputePlan{}
	err = AssetFromJSON(args, &inp)
	if err != nil {
		return
	}
	return createComputePlanInternal(db, inp)
}

func createComputePlanInternal(db *LedgerDB, inp inputComputePlan) (resp outputComputePlan, err error) {
	traintupleKeysByID := map[string]string{}

	resp.TraintupleKeys = []string{}

	DAG, err := createComputeDAG(inp)
	if err != nil {
		return resp, errors.BadRequest(err)
	}
	for i, task := range DAG.OrderTasks {
		switch task.TaskType {
		case TraintupleType:
			computeTraintuple := inp.Traintuples[task.InputIndex]
			inpTraintuple := inputTraintuple{
				Rank: strconv.Itoa(task.Depth),
			}
			if i != 0 {
				inpTraintuple.ComputePlanID = resp.ComputePlanID
			}
			err = inpTraintuple.Fill(computeTraintuple, traintupleKeysByID)
			if err != nil {
				return resp, errors.BadRequest("traintuple ID %s: "+err.Error(), computeTraintuple.ID)
			}

			// Intentionally skip the compute plan availability check: since the transaction hasn't been
			// committed yet, the index changes haven't been commited, so the check would always fail.
			traintupleKey, err := createTraintupleInternal(db, inpTraintuple, false)
			if err != nil {
				return resp, errors.BadRequest("traintuple ID %s: "+err.Error(), computeTraintuple.ID)
			}

			if i == 0 {
				resp.ComputePlanID = traintupleKey
			}

			traintupleKeysByID[computeTraintuple.ID] = traintupleKey
			resp.TraintupleKeys = append(resp.TraintupleKeys, traintupleKey)
		case CompositeTraintupleType:
			computeCompositeTraintuple := inp.CompositeTraintuples[task.InputIndex]
			inpCompositeTraintuple := inputCompositeTraintuple{
				Rank: strconv.Itoa(task.Depth),
			}
			if i != 0 {
				inpCompositeTraintuple.ComputePlanID = resp.ComputePlanID
			}
			err = inpCompositeTraintuple.Fill(computeCompositeTraintuple, traintupleKeysByID)
			if err != nil {
				return resp, errors.BadRequest("traintuple ID %s: "+err.Error(), computeCompositeTraintuple.ID)
			}
			_ = computeCompositeTraintuple
			// Intentionally skip the compute plan availability check: since the transaction hasn't been
			// committed yet, the index changes haven't been commited, so the check would always fail.
			compositeTraintupleKey, err := createCompositeTraintupleInternal(db, inpCompositeTraintuple, false)
			if err != nil {
				return resp, errors.BadRequest("traintuple ID %s: "+err.Error(), computeCompositeTraintuple.ID)
			}

			if i == 0 {
				resp.ComputePlanID = compositeTraintupleKey
			}

			traintupleKeysByID[computeCompositeTraintuple.ID] = compositeTraintupleKey
			resp.CompositeTraintupleKeys = append(resp.CompositeTraintupleKeys, compositeTraintupleKey)
		case AggregatetupleType:
			computeAggregatetuple := inp.Aggregatetuples[task.InputIndex]
			inpAggregatetuple := inputAggregatetuple{
				Rank: strconv.Itoa(task.Depth),
			}
			if i != 0 {
				inpAggregatetuple.ComputePlanID = resp.ComputePlanID
			}
			err = inpAggregatetuple.Fill(computeAggregatetuple, traintupleKeysByID)
			if err != nil {
				return resp, errors.BadRequest("traintuple ID %s: "+err.Error(), computeAggregatetuple.ID)
			}
			_ = computeAggregatetuple
			// Intentionally skip the compute plan availability check: since the transaction hasn't been
			// committed yet, the index changes haven't been commited, so the check would always fail.
			aggregatetupleKey, err := createAggregatetupleInternal(db, inpAggregatetuple, false)
			if err != nil {
				return resp, errors.BadRequest("traintuple ID %s: "+err.Error(), computeAggregatetuple.ID)
			}

			if i == 0 {
				resp.ComputePlanID = aggregatetupleKey
			}

			traintupleKeysByID[computeAggregatetuple.ID] = aggregatetupleKey
			resp.AggregatetupleKeys = append(resp.AggregatetupleKeys, aggregatetupleKey)
		}

	}

	resp.TesttupleKeys = []string{}
	for index, computeTesttuple := range inp.Testtuples {
		inpTesttuple := inputTesttuple{}
		err = inpTesttuple.Fill(computeTesttuple, traintupleKeysByID)
		if err != nil {
			return resp, errors.BadRequest("testtuple at index %s: "+err.Error(), index)
		}

		testtupleKey, err := createTesttupleInternal(db, inpTesttuple)
		if err != nil {
			return resp, errors.BadRequest("testtuple at index %s: "+err.Error(), index)
		}

		resp.TesttupleKeys = append(resp.TesttupleKeys, testtupleKey)
	}

	resp.Status, err = getComputePlanStatus(db, resp)
	return resp, err
}

func queryComputePlan(db *LedgerDB, args []string) (resp outputComputePlan, err error) {
	inp := inputHash{}
	err = AssetFromJSON(args, &inp)
	if err != nil {
		return
	}
	return getComputePlan(db, inp.Key)
}

func queryComputePlans(db *LedgerDB, args []string) (resp []outputComputePlan, err error) {
	resp = []outputComputePlan{}
	computePlanIDs, err := db.GetIndexKeys("computeplan~id", []string{"computeplan"})
	if err != nil {
		return
	}
	for _, key := range computePlanIDs {
		var computePlan outputComputePlan
		computePlan, err = getComputePlan(db, key)
		if err != nil {
			return
		}
		resp = append(resp, computePlan)
	}
	return resp, err
}

// getComputePlan returns details for a compute plan id.
// Traintuples, CompositeTraintuples and Aggregatetuples are ordered by ascending rank.
func getComputePlan(db *LedgerDB, key string) (resp outputComputePlan, err error) {

	// 1. Get tuples (regular, composite, aggregate)
	tupleKeys, err := db.GetIndexKeys("computePlan~computeplanid~worker~rank~key", []string{"computePlan", key})
	if err != nil {
		return
	}
	if len(tupleKeys) == 0 {
		err = errors.E("No traintuple found for compute plan %s", key)
		return
	}
	tuples := map[string]GenericTuple{}
	for _, tupleKey := range tupleKeys {
		var tuple GenericTuple
		tuple, err = db.GetGenericTuple(tupleKey)
		if err != nil {
			return
		}
		tuples[tupleKey] = tuple
	}

	// 2. Sort tuples by ascending rank
	sort.SliceStable(tupleKeys, func(i, j int) bool {
		return tuples[tupleKeys[i]].Rank < tuples[tupleKeys[j]].Rank
	})

	// 3. Get Testtuples associated with each tuple
	testtupleKeys := []string{}
	for _, tupleKey := range tupleKeys {
		var toAdd []string
		toAdd, err = db.GetIndexKeys("testtuple~traintuple~certified~key", []string{"testtuple", tupleKey})
		if err != nil {
			return
		}
		testtupleKeys = append(testtupleKeys, toAdd...)
	}

	// 4. Split tuple keys depending on their type
	traintupleKeys := []string{}
	compositeTraintupleKeys := []string{}
	aggregatetupleKeys := []string{}
	for _, tupleKey := range tupleKeys { // iterate over keys (sorted by rank) so that each output array is also sorted by rank
		tuple := tuples[tupleKey]
		switch tuple.AssetType {
		case TraintupleType:
			traintupleKeys = append(traintupleKeys, tupleKey)
		case CompositeTraintupleType:
			compositeTraintupleKeys = append(compositeTraintupleKeys, tupleKey)
		case AggregatetupleType:
			aggregatetupleKeys = append(aggregatetupleKeys, tupleKey)
		default:
			err = fmt.Errorf("Unknown tuple type: %v", tuple.AssetType)
			return
		}
	}

	resp = outputComputePlan{
		ComputePlanID:           key,
		TraintupleKeys:          traintupleKeys,
		CompositeTraintupleKeys: compositeTraintupleKeys,
		AggregatetupleKeys:      aggregatetupleKeys,
		TesttupleKeys:           testtupleKeys,
	}

	resp.Status, err = getComputePlanStatus(db, resp)
	if err != nil {
		return
	}

	return
}

func getComputePlanStatusByComputePlanID(db *LedgerDB, computePlanID string) (status string, err error) {
	computePlan, err := getComputePlan(db, computePlanID)
	if err != nil {
		return "", err
	}

	return getComputePlanStatus(db, computePlan)
}

func getComputePlanStatus(db *LedgerDB, computePlan outputComputePlan) (status string, err error) {
	// get all tuples like status in one slice
	statusCollection := []string{}

	keys := []string{}
	keys = append(keys, computePlan.TraintupleKeys...)
	keys = append(keys, computePlan.CompositeTraintupleKeys...)
	keys = append(keys, computePlan.AggregatetupleKeys...)
	keys = append(keys, computePlan.TesttupleKeys...)

	for _, key := range keys {
		tuple, err := db.GetGenericTuple(key)
		if err != nil {
			return "", err
		}

		statusCollection = append(statusCollection, tuple.Status)
	}

	return determineComputePlanStatus(statusCollection)
}

func determineComputePlanStatus(statusCollection []string) (status string, err error) {
	// this status order matters
	sts := []string{StatusCanceled, StatusFailed, StatusDoing, StatusTodo, StatusWaiting, StatusDone}
	for _, s := range sts {
		if stringInSlice(s, statusCollection) {
			return s, nil
		}
	}

	return StatusUndefined, nil
}

func cancelComputePlan(db *LedgerDB, args []string) (resp outputComputePlan, err error) {
	computeplan, err := queryComputePlan(db, args)
	if err != nil {
		return outputComputePlan{}, err
	}

	if stringInSlice(computeplan.Status, []string{StatusCanceled, StatusDone}) {
		return computeplan, nil
	}

	var tupleKeys []string
	tupleKeys = append(tupleKeys, computeplan.TraintupleKeys...)
	tupleKeys = append(tupleKeys, computeplan.CompositeTraintupleKeys...)
	tupleKeys = append(tupleKeys, computeplan.AggregatetupleKeys...)
	tupleKeys = append(tupleKeys, computeplan.TesttupleKeys...)

	for _, key := range tupleKeys {

		tuple, err := db.GetStatusUpdater(key)
		if err != nil {
			return outputComputePlan{}, err
		}
		err = tuple.commitStatusUpdate(db, key, StatusCanceled)
		if err != nil {
			return outputComputePlan{}, err
		}
	}

	computeplan.Status = StatusCanceled

	return computeplan, nil
}
