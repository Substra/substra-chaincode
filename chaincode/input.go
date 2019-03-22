package main

// -------------------------------------------------------------------------------------------
// Struct used to represent inputs for smart contracts. In Hyperledger Fabric, we get as input
// arg  [][]byte or []string, and it is not possible to input a string looking like a json
// -------------------------------------------------------------------------------------------

// inputObjective is the representation of input args to register a Objective
type inputObjective struct {
	Name                      string `validate:"required,gte=1,lte=100"`
	DescriptionHash           string `validate:"required,gte=64,lte=64,hexadecimal"`
	DescriptionStorageAddress string `validate:"required,url"`
	MetricsName               string `validate:"required,gte=1,lte=100"`
	MetricsHash               string `validate:"required,gte=64,lte=64,hexadecimal"`
	MetricsStorageAddress     string `validate:"required,url"`
	TestData                  string `validate:"required"`
	Permissions               string `validate:"required,oneof=all"`
}

// inputAlgo is the representation of input args to register an Algo
type inputAlgo struct {
	Name                      string `validate:"required,gte=1,lte=100"`
	Hash                      string `validate:"required,gte=64,lte=64,hexadecimal"`
	StorageAddress            string `validate:"required,url"`
	DescriptionHash           string `validate:"required,gte=64,lte=64,hexadecimal"`
	DescriptionStorageAddress string `validate:"required,url"`
	ObjectiveKey              string `validate:"required,gte=64,lte=64,hexadecimal"`
	Permissions               string `validate:"required,oneof=all"`
}

// inputDataManager is the representation of input args to register a DataManager
type inputDataManager struct {
	Name                      string `validate:"required,gte=1,lte=100"`
	OpenerHash                string `validate:"required,gte=64,lte=64,hexadecimal"`
	OpenerStorageAddress      string `validate:"required,url"`
	Type                      string `validate:"required,gte=1,lte=30"`
	DescriptionHash           string `validate:"required,gte=64,lte=64,hexadecimal"`
	DescriptionStorageAddress string `validate:"required,url"`
	ObjectiveKey              string //`validate:"required"`
	Permissions               string `validate:"required,oneof=all"`
}

// inputUpdateDataManager is the representation of input args to update a dataManager with a objective
type inputUpdateDataManager struct {
	DataManagerKey   string `validate:"required,gte=64,lte=64,hexadecimal"`
	ObjectiveKey string `validate:"required,gte=64,lte=64,hexadecimal"`
}

// inputData is the representation of input args to register one or more data
type inputData struct {
	Hashes      string `validate:"required"`
	DataManagerKeys string
	TestOnly    string `validate:"required,oneof=true false"`
}

// inputUpdateData is the representation of input args to update one or more data
type inputUpdateData struct {
	Hashes      string `validate:"required"`
	DataManagerKeys string `validate:"required"`
}

// inputTraintuple is the representation of input args to register a Traintuple
type inputTraintuple struct {
	AlgoKey    string `validate:"required,gte=64,lte=64,hexadecimal"`
	InModels   string //`validate:"omitEmpty"
	DataManagerKey string `validate:"required,gte=64,lte=64,hexadecimal"`
	DataKeys   string `validate:"required"`
	FLtask     string //`validate:"omitEmpty"`
	Rank       string //`validate:"omitEmpty"`
}

// inputTestuple is the representation of input args to register a Testtuple
type inputTesttuple struct {
	TraintupleKey string `validate:"required,gte=64,lte=64,hexadecimal"`
	DataManagerKey    string `validate:"omitempty,gte=64,lte=64,hexadecimal"`
	DataKeys      string
}
