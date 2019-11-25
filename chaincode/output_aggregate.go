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

import "fmt"

type outputAggregateTuple struct {
	Key           string            `json:"key"`
	Algo          *HashDressName    `json:"algo"`
	Creator       string            `json:"creator"`
	ComputePlanID string            `json:"computePlanID"`
	Log           string            `json:"log"`
	Objective     *TtObjective      `json:"objective"`
	InModels      []*Model          `json:"inModels"`
	OutModel      *HashDress        `json:"outModel"`
	Rank          int               `json:"rank"`
	Status        string            `json:"status"`
	Tag           string            `json:"tag"`
	Permissions   outputPermissions `json:"permissions"`
	Worker        string            `json:"worker"`
}

type outputAggregateAlgo struct {
	outputAlgo
}

func (out *outputAggregateAlgo) Fill(key string, in AggregateAlgo) {
	out.outputAlgo.Fill(key, in.Algo)
}

// Fill is a method of the receiver outputAggregateTuple. It returns all elements necessary to do a training task from an aggregate trainuple stored in the ledger
func (outputAggregateTuple *outputAggregateTuple) Fill(db LedgerDB, traintuple AggregateTuple, traintupleKey string) (err error) {
	outputAggregateTuple.Key = traintupleKey
	outputAggregateTuple.Creator = traintuple.Creator
	outputAggregateTuple.Log = traintuple.Log
	outputAggregateTuple.Status = traintuple.Status
	outputAggregateTuple.Rank = traintuple.Rank
	outputAggregateTuple.ComputePlanID = traintuple.ComputePlanID
	outputAggregateTuple.OutModel = traintuple.OutModel
	outputAggregateTuple.Tag = traintuple.Tag
	algo, err := db.GetAggregateAlgo(traintuple.AlgoKey)
	if err != nil {
		err = fmt.Errorf("could not retrieve aggregate algo with key %s - %s", traintuple.AlgoKey, err.Error())
		return
	}
	outputAggregateTuple.Algo = &HashDressName{
		Name:           algo.Name,
		Hash:           traintuple.AlgoKey,
		StorageAddress: algo.StorageAddress}

	// fill objective
	objective, err := db.GetObjective(traintuple.ObjectiveKey)
	if err != nil {
		err = fmt.Errorf("could not retrieve associated objective with key %s- %s", traintuple.ObjectiveKey, err.Error())
		return
	}
	if objective.Metrics == nil {
		err = fmt.Errorf("objective %s is missing metrics values", traintuple.ObjectiveKey)
		return
	}
	metrics := HashDress{
		Hash:           objective.Metrics.Hash,
		StorageAddress: objective.Metrics.StorageAddress,
	}
	outputAggregateTuple.Objective = &TtObjective{
		Key:     traintuple.ObjectiveKey,
		Metrics: &metrics,
	}

	// fill inModels
	for _, inModelKey := range traintuple.InModelKeys {
		if inModelKey == "" {
			break
		}
		hashDress, _err := db.GetOutModelHashDress(inModelKey, TrunkType, []AssetType{TraintupleType, CompositeTraintupleType, AggregateTupleType})
		if _err != nil {
			err = fmt.Errorf("could not fill in-model with key \"%s\": %s", inModelKey, _err.Error())
			return
		}
		inModel := &Model{
			TraintupleKey: inModelKey,
		}
		if hashDress != nil {
			inModel.Hash = hashDress.Hash
			inModel.StorageAddress = hashDress.StorageAddress
		}
		outputAggregateTuple.InModels = append(outputAggregateTuple.InModels, inModel)
	}

	outputAggregateTuple.Worker = traintuple.Worker

	return
}

// AddAggregateTuple adds one aggregate tuple to the event struct
func (te *TuplesEvent) AddAggregateTuple(out outputAggregateTuple) {
	te.AggregateTuples = append(te.AggregateTuples, out)
}

// SetAggregateTuples adds one or several tuples to the event struct
func (te *TuplesEvent) SetAggregateTuples(otuples ...outputAggregateTuple) {
	te.AggregateTuples = otuples
}
