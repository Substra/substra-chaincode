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

// -------------------------------------------------------------------------------------------
// Struct used to represent inputs for smart contracts. In Hyperledger Fabric, we get as input
// arg  [][]byte or []string, and it is not possible to input a string looking like a json
// -------------------------------------------------------------------------------------------

// inputObjective is the representation of input args to register a Objective
type inputObjective struct {
	Name                      string           `validate:"required,gte=1,lte=100" json:"name"`
	DescriptionHash           string           `validate:"required,len=64,hexadecimal" json:"descriptionHash"`
	DescriptionStorageAddress string           `validate:"required,url" json:"descriptionStorageAddress"`
	MetricsName               string           `validate:"required,gte=1,lte=100" json:"metricsName"`
	MetricsHash               string           `validate:"required,len=64,hexadecimal" json:"metricsHash"`
	MetricsStorageAddress     string           `validate:"required,url" json:"metricsStorageAddress"`
	TestDataset               inputDataset     `validate:"omitempty" json:"testDataset"`
	Permissions               inputPermissions `validate:"required" json:"permissions"`
}

// inputDataset is the representation in input args to register a dataset
type inputDataset struct {
	DataManagerKey string   `validate:"omitempty,len=64,hexadecimal" json:"dataManagerKey"`
	DataSampleKeys []string `validate:"omitempty,dive,len=64,hexadecimal" json:"dataSampleKeys"`
}

// inputAlgo is the representation of input args to register an Algo
type inputAlgo struct {
	Name                      string           `validate:"required,gte=1,lte=100" json:"name"`
	Hash                      string           `validate:"required,len=64,hexadecimal" json:"hash"`
	StorageAddress            string           `validate:"required,url" json:"storageAddress"`
	DescriptionHash           string           `validate:"required,len=64,hexadecimal" json:"descriptionHash"`
	DescriptionStorageAddress string           `validate:"required,url" json:"descriptionStorageAddress"`
	Permissions               inputPermissions `validate:"required" json:"permissions"`
}

// inputDataManager is the representation of input args to register a DataManager
type inputDataManager struct {
	Name                      string           `validate:"required,gte=1,lte=100" json:"name"`
	OpenerHash                string           `validate:"required,len=64,hexadecimal" json:"openerHash"`
	OpenerStorageAddress      string           `validate:"required,url" json:"openerStorageAddress"`
	Type                      string           `validate:"required,gte=1,lte=30" json:"type"`
	DescriptionHash           string           `validate:"required,len=64,hexadecimal" json:"descriptionHash"`
	DescriptionStorageAddress string           `validate:"required,url" json:"descriptionStorageAddress"`
	ObjectiveKey              string           `validate:"omitempty" json:"objectiveKey"` //`validate:"required"`
	Permissions               inputPermissions `validate:"required" json:"permissions"`
}

// inputUpdateDataManager is the representation of input args to update a dataManager with a objective
type inputUpdateDataManager struct {
	DataManagerKey string `validate:"required,len=64,hexadecimal" json:"dataManagerKey"`
	ObjectiveKey   string `validate:"required,len=64,hexadecimal" json:"objectiveKey"`
}

// inputDataSample is the representation of input args to register one or more dataSample
type inputDataSample struct {
	Hashes          []string `validate:"required,dive,len=64,hexadecimal" json:"hashes"`
	DataManagerKeys []string `validate:"omitempty,dive,len=64,hexadecimal" json:"dataManagerKeys"`
	TestOnly        string   `validate:"required,oneof=true false" json:"testOnly"`
}

// inputUpdateDataSample is the representation of input args to update one or more dataSample
type inputUpdateDataSample struct {
	Hashes          []string `validate:"required,dive,len=64,hexadecimal" json:"hashes"`
	DataManagerKeys []string `validate:"required,dive,len=64,hexadecimal" json:"dataManagerKeys"`
}

// inputTraintuple is the representation of input args to register a Traintuple
type inputTraintuple struct {
	AlgoKey        string   `validate:"required,len=64,hexadecimal" json:"algoKey"`
	ObjectiveKey   string   `validate:"required,len=64,hexadecimal" json:"objectiveKey"`
	InModels       []string `validate:"omitempty,dive,len=64,hexadecimal" json:"inModels"`
	DataManagerKey string   `validate:"required,len=64,hexadecimal" json:"dataManagerKey"`
	DataSampleKeys []string `validate:"required,unique,gt=0,dive,len=64,hexadecimal" json:"dataSampleKeys"`
	ComputePlanID  string   `validate:"omitempty" json:"computePlanID"`
	Rank           string   `validate:"omitempty" json:"rank"`
	Tag            string   `validate:"omitempty,lte=64" json:"tag"`
}

// inputTraintupleComposite is the representation of input args to register a Traintuple
type inputTraintupleComposite struct {
	AlgoKey        string              `validate:"required,len=64,hexadecimal" json:"algoKey"`
	ObjectiveKey   string              `validate:"required,len=64,hexadecimal" json:"objectiveKey"`
	InModelHead    inputModelComposite `validate:"omitempty" json:"inModelHead"`
	InModelTrunk   inputModelComposite `validate:"omitempty" json:"inModelTrunk"`
	DataManagerKey string              `validate:"required,len=64,hexadecimal" json:"dataManagerKey"`
	DataSampleKeys []string            `validate:"required,unique,gt=0,dive,len=64,hexadecimal" json:"dataSampleKeys"`
	ComputePlanID  string              `validate:"omitempty" json:"computePlanID"`
	Rank           string              `validate:"omitempty" json:"rank"`
	Tag            string              `validate:"omitempty,lte=64" json:"tag"`
}

type inputModelComposite struct {
	Key            string           `validate:"required,len=64,hexadecimal" json:"key"`
	OutPermissions inputPermissions `validate:"required" json:"outPermissions"`
}

// inputTestuple is the representation of input args to register a Testtuple
type inputTesttuple struct {
	TraintupleKey  string   `validate:"required,len=64,hexadecimal" json:"traintupleKey"`
	DataManagerKey string   `validate:"omitempty,len=64,hexadecimal" json:"dataManagerKey"`
	DataSampleKeys []string `validate:"omitempty,dive,len=64,hexadecimal" json:"dataSampleKeys"`
	Tag            string   `validate:"omitempty,lte=64" json:"tag"`
}

type inputHash struct {
	Key string `validate:"required,len=64,hexadecimal" json:"key"`
}

type inputLogSuccessTrain struct {
	inputLog
	OutModel inputHashDress `validate:"required" json:"outModel"`
	Perf     float32        `validate:"omitempty" json:"perf"`
}
type inputLogSuccessTrainComposite struct {
	inputLog
	OutModelHead  inputHashDress `validate:"required" json:"outModelHead"`
	OutModelTrunk inputHashDress `validate:"required" json:"outModelTrunk"`
	Perf          float32        `validate:"omitempty" json:"perf"`
}
type inputLogSuccessTest struct {
	inputLog
	Perf float32 `validate:"omitempty" json:"perf"`
}
type inputLogFailTrain struct {
	inputLog
}
type inputLogFailTest struct {
	inputLog
}
type inputLog struct {
	Key string `validate:"required,len=64,hexadecimal" json:"key"`
	Log string `validate:"lte=200" json:"log"`
}

type inputHashDress struct {
	Hash           string `validate:"required,len=64,hexadecimal" json:"hash"`
	StorageAddress string `validate:"required" json:"storageAddress"`
}

type inputQueryFilter struct {
	IndexName string `validate:"required" json:"indexName"`
	//TODO : Make Attributes a real list
	Attributes string `validate:"required" json:"attributes"`
}

// inputConputePlan represent a coherent set of tuples uploaded together.
// They share the same Algo and Objective represented by their respective keys.
// Traintuples is the list of all the traintuples planed by the compute plan
// Beware, it's order sensitive since the `InModelsIDs` can only be interpreted
// if the traintuples matching those IDs have already been created.
type inputComputePlan struct {
	AlgoKey      string                       `validate:"required,len=64,hexadecimal" json:"algoKey"`
	ObjectiveKey string                       `validate:"required,len=64,hexadecimal" json:"objectiveKey"`
	Traintuples  []inputComputePlanTraintuple `validate:"required,gt=0" json:"traintuples"`
	Testtuples   []inputComputePlanTesttuple  `validate:"omitempty" json:"testtuples"`
}

type inputComputePlanTraintuple struct {
	DataManagerKey string   `validate:"required,len=64,hexadecimal" json:"dataManagerKey"`
	DataSampleKeys []string `validate:"required,dive,len=64,hexadecimal" json:"dataSampleKeys"`
	ID             string   `validate:"required,lte=64" json:"id"`
	InModelsIDs    []string `validate:"omitempty,dive,lte=64" json:"inModelsIDs"`
	Tag            string   `validate:"omitempty,lte=64" json:"tag"`
}

type inputComputePlanTesttuple struct {
	DataManagerKey string   `validate:"omitempty,len=64,hexadecimal" json:"dataManagerKey"`
	DataSampleKeys []string `validate:"omitempty,dive,len=64,hexadecimal" json:"dataSampleKeys"`
	Tag            string   `validate:"omitempty,lte=64" json:"tag"`
	TraintupleID   string   `validate:"required,lte=64" json:"traintupleID"`
}

type inputLeaderboard struct {
	ObjectiveKey   string `validate:"omitempty,len=64,hexadecimal" json:"objectiveKey"`
	AscendingOrder bool   `json:"ascendingOrder,required"`
}

type inputPermissions struct {
	Process inputPermission `validate:"required" json:"process"`
}

type inputPermission struct {
	Public        bool     `json:"public,required"`
	AuthorizedIDs []string `validate:"required" json:"authorizedIDs"`
}
