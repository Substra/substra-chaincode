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
const OutputAssetPaginationHardLimit = 500

// Struct use as output representation of ledger data

type outputObjective struct {
	Key         string            `json:"key"`
	Name        string            `json:"name"`
	Description *HashDress        `json:"description"`
	Metrics     *HashDressName    `json:"metrics"`
	Owner       string            `json:"owner"`
	TestDataset *Dataset          `json:"test_dataset"`
	Permissions outputPermissions `json:"permissions"`
	Metadata    map[string]string `json:"metadata"`
}

func (out *outputObjective) Fill(in Objective) {
	out.Key = in.Key
	out.Name = in.Name
	out.Description = in.Description
	out.Metrics = in.Metrics
	out.Owner = in.Owner
	out.TestDataset = in.TestDataset
	if out.TestDataset != nil {
		out.TestDataset.Metadata = initMapOutput(in.TestDataset.Metadata)
	}
	out.Permissions.Fill(in.Permissions)
	out.Metadata = initMapOutput(in.Metadata)
}

// outputDataManager is the return representation of the DataManager type stored in the ledger
type outputDataManager struct {
	ObjectiveKey string            `json:"objective_key"`
	Description  *HashDress        `json:"description"`
	Key          string            `json:"key"`
	Metadata     map[string]string `json:"metadata"`
	Name         string            `json:"name"`
	Opener       *HashDress        `json:"opener"`
	Owner        string            `json:"owner"`
	Permissions  outputPermissions `json:"permissions"`
	Type         string            `json:"type"`
}

func (out *outputDataManager) Fill(in DataManager) {
	out.ObjectiveKey = in.ObjectiveKey
	out.Description = in.Description
	out.Key = in.Key
	out.Metadata = initMapOutput(in.Metadata)
	out.Name = in.Name
	out.Opener = in.Opener
	out.Owner = in.Owner
	out.Permissions.Fill(in.Permissions)
	out.Type = in.Type
}

type outputDataSample struct {
	DataManagerKeys []string `json:"data_manager_keys"`
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
	Metadata            map[string]string `json:"metadata"`
	TrainDataSampleKeys []string          `json:"train_data_sample_keys"`
	TestDataSampleKeys  []string          `json:"test_data_sample_keys"`
}

func (out *outputDataset) Fill(in DataManager, trainKeys []string, testKeys []string) {
	out.outputDataManager.Fill(in)
	out.TrainDataSampleKeys = trainKeys
	out.TestDataSampleKeys = testKeys
	out.Metadata = initMapOutput(in.Metadata)
}

type outputAlgo struct {
	Key         string            `json:"key"`
	Name        string            `json:"name"`
	Content     *HashDress        `json:"content"`
	Description *HashDress        `json:"description"`
	Owner       string            `json:"owner"`
	Permissions outputPermissions `json:"permissions"`
	Metadata    map[string]string `json:"metadata"`
}

func (out *outputAlgo) Fill(in Algo) {
	out.Key = in.Key
	out.Name = in.Name
	out.Content = &HashDress{
		Hash:           in.Hash,
		StorageAddress: in.StorageAddress,
	}
	out.Description = in.Description
	out.Owner = in.Owner
	out.Permissions.Fill(in.Permissions)
	out.Metadata = initMapOutput(in.Metadata)
}

// outputTtDataset is the representation of a Traintuple Dataset
type outputTtDataset struct {
	Key            string            `json:"key"`
	Worker         string            `json:"worker"`
	DataSampleKeys []string          `json:"keys"`
	OpenerHash     string            `json:"opener_hash"`
	Metadata       map[string]string `json:"metadata"`
}

// outputTraintuple is the representation of one the element type stored in the
// ledger. It describes a training task occuring on the platform
type outputTraintuple struct {
	Key           string            `json:"key"`
	Algo          *HashDressNameKey `json:"algo"`
	Creator       string            `json:"creator"`
	Dataset       *outputTtDataset  `json:"dataset"`
	ComputePlanID string            `json:"compute_plan_id"`
	InModels      []*Model          `json:"in_models"`
	Log           string            `json:"log"`
	Metadata      map[string]string `json:"metadata"`
	OutModel      *HashDress        `json:"out_model"`
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
	outputTraintuple.Metadata = initMapOutput(traintuple.Metadata)
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
	outputTraintuple.Algo = &HashDressNameKey{
		Key:            algo.Key,
		Name:           algo.Name,
		Hash:           algo.Hash,
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

	dataManager, err := db.GetDataManager(traintuple.Dataset.DataManagerKey)
	if err != nil {
		err = errors.Internal("could not retrieve data manager with key %s - %s", traintuple.Dataset.DataManagerKey, err.Error())
		return
	}

	// fill dataset
	outputTraintuple.Dataset = &outputTtDataset{
		Key:            dataManager.Key,
		Worker:         traintuple.Dataset.Worker,
		DataSampleKeys: traintuple.Dataset.DataSampleKeys,
		OpenerHash:     dataManager.Opener.Hash,
		Metadata:       initMapOutput(traintuple.Dataset.Metadata),
	}

	return
}

type outputTesttuple struct {
	Algo           *HashDressNameKey `json:"algo"`
	Certified      bool              `json:"certified"`
	ComputePlanID  string            `json:"compute_plan_id"`
	Creator        string            `json:"creator"`
	Dataset        *TtDataset        `json:"dataset"`
	Key            string            `json:"key"`
	Log            string            `json:"log"`
	Metadata       map[string]string `json:"metadata"`
	Objective      *TtObjective      `json:"objective"`
	Rank           int               `json:"rank"`
	Status         string            `json:"status"`
	Tag            string            `json:"tag"`
	TraintupleKey  string            `json:"traintuple_key"`
	TraintupleType string            `json:"traintuple_type"`
}

func (out *outputTesttuple) Fill(db *LedgerDB, key string, in Testtuple) error {
	out.Certified = in.Certified
	out.ComputePlanID = in.ComputePlanID
	out.Creator = in.Creator
	out.Dataset = in.Dataset
	out.Key = key
	out.Log = in.Log
	out.Metadata = initMapOutput(in.Metadata)
	out.Rank = in.Rank
	out.Status = in.Status
	out.Tag = in.Tag
	out.TraintupleKey = in.TraintupleKey

	// fill type
	traintupleType, err := db.GetAssetType(in.TraintupleKey)
	if err != nil {
		return errors.Internal("could not retrieve traintuple type with key %s - %s", in.TraintupleKey, err.Error())
	}
	out.TraintupleType = traintupleType.String()

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
	out.Algo = &HashDressNameKey{
		Key:            algo.Key,
		Name:           algo.Name,
		Hash:           algo.Hash,
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
	CompositeTraintuple    *outputCompositeTraintuple `json:"composite_traintuple,omitempty"`
	Traintuple             *outputTraintuple          `json:"traintuple,omitempty"`
	Testtuple              outputTesttuple            `json:"testtuple"`
	NonCertifiedTesttuples []outputTesttuple          `json:"non_certified_testtuples"`
}

type outputModel struct {
	Aggregatetuple      *outputAggregatetuple      `json:"aggregatetuple,omitempty"`
	CompositeTraintuple *outputCompositeTraintuple `json:"composite_traintuple,omitempty"`
	Traintuple          *outputTraintuple          `json:"traintuple,omitempty"`
}

// Event is the collection of tuples sent in an event
type Event struct {
	Testtuples           []outputTesttuple           `json:"testtuple"`
	Traintuples          []outputTraintuple          `json:"traintuple"`
	CompositeTraintuples []outputCompositeTraintuple `json:"composite_traintuple"`
	Aggregatetuples      []outputAggregatetuple      `json:"aggregatetuple"`
	ComputePlans         []eventComputePlan          `json:"compute_plan"`
}

type eventComputePlan struct {
	AlgoKeys       []string `json:"algo_keys"`
	ComputePlanID  string   `json:"compute_plan_id"`
	ModelsToDelete []string `json:"models_to_delete"`
	Status         string   `json:"status"`
}

type outputComputePlan struct {
	ComputePlanID           string            `json:"compute_plan_id"`
	TraintupleKeys          []string          `json:"traintuple_keys"`
	AggregatetupleKeys      []string          `json:"aggregatetuple_keys"`
	CompositeTraintupleKeys []string          `json:"composite_traintuple_keys"`
	TesttupleKeys           []string          `json:"testtuple_keys"`
	CleanModels             bool              `json:"clean_models"`
	Tag                     string            `json:"tag"`
	Metadata                map[string]string `json:"metadata"`
	Status                  string            `json:"status"`
	TupleCount              int               `json:"tuple_count"`
	DoneCount               int               `json:"done_count"`
	IDToKey                 map[string]string `json:"id_to_key"`
}

func (out *outputComputePlan) Fill(key string, in ComputePlan, newIDs []string) {
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
	out.Metadata = initMapOutput(in.Metadata)
	out.TupleCount = in.State.TupleCount
	out.DoneCount = in.State.DoneCount
	IDToKey := map[string]string{}
	for _, ID := range newIDs {
		IDToKey[ID] = in.IDToTrainTask[ID].Key
	}
	out.IDToKey = IDToKey
	out.CleanModels = in.CleanModels
}

type outputPermissions struct {
	Process Permission `json:"process"`
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
	TraintupleKey string         `json:"traintuple_key"`
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
		Hash:           algo.Hash,
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

type outputKey struct {
	Key string `json:"key"`
}

type outputMetrics struct {
	Duration int `json:"duration"`
}
