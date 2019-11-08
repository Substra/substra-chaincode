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
	"sort"
	"strconv"
)

// createComputePlan is the wrapper for the substra smartcontract CreateComputePlan
func createComputePlan(db LedgerDB, args []string) (resp outputComputePlan, err error) {
	inp := inputComputePlan{}
	err = AssetFromJSON(args, &inp)
	if err != nil {
		return
	}
	traintupleKeysByID := map[string]string{}
	resp.TraintupleKeys = []string{}
	var traintuplesTodo []outputTraintuple
	for i, computeTraintuple := range inp.Traintuples {
		inpTraintuple := inputTraintuple{}
		inpTraintuple.AlgoKey = inp.AlgoKey
		inpTraintuple.ObjectiveKey = inp.ObjectiveKey
		inpTraintuple.DataManagerKey = computeTraintuple.DataManagerKey
		inpTraintuple.DataSampleKeys = computeTraintuple.DataSampleKeys
		inpTraintuple.Tag = computeTraintuple.Tag
		inpTraintuple.Rank = strconv.Itoa(i)

		traintuple := Traintuple{}
		err := traintuple.SetFromInput(db, inpTraintuple)
		if err != nil {
			return resp, err
		}

		// Set the inModels by matching the id to traintuples key previously
		// encontered in this compute plan
		for _, InModelID := range computeTraintuple.InModelsIDs {
			inModelKey, ok := traintupleKeysByID[InModelID]
			if !ok {
				return resp, errors.BadRequest("traintuple ID %s: model ID %s not found, check traintuple list order", computeTraintuple.ID, InModelID)
			}
			traintuple.InModelKeys = append(traintuple.InModelKeys, inModelKey)
		}

		traintupleKey := traintuple.GetKey()

		// Set the ComputePlanID
		if i == 0 {
			traintuple.ComputePlanID = traintupleKey
			resp.ComputePlanID = traintuple.ComputePlanID
		} else {
			traintuple.ComputePlanID = resp.ComputePlanID
		}

		// Set status: if it has parents it's waiting
		// if not it's todo and it has to be included in the event
		if len(computeTraintuple.InModelsIDs) > 0 {
			traintuple.Status = StatusWaiting
		} else {
			traintuple.Status = StatusTodo
			out := outputTraintuple{}
			err = out.Fill(db, traintuple, traintupleKey)
			if err != nil {
				return resp, err
			}
			traintuplesTodo = append(traintuplesTodo, out)
		}

		err = traintuple.Save(db, traintupleKey)
		if err != nil {
			return resp, err
		}
		traintupleKeysByID[computeTraintuple.ID] = traintupleKey
		resp.TraintupleKeys = append(resp.TraintupleKeys, traintupleKey)
	}

	resp.TesttupleKeys = []string{}
	for index, computeTesttuple := range inp.Testtuples {
		traintupleKey, ok := traintupleKeysByID[computeTesttuple.TraintupleID]
		if !ok {
			return resp, errors.BadRequest("testtuple index %s: traintuple ID %s not found", index, computeTesttuple.TraintupleID)
		}
		testtuple := Testtuple{}
		testtuple.Model = &Model{TraintupleKey: traintupleKey}
		testtuple.ObjectiveKey = inp.ObjectiveKey
		testtuple.AlgoKey = inp.AlgoKey

		inputTesttuple := inputTesttuple{}
		inputTesttuple.DataManagerKey = computeTesttuple.DataManagerKey
		inputTesttuple.DataSampleKeys = computeTesttuple.DataSampleKeys
		inputTesttuple.Tag = computeTesttuple.Tag
		err = testtuple.SetFromInput(db, inputTesttuple)
		if err != nil {
			return resp, err
		}
		testtuple.Status = StatusWaiting
		testtupleKey := testtuple.GetKey()
		err = testtuple.Save(db, testtupleKey)
		if err != nil {
			return resp, err
		}
		resp.TesttupleKeys = append(resp.TesttupleKeys, testtupleKey)
	}

	event := TuplesEvent{}
	event.SetTraintuples(traintuplesTodo...)
	err = SendTuplesEvent(db.cc, event)
	if err != nil {
		return resp, err
	}

	return resp, err
}

func queryComputePlan(db LedgerDB, args []string) (resp outputComputePlanDetails, err error) {
	inp := inputHash{}
	err = AssetFromJSON(args, &inp)
	if err != nil {
		return
	}
	return getComputePlan(db, inp.Key)
}

func queryComputePlans(db LedgerDB, args []string) (resp []outputComputePlanDetails, err error) {
	computePlanIDs, err := db.GetIndexKeys("computeplan~id", []string{"computeplan"})
	if err != nil {
		return
	}
	for _, key := range computePlanIDs {
		var computePlan outputComputePlanDetails
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
func getComputePlan(db LedgerDB, key string) (resp outputComputePlanDetails, err error) {
	// 1. Get Traintuples and sort them by ascending rank
	var firstTt *Traintuple
	ttKeys, err := db.GetIndexKeys("traintuple~computeplanid~worker~rank~key", []string{"traintuple", key})
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
	var tstKeys []string
	for _, traintupleKey := range ttKeys {
		var toAdd []string
		toAdd, err = db.GetIndexKeys("testtuple~traintuple~certified~key", []string{"testtuple", traintupleKey})
		if err != nil {
			return
		}
		tstKeys = append(tstKeys, toAdd...)
	}

	// 3. Create response
	resp = outputComputePlanDetails{
		ComputePlanID: key,
		AlgoKey:       (*firstTt).AlgoKey,
		ObjectiveKey:  (*firstTt).ObjectiveKey,
		Traintuples:   ttKeys,
		Testtuples:    tstKeys,
	}
	return
}
