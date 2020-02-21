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
	"math"
)

// OutputAssetPaginationHardLimit is a used to avoid issues listing assets
const OutputAssetPaginationHardLimit = 250

// Struct use as output representation of ledger data

type outputObjective struct {
	Key         string            `json:"key"`
	Name        string            `json:"name"`
	Description HashDress         `json:"description"`
	Metrics     *HashDressName    `json:"metrics"`
	Owner       string            `json:"owner"`
	TestDataset *Dataset          `json:"testDataset"`
	Permissions outputPermissions `json:"permissions"`
}

func (out *outputObjective) Fill(key string, in Objective) {
	out.Key = key
	out.Name = in.Name
	out.Description.StorageAddress = in.DescriptionStorageAddress
	out.Description.Hash = key
	out.Metrics = in.Metrics
	out.Owner = in.Owner
	out.TestDataset = in.TestDataset
	out.Permissions.Fill(in.Permissions)
}

// outputDataManager is the return representation of the DataManager type stored in the ledger
type outputDataManager struct {
	ObjectiveKey string            `json:"objectiveKey"`
	Description  *HashDress        `json:"description"`
	Key          string            `json:"key"`
	Name         string            `json:"name"`
	Opener       HashDress         `json:"opener"`
	Owner        string            `json:"owner"`
	Permissions  outputPermissions `json:"permissions"`
	Type         string            `json:"type"`
}

func (out *outputDataManager) Fill(key string, in DataManager) {
	out.ObjectiveKey = in.ObjectiveKey
	out.Description = in.Description
	out.Key = key
	out.Name = in.Name
	out.Opener.Hash = key
	out.Opener.StorageAddress = in.OpenerStorageAddress
	out.Owner = in.Owner
	out.Permissions.Fill(in.Permissions)
	out.Type = in.Type
}

type outputDataSample struct {
	DataManagerKeys []string `json:"dataManagerKeys"`
	Owner           string   `json:"owner"`
	Key             string   `json:"key"`
}

func (out *outputDataSample) Fill(key string, in DataSample) {
	out.Key = key
	out.DataManagerKeys = in.DataManagerKeys
	out.Owner = in.Owner
}

type outputDataset struct {
	outputDataManager
	TrainDataSampleKeys []string `json:"trainDataSampleKeys"`
	TestDataSampleKeys  []string `json:"testDataSampleKeys"`
}

func (out *outputDataset) Fill(key string, in DataManager, trainKeys []string, testKeys []string) {
	out.outputDataManager.Fill(key, in)
	out.TrainDataSampleKeys = trainKeys
	out.TestDataSampleKeys = testKeys
}

type outputAlgo struct {
	Key         string            `json:"key"`
	Name        string            `json:"name"`
	Content     HashDress         `json:"content"`
	Description *HashDress        `json:"description"`
	Owner       string            `json:"owner"`
	Permissions outputPermissions `json:"permissions"`
}

func (out *outputAlgo) Fill(key string, in Algo) {
	out.Key = key
	out.Name = in.Name
	out.Content.Hash = key
	out.Content.StorageAddress = in.StorageAddress
	out.Description = in.Description
	out.Owner = in.Owner
	out.Permissions.Fill(in.Permissions)
}

// outputTtDataset is the representation of a Traintuple Dataset
type outputTtDataset struct {
	Worker         string   `json:"worker"`
	DataSampleKeys []string `json:"keys"`
	OpenerHash     string   `json:"openerHash"`
}

// outputTraintuple is the representation of one the element type stored in the
// ledger. It describes a training task occuring on the platform
type outputTraintuple struct {
	Key           string            `json:"key"`
	Algo          *HashDressName    `json:"algo"`
	Creator       string            `json:"creator"`
	Dataset       *outputTtDataset  `json:"dataset"`
	ComputePlanID string            `json:"computePlanID"`
	InModels      []*Model          `json:"inModels"`
	Log           string            `json:"log"`
	OutModel      *HashDress        `json:"outModel"`
	Permissions   outputPermissions `json:"permissions"`
	Rank          int               `json:"rank"`
	Status        string            `json:"status"`
	Tag           string            `json:"tag"`
}

//Fill is a method of the receiver outputTraintuple. It returns all elements necessary to do a training task from a trainuple stored in the ledger
func (outputTraintuple *outputTraintuple) Fill(db *LedgerDB, traintuple Traintuple, traintupleKey string) (err error) {

	outputTraintuple.Key = traintupleKey
	outputTraintuple.Creator = traintuple.Creator
	outputTraintuple.Permissions.Fill(traintuple.Permissions)
	outputTraintuple.Log = traintuple.Log
	outputTraintuple.Status = traintuple.Status
	outputTraintuple.Rank = traintuple.Rank
	outputTraintuple.ComputePlanID = traintuple.ComputePlanID
	outputTraintuple.OutModel = traintuple.OutModel
	outputTraintuple.Tag = traintuple.Tag
	// fill algo
	algo, err := db.GetAlgo(traintuple.AlgoKey)
	if err != nil {
		err = errors.Internal("could not retrieve algo with key %s - %s", traintuple.AlgoKey, err.Error())
		return
	}
	outputTraintuple.Algo = &HashDressName{
		Name:           algo.Name,
		Hash:           traintuple.AlgoKey,
		StorageAddress: algo.StorageAddress}

	// fill inModels
	for _, inModelKey := range traintuple.InModelKeys {
		if inModelKey == "" {
			break
		}
		parentTraintuple, err := db.GetTraintuple(inModelKey)
		if err != nil {
			return errors.Internal("could not retrieve parent traintuple with key %s - %s", inModelKey, err.Error())
		}
		inModel := &Model{
			TraintupleKey: inModelKey,
		}
		if parentTraintuple.OutModel != nil {
			inModel.Hash = parentTraintuple.OutModel.Hash
			inModel.StorageAddress = parentTraintuple.OutModel.StorageAddress
		}
		outputTraintuple.InModels = append(outputTraintuple.InModels, inModel)
	}

	// fill dataset
	outputTraintuple.Dataset = &outputTtDataset{
		Worker:         traintuple.Dataset.Worker,
		DataSampleKeys: traintuple.Dataset.DataSampleKeys,
		OpenerHash:     traintuple.Dataset.DataManagerKey,
	}

	return
}

type outputTesttuple struct {
	Algo           *HashDressName `json:"algo"`
	Certified      bool           `json:"certified"`
	ComputePlanID  string         `json:"computePlanID"`
	Creator        string         `json:"creator"`
	Dataset        *TtDataset     `json:"dataset"`
	Key            string         `json:"key"`
	Log            string         `json:"log"`
	Objective      *TtObjective   `json:"objective"`
	Rank           int            `json:"rank"`
	Status         string         `json:"status"`
	Tag            string         `json:"tag"`
	TraintupleKey  string         `json:"traintupleKey"`
	TraintupleType string         `json:"traintupleType"`
}

func (out *outputTesttuple) Fill(db *LedgerDB, key string, in Testtuple) error {
	out.Certified = in.Certified
	out.ComputePlanID = in.ComputePlanID
	out.Creator = in.Creator
	out.Dataset = in.Dataset
	out.Key = key
	out.Log = in.Log
	out.Rank = in.Rank
	out.Status = in.Status
	out.Tag = in.Tag
	out.TraintupleKey = in.TraintupleKey

	// fill type
	traintupleType, err := db.GetAssetType(in.TraintupleKey)
	if err != nil {
		return errors.Internal("could not retrieve traintuple type with key %s - %s", in.TraintupleKey, err.Error())
	}
	out.TraintupleType = LowerFirst(traintupleType.String())

	// fill algo
	var algo Algo
	switch traintupleType {
	case TraintupleType:
		algo, err = db.GetAlgo(in.AlgoKey)
		if err != nil {
			return errors.Internal("could not retrieve algo with key %s - %s", in.AlgoKey, err.Error())
		}
	case CompositeTraintupleType:
		compositeAlgo, err := db.GetCompositeAlgo(in.AlgoKey)
		if err != nil {
			return errors.Internal("could not retrieve composite algo with key %s - %s", in.AlgoKey, err.Error())
		}
		algo = compositeAlgo.Algo
	case AggregatetupleType:
		aggregateAlgo, err := db.GetAggregateAlgo(in.AlgoKey)
		if err != nil {
			return errors.Internal("could not retrieve aggregate algo with key %s - %s", in.AlgoKey, err.Error())
		}
		algo = aggregateAlgo.Algo
	}
	out.Algo = &HashDressName{
		Name:           algo.Name,
		Hash:           in.AlgoKey,
		StorageAddress: algo.StorageAddress}

	// fill objective
	objective, err := db.GetObjective(in.ObjectiveKey)
	if err != nil {
		return errors.Internal("could not retrieve associated objective with key %s- %s", in.ObjectiveKey, err.Error())
	}
	if objective.Metrics == nil {
		return errors.Internal("objective %s is missing metrics values", in.ObjectiveKey)
	}
	metrics := HashDress{
		Hash:           objective.Metrics.Hash,
		StorageAddress: objective.Metrics.StorageAddress,
	}
	out.Objective = &TtObjective{
		Key:     in.ObjectiveKey,
		Metrics: &metrics,
	}
	return nil
}

type outputModelDetails struct {
	Aggregatetuple         *outputAggregatetuple      `json:"aggregatetuple,omitempty"`
	CompositeTraintuple    *outputCompositeTraintuple `json:"compositeTraintuple,omitempty"`
	Traintuple             *outputTraintuple          `json:"traintuple,omitempty"`
	Testtuple              outputTesttuple            `json:"testtuple"`
	NonCertifiedTesttuples []outputTesttuple          `json:"nonCertifiedTesttuples"`
}

type outputModel struct {
	Aggregatetuple      *outputAggregatetuple      `json:"aggregatetuple,omitempty"`
	CompositeTraintuple *outputCompositeTraintuple `json:"compositeTraintuple,omitempty"`
	Traintuple          *outputTraintuple          `json:"traintuple,omitempty"`
	Testtuple           outputTesttuple            `json:"testtuple"`
}

// TuplesEvent is the collection of tuples sent in an event
type TuplesEvent struct {
	Testtuples           []outputTesttuple           `json:"testtuple"`
	Traintuples          []outputTraintuple          `json:"traintuple"`
	CompositeTraintuples []outputCompositeTraintuple `json:"compositeTraintuple"`
	Aggregatetuples      []outputAggregatetuple      `json:"aggregatetuple"`
}

type outputComputePlan struct {
	ComputePlanID           string            `json:"computePlanID"`
	TraintupleKeys          []string          `json:"traintupleKeys"`
	AggregatetupleKeys      []string          `json:"aggregatetupleKeys"`
	CompositeTraintupleKeys []string          `json:"compositeTraintupleKeys"`
	TesttupleKeys           []string          `json:"testtupleKeys"`
	Tag                     string            `json:"tag"`
	Status                  string            `json:"status"`
	TupleCount              int               `json:"tupleCount"`
	DoneCount               int               `json:"doneCount"`
	IDToKey                 map[string]string `json:"IDToKey"`
}

func (out *outputComputePlan) Fill(key string, in ComputePlan) {
	out.ComputePlanID = key
	nb := getLimitedNbSliceElements(in.TraintupleKeys)
	out.TraintupleKeys = in.TraintupleKeys[:nb]
	nb = getLimitedNbSliceElements(in.AggregatetupleKeys)
	out.AggregatetupleKeys = in.AggregatetupleKeys[:nb]
	nb = getLimitedNbSliceElements(in.CompositeTraintupleKeys)
	out.CompositeTraintupleKeys = in.CompositeTraintupleKeys[:nb]
	out.TesttupleKeys = in.TesttupleKeys
	out.Status = in.State.Status
	out.Tag = in.Tag
	out.TupleCount = in.State.TupleCount
	out.DoneCount = in.State.DoneCount
	IDToKey := map[string]string{}
	for ID, trainTask := range in.IDToTrainTask {
		IDToKey[ID] = trainTask.Key
	}
	out.IDToKey = IDToKey
}

type outputPermissions struct {
	Process Permission `validate:"required" json:"process"`
}

func (out *outputPermissions) Fill(in Permissions) {
	out.Process.Public = in.Process.Public
	out.Process.AuthorizedIDs = []string{}
	if !in.Process.Public {
		out.Process.AuthorizedIDs = in.Process.AuthorizedIDs
	}
}

type outputLeaderboard struct {
	Objective  outputObjective   `json:"objective"`
	Testtuples outputBoardTuples `json:"testtuples"`
}

type outputBoardTuples []outputBoardTuple

func (out outputBoardTuples) Len() int {
	return len(out)
}

func (out outputBoardTuples) Swap(i, j int) {
	out[i], out[j] = out[j], out[i]
}

func (out outputBoardTuples) Less(i, j int) bool {
	return out[i].Perf < out[j].Perf
}

type outputBoardTuple struct {
	Algo          *HashDressName `json:"algo"`
	Creator       string         `json:"creator"`
	Key           string         `json:"key"`
	TraintupleKey string         `json:"traintupleKey"`
	Perf          float32        `json:"perf"`
	Tag           string         `json:"tag"`
}

func (out *outputBoardTuple) Fill(db *LedgerDB, in Testtuple, testtupleKey string) error {
	out.Key = testtupleKey
	out.Creator = in.Creator
	algo, err := db.GetAlgo(in.AlgoKey)
	if err != nil {
		return err
	}
	out.Algo = &HashDressName{
		Name:           algo.Name,
		Hash:           in.AlgoKey,
		StorageAddress: algo.StorageAddress,
	}
	out.TraintupleKey = in.TraintupleKey
	out.Perf = in.Dataset.Perf
	out.Tag = in.Tag

	return nil
}

func getLimitedNbSliceElements(s []string) int {
	return int(math.Min(float64(len(s)), OutputAssetPaginationHardLimit))
}
