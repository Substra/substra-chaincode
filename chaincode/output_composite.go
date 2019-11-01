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
	"fmt"
)

type outputCompositeTraintuple struct {
	Key           string            `json:"key"`
	Algo          *HashDressName    `json:"algo"`
	Creator       string            `json:"creator"`
	Dataset       *TtDataset        `json:"dataset"`
	ComputePlanID string            `json:"computePlanID"`
	InModelHead   *Model            `json:"inModelHead"`
	InModelTrunk  *Model            `json:"inModelTrunk"`
	Log           string            `json:"log"`
	Objective     *TtObjective      `json:"objective"`
	OutHeadModel  outModelComposite `json:"outHeadModel"`
	OutTrunkModel outModelComposite `json:"outTrunkModel"`
	Rank          int               `json:"rank"`
	Status        string            `json:"status"`
	Tag           string            `json:"tag"`
}

type outModelComposite struct {
	OutModel    *HashDress        `json:"outModel"`
	Permissions outputPermissions `json:"permissions"`
}

//Fill is a method of the receiver outputCompositeTraintuple. It returns all elements necessary to do a training task from a trainuple stored in the ledger
func (outputCompositeTraintuple *outputCompositeTraintuple) Fill(db LedgerDB, traintuple CompositeTraintuple, traintupleKey string) (err error) {

	outputCompositeTraintuple.Key = traintupleKey
	outputCompositeTraintuple.Creator = traintuple.Creator
	outputCompositeTraintuple.Log = traintuple.Log
	outputCompositeTraintuple.Status = traintuple.Status
	outputCompositeTraintuple.Rank = traintuple.Rank
	outputCompositeTraintuple.ComputePlanID = traintuple.ComputePlanID
	outputCompositeTraintuple.OutHeadModel = outModelComposite{
		OutModel:    traintuple.OutHeadModel.OutModel,
		Permissions: getOutPermissions(traintuple.OutHeadModel.Permissions)}
	outputCompositeTraintuple.OutTrunkModel = outModelComposite{
		OutModel:    traintuple.OutTrunkModel.OutModel,
		Permissions: getOutPermissions(traintuple.OutTrunkModel.Permissions)}
	outputCompositeTraintuple.Tag = traintuple.Tag
	// fill algo
	algo, err := db.GetAlgo(traintuple.AlgoKey)
	if err != nil {
		err = fmt.Errorf("could not retrieve algo with key %s - %s", traintuple.AlgoKey, err.Error())
		return
	}
	outputCompositeTraintuple.Algo = &HashDressName{
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
	outputCompositeTraintuple.Objective = &TtObjective{
		Key:     traintuple.ObjectiveKey,
		Metrics: &metrics,
	}

	// fill in-model (head)
	if traintuple.InModelHead != "" {
		hashDress, _err := db.GetOutModelHashDress(traintuple.InModelHead, HeadType, []AssetType{TraintupleType, CompositeTraintupleType})
		if _err != nil {
			err = fmt.Errorf("could not fill (head) in-model with key \"%s\": %s", traintuple.InModelHead, _err.Error())
			return
		}
		outputCompositeTraintuple.InModelHead = &Model{
			TraintupleKey: traintuple.InModelHead}

		if hashDress != nil {
			outputCompositeTraintuple.InModelHead.Hash = hashDress.Hash
			outputCompositeTraintuple.InModelHead.StorageAddress = hashDress.StorageAddress
		}
	}

	// fill in-model (trunk)
	if traintuple.InModelTrunk != "" {
		hashDress, _err := db.GetOutModelHashDress(traintuple.InModelTrunk, TrunkType, []AssetType{TraintupleType, CompositeTraintupleType /* TODO: add AggregateTraintupleTYpe */})
		if _err != nil {
			err = fmt.Errorf("could not fill (trunk) in-model with key \"%s\": %s", traintuple.InModelTrunk, _err.Error())
			return
		}
		outputCompositeTraintuple.InModelTrunk = &Model{
			TraintupleKey: traintuple.InModelTrunk}

		if hashDress != nil {
			outputCompositeTraintuple.InModelTrunk.Hash = hashDress.Hash
			outputCompositeTraintuple.InModelTrunk.StorageAddress = hashDress.StorageAddress
		}
	}

	// fill dataset
	outputCompositeTraintuple.Dataset = &TtDataset{
		Worker:         traintuple.Dataset.Worker,
		DataSampleKeys: traintuple.Dataset.DataSampleKeys,
		OpenerHash:     traintuple.Dataset.DataManagerKey,
		Perf:           traintuple.Perf,
	}

	return
}

func getOutPermissions(in Permissions) (out outputPermissions) {
	out = outputPermissions{}
	out.Fill(in)
	return out
}
