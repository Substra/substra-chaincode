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

	return nil
}

// createComputePlan is the wrapper for the substra smartcontract CreateComputePlan
func createComputePlan(db LedgerDB, args []string) (resp outputComputePlan, err error) {
	inp := inputComputePlan{}
	err = AssetFromJSON(args, &inp)
	if err != nil {
		return
	}
	return createComputePlanInternal(db, inp)
}

func createComputePlanInternal(db LedgerDB, inp inputComputePlan) (resp outputComputePlan, err error) {
	traintupleKeysByID := map[string]string{}

	resp.ObjectiveKey = inp.ObjectiveKey
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
				Rank:         strconv.Itoa(i),
				ObjectiveKey: inp.ObjectiveKey,
			}
			if i != 0 {
				inpTraintuple.ComputePlanID = resp.ComputePlanID
			}
			err = inpTraintuple.Fill(computeTraintuple, traintupleKeysByID)
			if err != nil {
				return resp, errors.BadRequest("traintuple ID %s: "+err.Error(), computeTraintuple.ID)
			}

			traintupleKey, err := createTraintupleInternal(db, inpTraintuple)
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
				Rank:         strconv.Itoa(i),
				ObjectiveKey: inp.ObjectiveKey,
			}
			if i != 0 {
				inpCompositeTraintuple.ComputePlanID = resp.ComputePlanID
			}
			err = inpCompositeTraintuple.Fill(computeCompositeTraintuple, traintupleKeysByID)
			if err != nil {
				return resp, errors.BadRequest("traintuple ID %s: "+err.Error(), computeCompositeTraintuple.ID)
			}
			_ = computeCompositeTraintuple
			compositeTraintupleKey, err := createCompositeTraintupleInternal(db, inpCompositeTraintuple)
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
				Rank:         strconv.Itoa(i),
				ObjectiveKey: inp.ObjectiveKey,
			}
			if i != 0 {
				inpAggregatetuple.ComputePlanID = resp.ComputePlanID
			}
			err = inpAggregatetuple.Fill(computeAggregatetuple, traintupleKeysByID)
			if err != nil {
				return resp, errors.BadRequest("traintuple ID %s: "+err.Error(), computeAggregatetuple.ID)
			}
			_ = computeAggregatetuple
			aggregatetupleKey, err := createAggregatetupleInternal(db, inpAggregatetuple)
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

	// event := TuplesEvent{}
	//	event.SetTraintuples(traintuplesTodo...)
	// err = SendTuplesEvent(db.cc, event)
	// if err != nil {
	// return resp, err
	// }

	return resp, err
}

func queryComputePlan(db LedgerDB, args []string) (resp outputComputePlan, err error) {
	inp := inputHash{}
	err = AssetFromJSON(args, &inp)
	if err != nil {
		return
	}
	return getComputePlan(db, inp.Key)
}

func queryComputePlans(db LedgerDB, args []string) (resp []outputComputePlan, err error) {
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
// Traintuples are ordered by ascending rank.
func getComputePlan(db LedgerDB, key string) (resp outputComputePlan, err error) {
	// 1. Get Traintuples and sort them by ascending rank
	var firstTt *Traintuple
	ttKeys, err := db.GetIndexKeys("computePlan~computeplanid~worker~rank~key", []string{"computePlan", key})
	if err != nil {
		return
	}
	if len(ttKeys) == 0 {
		err = errors.E("No traintuple found for compute plan %s", key)
		return
	}
	tts := map[string]Traintuple{}
	for _, ttKey := range ttKeys {
		var tt Traintuple
		tt, err = db.GetTraintuple(ttKey)
		if err != nil {
			return
		}
		if firstTt == nil {
			firstTt = &tt
		}
		tts[ttKey] = tt
	}
	sort.SliceStable(ttKeys, func(i, j int) bool {
		return tts[ttKeys[i]].Rank < tts[ttKeys[j]].Rank
	})

	// 2. Get Testtuples associated with each Traintuple
	tstKeys := []string{}
	for _, traintupleKey := range ttKeys {
		var toAdd []string
		toAdd, err = db.GetIndexKeys("testtuple~traintuple~certified~key", []string{"testtuple", traintupleKey})
		if err != nil {
			return
		}
		tstKeys = append(tstKeys, toAdd...)
	}

	// 3. Create response
	resp = outputComputePlan{
		ComputePlanID:  key,
		ObjectiveKey:   (*firstTt).ObjectiveKey,
		TraintupleKeys: ttKeys,
		TesttupleKeys:  tstKeys,
	}
	return
}
