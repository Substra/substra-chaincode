package main

// -------------------------------------------------------------------------------------------
// Struct used to represent inputs for smart contracts. In Hyperledger Fabric, we get as input
// arg  [][]byte or []string, and it is not possible to input a string looking like a json
// -------------------------------------------------------------------------------------------

// inputObjective is the representation of input args to register a Objective
type inputObjective struct {
	Name                      string `validate:"required,gte=1,lte=100" json:"name"`
	DescriptionHash           string `validate:"required,len=64,hexadecimal" json:"descriptionHash"`
	DescriptionStorageAddress string `validate:"required,url" json:"descriptionStorageAddress"`
	MetricsName               string `validate:"required,gte=1,lte=100" json:"metricsName"`
	MetricsHash               string `validate:"required,len=64,hexadecimal" json:"metricsHash"`
	MetricsStorageAddress     string `validate:"required,url" json:"metricsStorageAddress"`
	TestDataset               string `validate:"required" json:"testDataset"`
	Permissions               string `validate:"required,oneof=all" json:"permissions"`
}

// inputAlgo is the representation of input args to register an Algo
type inputAlgo struct {
	Name                      string `validate:"required,gte=1,lte=100" json:"name"`
	Hash                      string `validate:"required,len=64,hexadecimal" json:"hash"`
	StorageAddress            string `validate:"required,url" json:"storageAddress"`
	DescriptionHash           string `validate:"required,len=64,hexadecimal" json:"descriptionHash"`
	DescriptionStorageAddress string `validate:"required,url" json:"descriptionStorageAddress"`
	Permissions               string `validate:"required,oneof=all" json:"permissions"`
}

// inputDataManager is the representation of input args to register a DataManager
type inputDataManager struct {
	Name                      string `validate:"required,gte=1,lte=100" json:"name"`
	OpenerHash                string `validate:"required,len=64,hexadecimal" json:"openerHash"`
	OpenerStorageAddress      string `validate:"required,url" json:"openerStorageAddress"`
	Type                      string `validate:"required,gte=1,lte=30" json:"type"`
	DescriptionHash           string `validate:"required,len=64,hexadecimal" json:"descriptionHash"`
	DescriptionStorageAddress string `validate:"required,url" json:"descriptionStorageAddress"`
	ObjectiveKey              string `validate:"omitempty" json:"objectiveKey"` //`validate:"required"`
	Permissions               string `validate:"required,oneof=all" json:"permissions"`
}

// inputUpdateDataManager is the representation of input args to update a dataManager with a objective
type inputUpdateDataManager struct {
	DataManagerKey string `validate:"required,len=64,hexadecimal" json:"dataManagerKey"`
	ObjectiveKey   string `validate:"required,len=64,hexadecimal" json:"objectiveKey"`
}

// inputDataSample is the representation of input args to register one or more dataSample
type inputDataSample struct {
	Hashes          string   `validate:"required" json:"hashes"`
	DataManagerKeys []string `validate:"omitempty,dive,len=64,hexadecimal" json:"dataManagerKeys"`
	TestOnly        string   `validate:"required,oneof=true false" json:"testOnly"`
}

// inputUpdateDataSample is the representation of input args to update one or more dataSample
type inputUpdateDataSample struct {
	Hashes          string   `validate:"required" json:"hashes"`
	DataManagerKeys []string `validate:"required,dive,len=64,hexadecimal" json:"dataManagerKeys"`
}

// inputTraintuple is the representation of input args to register a Traintuple
type inputTraintuple struct {
	AlgoKey        string   `validate:"required,len=64,hexadecimal" json:"algoKey"`
	ObjectiveKey   string   `validate:"required,len=64,hexadecimal" json:"objectiveKey"`
	InModels       []string `validate:"omitempty,dive,len=64,hexadecimal" json:"inModels"`
	DataManagerKey string   `validate:"required,len=64,hexadecimal" json:"dataManagerKey"`
	DataSampleKeys []string `validate:"required,dive,len=64,hexadecimal" json:"dataSampleKeys"`
	FLTask         string   `validate:"omitempty" json:"flTask"`
	Rank           string   `validate:"omitempty" json:"rank"`
	Tag            string   `validate:"omitempty,lte=64" json:"tag"`
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
	Perf     float32        `validate:"required" json:"perf"`
}
type inputLogSuccessTest struct {
	inputLog
	Perf float32 `validate:"required" json:"perf"`
}
type inputLogFailTrain struct {
	inputLog
}
type inputLogFailTest struct {
	inputLog
}
type inputLog struct {
	Key string `validate:"required,len=64,hexadecimal" json:"key"`
	Log string `validate:"required,lte=200" json:"log"`
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
type inputComputePlan struct {
	AlgoKey      string            `validate:"required,len=64,hexadecimal" json:"algoKey"`
	ObjectiveKey string            `validate:"required,len=64,hexadecimal" json:"objectiveKey"`
	Traintuples  []inputSubstruple `validate:"required,gt=0" json:"traintuples"`
}

type inputSubstruple struct {
	DataManagerKey string   `validate:"required,len=64,hexadecimal" json:"dataManagerKey"`
	DataSampleKeys []string `validate:"required,dive,len=64,hexadecimal" json:"dataSampleKeys"`
	InModelsUUID   []string `validate:"omitempty,dive,lte=64" json:"inModels"`
	Tag            string   `validate:"omitempty,lte=64" json:"tag"`
	UUID           string   `validate:"required,lte=64" json:"uuid"`
}
