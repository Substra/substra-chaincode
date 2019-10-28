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

// outputTraintuple is the representation of one the element type stored in the
// ledger. It describes a training task occuring on the platform
type outputTraintuple struct {
	Key           string            `json:"key"`
	Algo          *HashDressName    `json:"algo"`
	Creator       string            `json:"creator"`
	Dataset       *TtDataset        `json:"dataset"`
	ComputePlanID string            `json:"computePlanID"`
	InModels      []*Model          `json:"inModels"`
	Log           string            `json:"log"`
	Objective     *TtObjective      `json:"objective"`
	OutModel      *HashDress        `json:"outModel"`
	Permissions   outputPermissions `json:"permissions"`
	Rank          int               `json:"rank"`
	Status        string            `json:"status"`
	Tag           string            `json:"tag"`
}

//Fill is a method of the receiver outputTraintuple. It returns all elements necessary to do a training task from a trainuple stored in the ledger
func (outputTraintuple *outputTraintuple) Fill(db LedgerDB, traintuple Traintuple, traintupleKey string) (err error) {

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
		err = fmt.Errorf("could not retrieve algo with key %s - %s", traintuple.AlgoKey, err.Error())
		return
	}
	outputTraintuple.Algo = &HashDressName{
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
	outputTraintuple.Objective = &TtObjective{
		Key:     traintuple.ObjectiveKey,
		Metrics: &metrics,
	}

	// fill inModels
	for _, inModelKey := range traintuple.InModelKeys {
		if inModelKey == "" {
			break
		}
		parentTraintuple, err := db.GetTraintuple(inModelKey)
		if err != nil {
			return fmt.Errorf("could not retrieve parent traintuple with key %s - %s", inModelKey, err.Error())
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
	outputTraintuple.Dataset = &TtDataset{
		Worker:         traintuple.Dataset.Worker,
		DataSampleKeys: traintuple.Dataset.DataSampleKeys,
		OpenerHash:     traintuple.Dataset.DataManagerKey,
		Perf:           traintuple.Perf,
	}

	return
}

type outputTesttuple struct {
	Key       string         `json:"key"`
	Algo      *HashDressName `json:"algo"`
	Certified bool           `json:"certified"`
	Creator   string         `json:"creator"`
	Dataset   *TtDataset     `json:"dataset"`
	Log       string         `json:"log"`
	Model     *Model         `json:"model"`
	Objective *TtObjective   `json:"objective"`
	Status    string         `json:"status"`
	Tag       string         `json:"tag"`
}

func (out *outputTesttuple) Fill(db LedgerDB, key string, in Testtuple) error {
	out.Key = key
	out.Certified = in.Certified
	out.Creator = in.Creator
	out.Dataset = in.Dataset
	out.Log = in.Log
	out.Model = in.Model
	out.Status = in.Status
	out.Tag = in.Tag

	// fill algo
	algo, err := db.GetAlgo(in.AlgoKey)
	if err != nil {
		return fmt.Errorf("could not retrieve algo with key %s - %s", in.AlgoKey, err.Error())
	}
	out.Algo = &HashDressName{
		Name:           algo.Name,
		Hash:           in.AlgoKey,
		StorageAddress: algo.StorageAddress}

	// fill objective
	objective, err := db.GetObjective(in.ObjectiveKey)
	if err != nil {
		return fmt.Errorf("could not retrieve associated objective with key %s- %s", in.ObjectiveKey, err.Error())
	}
	if objective.Metrics == nil {
		return fmt.Errorf("objective %s is missing metrics values", in.ObjectiveKey)
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
	Traintuple             outputTraintuple  `json:"traintuple"`
	Testtuple              outputTesttuple   `json:"testtuple"`
	NonCertifiedTesttuples []outputTesttuple `json:"nonCertifiedTesttuples"`
}

type outputModel struct {
	Traintuple outputTraintuple `json:"traintuple"`
	Testtuple  outputTesttuple  `json:"testtuple"`
}

// TuplesEvent is the collection of tuples sent in an event
type TuplesEvent struct {
	Testtuples           []outputTesttuple           `json:"testtuple"`
	Traintuples          []outputTraintuple          `json:"traintuple"`
	TraintuplesComposite []outputTraintupleComposite `json:"traintuplesComposite"`
}

// SetTraintuplesComposite add one or several testtuples to the event struct
func (te *TuplesEvent) SetTraintuplesComposite(otuples ...outputTraintupleComposite) {
	te.TraintuplesComposite = otuples
}

// SetTesttuples add one or several testtuples to the event struct
func (te *TuplesEvent) SetTesttuples(otuples ...outputTesttuple) {
	te.Testtuples = otuples
}

// SetTraintuples add one or several traintuples to the event struct
func (te *TuplesEvent) SetTraintuples(otuples ...outputTraintuple) {
	te.Traintuples = otuples
}

// AddTraintuple add one traintuple to the event struct
func (te *TuplesEvent) AddTraintuple(out outputTraintuple) {
	te.Traintuples = append(te.Traintuples, out)
}

// AddTesttuple add one testtuple to the event struct
func (te *TuplesEvent) AddTesttuple(out outputTesttuple) {
	te.Testtuples = append(te.Testtuples, out)
}

type outputComputePlan struct {
	ComputePlanID  string   `json:"computePlanID"`
	TraintupleKeys []string `json:"traintupleKeys"`
	TesttupleKeys  []string `json:"testtupleKeys"`
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
	Algo    *HashDressName `json:"algo"`
	Creator string         `json:"creator"`
	Key     string         `json:"key"`
	Model   *Model         `json:"model"`
	Perf    float32        `json:"perf"`
	Tag     string         `json:"tag"`
}

func (out *outputBoardTuple) Fill(db LedgerDB, in Testtuple, testtupleKey string) error {
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
	out.Model = in.Model
	out.Perf = in.Dataset.Perf
	out.Tag = in.Tag
	return nil
}
