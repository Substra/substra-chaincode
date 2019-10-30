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

type outputTraintupleComposite struct {
	Key           string            `json:"key"`
	Algo          *HashDressName    `json:"algo"`
	Creator       string            `json:"creator"`
	Dataset       *TtDataset        `json:"dataset"`
	ComputePlanID string            `json:"computePlanID"`
	InModelTrunk  *Model            `json:"inModelTrunk"`
	InModelHead   *Model            `json:"inModelHead"`
	Log           string            `json:"log"`
	Objective     *TtObjective      `json:"objective"`
	OutModelTrunk outModelComposite `json:"outModelTrunk"`
	OutModelHead  outModelComposite `json:"outModelHead"`
	Rank          int               `json:"rank"`
	Status        string            `json:"status"`
	Tag           string            `json:"tag"`
}

type outModelComposite struct {
	OutModel    *HashDress        `json:"outModel"`
	Permissions outputPermissions `json:"permissions"`
}

//Fill is a method of the receiver outputTraintupleComposite. It returns all elements necessary to do a training task from a trainuple stored in the ledger
func (outputTraintupleComposite *outputTraintupleComposite) Fill(db LedgerDB, traintuple TraintupleComposite, traintupleKey string) (err error) {

	outputTraintupleComposite.Key = traintupleKey
	outputTraintupleComposite.Creator = traintuple.Creator
	outputTraintupleComposite.Log = traintuple.Log
	outputTraintupleComposite.Status = traintuple.Status
	outputTraintupleComposite.Rank = traintuple.Rank
	outputTraintupleComposite.ComputePlanID = traintuple.ComputePlanID
	outputTraintupleComposite.OutModelHead = outModelComposite{
		OutModel:    traintuple.OutModelHead.OutModel,
		Permissions: getOutPermissions(traintuple.OutModelHead.Permissions)}
	outputTraintupleComposite.OutModelTrunk = outModelComposite{
		OutModel:    traintuple.OutModelTrunk.OutModel,
		Permissions: getOutPermissions(traintuple.OutModelTrunk.Permissions)}
	outputTraintupleComposite.Tag = traintuple.Tag
	// fill algo
	algo, err := db.GetAlgo(traintuple.AlgoKey)
	if err != nil {
		err = fmt.Errorf("could not retrieve algo with key %s - %s", traintuple.AlgoKey, err.Error())
		return
	}
	outputTraintupleComposite.Algo = &HashDressName{
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
	outputTraintupleComposite.Objective = &TtObjective{
		Key:     traintuple.ObjectiveKey,
		Metrics: &metrics,
	}

	// fill inModel (trunk)
	trunk, err := db.GetTraintuple(traintuple.InModelTrunk)
	if err != nil {
		return fmt.Errorf("could not retrieve parent traintuple (trunk) with key %s - %s", traintuple.InModelTrunk, err.Error())
	}
	outputTraintupleComposite.InModelTrunk = &Model{
		TraintupleKey: traintuple.InModelTrunk,
	}
	if trunk.OutModel != nil {
		outputTraintupleComposite.InModelTrunk.Hash = trunk.OutModel.Hash
		outputTraintupleComposite.InModelTrunk.StorageAddress = trunk.OutModel.StorageAddress
	}

	// fill inModel (head)
	head, err := db.GetTraintuple(traintuple.InModelHead)
	if err != nil {
		return fmt.Errorf("could not retrieve parent traintuple (head) with key %s - %s", traintuple.InModelHead, err.Error())
	}
	outputTraintupleComposite.InModelHead = &Model{
		TraintupleKey: traintuple.InModelHead,
	}
	if head.OutModel != nil {
		outputTraintupleComposite.InModelHead.Hash = head.OutModel.Hash
		outputTraintupleComposite.InModelHead.StorageAddress = head.OutModel.StorageAddress
	}

	// fill dataset
	outputTraintupleComposite.Dataset = &TtDataset{
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
