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

var (
	// OpenPermissions represent struct for default public permissions that could apply to assets
	OpenPermissions = inputPermissions{
		Process: inputPermission{
			Public:        true,
			AuthorizedIDs: []string{},
		},
	}
)

// -------------------------------------------------------------------------------------------
// Struct used to represent inputs for smart contracts. In Hyperledger Fabric, we get as input
// arg  [][]byte or []string, and it is not possible to input a string looking like a json
// -------------------------------------------------------------------------------------------

// inputObjective is the representation of input args to register a Objective
type inputObjective struct {
	Key                       string            `validate:"required,len=32" json:"key"`
	Name                      string            `validate:"required,gte=1,lte=100" json:"name"`
	DescriptionHash           string            `validate:"required,len=64,hexadecimal" json:"description_hash"`
	DescriptionStorageAddress string            `validate:"required,url" json:"description_storage_address"`
	MetricsName               string            `validate:"required,gte=1,lte=100" json:"metrics_name"`
	MetricsHash               string            `validate:"required,len=64,hexadecimal" json:"metrics_hash"`
	MetricsStorageAddress     string            `validate:"required,url" json:"metrics_storage_address"`
	TestDataset               inputDataset      `validate:"omitempty" json:"test_dataset"`
	Permissions               inputPermissions  `validate:"required" json:"permissions"`
	Metadata                  map[string]string `validate:"lte=100,dive,keys,lte=50,endkeys,lte=100" json:"metadata"`
}

// inputDataset is the representation in input args to register a dataset
type inputDataset struct {
	DataManagerKey string   `validate:"omitempty,len=64,hexadecimal" json:"data_manager_key"`
	DataSampleKeys []string `validate:"omitempty,dive,len=64,hexadecimal" json:"data_sample_keys"`
}

// inputAlgo is the representation of input args to register an Algo
type inputAlgo struct {
	Name                      string            `validate:"required,gte=1,lte=100" json:"name"`
	Hash                      string            `validate:"required,len=64,hexadecimal" json:"hash"`
	StorageAddress            string            `validate:"required,url" json:"storage_address"`
	DescriptionHash           string            `validate:"required,len=64,hexadecimal" json:"description_hash"`
	DescriptionStorageAddress string            `validate:"required,url" json:"description_storage_address"`
	Permissions               inputPermissions  `validate:"required" json:"permissions"`
	Metadata                  map[string]string `validate:"lte=100,dive,keys,lte=50,endkeys,lte=100" json:"metadata"`
}

// inputDataManager is the representation of input args to register a DataManager
type inputDataManager struct {
	Name                      string            `validate:"required,gte=1,lte=100" json:"name"`
	OpenerHash                string            `validate:"required,len=64,hexadecimal" json:"opener_hash"`
	OpenerStorageAddress      string            `validate:"required,url" json:"opener_storage_address"`
	Type                      string            `validate:"required,gte=1,lte=30" json:"type"`
	DescriptionHash           string            `validate:"required,len=64,hexadecimal" json:"description_hash"`
	DescriptionStorageAddress string            `validate:"required,url" json:"description_storage_address"`
	ObjectiveKey              string            `validate:"omitempty,len=32,hexadecimal" json:"objective_key"` //`validate:"required"`
	Permissions               inputPermissions  `validate:"required" json:"permissions"`
	Metadata                  map[string]string `validate:"lte=100,dive,keys,lte=50,endkeys,lte=100" json:"metadata"`
}

// inputUpdateDataManager is the representation of input args to update a dataManager with a objective
type inputUpdateDataManager struct {
	DataManagerKey string `validate:"required,len=64,hexadecimal" json:"data_manager_key"`
	ObjectiveKey   string `validate:"required,len=32,hexadecimal" json:"objective_key"`
}

// inputDataSample is the representation of input args to register one or more dataSample
type inputDataSample struct {
	Hashes          []string `validate:"required,dive,len=64,hexadecimal" json:"hashes"`
	DataManagerKeys []string `validate:"omitempty,dive,len=64,hexadecimal" json:"data_manager_keys"`
	TestOnly        string   `validate:"required,oneof=true false" json:"testOnly"`
}

// inputUpdateDataSample is the representation of input args to update one or more dataSample
type inputUpdateDataSample struct {
	Hashes          []string `validate:"required,dive,len=64,hexadecimal" json:"hashes"`
	DataManagerKeys []string `validate:"required,dive,len=64,hexadecimal" json:"data_manager_keys"`
}

// inputTraintuple is the representation of input args to register a Traintuple
type inputTraintuple struct {
	AlgoKey        string            `validate:"required,len=64,hexadecimal" json:"algo_key"`
	InModels       []string          `validate:"omitempty,dive,len=64,hexadecimal" json:"in_models"`
	DataManagerKey string            `validate:"required,len=64,hexadecimal" json:"data_manager_key"`
	DataSampleKeys []string          `validate:"required,unique,gt=0,dive,len=64,hexadecimal" json:"data_sample_keys"`
	ComputePlanID  string            `validate:"omitempty" json:"compute_plan_id"`
	Rank           string            `validate:"omitempty" json:"rank"`
	Tag            string            `validate:"omitempty,lte=64" json:"tag"`
	Metadata       map[string]string `validate:"lte=100,dive,keys,lte=50,endkeys,lte=100" json:"metadata"`
}

// inputTestuple is the representation of input args to register a Testtuple
type inputTesttuple struct {
	DataManagerKey string            `validate:"omitempty,len=64,hexadecimal" json:"data_manager_key"`
	DataSampleKeys []string          `validate:"omitempty,dive,len=64,hexadecimal" json:"data_sample_keys"`
	ObjectiveKey   string            `validate:"required,len=32,hexadecimal" json:"objective_key"`
	Tag            string            `validate:"omitempty,lte=64" json:"tag"`
	Metadata       map[string]string `validate:"omitempty,lte=100,dive,keys,lte=50,endkeys,lte=100" json:"metadata"`
	TraintupleKey  string            `validate:"required,len=64,hexadecimal" json:"traintuple_key"`
}

type inputKeyOld struct {
	Key string `validate:"required,len=64,hexadecimal" json:"key"`
}

type inputKey struct {
	Key string `validate:"required,len=32,hexadecimal" json:"key"`
}

type inputLogSuccessTrain struct {
	inputLog
	OutModel inputHashDress `validate:"required" json:"out_model"`
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

type inputHash struct {
	Hash string `validate:"required,len=64,hexadecimal" json:"hash"`
}

type inputHashDress struct {
	Hash           string `validate:"required,len=64,hexadecimal" json:"hash"`
	StorageAddress string `validate:"required" json:"storage_address"`
}

type inputQueryFilter struct {
	IndexName string `validate:"required" json:"indexName"`
	//TODO : Make Attributes a real list
	Attributes string `validate:"required" json:"attributes"`
}

// inputConputePlan represent a coherent set of tuples uploaded together.
type inputComputePlan struct {
	Traintuples          []inputComputePlanTraintuple          `validate:"omitempty" json:"traintuples"`
	Aggregatetuples      []inputComputePlanAggregatetuple      `validate:"omitempty" json:"aggregatetuples"`
	CompositeTraintuples []inputComputePlanCompositeTraintuple `validate:"omitempty" json:"composite_traintuples"`
	Testtuples           []inputComputePlanTesttuple           `validate:"omitempty" json:"testtuples"`
}

// inputNewComputePlan represent the set of tuples to be added to the compute
// plan matching the ID
type inputNewComputePlan struct {
	CleanModels bool              `json:"clean_models"` // whether or not to delete intermediary models
	Tag         string            `validate:"omitempty,lte=64" json:"tag"`
	Metadata    map[string]string `validate:"lte=100,dive,keys,lte=50,endkeys,lte=100" json:"metadata"`
	inputComputePlan
}

// inputUpdateComputePlan represent the set of tuples to be added to the compute
// plan matching the ID
type inputUpdateComputePlan struct {
	ComputePlanID string `validate:"required,required,len=64,hexadecimal" json:"compute_plan_id"`
	inputComputePlan
}

type inputComputePlanTraintuple struct {
	DataManagerKey string            `validate:"required,len=64,hexadecimal" json:"data_manager_key"`
	DataSampleKeys []string          `validate:"required,dive,len=64,hexadecimal" json:"data_sample_keys"`
	AlgoKey        string            `validate:"required,len=64,hexadecimal" json:"algo_key"`
	ID             string            `validate:"required,lte=64" json:"id"`
	InModelsIDs    []string          `validate:"omitempty,dive,lte=64" json:"in_models_ids"`
	Tag            string            `validate:"omitempty,lte=64" json:"tag"`
	Metadata       map[string]string `validate:"omitempty,lte=100,dive,keys,lte=50,endkeys,lte=100" json:"metadata"`
}

type inputComputePlanAggregatetuple struct {
	AlgoKey     string            `validate:"required,len=64,hexadecimal" json:"algo_key"`
	ID          string            `validate:"required,lte=64" json:"id"`
	InModelsIDs []string          `validate:"omitempty,dive,lte=64" json:"in_models_ids"`
	Tag         string            `validate:"omitempty,lte=64" json:"tag"`
	Metadata    map[string]string `validate:"omitempty,lte=100,dive,keys,lte=50,endkeys,lte=100" json:"metadata"`
	Worker      string            `validate:"required" json:"worker"`
}

type inputComputePlanCompositeTraintuple struct {
	DataManagerKey           string            `validate:"required,len=64,hexadecimal" json:"data_manager_key"`
	DataSampleKeys           []string          `validate:"required,dive,len=64,hexadecimal" json:"data_sample_keys"`
	AlgoKey                  string            `validate:"required,len=64,hexadecimal" json:"algo_key"`
	ID                       string            `validate:"required,lte=64" json:"id"`
	InHeadModelID            string            `validate:"required_with=InTrunkModelID,omitempty,len=64,hexadecimal" json:"in_head_model_id"`
	InTrunkModelID           string            `validate:"required_with=InHeadModelID,omitempty,len=64,hexadecimal" json:"in_trunk_model_id"`
	OutTrunkModelPermissions inputPermissions  `validate:"required" json:"out_trunk_model_permissions"`
	Tag                      string            `validate:"omitempty,lte=64" json:"tag"`
	Metadata                 map[string]string `validate:"omitempty,lte=100,dive,keys,lte=50,endkeys,lte=100" json:"metadata"`
}

type inputComputePlanTesttuple struct {
	DataManagerKey string            `validate:"omitempty,len=64,hexadecimal" json:"data_manager_key"`
	DataSampleKeys []string          `validate:"omitempty,dive,len=64,hexadecimal" json:"data_sample_keys"`
	ObjectiveKey   string            `validate:"required,len=32,hexadecimal" json:"objective_key"`
	Tag            string            `validate:"omitempty,lte=64" json:"tag"`
	Metadata       map[string]string `validate:"omitempty,lte=100,dive,keys,lte=50,endkeys,lte=100" json:"metadata"`
	TraintupleID   string            `validate:"required,lte=64" json:"traintuple_id"`
}

type inputLeaderboard struct {
	ObjectiveKey   string `validate:"omitempty,len=32,hexadecimal" json:"objective_key"`
	AscendingOrder bool   `json:"ascendingOrder,required"`
}

type inputPermissions struct {
	Process inputPermission `validate:"required" json:"process"`
}

type inputPermission struct {
	Public        bool     `json:"public,required"`
	AuthorizedIDs []string `validate:"required" json:"authorized_ids"`
}
