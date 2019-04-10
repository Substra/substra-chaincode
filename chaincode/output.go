package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// Struct use as output representation of ledger data

type outputObjective struct {
	Key         string         `json:"key"`
	Name        string         `json:"name"`
	Description HashDress      `json:"description"`
	Metrics     *HashDressName `json:"metrics"`
	Owner       string         `json:"owner"`
	TestDataset *Dataset       `json:"testDataset"`
	Permissions string         `json:"permissions"`
}

func (out *outputObjective) Fill(key string, in Objective) {
	out.Key = key
	out.Name = in.Name
	out.Description.StorageAddress = in.DescriptionStorageAddress
	out.Description.Hash = key
	out.Metrics = in.Metrics
	out.Owner = in.Owner
	out.TestDataset = in.TestDataset
	out.Permissions = in.Permissions
}

// outputDataManager is the return representation of the DataManager type stored in the ledger
type outputDataManager struct {
	ObjectiveKey string     `json:"objectiveKey"`
	Description  *HashDress `json:"description"`
	Key          string     `json:"key"`
	Name         string     `json:"name"`
	Opener       HashDress  `json:"opener"`
	Owner        string     `json:"owner"`
	Permissions  string     `json:"permissions"`
	Type         string     `json:"type"`
}

func (out *outputDataManager) Fill(key string, in DataManager) {
	out.ObjectiveKey = in.ObjectiveKey
	out.Description = in.Description
	out.Key = key
	out.Name = in.Name
	out.Opener.Hash = key
	out.Opener.StorageAddress = in.OpenerStorageAddress
	out.Owner = in.Owner
	out.Permissions = in.Permissions
	out.Type = in.Type
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
	Key          string     `json:"key"`
	Name         string     `json:"name"`
	Content      HashDress  `json:"content"`
	Description  *HashDress `json:"description"`
	Owner        string     `json:"owner"`
	ObjectiveKey string     `json:"objectiveKey"`
	Permissions  string     `json:"permissions"`
}

func (out *outputAlgo) Fill(key string, in Algo) {
	out.Key = key
	out.Name = in.Name
	out.Content.Hash = key
	out.Content.StorageAddress = in.StorageAddress
	out.Description = in.Description
	out.Owner = in.Owner
	out.ObjectiveKey = in.ObjectiveKey
	out.Permissions = in.Permissions
}

// outputTraintuple is the representation of one the element type stored in the
// ledger. It describes a training task occuring on the platform
type outputTraintuple struct {
	Objective   *TtObjective   `json:"objective"`
	Algo        *HashDressName `json:"algo"`
	InModels    []*Model       `json:"inModels"`
	OutModel    *HashDress     `json:"outModel"`
	Dataset     *TtDataset     `json:"dataset"`
	FLtask      string         `json:"fltask"`
	Rank        int            `json:"rank"`
	Status      string         `json:"status"`
	Log         string         `json:"log"`
	Permissions string         `json:"permissions"`
	Creator     string         `json:"creator"`
}

//Fill is a method of the receiver outputTraintuple. It returns all elements necessary to do a training task from a trainuple stored in the ledger
func (outputTraintuple *outputTraintuple) Fill(stub shim.ChaincodeStubInterface, traintuple Traintuple) (err error) {

	outputTraintuple.Creator = traintuple.Creator
	outputTraintuple.Permissions = traintuple.Permissions
	outputTraintuple.Log = traintuple.Log
	outputTraintuple.Status = traintuple.Status
	outputTraintuple.Rank = traintuple.Rank
	outputTraintuple.FLtask = traintuple.FLtask
	outputTraintuple.OutModel = traintuple.OutModel
	// fill algo
	algo := Algo{}
	if err = getElementStruct(stub, traintuple.AlgoKey, &algo); err != nil {
		err = fmt.Errorf("could not retrieve algo with key %s - %s", traintuple.AlgoKey, err.Error())
		return
	}
	outputTraintuple.Algo = &HashDressName{
		Name:           algo.Name,
		Hash:           traintuple.AlgoKey,
		StorageAddress: algo.StorageAddress}

	// fill objective
	objective := Objective{}
	if err = getElementStruct(stub, algo.ObjectiveKey, &objective); err != nil {
		err = fmt.Errorf("could not retrieve associated objective with key %s- %s", algo.ObjectiveKey, err.Error())
		return
	}
	if objective.Metrics == nil {
		err = fmt.Errorf("objective %s is missing metrics values", algo.ObjectiveKey)
		return
	}
	metrics := HashDress{
		Hash:           objective.Metrics.Hash,
		StorageAddress: objective.Metrics.StorageAddress,
	}
	outputTraintuple.Objective = &TtObjective{
		Key:     algo.ObjectiveKey,
		Metrics: &metrics,
	}

	// fill inModels
	for _, inModelKey := range traintuple.InModelKeys {
		if inModelKey == "" {
			break
		}
		parentTraintuple := Traintuple{}
		if err = getElementStruct(stub, inModelKey, &parentTraintuple); err != nil {
			err = fmt.Errorf("could not retrieve parent traintuple with key %s - %s", inModelKey, err.Error())
			return
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

	// fill dataset from dataManager and dataSample
	dataManager := DataManager{}
	if err = getElementStruct(stub, traintuple.Dataset.DataManagerKey, &dataManager); err != nil {
		err = fmt.Errorf("could not retrieve dataManager with key %s - %s", traintuple.Dataset.DataManagerKey, err.Error())
		return
	}
	outputTraintuple.Dataset = &TtDataset{
		Worker:         dataManager.Owner,
		DataSampleKeys: traintuple.Dataset.DataSampleKeys,
		OpenerHash:     traintuple.Dataset.DataManagerKey,
		Perf:           traintuple.Perf,
	}

	return
}

type outputTesttuple struct {
	Key string `json:"key"`
	Testtuple
}

func (out *outputTesttuple) Fill(key string, in Testtuple) {
	out.Testtuple = in
	out.Key = key
}
