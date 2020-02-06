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

func (inpTraintuple *inputTraintuple) Fill(inpCP inputComputePlanTraintuple, IDToCPItem map[string]CPItem) error {
	inpTraintuple.DataManagerKey = inpCP.DataManagerKey
	inpTraintuple.DataSampleKeys = inpCP.DataSampleKeys
	inpTraintuple.AlgoKey = inpCP.AlgoKey
	inpTraintuple.Tag = inpCP.Tag

	// Set the inModels by matching the id to tuples key previously
	// encontered in this compute plan
	for _, InModelID := range inpCP.InModelsIDs {
		item, ok := IDToCPItem[InModelID]
		if !ok {
			return errors.BadRequest("model ID %s not found", InModelID)
		}
		inpTraintuple.InModels = append(inpTraintuple.InModels, item.Key)
	}

	return nil

}
func (inpAggregatetuple *inputAggregatetuple) Fill(inpCP inputComputePlanAggregatetuple, IDToCPItem map[string]CPItem) error {
	inpAggregatetuple.AlgoKey = inpCP.AlgoKey
	inpAggregatetuple.Tag = inpCP.Tag
	inpAggregatetuple.Worker = inpCP.Worker

	// Set the inModels by matching the id to tuples key previously
	// encontered in this compute plan
	for _, InModelID := range inpCP.InModelsIDs {
		item, ok := IDToCPItem[InModelID]
		if !ok {
			return errors.BadRequest("model ID %s not found", InModelID)
		}
		inpAggregatetuple.InModels = append(inpAggregatetuple.InModels, item.Key)
	}

	return nil

}
func (inpCompositeTraintuple *inputCompositeTraintuple) Fill(inpCP inputComputePlanCompositeTraintuple, IDToCPItem map[string]CPItem) error {
	inpCompositeTraintuple.DataManagerKey = inpCP.DataManagerKey
	inpCompositeTraintuple.DataSampleKeys = inpCP.DataSampleKeys
	inpCompositeTraintuple.AlgoKey = inpCP.AlgoKey
	inpCompositeTraintuple.Tag = inpCP.Tag
	inpCompositeTraintuple.OutTrunkModelPermissions = inpCP.OutTrunkModelPermissions

	// Set the inModels by matching the id to traintuples key previously
	// encontered in this compute plan
	if inpCP.InHeadModelID != "" {
		var ok bool
		item, ok := IDToCPItem[inpCP.InHeadModelID]
		if !ok {
			return errors.BadRequest("head model ID %s not found", inpCP.InHeadModelID)
		}
		inpCompositeTraintuple.InHeadModelKey = item.Key
	}
	if inpCP.InTrunkModelID != "" {
		var ok bool
		item, ok := IDToCPItem[inpCP.InTrunkModelID]
		if !ok {
			return errors.BadRequest("trunk model ID %s not found", inpCP.InTrunkModelID)
		}
		inpCompositeTraintuple.InTrunkModelKey = item.Key
	}
	return nil
}

func (inpTesttuple *inputTesttuple) Fill(inpCP inputComputePlanTesttuple, IDToCPItem map[string]CPItem) error {
	item, ok := IDToCPItem[inpCP.TraintupleID]
	if !ok {
		return errors.BadRequest("traintuple ID %s not found", inpCP.TraintupleID)
	}
	inpTesttuple.TraintupleKey = item.Key
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

	count := len(inp.Traintuples) + len(inp.Aggregatetuples) + len(inp.CompositeTraintuples)
	if count == 0 {
		return resp, errors.BadRequest("empty update for compute plan %s", inp.ID)
	}
	return updateComputePlanInternal(db, inp.ID, inp.inputComputePlan)
}

func createComputePlanInternal(db *LedgerDB, inp inputComputePlan, tag string) (resp outputComputePlan, err error) {
	var computePlan ComputePlan
	computePlan.Status = StatusWaiting
	computePlan.Tag = tag
	ID, err := computePlan.Create(db)
	if err != nil {
		return resp, err
	}
	count := len(inp.Traintuples) + len(inp.Aggregatetuples) + len(inp.CompositeTraintuples)
	if count == 0 {
		resp.Fill(ID, computePlan)
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
	IDToItem := map[string]CPItem{}
	for ID, item := range computePlan.IDToItem {
		IDToItem[ID] = item
	}
	DAG, err := createComputeDAG(inp, computePlan.IDToItem)
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
			err = inpTraintuple.Fill(computeTraintuple, IDToItem)
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
			err = inpCompositeTraintuple.Fill(computeCompositeTraintuple, IDToItem)
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
			err = inpAggregatetuple.Fill(computeAggregatetuple, IDToItem)
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
		IDToItem[task.ID] = CPItem{Depth: task.Depth, Key: tupleKey}
	}

	for index, computeTesttuple := range inp.Testtuples {
		inpTesttuple := inputTesttuple{}
		err = inpTesttuple.Fill(computeTesttuple, IDToItem)
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
	computePlan.IDToItem = IDToItem
	err = computePlan.Save(db, computePlanID)
	if err != nil {
		return resp, err
	}
	resp.Fill(computePlanID, computePlan)
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

	resp.Fill(key, computePlan)
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

	computeplan.Status = StatusCanceled
	err = db.Put(inp.Key, computeplan)
	if err != nil {
		return outputComputePlan{}, err
	}

	resp.Fill(inp.Key, computeplan)
	return resp, nil
}

// Create generate on ID for the compute plan, add it to the ledger
// and register it in the compute plan index
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

// Save add or update the compute plan in the ledger
func (cp *ComputePlan) Save(db *LedgerDB, ID string) error {
	err := db.Put(ID, cp)
	if err != nil {
		return err
	}
	return nil
}

// CheckNewTupleStatus check the tuple status (from an updated tuple or a new one)
// and, if required, it updates the compute plan' status and/or its doneCount.
// It returns true if there is any change to the compute plan, false otherwise.
func (cp *ComputePlan) CheckNewTupleStatus(tupleStatus string) bool {
	switch cp.Status {
	case StatusFailed, StatusCanceled:
	case StatusDone:
		// We might add tuples to a done compute plan
		if stringInSlice(tupleStatus, []string{StatusWaiting, StatusTodo}) {
			cp.Status = tupleStatus
			return true
		}
	case StatusDoing:
		switch tupleStatus {
		case StatusFailed:
			cp.Status = tupleStatus
			return true
		case StatusDone:
			cp.DoneCount++
			if cp.DoneCount == cp.TupleCount {
				cp.Status = tupleStatus
			}
			return true
		}
	case StatusTodo:
		if tupleStatus == StatusDoing {
			cp.Status = tupleStatus
			return true
		}
	case StatusWaiting:
		if tupleStatus == StatusTodo {
			cp.Status = tupleStatus
			return true
		}
	case "":
		cp.Status = tupleStatus
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
	cp.TupleCount++
	cp.CheckNewTupleStatus(status)
}

// UpdateComputePlan retreive the compute plan if the ID is not empty,
// check if the updated status change anything and save it if it's the case
func UpdateComputePlan(db *LedgerDB, ComputePlanID, tupleStatus, tupleKey string) error {
	if ComputePlanID == "" {
		return nil
	}
	cp, err := db.GetComputePlan(ComputePlanID)
	if err != nil {
		return err
	}
	if cp.CheckNewTupleStatus(tupleStatus) {
		return cp.Save(db, ComputePlanID)
	}
	return nil
}
