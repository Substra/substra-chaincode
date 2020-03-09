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
	inpTraintuple.DataManagerKey = inpCP.DataManagerKey
	inpTraintuple.DataSampleKeys = inpCP.DataSampleKeys
	inpTraintuple.AlgoKey = inpCP.AlgoKey
	inpTraintuple.Tag = inpCP.Tag

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
	inpAggregatetuple.AlgoKey = inpCP.AlgoKey
	inpAggregatetuple.Tag = inpCP.Tag
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
	inpCompositeTraintuple.DataManagerKey = inpCP.DataManagerKey
	inpCompositeTraintuple.DataSampleKeys = inpCP.DataSampleKeys
	inpCompositeTraintuple.AlgoKey = inpCP.AlgoKey
	inpCompositeTraintuple.Tag = inpCP.Tag
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
	inpTesttuple.TraintupleKey = trainTask.Key
	inpTesttuple.DataManagerKey = inpCP.DataManagerKey
	inpTesttuple.DataSampleKeys = inpCP.DataSampleKeys
	inpTesttuple.Tag = inpCP.Tag
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
	return createComputePlanInternal(db, inp.inputComputePlan, inp.Tag)
}

func updateComputePlan(db *LedgerDB, args []string) (resp outputComputePlan, err error) {
	inp := inputUpdateComputePlan{}
	err = AssetFromJSON(args, &inp)
	if err != nil {
		return
	}

	count := len(inp.Traintuples) +
		len(inp.Aggregatetuples) +
		len(inp.CompositeTraintuples) +
		len(inp.Testtuples)
	if count == 0 {
		return resp, errors.BadRequest("empty update for compute plan %s", inp.ComputePlanID)
	}
	return updateComputePlanInternal(db, inp.ComputePlanID, inp.inputComputePlan)
}

func createComputePlanInternal(db *LedgerDB, inp inputComputePlan, tag string) (resp outputComputePlan, err error) {
	var computePlan ComputePlan
	computePlan.State.Status = StatusWaiting
	computePlan.Tag = tag
	ID, err := computePlan.Create(db)
	if err != nil {
		return resp, err
	}
	count := len(inp.Traintuples) +
		len(inp.Aggregatetuples) +
		len(inp.CompositeTraintuples) +
		len(inp.Testtuples)
	if count == 0 {
		resp.Fill(ID, computePlan, []string{})
		return resp, nil
	}
	return updateComputePlanInternal(db, ID, inp)
}

func updateComputePlanInternal(db *LedgerDB, computePlanID string, inp inputComputePlan) (resp outputComputePlan, err error) {
	var tupleKey string
	computePlan, err := db.GetComputePlan(computePlanID)
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
			inpTraintuple.ComputePlanID = computePlanID
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
			inpCompositeTraintuple.ComputePlanID = computePlanID
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
			inpAggregatetuple.ComputePlanID = computePlanID
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

	computePlan, err = db.GetComputePlan(computePlanID)
	if err != nil {
		return resp, err
	}
	computePlan.IDToTrainTask = IDToTrainTask
	err = computePlan.Save(db, computePlanID)
	if err != nil {
		return resp, err
	}
	resp.Fill(computePlanID, computePlan, NewIDs)
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

func queryComputePlans(db *LedgerDB, args []string) (resp []outputComputePlan, err error) {
	resp = []outputComputePlan{}
	computePlanIDs, err := db.GetIndexKeys("computePlan~id", []string{"computePlan"})
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
		return resp, err
	}

	resp.Fill(key, computePlan, []string{})
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

	err = db.AddComputePlanEvent(inp.Key, computeplan.State.Status, computeplan.State.IntermediaryModel)
	if err != nil {
		return outputComputePlan{}, err
	}
	resp.Fill(inp.Key, computeplan, []string{})
	return resp, nil
}

// Create generate on ID for the compute plan, add it to the ledger
// and register it in the compute plan index
func (cp *ComputePlan) Create(db *LedgerDB) (string, error) {
	ID := GetRandomHash()
	cp.StateKey = GetRandomHash()
	cp.AssetType = ComputePlanType
	err := db.Add(ID, cp)
	if err != nil {
		return "", err
	}
	err = db.Add(cp.StateKey, cp.State)
	if err != nil {
		return "", err
	}
	if err := db.CreateIndex("computePlan~id", []string{"computePlan", ID}); err != nil {
		return "", err
	}
	return ID, nil
}

// Save add or update the compute plan in the ledger
func (cp *ComputePlan) Save(db *LedgerDB, ID string) error {
	err := db.Put(ID, cp)
	if err != nil {
		return err
	}
	return cp.SaveState(db)
}

// SaveState add or update the compute plan in the ledger
func (cp *ComputePlan) SaveState(db *LedgerDB) error {
	return db.Put(cp.StateKey, cp.State)
}

// UpdateStatus check the tuple status (from an updated tuple or a new one)
// and, if required, it updates the compute plan' status and/or its doneCount.
// It returns true if there is any change to the compute plan, false otherwise.
func (cp *ComputePlan) UpdateStatus(tupleStatus string) bool {
	switch cp.State.Status {
	case StatusFailed, StatusCanceled:
	case StatusDone:
		// We might add tuples to a done compute plan
		if stringInSlice(tupleStatus, []string{StatusWaiting, StatusTodo}) {
			cp.State.Status = tupleStatus
			return true
		}
	case StatusDoing:
		switch tupleStatus {
		case StatusFailed:
			cp.State.Status = tupleStatus
			return true
		case StatusDone:
			cp.State.DoneCount++
			if cp.State.DoneCount == cp.State.TupleCount {
				cp.State.Status = tupleStatus
			}
			return true
		}
	case StatusTodo:
		if tupleStatus == StatusDoing {
			cp.State.Status = tupleStatus
			return true
		}
	case StatusWaiting:
		if tupleStatus == StatusTodo {
			cp.State.Status = tupleStatus
			return true
		}
	case "":
		cp.State.Status = tupleStatus
		return true
	}
	return false
}

// AddTuple add the tuple key to the compute plan and update it accordingly
func (cp *ComputePlan) AddTuple(tupleType AssetType, key, status string) {
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
	cp.State.TupleCount++
	cp.UpdateStatus(status)
}

// UpdateComputePlanState retreive the compute plan if the ID is not empty,
// check if the updated status change anything and save it if it's the case
func UpdateComputePlanState(db *LedgerDB, ComputePlanID, tupleStatus, tupleKey string) error {
	if ComputePlanID == "" {
		return nil
	}
	cp, err := db.GetComputePlan(ComputePlanID)
	if err != nil {
		return err
	}
	statusUpdated := cp.UpdateStatus(tupleStatus)
	doneModels, err := cp.HandleIntermediaryModel(db)
	if err != nil {
		return err
	}
	if statusUpdated || len(doneModels) != 0 {
		db.AddComputePlanEvent(ComputePlanID, cp.State.Status, doneModels)
		return cp.SaveState(db)
	}
	return nil
}

// ListModelIfIntermediary will reference the hash model if the compute plan ID
// is not empty and if it's an intermediary model meaning without any children
func ListModelIfIntermediary(db *LedgerDB, ComputePlanID, tupleKey, modelHash string) error {
	if ComputePlanID == "" {
		return nil
	}
	allChildKeys, err := db.GetIndexKeys("tuple~inModel~key", []string{"tuple", tupleKey})
	if err != nil {
		return err
	}
	if len(allChildKeys) == 0 {
		// If a tuple has no children it's concidered final and should not be
		// listed in the index
		return nil
	}
	cp, err := db.GetComputePlan(ComputePlanID)
	if err != nil {
		return err
	}
	cp.State.IntermediaryModel = append(cp.State.IntermediaryModel, modelHash)

	return cp.SaveState(db)
}

// HandleIntermediaryModel is a function
func (cp *ComputePlan) HandleIntermediaryModel(db *LedgerDB) ([]string, error) {
	var doneModels, inUseModels []string
	for _, hash := range cp.State.IntermediaryModel {
		done := true
		keys, err := db.GetIndexKeys("tuple~modelHash~key", []string{"tuple", hash})
		if err != nil {
			return []string{}, err
		}
		tupleKey := keys[0]
		tupleChildKeys, err := db.GetIndexKeys("tuple~inModel~key", []string{"tuple", tupleKey})
		if err != nil {
			return []string{}, err
		}
		for _, key := range tupleChildKeys {
			tuple, err := db.GetGenericTuple(key)
			if err != nil {
				return []string{}, err
			}
			if tuple.Status != StatusDone {
				inUseModels = append(inUseModels, hash)
				done = false
				break
			}
		}
		if done {
			doneModels = append(doneModels, hash)
		}
	}
	cp.State.IntermediaryModel = inUseModels
	return doneModels, nil
}
