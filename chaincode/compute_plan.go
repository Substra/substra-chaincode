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
	"strconv"
)

func (inpTraintuple *inputTraintuple) Fill(inpCP inputComputePlanTraintuple, IDToTrainTask map[string]TrainTask) error {
	inpTraintuple.Key = inpCP.Key
	inpTraintuple.DataManagerKey = inpCP.DataManagerKey
	inpTraintuple.DataSampleKeys = inpCP.DataSampleKeys
	inpTraintuple.AlgoKey = inpCP.AlgoKey
	inpTraintuple.Tag = inpCP.Tag
	inpTraintuple.Metadata = inpCP.Metadata

	// Set the inModels by matching the id to tuples key previously
	// encontered in this compute plan
	for _, InModelID := range inpCP.InModelsIDs {
		trainTask, ok := IDToTrainTask[InModelID]
		if !ok {
			return errors.BadRequest("model ID %s not found", InModelID)
		}
		inpTraintuple.InModels = append(inpTraintuple.InModels, trainTask.Key)
	}

	return nil

}
func (inpAggregatetuple *inputAggregatetuple) Fill(inpCP inputComputePlanAggregatetuple, IDToTrainTask map[string]TrainTask) error {
	inpAggregatetuple.Key = inpCP.Key
	inpAggregatetuple.AlgoKey = inpCP.AlgoKey
	inpAggregatetuple.Tag = inpCP.Tag
	inpAggregatetuple.Metadata = inpCP.Metadata
	inpAggregatetuple.Worker = inpCP.Worker

	// Set the inModels by matching the id to tuples key previously
	// encontered in this compute plan
	for _, InModelID := range inpCP.InModelsIDs {
		trainTask, ok := IDToTrainTask[InModelID]
		if !ok {
			return errors.BadRequest("model ID %s not found", InModelID)
		}
		inpAggregatetuple.InModels = append(inpAggregatetuple.InModels, trainTask.Key)
	}

	return nil

}
func (inpCompositeTraintuple *inputCompositeTraintuple) Fill(inpCP inputComputePlanCompositeTraintuple, IDToTrainTask map[string]TrainTask) error {
	inpCompositeTraintuple.Key = inpCP.Key
	inpCompositeTraintuple.DataManagerKey = inpCP.DataManagerKey
	inpCompositeTraintuple.DataSampleKeys = inpCP.DataSampleKeys
	inpCompositeTraintuple.AlgoKey = inpCP.AlgoKey
	inpCompositeTraintuple.Tag = inpCP.Tag
	inpCompositeTraintuple.Metadata = inpCP.Metadata
	inpCompositeTraintuple.OutTrunkModelPermissions = inpCP.OutTrunkModelPermissions

	// Set the inModels by matching the id to traintuples key previously
	// encontered in this compute plan
	if inpCP.InHeadModelID != "" {
		var ok bool
		trainTask, ok := IDToTrainTask[inpCP.InHeadModelID]
		if !ok {
			return errors.BadRequest("head model ID %s not found", inpCP.InHeadModelID)
		}
		inpCompositeTraintuple.InHeadModelKey = trainTask.Key
	}
	if inpCP.InTrunkModelID != "" {
		var ok bool
		trainTask, ok := IDToTrainTask[inpCP.InTrunkModelID]
		if !ok {
			return errors.BadRequest("trunk model ID %s not found", inpCP.InTrunkModelID)
		}
		inpCompositeTraintuple.InTrunkModelKey = trainTask.Key
	}
	return nil
}

func (inpTesttuple *inputTesttuple) Fill(inpCP inputComputePlanTesttuple, IDToTrainTask map[string]TrainTask) error {
	trainTask, ok := IDToTrainTask[inpCP.TraintupleID]
	if !ok {
		return errors.BadRequest("traintuple ID %s not found", inpCP.TraintupleID)
	}
	inpTesttuple.Key = inpCP.Key
	inpTesttuple.TraintupleKey = trainTask.Key
	inpTesttuple.DataManagerKey = inpCP.DataManagerKey
	inpTesttuple.DataSampleKeys = inpCP.DataSampleKeys
	inpTesttuple.Tag = inpCP.Tag
	inpTesttuple.Metadata = inpCP.Metadata
	inpTesttuple.ObjectiveKey = inpCP.ObjectiveKey

	return nil
}

// createComputePlan is the wrapper for the substra smartcontract CreateComputePlan
func createComputePlan(db *LedgerDB, args []string) (resp outputComputePlan, err error) {
	inp := inputNewComputePlan{}
	err = AssetFromJSON(args, &inp)
	if err != nil {
		return
	}
	return createComputePlanInternal(db, inp.inputComputePlan, inp.Tag, inp.Metadata, inp.CleanModels)
}

func updateComputePlan(db *LedgerDB, args []string) (resp outputComputePlan, err error) {
	inp := inputComputePlan{}
	err = AssetFromJSON(args, &inp)
	if err != nil {
		return
	}

	count := len(inp.Traintuples) +
		len(inp.Aggregatetuples) +
		len(inp.CompositeTraintuples) +
		len(inp.Testtuples)
	if count == 0 {
		return resp, errors.BadRequest("empty update for compute plan %s", inp.Key)
	}
	return updateComputePlanInternal(db, inp)
}

func createComputePlanInternal(db *LedgerDB, inp inputComputePlan, tag string, metadata map[string]string, cleanModels bool) (resp outputComputePlan, err error) {
	var computePlan ComputePlan
	computePlan.State.Status = StatusWaiting
	computePlan.Tag = tag
	computePlan.Metadata = metadata
	computePlan.CleanModels = cleanModels
	err = computePlan.Create(db, inp.Key)
	if err != nil {
		return resp, err
	}
	count := len(inp.Traintuples) +
		len(inp.Aggregatetuples) +
		len(inp.CompositeTraintuples) +
		len(inp.Testtuples)
	if count == 0 {
		resp.Fill(inp.Key, computePlan, []string{}, 0, 0)
		return resp, nil
	}
	return updateComputePlanInternal(db, inp)
}

func updateComputePlanInternal(db *LedgerDB, inp inputComputePlan) (resp outputComputePlan, err error) {
	var tupleKey string
	computePlan, err := db.GetComputePlan(inp.Key)
	if err != nil {
		return resp, err
	}
	IDToTrainTask := map[string]TrainTask{}
	for ID, trainTask := range computePlan.IDToTrainTask {
		IDToTrainTask[ID] = trainTask
	}
	NewIDs := []string{}
	DAG, err := createComputeDAG(inp, computePlan.IDToTrainTask)
	if err != nil {
		return resp, errors.BadRequest(err)
	}
	for _, task := range DAG.OrderTasks {
		switch task.TaskType {
		case TraintupleType:
			computeTraintuple := inp.Traintuples[task.InputIndex]
			inpTraintuple := inputTraintuple{
				Rank: strconv.Itoa(task.Depth),
			}
			inpTraintuple.ComputePlanKey = inp.Key
			err = inpTraintuple.Fill(computeTraintuple, IDToTrainTask)
			if err != nil {
				return resp, errors.BadRequest("traintuple ID %s: "+err.Error(), computeTraintuple.ID)
			}

			// Intentionally skip the compute plan availability check: since the transaction hasn't been
			// committed yet, the index changes haven't been commited, so the check would always fail.
			tupleKey, err = createTraintupleInternal(db, inpTraintuple, false)
			if err != nil {
				return resp, errors.BadRequest("traintuple ID %s: "+err.Error(), computeTraintuple.ID)
			}
		case CompositeTraintupleType:
			computeCompositeTraintuple := inp.CompositeTraintuples[task.InputIndex]
			inpCompositeTraintuple := inputCompositeTraintuple{
				Rank: strconv.Itoa(task.Depth),
			}
			inpCompositeTraintuple.ComputePlanKey = inp.Key
			err = inpCompositeTraintuple.Fill(computeCompositeTraintuple, IDToTrainTask)
			if err != nil {
				return resp, errors.BadRequest("traintuple ID %s: "+err.Error(), computeCompositeTraintuple.ID)
			}
			// Intentionally skip the compute plan availability check: since the transaction hasn't been
			// committed yet, the index changes haven't been commited, so the check would always fail.
			tupleKey, err = createCompositeTraintupleInternal(db, inpCompositeTraintuple, false)
			if err != nil {
				return resp, errors.BadRequest("traintuple ID %s: "+err.Error(), computeCompositeTraintuple.ID)
			}
		case AggregatetupleType:
			computeAggregatetuple := inp.Aggregatetuples[task.InputIndex]
			inpAggregatetuple := inputAggregatetuple{
				Rank: strconv.Itoa(task.Depth),
			}
			inpAggregatetuple.ComputePlanKey = inp.Key
			err = inpAggregatetuple.Fill(computeAggregatetuple, IDToTrainTask)
			if err != nil {
				return resp, errors.BadRequest("traintuple ID %s: "+err.Error(), computeAggregatetuple.ID)
			}
			// Intentionally skip the compute plan availability check: since the transaction hasn't been
			// committed yet, the index changes haven't been commited, so the check would always fail.
			tupleKey, err = createAggregatetupleInternal(db, inpAggregatetuple, false)
			if err != nil {
				return resp, errors.BadRequest("traintuple ID %s: "+err.Error(), computeAggregatetuple.ID)
			}
		}
		IDToTrainTask[task.ID] = TrainTask{Depth: task.Depth, Key: tupleKey}
		NewIDs = append(NewIDs, task.ID)
	}

	for index, computeTesttuple := range inp.Testtuples {
		inpTesttuple := inputTesttuple{}
		err = inpTesttuple.Fill(computeTesttuple, IDToTrainTask)
		if err != nil {
			return resp, errors.BadRequest("testtuple at index %s: "+err.Error(), index)
		}

		_, err := createTesttupleInternal(db, inpTesttuple)
		if err != nil {
			return resp, errors.BadRequest("testtuple at index %s: "+err.Error(), index)
		}

	}

	computePlan, err = db.GetComputePlan(inp.Key)
	if err != nil {
		return resp, err
	}
	computePlan.IDToTrainTask = IDToTrainTask
	err = computePlan.Save(db, inp.Key)
	if err != nil {
		return resp, err
	}
	doneCount, tupleCount, err := computePlan.getTupleCounts(db)
	if err != nil {
		return resp, err
	}

	resp.Fill(inp.Key, computePlan, NewIDs, doneCount, tupleCount)
	return resp, err
}

func queryComputePlan(db *LedgerDB, args []string) (resp outputComputePlan, err error) {
	inp := inputKey{}
	err = AssetFromJSON(args, &inp)
	if err != nil {
		return
	}
	return getOutComputePlan(db, inp.Key)
}

func queryComputePlans(db *LedgerDB, args []string) (outComputePlans []outputComputePlan, bookmark string, err error) {
	inp := inputBookmark{}
	outComputePlans = []outputComputePlan{}

	if len(args) > 1 {
		err = errors.BadRequest("incorrect number of arguments, expecting at most one argument")
		return
	}

	if len(args) == 1 && args[0] != "" {
		err = AssetFromJSON(args, &inp)
		if err != nil {
			return
		}
	}

	computePlanKeys, bookmark, err := db.GetIndexKeysWithPagination("computePlan~key", []string{"computePlan"}, OutputPageSize, inp.Bookmark)

	if err != nil {
		return
	}

	for _, key := range computePlanKeys {
		var computePlan outputComputePlan
		computePlan, err = getOutComputePlan(db, key)
		if err != nil {
			return
		}
		outComputePlans = append(outComputePlans, computePlan)
	}
	return
}

// getComputePlan returns details for a compute plan key.
// Traintuples, CompositeTraintuples and Aggregatetuples are ordered by ascending rank.
func getOutComputePlan(db *LedgerDB, key string) (resp outputComputePlan, err error) {

	computePlan, err := db.GetComputePlan(key)
	if err != nil {
		return resp, err
	}

	doneCount, tupleCount, err := computePlan.getTupleCounts(db)
	if err != nil {
		return resp, err
	}

	resp.Fill(key, computePlan, []string{}, doneCount, tupleCount)
	return resp, err
}

func cancelComputePlan(db *LedgerDB, args []string) (resp outputComputePlan, err error) {
	inp := inputKey{}
	err = AssetFromJSON(args, &inp)
	if err != nil {
		return
	}
	computeplan, err := db.GetComputePlan(inp.Key)
	if err != nil {
		return outputComputePlan{}, err
	}

	computeplan.State.Status = StatusCanceled
	err = computeplan.SaveState(db)
	if err != nil {
		return outputComputePlan{}, err
	}

	models, err := computeplan.removeAllIntermediaryModels(db)
	if err != nil {
		return outputComputePlan{}, err
	}

	err = db.AddComputePlanEvent(inp.Key, computeplan.State.Status, models)
	if err != nil {
		return outputComputePlan{}, err
	}

	doneCount, tupleCount, err := computeplan.getTupleCounts(db)
	if err != nil {
		return resp, err
	}
	resp.Fill(inp.Key, computeplan, []string{}, doneCount, tupleCount)
	return resp, nil
}

// Create adds a Compute Plan to the ledger and registers it in the compute plan index
func (cp *ComputePlan) Create(db *LedgerDB, key string) error {
	cp.Key = key
	cp.StateKey = GetRandomHash()
	cp.AssetType = ComputePlanType
	cp.Workers = []string{}
	err := db.Add(key, cp)
	if err != nil {
		return err
	}
	err = db.Add(cp.StateKey, cp.State)
	if err != nil {
		return err
	}
	if err := db.CreateIndex("computePlan~key", []string{"computePlan", key}); err != nil {
		return err
	}
	return nil
}

// Save add or update the compute plan in the ledger
func (cp *ComputePlan) Save(db *LedgerDB, key string) error {
	err := db.Put(key, cp)
	if err != nil {
		return err
	}
	return cp.SaveState(db)
}

// SaveState add or update the compute plan in the ledger
func (cp *ComputePlan) SaveState(db *LedgerDB) error {
	return db.Put(cp.StateKey, cp.State)
}

// UpdateState check the tuple status (from an updated tuple or a new one)
// and, if required, it updates the compute plan' status and/or its doneCount.
// It returns true and the list of models to delete if there is any change to the compute plan, false and empty list otherwise.
func (cp *ComputePlan) UpdateState(db *LedgerDB, tupleStatus string, worker string) (bool, []string, error) {
	switch cp.State.Status {
	case StatusFailed, StatusCanceled:
	case StatusDone:
		// We might add tuples to a done compute plan
		if stringInSlice(tupleStatus, []string{StatusWaiting, StatusTodo}) {
			cp.State.Status = tupleStatus
			return true, []string{}, nil
		}
	case StatusDoing:
		switch tupleStatus {
		case StatusFailed:
			cp.State.Status = tupleStatus
			return true, []string{}, nil
		case StatusDone:
			// In order for the CP to transition to the "done" state, each worker must have all
			// its tuples in the "done" state. Checking the state of all the workers is
			// expensive and can lead to MVCC conflicts because workers each write to their
			// respective states concurrently. To mitigate this issue, we first check if the
			//  *current* worker has finished processing all of its tuples. Only if that's the
			// case do we inspect the state of other workers.
			cp.incrementWorkerDoneCount(db, worker)
			wStateKey := cp.getCPWorkerStateKey(worker)
			wState, err := db.GetCPWorkerState(wStateKey)
			if err != nil {
				return false, []string{}, err
			}
			if wState.DoneCount == wState.TupleCount {
				doneCount, tupleCount, err := cp.getTupleCounts(db)
				if err != nil {
					return false, []string{}, err
				}
				if doneCount == tupleCount {
					modelsToDelete, err := cp.removeAllIntermediaryModels(db)
					if err != nil {
						return false, []string{}, err
					}
					cp.State.Status = StatusDone
					return true, modelsToDelete, nil
				}
			}
			return false, []string{}, nil
		}
	case StatusTodo:
		if tupleStatus == StatusDoing {
			cp.State.Status = tupleStatus
			return true, []string{}, nil
		}
	case StatusWaiting:
		if tupleStatus == StatusTodo {
			cp.State.Status = tupleStatus
			return true, []string{}, nil
		}
	case "":
		cp.State.Status = tupleStatus
		return true, []string{}, nil
	}
	return false, []string{}, nil
}

// AddTuple add the tuple key to the compute plan and update it accordingly
func (cp *ComputePlan) AddTuple(db *LedgerDB, tupleType AssetType, key, status string, worker string) error {
	switch tupleType {
	case TraintupleType:
		cp.TraintupleKeys = append(cp.TraintupleKeys, key)
	case CompositeTraintupleType:
		cp.CompositeTraintupleKeys = append(cp.CompositeTraintupleKeys, key)
	case AggregatetupleType:
		cp.AggregatetupleKeys = append(cp.AggregatetupleKeys, key)
	case TesttupleType:
		cp.TesttupleKeys = append(cp.TesttupleKeys, key)
	}
	cp.incrementWorkerTupleCount(db, worker)
	_, _, err := cp.UpdateState(db, status, worker)
	return err
}

// UpdateComputePlanState retreive the compute plan if the ID is not empty,
// check if the updated status change anything and save it if it's the case
func UpdateComputePlanState(db *LedgerDB, ComputePlanKey, tupleStatus, tupleKey string, worker string) error {
	if ComputePlanKey == "" {
		return nil
	}
	cp, err := db.GetComputePlan(ComputePlanKey)
	if err != nil {
		return err
	}
	stateUpdated, modelsToDelete, err := cp.UpdateState(db, tupleStatus, worker)
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}
	if stateUpdated || len(modelsToDelete) > 0 {
		err = db.AddComputePlanEvent(ComputePlanKey, cp.State.Status, modelsToDelete)
		if err != nil {
			return err
		}
		return cp.SaveState(db)
	}
	return nil
}
