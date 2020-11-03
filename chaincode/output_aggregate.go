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

import "chaincode/errors"

type outputAggregatetuple struct {
	Key           string            `json:"key"`
	Algo          *KeyHashDressName `json:"algo"`
	Creator       string            `json:"creator"`
	ComputePlanID string            `json:"compute_plan_id"`
	Log           string            `json:"log"`
	Metadata      map[string]string `json:"metadata"`
	InModels      []*Model          `json:"in_models"`
	OutModel      *KeyHashDress     `json:"out_model"`
	Rank          int               `json:"rank"`
	Status        string            `json:"status"`
	Tag           string            `json:"tag"`
	Permissions   outputPermissions `json:"permissions"`
	Worker        string            `json:"worker"`
}

type outputAggregateAlgo struct {
	outputAlgo
}

func (out *outputAggregateAlgo) Fill(in AggregateAlgo) {
	out.outputAlgo.Fill(in.Algo)
}

// Fill is a method of the receiver outputAggregatetuple. It returns all elements necessary to do a training task from an aggregate trainuple stored in the ledger
func (outputAggregatetuple *outputAggregatetuple) Fill(db *LedgerDB, traintuple Aggregatetuple) (err error) {
	outputAggregatetuple.Key = traintuple.Key
	outputAggregatetuple.Creator = traintuple.Creator
	outputAggregatetuple.Log = traintuple.Log
	outputAggregatetuple.Metadata = initMapOutput(traintuple.Metadata)
	outputAggregatetuple.Status = traintuple.Status
	outputAggregatetuple.Rank = traintuple.Rank
	outputAggregatetuple.ComputePlanID = traintuple.ComputePlanID
	outputAggregatetuple.OutModel = traintuple.OutModel
	outputAggregatetuple.Tag = traintuple.Tag
	algo, err := db.GetAggregateAlgo(traintuple.AlgoKey)
	if err != nil {
		err = errors.Internal("could not retrieve aggregate algo with key %s - %s", traintuple.AlgoKey, err.Error())
		return
	}
	outputAggregatetuple.Algo = &KeyHashDressName{
		Key:            algo.Key,
		Name:           algo.Name,
		Hash:           algo.Hash,
		StorageAddress: algo.StorageAddress}

	// fill inModels
	for _, inModelKey := range traintuple.InModelKeys {
		if inModelKey == "" {
			break
		}
		keyHashDress, _err := db.GetOutModelKeyHashDress(inModelKey, []AssetType{TraintupleType, CompositeTraintupleType, AggregatetupleType})
		if _err != nil {
			err = errors.Internal("could not fill in-model with key \"%s\": %s", inModelKey, _err.Error())
			return
		}
		inModel := &Model{
			TraintupleKey: inModelKey,
		}
		if keyHashDress != nil {
			inModel.Key = keyHashDress.Key
			inModel.Hash = keyHashDress.Hash
			inModel.StorageAddress = keyHashDress.StorageAddress
		}
		outputAggregatetuple.InModels = append(outputAggregatetuple.InModels, inModel)
	}

	outputAggregatetuple.Worker = traintuple.Worker
	outputAggregatetuple.Permissions.Fill(traintuple.Permissions)

	return
}
