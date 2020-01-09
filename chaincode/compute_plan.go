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
	var computePlanID, tupleKey string
	traintupleKeysByID := map[string]string{}
	computePlan := ComputePlan{}
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
			inpTraintuple.ComputePlanID = computePlanID
			err = inpTraintuple.Fill(computeTraintuple, traintupleKeysByID)
			if err != nil {
				return resp, errors.BadRequest("traintuple ID %s: "+err.Error(), computeTraintuple.ID)
			}

			// Intentionally skip the compute plan availability check: since the transaction hasn't been
			// committed yet, the index changes haven't been commited, so the check would always fail.
			tupleKey, err = createTraintupleInternal(db, inpTraintuple, false)
			if err != nil {
				return resp, errors.BadRequest("traintuple ID %s: "+err.Error(), computeTraintuple.ID)
			}

			traintupleKeysByID[computeTraintuple.ID] = tupleKey
			computePlan.TraintupleKeys = append(computePlan.TraintupleKeys, tupleKey)
		case CompositeTraintupleType:
			computeCompositeTraintuple := inp.CompositeTraintuples[task.InputIndex]
			inpCompositeTraintuple := inputCompositeTraintuple{
				Rank: strconv.Itoa(task.Depth),
			}
			inpCompositeTraintuple.ComputePlanID = computePlanID
			err = inpCompositeTraintuple.Fill(computeCompositeTraintuple, traintupleKeysByID)
			if err != nil {
				return resp, errors.BadRequest("traintuple ID %s: "+err.Error(), computeCompositeTraintuple.ID)
			}
			// Intentionally skip the compute plan availability check: since the transaction hasn't been
			// committed yet, the index changes haven't been commited, so the check would always fail.
			tupleKey, err = createCompositeTraintupleInternal(db, inpCompositeTraintuple, false)
			if err != nil {
				return resp, errors.BadRequest("traintuple ID %s: "+err.Error(), computeCompositeTraintuple.ID)
			}

			traintupleKeysByID[computeCompositeTraintuple.ID] = tupleKey
			computePlan.CompositeTraintupleKeys = append(computePlan.CompositeTraintupleKeys, tupleKey)
		case AggregatetupleType:
			computeAggregatetuple := inp.Aggregatetuples[task.InputIndex]
			inpAggregatetuple := inputAggregatetuple{
				Rank: strconv.Itoa(task.Depth),
			}
			inpAggregatetuple.ComputePlanID = computePlanID
			err = inpAggregatetuple.Fill(computeAggregatetuple, traintupleKeysByID)
			if err != nil {
				return resp, errors.BadRequest("traintuple ID %s: "+err.Error(), computeAggregatetuple.ID)
			}
			// Intentionally skip the compute plan availability check: since the transaction hasn't been
			// committed yet, the index changes haven't been commited, so the check would always fail.
			tupleKey, err = createAggregatetupleInternal(db, inpAggregatetuple, false)
			if err != nil {
				return resp, errors.BadRequest("traintuple ID %s: "+err.Error(), computeAggregatetuple.ID)
			}

			traintupleKeysByID[computeAggregatetuple.ID] = tupleKey
			computePlan.AggregatetupleKeys = append(computePlan.AggregatetupleKeys, tupleKey)
		}
		if i == 0 {
			tuple, err := db.GetGenericTuple(tupleKey)
			if err != nil {
				return resp, err
			}
			computePlanID = tuple.ComputePlanID
		}

	}

	computePlan.TesttupleKeys = []string{}
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

		computePlan.TesttupleKeys = append(computePlan.TesttupleKeys, testtupleKey)
	}

	computePlan.Status = StatusTodo
	resp.Fill(computePlanID, computePlan)
	return resp, err
}

func queryComputePlan(db *LedgerDB, args []string) (resp outputComputePlan, err error) {
	inp := inputHash{}
	err = AssetFromJSON(args, &inp)
	if err != nil {
		return
	}
	return getOutComputePlan(db, inp.Key)
}

func queryComputePlans(db *LedgerDB, args []string) (resp []outputComputePlan, err error) {
	resp = []outputComputePlan{}
	computePlanIDs, err := db.GetIndexKeys("computePlan~id", []string{"computeplan"})
	if err != nil {
		return
	}
	for _, key := range computePlanIDs {
		var computePlan outputComputePlan
		computePlan, err = getOutComputePlan(db, key)
		if err != nil {
			return
		}
		resp = append(resp, computePlan)
	}
	return resp, err
}

// getComputePlan returns details for a compute plan id.
// Traintuples, CompositeTraintuples and Aggregatetuples are ordered by ascending rank.
func getOutComputePlan(db *LedgerDB, key string) (resp outputComputePlan, err error) {

	computePlan, err := db.GetComputePlan(key)
	if err != nil {
		return
	}

	resp.Fill(key, computePlan)
	return
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

	return "", fmt.Errorf("unknown compute plan status")
}

func cancelComputePlan(db *LedgerDB, args []string) (resp outputComputePlan, err error) {
	inp := inputHash{}
	err = AssetFromJSON(args, &inp)
	if err != nil {
		return
	}
	computeplan, err := db.GetComputePlan(inp.Key)
	if err != nil {
		return outputComputePlan{}, err
	}

	computeplan.Status = StatusCanceled
	err = db.Put(inp.Key, computeplan)
	if err != nil {
		return outputComputePlan{}, err
	}

	resp.Fill(inp.Key, computeplan)
	return resp, nil
}

func (cp *ComputePlan) Create(db *LedgerDB) (string, error) {
	ID := GetRandomHash()
	cp.AssetType = ComputePlanType
	err := db.Add(ID, cp)
	if err != nil {
		return "", err
	}
	if err := db.CreateIndex("computePlan~id", []string{"computePlan", ID}); err != nil {
		return "", err
	}
	return ID, nil
}

func (cp *ComputePlan) Save(db *LedgerDB, ID string) error {
	err := db.Put(ID, cp)
	if err != nil {
		return err
	}
	return nil
}
