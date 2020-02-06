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
	Algo          *HashDressName        `json:"algo"`
	Creator       string                `json:"creator"`
	Dataset       *outputTtDataset      `json:"dataset"`
	ComputePlanID string                `json:"computePlanID"`
	InHeadModel   *Model                `json:"inHeadModel"`
	InTrunkModel  *Model                `json:"inTrunkModel"`
	Log           string                `json:"log"`
	OutHeadModel  outHeadModelComposite `json:"outHeadModel"`
	OutTrunkModel outModelComposite     `json:"outTrunkModel"`
	Rank          int                   `json:"rank"`
	Status        string                `json:"status"`
	Tag           string                `json:"tag"`
}

type outHeadModelComposite struct {
	OutModel    *Hash             `json:"outModel"`
	Permissions outputPermissions `json:"permissions"`
}

type outModelComposite struct {
	OutModel    *HashDress        `json:"outModel"`
	Permissions outputPermissions `json:"permissions"`
}

//Fill is a method of the receiver outputCompositeTraintuple. It returns all elements necessary to do a training task from a trainuple stored in the ledger
func (outputCompositeTraintuple *outputCompositeTraintuple) Fill(db *LedgerDB, traintuple CompositeTraintuple, traintupleKey string) (err error) {

	outputCompositeTraintuple.Key = traintupleKey
	outputCompositeTraintuple.Creator = traintuple.Creator
	outputCompositeTraintuple.Log = traintuple.Log
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
	outputCompositeTraintuple.Algo = &HashDressName{
		Name:           algo.Name,
		Hash:           traintuple.AlgoKey,
		StorageAddress: algo.StorageAddress}

	// fill in-model (head)
	if traintuple.InHeadModel != "" {
		// Head can only be a composite traintuple's head out model
		hash, _err := db.GetOutModelHash(traintuple.InHeadModel, HeadType, []AssetType{CompositeTraintupleType})
		if _err != nil {
			err = errors.Internal("could not fill (head) in-model with key \"%s\": %s", traintuple.InHeadModel, _err.Error())
			return
		}
		outputCompositeTraintuple.InHeadModel = &Model{
			TraintupleKey: traintuple.InHeadModel}

		if hash != nil {
			outputCompositeTraintuple.InHeadModel.Hash = hash.Hash
		}
	}

	// fill in-model (trunk)
	if traintuple.InTrunkModel != "" {
		// Trunk can be either:
		// - a traintuple's out model
		// - a composite traintuple's head out model
		// - an aggregate tuple's out model
		hashDress, _err := db.GetOutModelHashDress(traintuple.InTrunkModel, TrunkType, []AssetType{TraintupleType, CompositeTraintupleType, AggregatetupleType})
		if _err != nil {
			err = errors.Internal("could not fill (trunk) in-model with key \"%s\": %s", traintuple.InTrunkModel, _err.Error())
			return
		}
		outputCompositeTraintuple.InTrunkModel = &Model{
			TraintupleKey: traintuple.InTrunkModel}

		if hashDress != nil {
			outputCompositeTraintuple.InTrunkModel.Hash = hashDress.Hash
			outputCompositeTraintuple.InTrunkModel.StorageAddress = hashDress.StorageAddress
		}
	}

	// fill dataset
	outputCompositeTraintuple.Dataset = &outputTtDataset{
		Worker:         traintuple.Dataset.Worker,
		DataSampleKeys: traintuple.Dataset.DataSampleKeys,
		OpenerHash:     traintuple.Dataset.DataManagerKey,
	}

	return
}

func getOutPermissions(in Permissions) (out outputPermissions) {
	out = outputPermissions{}
	out.Fill(in)
	return out
}

func (out *outputCompositeAlgo) Fill(key string, in CompositeAlgo) {
	out.outputAlgo.Fill(key, in.Algo)
}
