package main

// ---------------------------------------------------------------------------------
// Representation of elements stored in the ledger
// ---------------------------------------------------------------------------------
type AssetType uint8

const (
	ObjectiveType AssetType = iota
	DataManagerType
	DataSampleType
	AlgoType
	TraintupleType
	TesttupleType
)

// Objective is the representation of one of the element type stored in the ledger
type Objective struct {
	Name                      string         `json:"name"`
	AssetType                 AssetType      `json:"assetType"`
	DescriptionStorageAddress string         `json:"descriptionStorageAddress"`
	Metrics                   *HashDressName `json:"metrics"`
	Owner                     string         `json:"owner"`
	TestDataset               *Dataset       `json:"testDataset"`
	Permissions               string         `json:"permissions"`
}

// DataManager is the representation of one of the elements type stored in the ledger
type DataManager struct {
	Name                 string     `json:"name"`
	AssetType            AssetType  `json:"assetType"`
	OpenerStorageAddress string     `json:"openerStorageAddress"`
	Type                 string     `json:"type"`
	Description          *HashDress `json:"description"`
	Owner                string     `json:"owner"`
	ObjectiveKey         string     `json:"objectiveKey"`
	Permissions          string     `json:"permissions"`
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
	Name           string     `json:"name"`
	AssetType      AssetType  `json:"assetType"`
	StorageAddress string     `json:"storageAddress"`
	Description    *HashDress `json:"description"`
	Owner          string     `json:"owner"`
	Permissions    string     `json:"permissions"`
}

// Traintuple is the representation of one the element type stored in the ledger. It describes a training task occuring on the platform
type Traintuple struct {
	AssetType    AssetType  `json:"assetType"`
	AlgoKey      string     `json:"algoKey"`
	Creator      string     `json:"creator"`
	Dataset      *Dataset   `json:"dataset"`
	FLTask       string     `json:"fltask"`
	InModelKeys  []string   `json:"inModels"`
	Log          string     `json:"log"`
	ObjectiveKey string     `json:"objectiveKey"`
	OutModel     *HashDress `json:"outModel"`
	Perf         float32    `json:"perf"`
	Permissions  string     `json:"permissions"`
	Rank         int        `json:"rank"`
	Status       string     `json:"status"`
	Tag          string     `json:"tag"`
}

// Testtuple is the representation of one the element type stored in the ledger. It describes a training task occuring on the platform
type Testtuple struct {
	AssetType   AssetType      `json:"assetType"`
	Algo        *HashDressName `json:"algo"`
	Certified   bool           `json:"certified"`
	Creator     string         `json:"creator"`
	Dataset     *TtDataset     `json:"dataset"`
	Log         string         `json:"log"`
	Model       *Model         `json:"model"`
	Objective   *TtObjective   `json:"objective"`
	Permissions string         `json:"permissions"`
	Status      string         `json:"status"`
	Tag         string         `json:"tag"`
}

// ---------------------------------------------------------------------------------
// Struct used in the representation of elements stored in the ledger
// ---------------------------------------------------------------------------------

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
	Worker         string   `json:"-"`
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
