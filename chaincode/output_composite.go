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

type outputCompositeAlgo struct {
	outputAlgo
}

type outputCompositeTraintuple struct {
	Key           string                `json:"key"`
	Algo          *KeyHashDressName     `json:"algo"`
	Creator       string                `json:"creator"`
	Dataset       *outputTtDataset      `json:"dataset"`
	ComputePlanID string                `json:"compute_plan_id"`
	InHeadModel   *Model                `json:"in_head_model"`
	InTrunkModel  *Model                `json:"in_trunk_model"`
	Log           string                `json:"log"`
	Metadata      map[string]string     `json:"metadata"`
	OutHeadModel  outHeadModelComposite `json:"out_head_model"`
	OutTrunkModel outModelComposite     `json:"out_trunk_model"`
	Rank          int                   `json:"rank"`
	Status        string                `json:"status"`
	Tag           string                `json:"tag"`
}

type outHeadModelComposite struct {
	OutModel    *KeyHash          `json:"out_model"`
	Permissions outputPermissions `json:"permissions"`
}

type outModelComposite struct {
	OutModel    *KeyHashDress     `json:"out_model"`
	Permissions outputPermissions `json:"permissions"`
}

//Fill is a method of the receiver outputCompositeTraintuple. It returns all elements necessary to do a training task from a trainuple stored in the ledger
func (outputCompositeTraintuple *outputCompositeTraintuple) Fill(db *LedgerDB, traintuple CompositeTraintuple) (err error) {

	outputCompositeTraintuple.Key = traintuple.Key
	outputCompositeTraintuple.Creator = traintuple.Creator
	outputCompositeTraintuple.Log = traintuple.Log
	outputCompositeTraintuple.Metadata = initMapOutput(traintuple.Metadata)
	outputCompositeTraintuple.Status = traintuple.Status
	outputCompositeTraintuple.Rank = traintuple.Rank
	outputCompositeTraintuple.ComputePlanID = traintuple.ComputePlanID
	outputCompositeTraintuple.OutHeadModel = outHeadModelComposite{
		OutModel:    traintuple.OutHeadModel.OutModel,
		Permissions: getOutPermissions(traintuple.OutHeadModel.Permissions)}
	outputCompositeTraintuple.OutTrunkModel = outModelComposite{
		OutModel:    traintuple.OutTrunkModel.OutModel,
		Permissions: getOutPermissions(traintuple.OutTrunkModel.Permissions)}
	outputCompositeTraintuple.Tag = traintuple.Tag
	// fill algo
	algo, err := db.GetCompositeAlgo(traintuple.AlgoKey)
	if err != nil {
		err = errors.Internal("could not retrieve composite algo with key %s - %s", traintuple.AlgoKey, err.Error())
		return
	}
	outputCompositeTraintuple.Algo = &KeyHashDressName{
		Key:            algo.Key,
		Name:           algo.Name,
		Hash:           algo.Hash,
		StorageAddress: algo.StorageAddress}

	// fill in-model (head)
	if traintuple.InHeadModel != "" {
		// Head can only be a composite traintuple's head out model
		hash, _err := db.GetOutHeadModelHashKey(traintuple.InHeadModel)
		if _err != nil {
			err = errors.Internal("could not fill (head) in-model with key \"%s\": %s", traintuple.InHeadModel, _err.Error())
			return
		}
		outputCompositeTraintuple.InHeadModel = &Model{
			TraintupleKey: traintuple.InHeadModel}

		if hash != nil {
			outputCompositeTraintuple.InHeadModel.Key = hash.Key
			outputCompositeTraintuple.InHeadModel.Hash = hash.Hash
		}
	}

	// fill in-model (trunk)
	if traintuple.InTrunkModel != "" {
		// Trunk can be either:
		// - a traintuple's out model
		// - a composite traintuple's head out model
		// - an aggregate tuple's out model
		hashDress, _err := db.GetOutModelHashDressKey(traintuple.InTrunkModel, []AssetType{TraintupleType, CompositeTraintupleType, AggregatetupleType})
		if _err != nil {
			err = errors.Internal("could not fill (trunk) in-model with key \"%s\": %s", traintuple.InTrunkModel, _err.Error())
			return
		}
		outputCompositeTraintuple.InTrunkModel = &Model{
			TraintupleKey: traintuple.InTrunkModel}

		if hashDress != nil {
			outputCompositeTraintuple.InTrunkModel.Key = hashDress.Key
			outputCompositeTraintuple.InTrunkModel.Hash = hashDress.Hash
			outputCompositeTraintuple.InTrunkModel.StorageAddress = hashDress.StorageAddress
		}
	}

	dataManager, err := db.GetDataManager(traintuple.Dataset.DataManagerKey)
	if err != nil {
		err = errors.Internal("could not retrieve data manager with key %s - %s", traintuple.Dataset.DataManagerKey, err.Error())
		return
	}

	// fill dataset
	outputCompositeTraintuple.Dataset = &outputTtDataset{
		Key:            dataManager.Key,
		Worker:         traintuple.Dataset.Worker,
		DataSampleKeys: traintuple.Dataset.DataSampleKeys,
		OpenerHash:     dataManager.Opener.Hash,
		Metadata:       initMapOutput(traintuple.Dataset.Metadata),
	}

	return
}

func getOutPermissions(in Permissions) (out outputPermissions) {
	out = outputPermissions{}
	out.Fill(in)
	return out
}

func (out *outputCompositeAlgo) Fill(in CompositeAlgo) {
	out.outputAlgo.Fill(in.Algo)
}
