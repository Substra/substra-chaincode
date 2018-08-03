package main

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"gopkg.in/go-playground/validator.v9"
	"strconv"
	"strings"
)

// ---------------------------------------------------------------------------------
// Representation of elements stored in the ledger
// ---------------------------------------------------------------------------------

// Challenge is the representation of one of the element type stored in the ledger
type Challenge struct {
	Name                      string   `json:"name"`
	DescriptionStorageAddress string   `json:"descriptionStorageAddress"`
	Metrics                   *Metrics `json:"metrics"`
	Owner                     string   `json:"owner"`
	TestDataKeys              []string `json:"testDataKeys"`
	Permissions               string   `json:"permissions"`
}

// Dataset is the representation of one of the elements type stored in the ledger
type Dataset struct {
	Name          string     `json:"name"`
	Size          float32    `json:"size"`
	NbData        int        `json:"nbData"`
	Type          string     `json:"type"`
	Description   *HashDress `json:"description"`
	Owner         string     `json:"owner"`
	ChallengeKeys []string   `json:"challengeKeys"`
	Permissions   string     `json:"permissions"`
}

// Data is the representation of one of the element type stored in the ledger
type Data struct {
	DatasetKey string `json:"datasetKey"`
	TestOnly   bool   `json:"testOnly"`
}

// Algo is the representation of one of the element type stored in the ledger
type Algo struct {
	Name           string     `json:"name"`
	StorageAddress string     `json:"storageAddress"`
	Description    *HashDress `json:"description"`
	Owner          string     `json:"owner"`
	ChallengeKey   string     `json:"challengeKey"`
	Permissions    string     `json:"permissions"`
}

// Traintuple is the representation of one the element type stored in the ledger. It describes a training task occuring on the platform
type Traintuple struct {
	Challenge   *TtChallenge `json:"challenge"`
	Algo        *HashDress   `json:"algo"`
	StartModel  *HashDress   `json:"startModel"`
	EndModel    *HashDress   `json:"endModel"`
	TrainData   *TtData      `json:"trainData"`
	TestData    *TtData      `json:"testData"`
	Status      string       `json:"status"`
	Rank        int          `json:"rank"`
	Perf        float32      `json:"perf"`
	Log         string       `json:"log"`
	Permissions string       `json:"permissions"`
	Creator     string       `json:"creator"`
}

// ---------------------------------------------------------------------------------
// Struct used in the representation of elements stored in the ledger
// ---------------------------------------------------------------------------------

// Metrics stores info about a metrics
type Metrics struct {
	Name           string `json:"name"`
	Hash           string `json:"hash"`
	StorageAddress string `json:"storageAddress"`
}

// HashDress stores a hash and a Storage Address
type HashDress struct {
	Hash           string `json:"hash"`
	StorageAddress string `json:"storageAddress"`
}

// TtChallenge stores info about a challenge in a Traintuple
type TtChallenge struct {
	Hash    string     `json:"hash"`
	Metrics *HashDress `json:"metrics"`
}

// TtData stores info about data in a Traintyple (train or test data) and in a PredTuple (later)
type TtData struct {
	Worker     string    `json:"worker"`
	Keys       []string  `json:"keys"`
	OpenerHash string    `json:"openerHash"`
	Perf       []float32 `json:"perf"`
}

// ---------------------------------------------------------------------------------
// Struct used to represent inputs for smart contracts. In Hyperledger Fabric, we
// get as input arg  [][]byte or []string, and it is not possible to input a string
//looking like a json
// ---------------------------------------------------------------------------------

// inputChallenge is the representation of input args to register a Challenge
type inputChallenge struct {
	Name                      string `validate:"required,gte=1,lte=100"`
	DescriptionHash           string `validate:"required,gte=64,lte=64"`
	DescriptionStorageAddress string `validate:"required,url"`
	MetricsName               string `validate:"required,gte=1,lte=100"`
	MetricsHash               string `validate:"required,gte=64,lte=64"`
	MetricsStorageAddress     string `validate:"required,url"`
	TestDataKeys              string `validate:"required"`
	Permissions               string `validate:"required,oneof=all"`
}

// set is a method of the receiver Challenge. It checks the validity of inputChallenge and uses its fields to set the Challenge.
// Returns the challengeKey and the datasetKey associated to test data
func (challenge *Challenge) Set(stub shim.ChaincodeStubInterface, inp inputChallenge) (challengeKey string, datasetKey string, err error) {
	// checking validity of submitted fields
	validate := validator.New()
	err = validate.Struct(inp)
	if err != nil {
		err = fmt.Errorf("invalid inputs %s", err.Error())
		return
	}
	testDataKeys := strings.Split(strings.Replace(inp.TestDataKeys, " ", "", -1), ",")
	datasetKey, err = checkSameDataset(stub, testDataKeys)
	if err != nil {
		err = fmt.Errorf("invalid test data %s", err.Error())
		return
	}
	challenge.Name = inp.Name
	challenge.DescriptionStorageAddress = inp.DescriptionStorageAddress
	challenge.Metrics = &Metrics{
		Name:           inp.MetricsName,
		Hash:           inp.MetricsHash,
		StorageAddress: inp.MetricsStorageAddress,
	}
	owner, err := getTxCreator(stub)
	if err != nil {
		return
	}
	challenge.Owner = owner
	challenge.TestDataKeys = testDataKeys
	challenge.Permissions = inp.Permissions
	// create challenge key
	challengeKey = "challenge_" + inp.DescriptionHash
	return
}

// inputDataset is the representation of input args to register a Dataset
type inputDataset struct {
	Name                      string `validate:"required,gte=1,lte=100"`
	OpenerHash                string `validate:"required,gte=64,lte=64"`
	OpenerStorageAddress      string `validate:"required,url"`
	Type                      string `validate:"required,gte=1,lte=30"`
	DescriptionHash           string `validate:"required,gte=64,lte=64"`
	DescriptionStorageAddress string `validate:"required,url"`
	ChallengeKeys             string //`validate:"required"`
	Permissions               string `validate:"required,oneof=all"`
}

// set is a method of the receiver Dataset. It checks the validity of inputDataset and uses its fields to set the Dataset
// Returns the datasetKey and associated challengeKeys
func (dataset *Dataset) Set(stub shim.ChaincodeStubInterface, inp inputDataset) (string, []string, error) {
	// checking validity of submitted fields
	validate := validator.New()
	err := validate.Struct(inp)
	if err != nil {
		err = fmt.Errorf("invalid inputs %s", err.Error())
		return "", nil, nil
	}
	var challengeKeys []string
	if len(inp.ChallengeKeys) > 0 {
		challengeKeys := strings.Split(strings.Replace(inp.ChallengeKeys, " ", "", -1), ",")
		for _, challengeKey := range challengeKeys {
			_, err = getElementBytes(stub, challengeKey)
			if err != nil {
				err = fmt.Errorf("error checking associated challenge(s) %s", err.Error())
				return "", nil, nil
			}
		}
	} else {
		challengeKeys = nil
	}
	dataset.Name = inp.Name
	dataset.Type = inp.Type
	dataset.Description = &HashDress{
		Hash:           inp.DescriptionHash,
		StorageAddress: inp.DescriptionStorageAddress,
	}
	owner, err := getTxCreator(stub)
	if err != nil {
		return "", nil, err
	}
	dataset.Owner = owner
	dataset.ChallengeKeys = challengeKeys
	dataset.Permissions = inp.Permissions
	// create datasetKey
	datasetKey := "dataset_" + inp.OpenerHash
	return datasetKey, challengeKeys, nil
}

// inputData is the representation of input args to register a Data
type inputData struct {
	Hashes     string `validate:"required"`
	DatasetKey string `validate:"required,gte=72,lte=72"`
	Size       string `validate:"required"`
	TestOnly   string `validate:"required,oneof=true false"`
}

// update is a method of the receiver Dataset. It checks the validity of inputData
// and uses its fields to update the associated dataset
func (dataset *Dataset) Update(stub shim.ChaincodeStubInterface, inp inputData) (datasetKey string, dataHashes []string, testOnly bool, err error) {

	// check if associated dataset exists
	datasetKey = inp.DatasetKey
	err = getElementStruct(stub, datasetKey, &dataset)
	if err != nil {
		return
	}
	// check data size can be converted to float
	size64, err := strconv.ParseFloat(inp.Size, 32)
	if err != nil {
		return
	}
	size := float32(size64)
	// convert input testOnly to boolean
	testOnly, err = strconv.ParseBool(inp.TestOnly)
	if err != nil {
		return
	}

	dataHashes = strings.Split(strings.Replace(inp.Hashes, " ", "", -1), ",")
	dataset.NbData += len(dataHashes)
	dataset.Size += size
	// check validity of dataHashes
	for _, dataHash := range dataHashes {
		// create data key
		if len(dataHash) != 64 {
			err = fmt.Errorf("invalid data hash %s", dataHash)
		}
	}
	return
}

// inputAlgo is the representation of input args to register an Algo
type inputAlgo struct {
	Name                      string `validate:"required,gte=1,lte=100"`
	Hash                      string `validate:"required,gte=64,lte=64"`
	StorageAddress            string `validate:"required,url"`
	DescriptionHash           string `validate:"required,gte=64,lte=64"`
	DescriptionStorageAddress string `validate:"required,url"`
	ChallengeKey              string `validate:"required,gte=74,lte=74"`
	Permissions               string `validate:"required,oneof=all"`
}

// set is a method of the receiver Algo. It checks the validity of inputAlgo and uses its fields to set the Algo
// Returns the algoKey
func (algo *Algo) Set(stub shim.ChaincodeStubInterface, inp inputAlgo) (algoKey string, err error) {
	// checking validity of submitted fields
	validate := validator.New()
	err = validate.Struct(inp)
	if err != nil {
		err = fmt.Errorf("invalid inputs %s", err.Error())
		return
	}
	// check associated challenge exists
	_, err = getElementBytes(stub, inp.ChallengeKey)
	if err != nil {
		return
	}
	// create algo key
	algoKey = "algo_" + inp.Hash
	// find associated owner
	owner, err := getTxCreator(stub)
	if err != nil {
		return
	}
	// set algo
	algo.Name = inp.Name
	algo.StorageAddress = inp.StorageAddress
	algo.Description = &HashDress{
		Hash:           inp.DescriptionHash,
		StorageAddress: inp.DescriptionStorageAddress,
	}
	algo.Owner = owner
	algo.ChallengeKey = inp.ChallengeKey
	algo.Permissions = inp.Permissions
	return
}

// inputTraintuple is the representation of input args to register a Traintuple
type inputTraintuple struct {
	ChallengeKey  string `validate:"required,gte=74,lte=74"`
	AlgoKey       string `validate:"required,gte=69,lte=69"`
	StartModelKey string `validate:"required"`
	TrainDataKeys string `validate:"required"`
}
