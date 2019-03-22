package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"gopkg.in/go-playground/validator.v9"
)

const lenKey int = 64

// -------------------------------------------------------------------------------------------
// Methods on receivers traintuple and testuples
// -------------------------------------------------------------------------------------------

// Set is a method of the receiver Traintuple. It checks the validity of inputTraintuple and uses its fields to set the Traintuple
func (traintuple *Traintuple) Set(stub shim.ChaincodeStubInterface, inp inputTraintuple) (traintupleKey string, err error) {

	validate := validator.New()
	if err = validate.Struct(inp); err != nil {
		err = fmt.Errorf("invalid inputs to update data %s", err.Error())
		return
	}

	// TODO later: check permissions
	// find associated creator and check permissions (TODO later)
	creator, err := getTxCreator(stub)
	if err != nil {
		return
	}
	traintuple.Creator = creator
	traintuple.Permissions = "all"

	// WARNING FOR NOW NO CHECK ABOUT RANK AND FLtask
	// The FL task should be returned by the ledger now TODO
	// and the rank should be checked also...
	if (inp.Rank != "") && (inp.FLtask != "") {
		rank, _ := strconv.Atoi(inp.Rank)
		traintuple.Rank = rank
		traintuple.FLtask = inp.FLtask
	}

	// check if algo exists
	if _, err = getElementBytes(stub, inp.AlgoKey); err != nil {
		err = fmt.Errorf("could not retrieve algo with key %s - %s", inp.AlgoKey, err.Error())
		return
	}
	traintuple.AlgoKey = inp.AlgoKey

	// create key: hash of algo + start model + train dataset + creator (keys)
	// certainly not be the most efficient key... but let's make it work and them try to make it better...
	tKey := sha256.Sum256([]byte(inp.AlgoKey + inp.InModels + inp.DataSampleKeys + creator))
	traintupleKey = hex.EncodeToString(tKey[:])
	// check if traintuple key already exist
	if elementBytes, _ := stub.GetState(traintupleKey); elementBytes != nil {
		err = fmt.Errorf("traintuple with these algo, in models, and train dataSample already exist - %s", traintupleKey)
		return
	}

	// check if InModels is empty or if mentionned models do exist and fill inModels
	status := "todo"
	parentTraintupleKeys := strings.Split(strings.Replace(inp.InModels, " ", "", -1), ",")
	for _, parentTraintupleKey := range parentTraintupleKeys {
		if parentTraintupleKey == "" {
			break
		}
		parentTraintuple := Traintuple{}
		if err = getElementStruct(stub, parentTraintupleKey, &parentTraintuple); err != nil {
			err = fmt.Errorf("could not retrieve parent traintuple with key %s - %s %d", parentTraintupleKeys, err.Error(), len(parentTraintupleKeys))
			return
		}
		if parentTraintuple.OutModel == nil {
			status = "waiting"
		}
		traintuple.InModelKeys = append(traintuple.InModelKeys, parentTraintupleKey)
	}
	traintuple.Status = status

	// check if DataSampleKeys are from the same dataManager and if they are not test only dataSample
	dataSampleKeys := strings.Split(strings.Replace(inp.DataSampleKeys, " ", "", -1), ",")
	_, trainOnly, err := checkSameDataManager(stub, inp.DataManagerKey, dataSampleKeys)
	if err != nil {
		return
	}
	if !trainOnly {
		err = fmt.Errorf("not possible to create a traintuple with test only data")
		return
	}

	// fill traintuple.Dataset from dataManager and dataSample
	if _, err = getElementBytes(stub, inp.DataManagerKey); err != nil {
		err = fmt.Errorf("could not retrieve dataManager with key %s - %s", inp.DataManagerKey, err.Error())
		return
	}
	traintuple.Dataset = &Dataset{
		DataManagerKey: inp.DataManagerKey,
		DataSampleKeys: dataSampleKeys,
	}
	return
}

// Set is a method of the receiver outputTraintuple. It returns all elements necessary to do a training task from a trainuple stored in the ledger
func (outputTraintuple *outputTraintuple) Set(stub shim.ChaincodeStubInterface, traintuple Traintuple) (err error) {

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
		Worker:     dataManager.Owner,
		Keys:       traintuple.Dataset.DataSampleKeys,
		OpenerHash: traintuple.Dataset.DataManagerKey,
		Perf:       traintuple.Perf,
	}

	return
}

// Set is a method of the receiver Testtuple. It checks the validity of inputTesttuple and uses its fields to set the Testtuple
func (testtuple *Testtuple) Set(stub shim.ChaincodeStubInterface, inp inputTesttuple) (testtupleKey string, err error) {

	validate := validator.New()
	if err = validate.Struct(inp); err != nil {
		return "", fmt.Errorf("invalid inputs to update data %s", err.Error())
	}

	// check associated traintuple
	traintuple := Traintuple{}
	if err = getElementStruct(stub, inp.TraintupleKey, &traintuple); err != nil {
		return testtupleKey, fmt.Errorf("could not retrieve traintuple with key %s - %s", inp.TraintupleKey, err.Error())
	}

	// TODO later: check permissions
	// find associated creator and check permissions (TODO later)
	creator, err := getTxCreator(stub)
	if err != nil {
		return testtupleKey, err
	}
	testtuple.Creator = creator
	testtuple.Permissions = "all"

	// fill info from associated traintuple
	outputTraintuple := &outputTraintuple{}
	outputTraintuple.Set(stub, traintuple)
	testtuple.Objective = outputTraintuple.Objective
	testtuple.Algo = outputTraintuple.Algo
	testtuple.Model = &Model{
		TraintupleKey: inp.TraintupleKey,
	}
	if traintuple.OutModel != nil {
		testtuple.Model.Hash = outputTraintuple.OutModel.Hash
		testtuple.Model.StorageAddress = outputTraintuple.OutModel.StorageAddress
	}

	switch status := traintuple.Status; status {
	case "done":
		testtuple.Status = "todo"
	case "failed":
		err = fmt.Errorf(
			"could not register this testtuple, the traintuple %s has a failed status",
			inp.TraintupleKey)
		return
	default:
		testtuple.Status = "waiting"
	}

	// Get test dataset from objective
	objective := Objective{}
	if err = getElementStruct(stub, testtuple.Objective.Key, &objective); err != nil {
		return testtupleKey, fmt.Errorf("could not retrieve objective with key %s - %s", testtuple.Objective.Key, err.Error())
	}
	objectiveDataManagerKey := objective.TestDataset.DataManagerKey
	objectiveDataSampleKeys := objective.TestDataset.DataSampleKeys
	// For now we need to sort it but in fine it should be save sorted
	// TODO
	sort.Strings(objectiveDataSampleKeys)

	var dataManagerKey string
	var dataSampleKeys []string
	if len(inp.DataManagerKey) > 0 && len(inp.DataSampleKeys) > 0 {
		// non-certified testtuple
		// test dataset are specified by the user
		dataSampleKeys = strings.Split(strings.Replace(inp.DataSampleKeys, " ", "", -1), ",")
		_, _, err = checkSameDataManager(stub, inp.DataManagerKey, dataSampleKeys)
		if err != nil {
			return testtupleKey, err
		}
		dataManagerKey = inp.DataManagerKey
		sort.Strings(dataSampleKeys)
		testtuple.Certified = objectiveDataManagerKey == dataManagerKey && reflect.DeepEqual(objectiveDataSampleKeys, dataSampleKeys)
	} else if len(inp.DataManagerKey) > 0 || len(inp.DataSampleKeys) > 0 {
		return testtupleKey, fmt.Errorf("invalid input: dataManagerKey and dataSampleKey should be provided together")
	} else {
		dataSampleKeys = objectiveDataSampleKeys
		dataManagerKey = objectiveDataManagerKey
		testtuple.Certified = true
	}
	// retrieve dataManager owner
	dataManager := DataManager{}
	if err = getElementStruct(stub, dataManagerKey, &dataManager); err != nil {
		return testtupleKey, fmt.Errorf("could not retrieve dataManager with key %s - %s", dataManagerKey, err.Error())
	}
	testtuple.Dataset = &TtDataset{
		Worker:     dataManager.Owner,
		Keys:       dataSampleKeys,
		OpenerHash: dataManagerKey,
	}

	// create testtuple key and check if it already exists
	toHash := "testtuple"
	toHash += testtuple.Model.TraintupleKey
	toHash += strings.Join(dataSampleKeys, ",")
	tKey := sha256.Sum256([]byte(toHash))
	testtupleKey = hex.EncodeToString(tKey[:])
	if testtupleBytes, err := stub.GetState(testtupleKey); testtupleBytes != nil {
		return testtupleKey, fmt.Errorf("this testtuple already exists (tkey: %s)", testtupleKey)
	} else if err != nil {
		return testtupleKey, err
	}

	return testtupleKey, err
}

// -------------------------------------------------------------------------------------------
// Smart contracts related to traintuples and testuples
// args  [][]byte or []string, it is not possible to input a string looking like a json
// -------------------------------------------------------------------------------------------

// createTraintuple adds a Traintuple in the ledger
func createTraintuple(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	expectedArgs := getFieldNames(&inputTraintuple{})
	if nbArgs := len(expectedArgs); (nbArgs != len(args)) && (nbArgs != len(args)-2) {
		return nil, fmt.Errorf("incorrect arguments, expecting %d args: %s", nbArgs, strings.Join(expectedArgs, ", "))
	}

	// convert input strings args to input struct inputTraintuplet
	inp := inputTraintuple{}
	stringToInputStruct(args, &inp)
	// check validity of input arg and set traintuples
	traintuple := Traintuple{}
	traintupleKey, err := traintuple.Set(stub, inp)
	if err != nil {
		return nil, err
	}

	// store in ledger
	traintupleBytes, _ := json.Marshal(traintuple)
	if err = stub.PutState(traintupleKey, traintupleBytes); err != nil {
		return nil, fmt.Errorf("could not put in ledger traintuple with algo %s inModels %s - %s", inp.AlgoKey, inp.InModels, err.Error())
	}
	// get worker
	worker, err := getDataManagerOwner(stub, traintuple.Dataset.DataManagerKey)
	if err != nil {
		return nil, err
	}

	// create composite keys
	if err = createCompositeKey(stub, "traintuple~algo~key", []string{"traintuple", traintuple.AlgoKey, traintupleKey}); err != nil {
		return nil, fmt.Errorf("issue creating composite keys - %s", err.Error())
	}
	if err = createCompositeKey(stub, "traintuple~worker~status~key", []string{"traintuple", worker, traintuple.Status, traintupleKey}); err != nil {
		return nil, fmt.Errorf("issue creating composite keys - %s", err.Error())
	}
	for _, inModelKey := range traintuple.InModelKeys {
		if err = createCompositeKey(stub, "traintuple~inModel~key", []string{"traintuple", inModelKey, traintupleKey}); err != nil {
			return nil, fmt.Errorf("issue creating composite keys - %s", err.Error())
		}
	}
	return []byte(traintupleKey), nil
}

// createTesttuple adds a Testtuple in the ledger
func createTesttuple(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	expectedArgs := getFieldNames(&inputTesttuple{})
	if nbArgs := len(expectedArgs); nbArgs != len(args) {
		return nil, fmt.Errorf("incorrect arguments, expecting %d args: %s", nbArgs, strings.Join(expectedArgs, ", "))
	}

	// convert input strings args to input struct inputTesttuple
	inp := inputTesttuple{}
	stringToInputStruct(args, &inp)
	// check validity of input arg and set testtuple
	testtuple := Testtuple{}
	testtupleKey, err := testtuple.Set(stub, inp)
	if err != nil {
		return nil, err
	}
	testtupleBytes, _ := json.Marshal(testtuple)
	if err = stub.PutState(testtupleKey, testtupleBytes); err != nil {
		return nil, fmt.Errorf("could not put in ledger testtuple associated with traintuple %s - %s", inp.TraintupleKey, err.Error())
	}

	// create composite keys
	if err = createCompositeKey(stub, "testtuple~algo~key", []string{"testtuple", testtuple.Algo.Hash, testtupleKey}); err != nil {
		return nil, fmt.Errorf("issue creating composite keys - %s", err.Error())
	}
	if err = createCompositeKey(stub, "testtuple~worker~status~key", []string{"testtuple", testtuple.Dataset.Worker, testtuple.Status, testtupleKey}); err != nil {
		return nil, fmt.Errorf("issue creating composite keys - %s", err.Error())
	}
	if err = createCompositeKey(stub, "testtuple~traintuple~certified~key", []string{"testtuple", inp.TraintupleKey, strconv.FormatBool(testtuple.Certified), testtupleKey}); err != nil {
		return nil, fmt.Errorf("issue creating composite keys - %s", err.Error())
	}
	return []byte(testtupleKey), nil
}

// logStartTrain modifies a traintuple by changing its status from todo to doing
func logStartTrain(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	expectedArgs := [1]string{"key of the traintuple to update"}
	if nbArgs := len(expectedArgs); nbArgs != len(args) {
		return nil, fmt.Errorf("incorrect arguments, expecting %d args: %s", nbArgs, strings.Join(expectedArgs[:], ", "))
	}

	oldStatus := "todo"
	status := "doing"
	traintupleKey := args[0]
	// get traintuple, check validity of the update, and update its status
	traintuple, oldStatus, err := updateStatusTraintuple(stub, traintupleKey, status)
	if err != nil {
		return nil, err
	}
	// save to ledger
	traintupleBytes, _ := json.Marshal(traintuple)
	if err := stub.PutState(traintupleKey, traintupleBytes); err != nil {
		return traintupleBytes, fmt.Errorf("failed to update traintuple status to %s with key %s", status, traintupleKey)
	}
	// get worker
	worker, err := getDataManagerOwner(stub, traintuple.Dataset.DataManagerKey)
	if err != nil {
		return nil, err
	}
	// update associated composite keys
	indexName := "traintuple~worker~status~key"
	oldAttributes := []string{"traintuple", worker, oldStatus, traintupleKey}
	newAttributes := []string{"traintuple", worker, status, traintupleKey}
	if err := updateCompositeKey(stub, indexName, oldAttributes, newAttributes); err != nil {
		return nil, err
	}
	return traintupleBytes, nil
}

// logStartTest modifies a testtuple by changing its status from todo to doing
func logStartTest(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	expectedArgs := [1]string{"key of the testtuple to update"}
	if nbArgs := len(expectedArgs); nbArgs != len(args) {
		return nil, fmt.Errorf("incorrect arguments, expecting %d args: %s", nbArgs, strings.Join(expectedArgs[:], ", "))
	}

	oldStatus := "todo"
	status := "doing"
	testtupleKey := args[0]
	// get testtuple, check validity of the update, and update its status
	testtuple, oldStatus, err := updateStatusTesttuple(stub, testtupleKey, status)
	if err != nil {
		return nil, err
	}
	// save to ledger
	testtupleBytes, _ := json.Marshal(testtuple)
	if err := stub.PutState(testtupleKey, testtupleBytes); err != nil {
		return testtupleBytes, fmt.Errorf("failed to update testtuple status to %s with key %s", status, testtupleKey)
	}
	// update associated composite keys
	indexName := "testtuple~worker~status~key"
	oldAttributes := []string{"testtuple", testtuple.Dataset.Worker, oldStatus, testtupleKey}
	newAttributes := []string{"testtuple", testtuple.Dataset.Worker, status, testtupleKey}
	if err := updateCompositeKey(stub, indexName, oldAttributes, newAttributes); err != nil {
		return nil, err
	}
	return testtupleBytes, nil
}

// logSuccessTrain modifies a traintuple by changing its status from doing to done
// reports logs and associated performances
func logSuccessTrain(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	expectedArgs := [4]string{
		"key of the traintuple to update",
		"end model hash and storage address (endModelHash, endModelStorageAddress)",
		"train perf (float)",
		"logs"}
	if nbArgs := len(expectedArgs); nbArgs != len(args) {
		return nil, fmt.Errorf("incorrect arguments, expecting %d args: %s", nbArgs, strings.Join(expectedArgs[:], ", "))
	}

	oldStatus := "doing"
	status := "done"
	// get and check inputs
	traintupleKey := args[0]
	outModel := strings.Split(strings.Replace(args[1], " ", "", -1), ",")
	if len(outModel[0]) != lenKey {
		return nil, fmt.Errorf("invalid length of model hash")
	}
	// get train perf, check validity
	perf, err := strconv.ParseFloat(args[2], 32)
	if err != nil {
		return nil, err
	}
	// get logs and check validity
	log := args[3]
	if err = checkLog(log); err != nil {
		return nil, err
	}

	// get traintuple, check validity of the update, and update its status
	traintuple, oldStatus, err := updateStatusTraintuple(stub, traintupleKey, status)
	if err != nil {
		return nil, err
	}

	// update traintuple
	traintuple.Perf = float32(perf)
	traintuple.OutModel = &HashDress{
		Hash:           outModel[0],
		StorageAddress: outModel[1]}
	traintuple.Log += log
	traintupleBytes, _ := json.Marshal(traintuple)
	if err = stub.PutState(traintupleKey, traintupleBytes); err != nil {
		return nil, fmt.Errorf("failed to update traintuple status to trained with key %s", traintupleKey)
	}
	// get worker
	worker, err := getDataManagerOwner(stub, traintuple.Dataset.DataManagerKey)
	if err != nil {
		return nil, err
	}
	// update associated composite keys
	indexName := "traintuple~worker~status~key"
	oldAttributes := []string{"traintuple", worker, oldStatus, traintupleKey}
	newAttributes := []string{"traintuple", worker, status, traintupleKey}
	if err = updateCompositeKey(stub, indexName, oldAttributes, newAttributes); err != nil {
		return nil, err
	}

	// update traintuples whose inModel use the outModel of this traintuple
	if err := updateWaitingTraintuples(stub, traintupleKey, "done"); err != nil {
		return nil, err
	}
	// update testtuple associated with this traintuple
	if err := updateWaitingTesttuples(stub, traintupleKey, traintuple.OutModel, "todo"); err != nil {
		return nil, err
	}
	return traintupleBytes, nil
}

// logSuccessTest modifies a testtuple by changing its status to done, reports perf and logs
func logSuccessTest(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	expectedArgs := [3]string{"key of the testtuple to update", "test perf (float)", "logs"}
	if nbArgs := len(expectedArgs); nbArgs != len(args) {
		return nil, fmt.Errorf("incorrect arguments, expecting %d args: %s", nbArgs, strings.Join(expectedArgs[:], ", "))
	}

	status := "done"
	// get testtuple
	testtupleKey := args[0]
	if len(testtupleKey) != lenKey {
		return nil, fmt.Errorf("invalid testtuple key")
	}
	// get testtuple, check validity of the update, and update its status
	testtuple, oldStatus, err := updateStatusTesttuple(stub, testtupleKey, status)
	if err != nil {
		return nil, err
	}

	// get test perf, check validity
	perf, err := strconv.ParseFloat(args[1], 32)
	if err != nil {
		return nil, err
	}
	testtuple.Dataset.Perf = float32(perf)

	// get logs and check validity
	log := args[2]
	if err = checkLog(log); err != nil {
		return nil, err
	}
	testtuple.Log += log

	// save to ledger
	testtupleBytes, _ := json.Marshal(testtuple)
	if err = stub.PutState(testtupleKey, testtupleBytes); err != nil {
		return nil, fmt.Errorf("failed to update testtuple status to trained with key %s", testtupleKey)
	}

	// update associated composite keys
	indexName := "testtuple~worker~status~key"
	oldAttributes := []string{"testtuple", testtuple.Dataset.Worker, oldStatus, testtupleKey}
	newAttributes := []string{"testtuple", testtuple.Dataset.Worker, status, testtupleKey}
	if err = updateCompositeKey(stub, indexName, oldAttributes, newAttributes); err != nil {
		return nil, err
	}
	return testtupleBytes, nil
}

// logFailTrain modifies a traintuple by changing its status to fail and reports associated logs
func logFailTrain(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	expectedArgs := [2]string{"the key of the traintuple to update", "logs"}
	if nbArgs := len(expectedArgs); nbArgs != len(args) {
		return nil, fmt.Errorf("incorrect arguments, expecting %d args: %s", nbArgs, strings.Join(expectedArgs[:], ", "))
	}
	// get input args
	traintupleKey := args[0]
	if len(traintupleKey) != lenKey {
		return nil, fmt.Errorf("invalid traintuple key")
	}
	log := args[1]
	if err := checkLog(log); err != nil {
		return nil, err
	}
	// get traintuple and updates its status
	traintuple, oldStatus, err := updateStatusTraintuple(stub, traintupleKey, "failed")
	if err != nil {
		return nil, err
	}
	traintuple.Log += log
	traintupleBytes, _ := json.Marshal(traintuple)
	if err := stub.PutState(traintupleKey, traintupleBytes); err != nil {
		return nil, fmt.Errorf("failed to update traintuple status to failed with key %s", traintupleKey)
	}
	// get worker
	worker, err := getDataManagerOwner(stub, traintuple.Dataset.DataManagerKey)
	if err != nil {
		return nil, err
	}
	// update associated composite keys
	indexName := "traintuple~worker~status~key"
	oldAttributes := []string{"traintuple", worker, oldStatus, traintupleKey}
	newAttributes := []string{"traintuple", worker, "failed", traintupleKey}
	if err := updateCompositeKey(stub, indexName, oldAttributes, newAttributes); err != nil {
		return nil, err
	}
	// update testtuple associated with this traintuple
	if err := updateWaitingTesttuples(stub, traintupleKey, &HashDress{}, "failed"); err != nil {
		return nil, err
	}
	// update traintuple having this traintuple as inModel
	if err := updateWaitingTraintuples(stub, traintupleKey, "failed"); err != nil {
		return nil, err
	}
	return traintupleBytes, nil
}

// logFailTest modifies a testtuple by changing its status to fail and reports associated logs
func logFailTest(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	expectedArgs := [2]string{"the key of the testtuple to update", "logs"}
	if nbArgs := len(expectedArgs); nbArgs != len(args) {
		return nil, fmt.Errorf("incorrect arguments, expecting %d args: %s", nbArgs, strings.Join(expectedArgs[:], ", "))
	}

	status := "failed"
	// get input args
	testtupleKey := args[0]
	if len(testtupleKey) != lenKey {
		return nil, fmt.Errorf("invalid testtuple key")
	}
	log := args[1]
	if err := checkLog(log); err != nil {
		return nil, err
	}
	// get testtuple and updates its status
	testtuple, oldStatus, err := updateStatusTesttuple(stub, testtupleKey, status)
	if err != nil {
		return nil, err
	}
	testtuple.Log += log
	testtupleBytes, _ := json.Marshal(testtuple)
	if err := stub.PutState(testtupleKey, testtupleBytes); err != nil {
		return nil, fmt.Errorf("failed to update testtuple status to failed with key %s", testtupleKey)
	}
	// update associated composite keys
	indexName := "testtuple~worker~status~key"
	oldAttributes := []string{"testtuple", testtuple.Dataset.Worker, oldStatus, testtupleKey}
	newAttributes := []string{"testtuple", testtuple.Dataset.Worker, status, testtupleKey}
	if err := updateCompositeKey(stub, indexName, oldAttributes, newAttributes); err != nil {
		return nil, err
	}
	return testtupleBytes, nil
}

// queryTraintuple returns info about a traintuple given its key
func queryTraintuple(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	expectedArgs := [1]string{"traintuple key"}
	if nbArgs := len(expectedArgs); nbArgs != len(args) {
		return nil, fmt.Errorf("incorrect arguments, expecting %d args: %s", nbArgs, strings.Join(expectedArgs[:], ", "))
	}
	traintupleKey := args[0]
	traintuple := Traintuple{}
	if err := getElementStruct(stub, traintupleKey, &traintuple); err != nil {
		return nil, err
	}
	outputTraintuple := outputTraintuple{}
	outputTraintuple.Set(stub, traintuple)
	// Marshal payload
	payload, err := json.Marshal(outputTraintuple)
	return payload, err

}

// queryTraintuples returns all traintuples
func queryTraintuples(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("incorrect number of arguments, expecting nothing")
	}
	elementsKeys, err := getKeysFromComposite(stub, "traintuple~algo~key", []string{"traintuple"})
	if err != nil {
		return nil, err
	}
	var elements []map[string]interface{}
	for _, key := range elementsKeys {
		var element map[string]interface{}
		outputTraintuple, err := getOutputTraintuple(stub, key)
		if err != nil {
			return nil, err
		}
		oo, _ := json.Marshal(outputTraintuple)
		json.Unmarshal(oo, &element)
		element["key"] = key
		elements = append(elements, element)
	}
	return json.Marshal(elements)
}

// queryTesttuple returns a testtuple of the ledger given its key
func queryTesttuple(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 || len(args[0]) != 64 {
		return nil, fmt.Errorf("incorrect arguments, expecting key, received: %s", args[0])
	}
	key := args[0]
	var testtuple Testtuple
	if err := getElementStruct(stub, key, &testtuple); err != nil {
		return nil, err
	}
	return json.Marshal(testtuple)
}

// queryTesttuples returns all testtuples of the ledger
func queryTesttuples(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("incorrect number of arguments, expecting nothing")
	}
	var indexName = "testtuple~traintuple~certified~key"
	elementsKeys, err := getKeysFromComposite(stub, indexName, []string{"testtuple"})
	if err != nil {
		return nil, fmt.Errorf("issue getting keys from composite key %s - %s", indexName, err.Error())
	}
	var outTesttuples []Testtuple
	for _, key := range elementsKeys {
		var testtuple Testtuple
		if err := getElementStruct(stub, key, &testtuple); err != nil {
			return nil, err
		}
		outTesttuples = append(outTesttuples, testtuple)
	}
	return json.Marshal(outTesttuples)
}

// queryModelDetails returns info about the testtuple and algo related to a traintuple
func queryModelDetails(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	expectedArgs := [1]string{"traintuple key"}
	if nbArgs := len(expectedArgs); nbArgs != len(args) {
		return nil, fmt.Errorf("incorrect arguments, expecting %d args: %s", nbArgs, strings.Join(expectedArgs[:], ", "))
	}
	traintupleKey := args[0]

	mPayload := make(map[string]interface{})
	// get associated traintuple
	var element map[string]interface{}
	outputTraintuple, err := getOutputTraintuple(stub, traintupleKey)
	if err != nil {
		return nil, err
	}
	oo, _ := json.Marshal(outputTraintuple)
	json.Unmarshal(oo, &element)
	element["key"] = traintupleKey
	mPayload["traintuple"] = element
	// get certified and non-certified testtuples related to traintuple
	var nonCertifiedTesttuples []map[string]interface{}
	testtupleKeys, err := getKeysFromComposite(stub, "testtuple~traintuple~certified~key", []string{"testtuple", traintupleKey})
	if err != nil {
		return nil, err
	}
	for _, testtupleKey := range testtupleKeys {
		// get testtuple and serialize it
		var testtuple map[string]interface{}
		if err = getElementStruct(stub, testtupleKey, &testtuple); err != nil {
			return nil, err
		}
		testtuple["key"] = testtupleKey
		if testtuple["certified"] == true {
			mPayload["testtuple"] = testtuple
		} else {
			nonCertifiedTesttuples = append(nonCertifiedTesttuples, testtuple)
		}
		mPayload["nonCertifiedTesttuples"] = nonCertifiedTesttuples
	}
	// Marshal payload
	payload, err := json.Marshal(mPayload)
	return payload, err
}

// queryModels returns all traintuples and associated testuples
// TODO
func queryModels(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("incorrect number of arguments, expecting nothing")
	}

	traintupleKeys, err := getKeysFromComposite(stub, "traintuple~algo~key", []string{"traintuple"})
	if err != nil {
		return nil, err
	}
	var elements []map[string]interface{}
	for _, traintupleKey := range traintupleKeys {
		element := make(map[string]interface{})
		traintuple := make(map[string]interface{})

		// get traintuple
		outputTraintuple, err := getOutputTraintuple(stub, traintupleKey)
		if err != nil {
			return nil, err
		}
		oo, _ := json.Marshal(outputTraintuple)
		json.Unmarshal(oo, &traintuple)
		traintuple["key"] = traintupleKey
		element["traintuple"] = traintuple

		// get testtuple related to traintuple
		testtupleKeys, err := getKeysFromComposite(stub, "testtuple~traintuple~certified~key", []string{"testtuple", traintupleKey, "true"})
		if err != nil {
			return nil, err
		}
		if len(testtupleKeys) == 1 {
			// get testtuple and serialize it
			testtupleKey := testtupleKeys[0]
			testtuple := make(map[string]interface{})
			if err = getElementStruct(stub, testtupleKey, &testtuple); err != nil {
				return nil, err
			}
			testtuple["key"] = testtupleKey
			element["testtuple"] = testtuple
		}
		elements = append(elements, element)
	}
	return json.Marshal(elements)
}

// --------------------------------------------------------------
// Utils for smartcontracts related to traintuples and testtuples
// --------------------------------------------------------------

// getOutputTraintuple takes as input a traintuple key and returns the outputTraintuple
func getOutputTraintuple(stub shim.ChaincodeStubInterface, traintupleKey string) (outputTraintuple, error) {
	traintuple := Traintuple{}
	outputTraintuple := outputTraintuple{}
	if err := getElementStruct(stub, traintupleKey, &traintuple); err != nil {
		return outputTraintuple, err
	}
	outputTraintuple.Set(stub, traintuple)
	return outputTraintuple, nil
}

// checkLog checks the validity of logs
func checkLog(log string) (err error) {
	maxLength := 200
	if length := len(log); length > maxLength {
		err = fmt.Errorf("too long log, is %d and should be %d ", length, maxLength)
	}
	return
}

// check validity of traintuple update: consistent status and agent submitting the transaction
func checkUpdateTuple(stub shim.ChaincodeStubInterface, worker string, oldStatus string, newStatus string) error {
	txCreator, err := getTxCreator(stub)
	if err != nil {
		return err
	}
	if txCreator != worker {
		return fmt.Errorf("%s is not allowed to update tuple", txCreator)
	}
	statusPossibilities := map[string]string{
		"todo":  "doing",
		"doing": "done"}
	if statusPossibilities[oldStatus] != newStatus && newStatus != "failed" {
		return fmt.Errorf("cannot change status from %s to %s", oldStatus, newStatus)
	}
	return nil
}

//
// TODO: change names for all functions related to updates, since it is not consistent
//

// updateStatusTraintuple retrieves a traintuple given its key, check the validity of the status update, changes the status of a traintuple, and returns the updated traintuple and its oldStatus
func updateStatusTraintuple(stub shim.ChaincodeStubInterface, traintupleKey string, status string) (Traintuple, string, error) {

	var oldStatus string
	traintuple := Traintuple{}
	if len(traintupleKey) != lenKey {
		return traintuple, oldStatus, fmt.Errorf("invalid traintuple key")
	}
	if err := getElementStruct(stub, traintupleKey, &traintuple); err != nil {
		return traintuple, oldStatus, err
	}
	oldStatus = traintuple.Status
	// get worker
	worker, err := getDataManagerOwner(stub, traintuple.Dataset.DataManagerKey)
	if err != nil {
		return traintuple, oldStatus, err
	}
	// check validity of worker and change of status
	if err := checkUpdateTuple(stub, worker, oldStatus, status); err != nil {
		return traintuple, oldStatus, err
	}
	// update traintuple
	traintuple.Status = status

	return traintuple, oldStatus, nil
}

// updateStatusTesttuple retrieves a testtuple given its key, check the validity of the status update, changes the status of a testtuple, and returns the updated tesntuple and its oldStatus
func updateStatusTesttuple(stub shim.ChaincodeStubInterface, testtupleKey string, status string) (Testtuple, string, error) {

	var oldStatus string
	testtuple := Testtuple{}
	if len(testtupleKey) != lenKey {
		return testtuple, oldStatus, fmt.Errorf("invalid testtuple key")
	}
	if err := getElementStruct(stub, testtupleKey, &testtuple); err != nil {
		return testtuple, oldStatus, err
	}
	oldStatus = testtuple.Status
	// check validity of worker and change of status
	if err := checkUpdateTuple(stub, testtuple.Dataset.Worker, oldStatus, status); err != nil {
		return testtuple, oldStatus, err
	}
	// update traintuple
	testtuple.Status = status

	return testtuple, oldStatus, nil
}

// updateWaitingTraintuples updates the status of waiting trainuples  InModels of traintuples once they have been trained (succesfully or failed)
// func updateWaitingInModels(stub shim.ChaincodeStubInterface, modelTraintupleKey string, model *HashDress) error {
func updateWaitingTraintuples(stub shim.ChaincodeStubInterface, modelTraintupleKey string, status string) error {

	indexName := "traintuple~inModel~key"
	// get traintuples having as inModels the input traintuple
	traintupleKeys, err := getKeysFromComposite(stub, indexName, []string{"traintuple", modelTraintupleKey})
	if err != nil {
		return fmt.Errorf("error while getting associated traintuples to update their inModel")
	}
	for _, traintupleKey := range traintupleKeys {
		// get and update traintuple
		traintuple := Traintuple{}
		if err := getElementStruct(stub, traintupleKey, &traintuple); err != nil {
			return err
		}
		oldStatus := traintuple.Status
		if status == "failed" {
			traintuple.Status = status
		} else if status == "done" {
			err := traintuple.CheckReady(stub, modelTraintupleKey)
			if err != nil {
				return err
			}
		}
		// remove associated composite key
		compositeKey, err := stub.CreateCompositeKey(indexName, []string{"traintuple", modelTraintupleKey, traintupleKey})
		if err != nil {
			return fmt.Errorf("failed to recreate composite key to update traintuple %s with inModel %s - %s", traintupleKey, modelTraintupleKey, err.Error())
		}
		if err := stub.DelState(compositeKey); err != nil {
			return fmt.Errorf("failed to delete associated composite key to update traintuple %s with inModel %s - %s", traintupleKey, modelTraintupleKey, err.Error())
		}
		if oldStatus != traintuple.Status {
			// save update in ledger
			traintupleBytes, _ := json.Marshal(traintuple)
			if err := stub.PutState(traintupleKey, traintupleBytes); err != nil {
				return fmt.Errorf("failed to update traintuple %s with inModel %s - %s", traintupleKey, modelTraintupleKey, err.Error())
			}
			// get worker
			worker, err := getDataManagerOwner(stub, traintuple.Dataset.DataManagerKey)
			if err != nil {
				return err
			}
			// update associated composite keys
			indexName := "traintuple~worker~status~key"
			oldAttributes := []string{"traintuple", worker, "waiting", traintupleKey}
			newAttributes := []string{"traintuple", worker, traintuple.Status, traintupleKey}
			if err = updateCompositeKey(stub, indexName, oldAttributes, newAttributes); err != nil {
				return err
			}
		}

	}
	return nil
}

// CheckReady checks if inModels of a traintuple have been trained, except the newDoneTraintupleKey (since the transaction is not commited)
// and updates the traintuple status if necessary
// TODO idem change name of the function
// func (traintuple *Traintuple) CheckReady(modelTraintupleKey string, model *HashDress) {
func (traintuple *Traintuple) CheckReady(stub shim.ChaincodeStubInterface, newDoneTraintupleKey string) error {
	status := "todo"
	// would be much easier if traintuple.InModels was map[string]HashDress
	for _, key := range traintuple.InModelKeys {
		if key != newDoneTraintupleKey {
			tt := Traintuple{}
			if err := getElementStruct(stub, key, &tt); err != nil {
				return err
			}
			if tt.Status == "waiting" || tt.Status == "todo" {
				status = "waiting"
			} else if tt.Status == "failed" {
				status = "failed"
			}
		}
	}
	traintuple.Status = status
	return nil

}

// updateWaitingTesttuple updates Status of testtuple whose associated traintuple has been trained or has failed
func updateWaitingTesttuples(stub shim.ChaincodeStubInterface, traintupleKey string, model *HashDress, testtupleStatus string) error {

	if !stringInSlice(testtupleStatus, []string{"todo", "failed"}) {
		return fmt.Errorf("invalid status of associated traintuple")
	}

	indexName := "testtuple~traintuple~certified~key"
	// get testtuple associated with this traintuple and updates its status
	testtupleKeys, err := getKeysFromComposite(stub, indexName, []string{"testtuple", traintupleKey})
	if err != nil || len(testtupleKeys) == 0 {
		return err
	}
	testtupleKey := testtupleKeys[0]
	// get and update testtuple
	testtuple := Testtuple{}
	if err := getElementStruct(stub, testtupleKey, &testtuple); err != nil {
		return err
	}
	oldStatus := testtuple.Status
	testtuple.Model = &Model{
		TraintupleKey:  traintupleKey,
		Hash:           model.Hash,
		StorageAddress: model.StorageAddress,
	}
	testtuple.Status = testtupleStatus
	testtupleBytes, _ := json.Marshal(testtuple)
	if err := stub.PutState(testtupleKey, testtupleBytes); err != nil {
		return fmt.Errorf("failed to update testtuple associated with traintuple %s - %s", traintupleKey, err.Error())
	}
	// update associated composite key
	indexName = "testtuple~worker~status~key"
	oldAttributes := []string{"testtuple", testtuple.Dataset.Worker, oldStatus, testtupleKey}
	newAttributes := []string{"testtuple", testtuple.Dataset.Worker, testtupleStatus, testtupleKey}
	if err := updateCompositeKey(stub, indexName, oldAttributes, newAttributes); err != nil {
		return err
	}
	return nil
}

// getTraintuplesPayload takes as input a list of keys and returns a paylaod containing a list of associated retrieved elements
func getTraintuplesPayload(stub shim.ChaincodeStubInterface, traintupleKeys []string) ([]byte, error) {

	var elements []map[string]interface{}
	for _, key := range traintupleKeys {
		var element map[string]interface{}
		outputTraintuple, err := getOutputTraintuple(stub, key)
		if err != nil {
			return nil, err
		}
		oo, err := json.Marshal(outputTraintuple)
		if err != nil {
			return nil, err
		}
		json.Unmarshal(oo, &element)
		element["key"] = key
		elements = append(elements, element)
	}
	return json.Marshal(elements)
}
