package main

// ---------------------------------------------------------------------------------
// Representation of elements stored in the ledger
// ---------------------------------------------------------------------------------

// Objective is the representation of one of the element type stored in the ledger
type Objective struct {
	Name                      string         `json:"name"`
	DescriptionStorageAddress string         `json:"descriptionStorageAddress"`
	Metrics                   *HashDressName `json:"metrics"`
	Owner                     string         `json:"owner"`
	TestData                  *Dataset   `json:"testData"`
	Permissions               string         `json:"permissions"`
}

// DataManager is the representation of one of the elements type stored in the ledger
type DataManager struct {
	Name                 string     `json:"name"`
	OpenerStorageAddress string     `json:"openerStorageAddress"`
	Type                 string     `json:"type"`
	Description          *HashDress `json:"description"`
	Owner                string     `json:"owner"`
	ObjectiveKey         string     `json:"objectiveKey"`
	Permissions          string     `json:"permissions"`
}

// Data is the representation of one of the element type stored in the ledger
type Data struct {
	DataManagerKeys []string `json:"dataManagerKeys"`
	Owner       string   `json:"owner"`
	TestOnly    bool     `json:"testOnly"`
}

// Algo is the representation of one of the element type stored in the ledger
type Algo struct {
	Name           string     `json:"name"`
	StorageAddress string     `json:"storageAddress"`
	Description    *HashDress `json:"description"`
	Owner          string     `json:"owner"`
	ObjectiveKey   string     `json:"objectiveKey"`
	Permissions    string     `json:"permissions"`
}

// Traintuple is the representation of one the element type stored in the ledger. It describes a training task occuring on the platform
type Traintuple struct {
	AlgoKey     string       `json:"algoKey"`
	InModelKeys []string     `json:"inModels"`
	OutModel    *HashDress   `json:"outModel"`
	Data        *Dataset `json:"data"`
	Perf        float32      `json:"perf"`
	FLtask      string       `json:"fltask"`
	Rank        int          `json:"rank"`
	Status      string       `json:"status"`
	Log         string       `json:"log"`
	Permissions string       `json:"permissions"`
	Creator     string       `json:"creator"`
}

// Testtuple is the representation of one the element type stored in the ledger. It describes a training task occuring on the platform
type Testtuple struct {
	Objective   *TtObjective   `json:"objective"`
	Algo        *HashDressName `json:"algo"`
	Model       *Model         `json:"model"`
	Data        *TtData        `json:"data"`
	Certified   bool           `json:"certified"`
	Status      string         `json:"status"`
	Log         string         `json:"log"`
	Permissions string         `json:"permissions"`
	Creator     string         `json:"creator"`
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

// Dataset stores info about a dataManagerKey and a list of associated data
type Dataset struct {
	DataManagerKey string   `json:"dataManagerKey"`
	DataKeys   []string `json:"dataKeys"`
}

// ----------------------------------------------------------------------------------------------
// Representation of output when querying elements if different from what is stored in the ledger
// ----------------------------------------------------------------------------------------------

// outputTraintuple is the representation of one the element type stored in the ledger. It describes a training task occuring on the platform
type outputTraintuple struct {
	Objective   *TtObjective   `json:"objective"`
	Algo        *HashDressName `json:"algo"`
	InModels    []*Model       `json:"inModels"`
	OutModel    *HashDress     `json:"outModel"`
	Data        *TtData        `json:"data"`
	FLtask      string         `json:"fltask"`
	Rank        int            `json:"rank"`
	Status      string         `json:"status"`
	Log         string         `json:"log"`
	Permissions string         `json:"permissions"`
	Creator     string         `json:"creator"`
}

// ---------------------------------------------------------------------------------
// Struct used in the representation of outputs when querying some elements
// ---------------------------------------------------------------------------------

// TtData stores info about data in a Traintyple (train or test data) and in a PredTuple (later)
type TtData struct {
	Worker     string   `json:"worker"`
	Keys       []string `json:"keys"`
	OpenerHash string   `json:"openerHash"`
	Perf       float32  `json:"perf"`
}

// TtObjective stores info about a objective in a Traintuple
type TtObjective struct {
	Key     string     `json:"hash"`
	Metrics *HashDress `json:"metrics"`
}
