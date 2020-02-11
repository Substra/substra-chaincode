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

// ---------------------------------------------------------------------------------
// Representation of elements stored in the ledger
// ---------------------------------------------------------------------------------

// AssetType is use to check the type of an asset
type AssetType uint8

// Const representing the types of asset findable in the ledger
const (
	ObjectiveType AssetType = iota
	DataManagerType
	DataSampleType
	AlgoType
	CompositeAlgoType
	AggregateAlgoType
	TraintupleType
	CompositeTraintupleType
	AggregatetupleType
	TesttupleType
	ComputePlanType
	// when adding a new type here, don't forget to update
	// the String() function in utils.go
)

// StatusUpdater is exported
type StatusUpdater interface {
	commitStatusUpdate(db *LedgerDB, key string, status string) error
}

// Objective is the representation of one of the element type stored in the ledger
type Objective struct {
	Name                      string         `json:"name"`
	AssetType                 AssetType      `json:"assetType"`
	DescriptionStorageAddress string         `json:"descriptionStorageAddress"`
	Metrics                   *HashDressName `json:"metrics"`
	Owner                     string         `json:"owner"`
	TestDataset               *Dataset       `json:"testDataset"`
	Permissions               Permissions    `json:"permissions"`
}

// DataManager is the representation of one of the elements type stored in the ledger
type DataManager struct {
	Name                 string      `json:"name"`
	AssetType            AssetType   `json:"assetType"`
	OpenerStorageAddress string      `json:"openerStorageAddress"`
	Type                 string      `json:"type"`
	Description          *HashDress  `json:"description"`
	Owner                string      `json:"owner"`
	ObjectiveKey         string      `json:"objectiveKey"`
	Permissions          Permissions `json:"permissions"`
}

// DataSample is the representation of one of the element type stored in the ledger
type DataSample struct {
	AssetType       AssetType `json:"assetType"`
	DataManagerKeys []string  `json:"dataManagerKeys"`
	Owner           string    `json:"owner"`
	TestOnly        bool      `json:"testOnly"`
}

// Algo is the representation of one of the element type stored in the ledger
type Algo struct {
	Name           string      `json:"name"`
	AssetType      AssetType   `json:"assetType"`
	StorageAddress string      `json:"storageAddress"`
	Description    *HashDress  `json:"description"`
	Owner          string      `json:"owner"`
	Permissions    Permissions `json:"permissions"`
}

// CompositeAlgo is the representation of one of the element type stored in the ledger
type CompositeAlgo struct {
	Algo
}

// AggregateAlgo is the representation of one of the element type stored in the ledger
type AggregateAlgo struct {
	Algo
}

// GenericTuple is a structure that contains the fields
// that are common to Traintuple, CompositeTraintuple and
// AggregateTuple
type GenericTuple struct {
	AssetType     AssetType `json:"assetType"`
	AlgoKey       string    `json:"algoKey"`
	ComputePlanID string    `json:"computePlanID"`
	Creator       string    `json:"creator"`
	Log           string    `json:"log"`
	Rank          int       `json:"rank"`
	Status        string    `json:"status"`
	Tag           string    `json:"tag"`
}

// Traintuple is the representation of one the element type stored in the ledger. It describes a training task occuring on the platform
type Traintuple struct {
	AssetType     AssetType   `json:"assetType"`
	AlgoKey       string      `json:"algoKey"`
	ComputePlanID string      `json:"computePlanID"`
	Creator       string      `json:"creator"`
	Log           string      `json:"log"`
	Rank          int         `json:"rank"`
	Status        string      `json:"status"`
	Tag           string      `json:"tag"`
	Dataset       *Dataset    `json:"dataset"`
	InModelKeys   []string    `json:"inModels"`
	OutModel      *HashDress  `json:"outModel"`
	Permissions   Permissions `json:"permissions"`
}

// CompositeTraintuple is like a traintuple, but for composite model composition
type CompositeTraintuple struct {
	AssetType     AssetType                       `json:"assetType"`
	AlgoKey       string                          `json:"algoKey"`
	ComputePlanID string                          `json:"computePlanID"`
	Creator       string                          `json:"creator"`
	Log           string                          `json:"log"`
	Rank          int                             `json:"rank"`
	Status        string                          `json:"status"`
	Tag           string                          `json:"tag"`
	Dataset       *Dataset                        `json:"dataset"`
	InHeadModel   string                          `json:"inHeadModel"`
	InTrunkModel  string                          `json:"inTrunkModel"`
	OutHeadModel  CompositeTraintupleOutHeadModel `json:"outHeadModel"`
	OutTrunkModel CompositeTraintupleOutModel     `json:"outTrunkModel"`
}

// Aggregatetuple is like a traintuple, but for aggregate model composition
type Aggregatetuple struct {
	AssetType     AssetType   `json:"assetType"`
	AlgoKey       string      `json:"algoKey"`
	ComputePlanID string      `json:"computePlanID"`
	Creator       string      `json:"creator"`
	Log           string      `json:"log"`
	Rank          int         `json:"rank"`
	Status        string      `json:"status"`
	Tag           string      `json:"tag"`
	InModelKeys   []string    `json:"inModels"`
	OutModel      *HashDress  `json:"outModel"`
	Permissions   Permissions `json:"permissions"` // TODO (aggregate): what do permissions mean here?
	Worker        string      `json:"worker"`
}

// CompositeTraintupleOutModel is the out-model of a CompositeTraintuple
type CompositeTraintupleOutModel struct {
	OutModel    *HashDress  `json:"outModel"`
	Permissions Permissions `json:"permissions"`
}

// CompositeTraintupleOutHeadModel is the out-model of a CompositeTraintuple
type CompositeTraintupleOutHeadModel struct {
	OutModel    *Hash       `json:"outModel"`
	Permissions Permissions `json:"permissions"`
}

// Testtuple is the representation of one the element type stored in the ledger. It describes a training task occuring on the platform
type Testtuple struct {
	AlgoKey       string      `json:"algo"`
	AssetType     AssetType   `json:"assetType"`
	Certified     bool        `json:"certified"`
	ComputePlanID string      `json:"computePlanID"`
	Creator       string      `json:"creator"`
	Dataset       *TtDataset  `json:"dataset"`
	Log           string      `json:"log"`
	TraintupleKey string      `json:"traintupleKey"`
	ObjectiveKey  string      `json:"objective"`
	Permissions   Permissions `json:"permissions"`
	Rank          int         `json:"rank"`
	Status        string      `json:"status"`
	Tag           string      `json:"tag"`
}

// ComputePlan is the ledger's representation of a compute plan.
type ComputePlan struct {
	AggregatetupleKeys      []string             `json:"aggregatetupleKeys"`
	AssetType               AssetType            `json:"assetType"`
	CompositeTraintupleKeys []string             `json:"compositeTraintupleKeys"`
	DoneCount               int                  `json:"DoneCount"`
	IDToTrainTask           map[string]TrainTask `json:"IDToTrainTask"`
	Status                  string               `json:"status"`
	Tag                     string               `json:"tag"`
	TesttupleKeys           []string             `json:"testtupleKeys"`
	TraintupleKeys          []string             `json:"traintupleKeys"`
	TupleCount              int                  `json:"tupleCount"`
}

// TrainTask is represent the information for one tuple in a Compute Plan
type TrainTask struct {
	Depth int    `json:"depth"`
	Key   string `json:"key"`
}

// ---------------------------------------------------------------------------------
// Struct used in the representation of elements stored in the ledger
// ---------------------------------------------------------------------------------

// Hash stores a hash
type Hash struct {
	Hash string `json:"hash"`
}

// HashDress stores a hash and a Storage Address
type HashDress struct {
	Hash           string `json:"hash"`
	StorageAddress string `json:"storageAddress"`
}

// HashDressName stores a hash, storage address and a name
type HashDressName struct {
	Name           string `json:"name"`
	Hash           string `json:"hash"`
	StorageAddress string `json:"storageAddress"`
}

// Model stores the traintupleKey leading to the model, its hash and storage addressl
type Model struct {
	TraintupleKey  string `json:"traintupleKey"`
	Hash           string `json:"hash"`
	StorageAddress string `json:"storageAddress"`
}

// Dataset stores info about a dataManagerKey and a list of associated dataSample
type Dataset struct {
	DataManagerKey string   `json:"dataManagerKey"`
	DataSampleKeys []string `json:"dataSampleKeys"`
	Worker         string   `json:"worker"`
}

// ---------------------------------------------------------------------------------
// Struct used in the representation of outputs when querying some elements
// ---------------------------------------------------------------------------------

// TtDataset stores info about dataset in a Traintyple (train or test data) and in a PredTuple (later)
type TtDataset struct {
	Worker         string   `json:"worker"`
	DataSampleKeys []string `json:"keys"`
	OpenerHash     string   `json:"openerHash"`
	Perf           float32  `json:"perf"`
}

// TtObjective stores info about a objective in a Traintuple
type TtObjective struct {
	Key     string     `json:"hash"`
	Metrics *HashDress `json:"metrics"`
}

// Node stores informations about node registered into the network,
// would be used to list authorized nodes for permissions
type Node struct {
	ID string `json:"id"`
}
