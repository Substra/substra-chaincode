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
	traintuple.Tag = inp.Tag
	// check if algo exists
	if _, err = getElementBytes(stub, inp.AlgoKey); err != nil {
		err = fmt.Errorf("could not retrieve algo with key %s - %s", inp.AlgoKey, err.Error())
		return
	}
	traintuple.AlgoKey = inp.AlgoKey

	// check objective and add it
	obj := Objective{}
	if err = getElementStruct(stub, inp.ObjectiveKey, &obj); err != nil {
		err = fmt.Errorf("could not retrieve objective with key %s - %s", inp.ObjectiveKey, err.Error())
		return
	}
	traintuple.ObjectiveKey = inp.ObjectiveKey

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
		// set traintuple to waiting if one of the parent traintuples is not done
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
	traintuple.Dataset.Worker, err = getDataManagerOwner(stub, traintuple.Dataset.DataManagerKey)
	if err != nil {
		return
	}

	hashKeys := []string{creator, traintuple.AlgoKey, traintuple.Dataset.DataManagerKey}
	hashKeys = append(hashKeys, traintuple.Dataset.DataSampleKeys...)
	hashKeys = append(hashKeys, traintuple.InModelKeys...)
	traintupleKey, err = HashForKey(stub, "traintuple", hashKeys...)
	if err != nil {
		return
	}

	// check FLTask and Rank and set it when required
	if inp.Rank == "" {
		if inp.FLTask != "" {
			err = fmt.Errorf("invalit inputs, a FLTask should have a rank")
			return
		}
	} else {
		traintuple.Rank, err = strconv.Atoi(inp.Rank)
		if err != nil {
			return
		}
		if inp.FLTask == "" {
			if traintuple.Rank != 0 {
				err = fmt.Errorf("invalid inputs, a new FLTask should have a rank 0")
				return
			}
			traintuple.FLTask = traintupleKey
		} else {
			var ttKeys []string
			attributes := []string{"traintuple", inp.FLTask}
			ttKeys, err = getKeysFromComposite(stub, "traintuple~fltask~worker~rank~key", attributes)
			if err != nil {
				return
			} else if len(ttKeys) == 0 {
				err = fmt.Errorf("cannot find the FLTask %s", inp.FLTask)
				return
			}
			for _, ttKey := range ttKeys {
				FLTraintuple := Traintuple{}
				err = getElementStruct(stub, ttKey, &FLTraintuple)
				if err != nil {
					return
				} else if FLTraintuple.AlgoKey != inp.AlgoKey {
					err = fmt.Errorf("previous traintuple for FLTask %s does not have the same algo key %s", inp.FLTask, inp.AlgoKey)
					return
				}
			}

			attributes = []string{"traintuple", inp.FLTask, traintuple.Dataset.Worker, inp.Rank}
			ttKeys, err = getKeysFromComposite(stub, "traintuple~fltask~worker~rank~key", attributes)
			if err != nil {
				return
			} else if len(ttKeys) > 0 {
				err = fmt.Errorf("FLTask %s with worker %s rank %d already exists", inp.FLTask, traintuple.Dataset.Worker, traintuple.Rank)
				return
			}

			traintuple.FLTask = inp.FLTask
		}
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
	testtuple.Tag = inp.Tag

	// fill info from associated traintuple
	outputTraintuple := &outputTraintuple{}
	outputTraintuple.Fill(stub, traintuple, inp.TraintupleKey)
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
	var objectiveDataManagerKey string
	var objectiveDataSampleKeys []string
	if objective.TestDataset != nil {
		objectiveDataManagerKey = objective.TestDataset.DataManagerKey
		objectiveDataSampleKeys = objective.TestDataset.DataSampleKeys
	}
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
	} else if objective.TestDataset != nil {
		dataSampleKeys = objectiveDataSampleKeys
		dataManagerKey = objectiveDataManagerKey
		testtuple.Certified = true
	} else {
		err = fmt.Errorf("can not create a certified testtuple, no data associated with objective %s", testtuple.Objective.Key)
		return
	}
	// retrieve dataManager owner
	dataManager := DataManager{}
	if err = getElementStruct(stub, dataManagerKey, &dataManager); err != nil {
		return testtupleKey, fmt.Errorf("could not retrieve dataManager with key %s - %s", dataManagerKey, err.Error())
	}
	testtuple.Dataset = &TtDataset{
		Worker:         dataManager.Owner,
		DataSampleKeys: dataSampleKeys,
		OpenerHash:     dataManagerKey,
	}

	// create testtuple key and check if it already exists
	hashKeys := []string{testtuple.Model.TraintupleKey, dataManagerKey, creator}
	hashKeys = append(hashKeys, dataSampleKeys...)
	testtupleKey, err = HashForKey(stub, "testtuple", hashKeys...)

	return testtupleKey, err
}

// -------------------------------------------------------------------------------------------
// Smart contracts related to traintuples and testuples
// args  [][]byte or []string, it is not possible to input a string looking like a json
// -------------------------------------------------------------------------------------------

// createTraintuple adds a Traintuple in the ledger
func createTraintuple(stub shim.ChaincodeStubInterface, args []string) (resp map[string]string, err error) {
	inp := inputTraintuple{}
	err = AssetFromJSON(args[0], &inp)
	if err != nil {
		return
	}

	// check validity of input arg and set traintuples
	traintuple := Traintuple{}
	traintupleKey, err := traintuple.Set(stub, inp)
	if err != nil {
		return
	}

	// store in ledger
	traintupleBytes, _ := json.Marshal(traintuple)
	if err = stub.PutState(traintupleKey, traintupleBytes); err != nil {
		err = fmt.Errorf("could not put in ledger traintuple with algo %s inModels %s - %s", inp.AlgoKey, inp.InModels, err.Error())
		return
	}

	// create composite keys
	if err = createCompositeKey(stub, "traintuple~algo~key", []string{"traintuple", traintuple.AlgoKey, traintupleKey}); err != nil {
		err = fmt.Errorf("issue creating composite keys - %s", err.Error())
		return
	}
	if err = createCompositeKey(stub, "traintuple~worker~status~key", []string{"traintuple", traintuple.Dataset.Worker, traintuple.Status, traintupleKey}); err != nil {
		err = fmt.Errorf("issue creating composite keys - %s", err.Error())
		return
	}
	for _, inModelKey := range traintuple.InModelKeys {
		if err = createCompositeKey(stub, "traintuple~inModel~key", []string{"traintuple", inModelKey, traintupleKey}); err != nil {
			err = fmt.Errorf("issue creating composite keys - %s", err.Error())
			return
		}
	}
	if traintuple.FLTask != "" {
		if err = createCompositeKey(stub, "traintuple~fltask~worker~rank~key", []string{"traintuple", traintuple.FLTask, traintuple.Dataset.Worker, inp.Rank, traintupleKey}); err != nil {
			err = fmt.Errorf("issue creating composite keys - %s", err.Error())
			return
		}
	}
	if traintuple.Tag != "" {
		err = createCompositeKey(stub, "traintuple~tag~key", []string{"traintuple", traintuple.Tag, traintupleKey})
		if err != nil {
			return nil, err
		}
	}

	out := outputTraintuple{}
	err = out.Fill(stub, traintuple, traintupleKey)
	if err != nil {
		return nil, err
	}

	// We can only send one event
	if traintuple.Status == "todo" {
		err = SetEvent(stub, "traintuple-ready", out)
	} else {
		err = SetEvent(stub, "traintuple-created", out)
	}

	if err != nil {
		return nil, err
	}

	return map[string]string{"key": traintupleKey}, nil
}

// createTesttuple adds a Testtuple in the ledger
func createTesttuple(stub shim.ChaincodeStubInterface, args []string) (resp map[string]string, err error) {
	inp := inputTesttuple{}
	err = AssetFromJSON(args[0], &inp)
	if err != nil {
		return
	}

	// check validity of input arg and set testtuple
	testtuple := Testtuple{}
	testtupleKey, err := testtuple.Set(stub, inp)
	if err != nil {
		return
	}
	testtupleBytes, _ := json.Marshal(testtuple)
	if err = stub.PutState(testtupleKey, testtupleBytes); err != nil {
		err = fmt.Errorf("could not put in ledger testtuple associated with traintuple %s - %s", inp.TraintupleKey, err.Error())
		return
	}

	// create composite keys
	if err = createCompositeKey(stub, "testtuple~algo~key", []string{"testtuple", testtuple.Algo.Hash, testtupleKey}); err != nil {
		err = fmt.Errorf("issue creating composite keys - %s", err.Error())
		return
	}
	if err = createCompositeKey(stub, "testtuple~worker~status~key", []string{"testtuple", testtuple.Dataset.Worker, testtuple.Status, testtupleKey}); err != nil {
		err = fmt.Errorf("issue creating composite keys - %s", err.Error())
		return
	}
	if err = createCompositeKey(stub, "testtuple~traintuple~certified~key", []string{"testtuple", inp.TraintupleKey, strconv.FormatBool(testtuple.Certified), testtupleKey}); err != nil {
		err = fmt.Errorf("issue creating composite keys - %s", err.Error())
		return
	}
	if testtuple.Tag != "" {
		err = createCompositeKey(stub, "testtuple~tag~key", []string{"traintuple", testtuple.Tag, testtupleKey})
		if err != nil {
			return nil, err
		}
	}

	out := outputTesttuple{}
	out.Fill(testtupleKey, testtuple)

	// We can only send one event
	if testtuple.Status == "todo" {
		err = SetEvent(stub, "testtuple-ready", out)
	} else {
		err = SetEvent(stub, "testtuple-created", out)
	}

	if err != nil {
		return nil, err
	}

	return map[string]string{"key": testtupleKey}, nil
}

// logStartTrain modifies a traintuple by changing its status from todo to doing
func logStartTrain(stub shim.ChaincodeStubInterface, args []string) (traintuple Traintuple, err error) {
	inp := inputHashe{}
	err = AssetFromJSON(args[0], &inp)
	if err != nil {
		return
	}

	status := "doing"
	// get traintuple, check validity of the update, and update its status
	traintuple = Traintuple{}

	if err = getElementStruct(stub, inp.Key, &traintuple); err != nil {
		return
	}
	err = traintuple.checkNewStatus(stub, status)
	if err != nil {
		return
	}

	// save to ledger
	if err = traintuple.commitUpdate(stub, inp.Key, status); err != nil {
		return
	}
	return
}

// logStartTest modifies a testtuple by changing its status from todo to doing
func logStartTest(stub shim.ChaincodeStubInterface, args []string) (testtuple Testtuple, err error) {
	inp := inputHashe{}
	err = AssetFromJSON(args[0], &inp)
	if err != nil {
		return
	}

	oldStatus := "todo"
	status := "doing"
	// get testtuple, check validity of the update, and update its status
	testtuple, oldStatus, err = updateStatusTesttuple(stub, inp.Key, status)
	if err != nil {
		return
	}
	// save to ledger
	testtupleBytes, _ := json.Marshal(testtuple)
	if err = stub.PutState(inp.Key, testtupleBytes); err != nil {
		err = fmt.Errorf("failed to update testtuple status to %s with key %s", status, inp.Key)
		return
	}
	// update associated composite keys
	indexName := "testtuple~worker~status~key"
	oldAttributes := []string{"testtuple", testtuple.Dataset.Worker, oldStatus, inp.Key}
	newAttributes := []string{"testtuple", testtuple.Dataset.Worker, status, inp.Key}
	if err = updateCompositeKey(stub, indexName, oldAttributes, newAttributes); err != nil {
		return
	}
	return
}

// logSuccessTrain modifies a traintuple by changing its status from doing to done
// reports logs and associated performances
func logSuccessTrain(stub shim.ChaincodeStubInterface, args []string) (traintuple Traintuple, err error) {
	inp := inputLogSuccessTrain{}
	err = AssetFromJSON(args[0], &inp)
	if err != nil {
		return
	}

	status := "done"
	// get traintuple, check validity of the update, and update its status
	traintuple = Traintuple{}
	if err = getElementStruct(stub, inp.Key, &traintuple); err != nil {
		return
	}
	err = traintuple.checkNewStatus(stub, status)
	if err != nil {
		return
	}

	// update traintuple
	traintuple.Perf = inp.Perf
	traintuple.OutModel = &HashDress{
		Hash:           inp.OutModel.Hash,
		StorageAddress: inp.OutModel.StorageAddress}
	traintuple.Log += inp.Log

	if err = traintuple.commitUpdate(stub, inp.Key, status); err != nil {
		return
	}

	// update traintuples whose inModel use the outModel of this traintuple
	if err = updateWaitingTraintuples(stub, inp.Key, "done"); err != nil {
		return
	}
	// update testtuple associated with this traintuple
	if err = updateWaitingTesttuples(stub, inp.Key, traintuple.OutModel, "todo"); err != nil {
		return
	}
	return
}

// logSuccessTest modifies a testtuple by changing its status to done, reports perf and logs
func logSuccessTest(stub shim.ChaincodeStubInterface, args []string) (testtuple Testtuple, err error) {
	inp := inputLogSuccessTest{}
	err = AssetFromJSON(args[0], &inp)
	if err != nil {
		return
	}

	status := "done"
	// get testtuple, check validity of the update, and update its status
	testtuple, oldStatus, err := updateStatusTesttuple(stub, inp.Key, status)
	if err != nil {
		return
	}
	testtuple.Dataset.Perf = inp.Perf
	testtuple.Log += inp.Log

	// save to ledger
	testtupleBytes, _ := json.Marshal(testtuple)
	if err = stub.PutState(inp.Key, testtupleBytes); err != nil {
		err = fmt.Errorf("failed to update testtuple status to trained with key %s", inp.Key)
		return
	}

	// update associated composite keys
	indexName := "testtuple~worker~status~key"
	oldAttributes := []string{"testtuple", testtuple.Dataset.Worker, oldStatus, inp.Key}
	newAttributes := []string{"testtuple", testtuple.Dataset.Worker, status, inp.Key}
	if err = updateCompositeKey(stub, indexName, oldAttributes, newAttributes); err != nil {
		return
	}
	return
}

// logFailTrain modifies a traintuple by changing its status to fail and reports associated logs
func logFailTrain(stub shim.ChaincodeStubInterface, args []string) (traintuple Traintuple, err error) {
	inp := inputLogFailTrain{}
	err = AssetFromJSON(args[0], &inp)
	if err != nil {
		return
	}

	// get traintuple and updates its status
	status := "failed"
	traintuple = Traintuple{}
	if err = getElementStruct(stub, inp.Key, &traintuple); err != nil {
		return
	}
	err = traintuple.checkNewStatus(stub, status)
	if err != nil {
		return
	}
	traintuple.Log += inp.Log

	if err = traintuple.commitUpdate(stub, inp.Key, status); err != nil {
		return
	}

	// update testtuple associated with this traintuple
	if err = updateWaitingTesttuples(stub, inp.Key, &HashDress{}, status); err != nil {
		return
	}
	// update traintuple having this traintuple as inModel
	if err = updateWaitingTraintuples(stub, inp.Key, status); err != nil {
		return
	}
	return
}

// logFailTest modifies a testtuple by changing its status to fail and reports associated logs
func logFailTest(stub shim.ChaincodeStubInterface, args []string) (testtuple Testtuple, err error) {
	inp := inputLogFailTest{}
	err = AssetFromJSON(args[0], &inp)
	if err != nil {
		return
	}

	// get testtuple and updates its status
	status := "failed"
	testtuple, oldStatus, err := updateStatusTesttuple(stub, inp.Key, status)
	if err != nil {
		return
	}
	testtuple.Log += inp.Log
	testtupleBytes, _ := json.Marshal(testtuple)
	if err = stub.PutState(inp.Key, testtupleBytes); err != nil {
		err = fmt.Errorf("failed to update testtuple status to failed with key %s", inp.Key)
		return
	}
	// update associated composite keys
	indexName := "testtuple~worker~status~key"
	oldAttributes := []string{"testtuple", testtuple.Dataset.Worker, oldStatus, inp.Key}
	newAttributes := []string{"testtuple", testtuple.Dataset.Worker, status, inp.Key}
	if err = updateCompositeKey(stub, indexName, oldAttributes, newAttributes); err != nil {
		return
	}
	return
}

// queryTraintuple returns info about a traintuple given its key
func queryTraintuple(stub shim.ChaincodeStubInterface, args []string) (outputTraintuple outputTraintuple, err error) {
	inp := inputHashe{}
	err = AssetFromJSON(args[0], &inp)
	if err != nil {
		return
	}
	traintuple := Traintuple{}
	if err = getElementStruct(stub, inp.Key, &traintuple); err != nil {
		return
	}
	outputTraintuple.Fill(stub, traintuple, inp.Key)
	return
}

// queryTraintuples returns all traintuples
func queryTraintuples(stub shim.ChaincodeStubInterface, args []string) (elements []map[string]interface{}, err error) {
	if len(args) != 0 {
		err = fmt.Errorf("incorrect number of arguments, expecting nothing")
		return
	}
	elementsKeys, err := getKeysFromComposite(stub, "traintuple~algo~key", []string{"traintuple"})
	if err != nil {
		return
	}
	for _, key := range elementsKeys {
		var element map[string]interface{}
		var outputTraintuple outputTraintuple
		outputTraintuple, err = getOutputTraintuple(stub, key)
		if err != nil {
			return
		}
		oo, _ := json.Marshal(outputTraintuple)
		json.Unmarshal(oo, &element)
		element["key"] = key
		elements = append(elements, element)
	}
	return
}

// queryTesttuple returns a testtuple of the ledger given its key
func queryTesttuple(stub shim.ChaincodeStubInterface, args []string) (out outputTesttuple, err error) {
	inp := inputHashe{}
	err = AssetFromJSON(args[0], &inp)
	if err != nil {
		return
	}
	var testtuple Testtuple
	if err = getElementStruct(stub, inp.Key, &testtuple); err != nil {
		return
	}
	out.Fill(inp.Key, testtuple)
	return
}

// queryTesttuples returns all testtuples of the ledger
func queryTesttuples(stub shim.ChaincodeStubInterface, args []string) (outTesttuples []outputTesttuple, err error) {
	if len(args) != 0 {
		err = fmt.Errorf("incorrect number of arguments, expecting nothing")
		return
	}
	var indexName = "testtuple~traintuple~certified~key"
	elementsKeys, err := getKeysFromComposite(stub, indexName, []string{"testtuple"})
	if err != nil {
		err = fmt.Errorf("issue getting keys from composite key %s - %s", indexName, err.Error())
		return
	}
	for _, key := range elementsKeys {
		var testtuple Testtuple
		var out outputTesttuple
		if err = getElementStruct(stub, key, &testtuple); err != nil {
			return
		}
		out.Fill(key, testtuple)
		outTesttuples = append(outTesttuples, out)
	}
	return
}

// queryModelDetails returns info about the testtuple and algo related to a traintuple
func queryModelDetails(stub shim.ChaincodeStubInterface, args []string) (mPayload map[string]interface{}, err error) {
	inp := inputHashe{}
	err = AssetFromJSON(args[0], &inp)
	if err != nil {
		return
	}

	// get associated traintuple
	var element map[string]interface{}
	outputTraintuple, err := getOutputTraintuple(stub, inp.Key)
	if err != nil {
		return
	}
	oo, _ := json.Marshal(outputTraintuple)
	json.Unmarshal(oo, &element)
	element["key"] = inp.Key
	mPayload = map[string]interface{}{"traintuple": element}
	// get certified and non-certified testtuples related to traintuple
	var nonCertifiedTesttuples []map[string]interface{}
	testtupleKeys, err := getKeysFromComposite(stub, "testtuple~traintuple~certified~key", []string{"testtuple", inp.Key})
	if err != nil {
		return
	}
	for _, testtupleKey := range testtupleKeys {
		// get testtuple and serialize it
		var testtuple map[string]interface{}
		if err = getElementStruct(stub, testtupleKey, &testtuple); err != nil {
			return
		}
		testtuple["key"] = testtupleKey
		if testtuple["certified"] == true {
			mPayload["testtuple"] = testtuple
		} else {
			nonCertifiedTesttuples = append(nonCertifiedTesttuples, testtuple)
		}
		mPayload["nonCertifiedTesttuples"] = nonCertifiedTesttuples
	}
	return
}

// queryModels returns all traintuples and associated testuples
// TODO
func queryModels(stub shim.ChaincodeStubInterface, args []string) (elements []map[string]interface{}, err error) {
	if len(args) != 0 {
		err = fmt.Errorf("incorrect number of arguments, expecting nothing")
		return
	}

	traintupleKeys, err := getKeysFromComposite(stub, "traintuple~algo~key", []string{"traintuple"})
	if err != nil {
		return
	}
	for _, traintupleKey := range traintupleKeys {
		var outputTraintuple outputTraintuple
		element := make(map[string]interface{})
		traintuple := make(map[string]interface{})
		// get traintuple
		outputTraintuple, err = getOutputTraintuple(stub, traintupleKey)
		if err != nil {
			return
		}
		oo, _ := json.Marshal(outputTraintuple)
		json.Unmarshal(oo, &traintuple)
		traintuple["key"] = traintupleKey
		element["traintuple"] = traintuple

		// get testtuple related to traintuple
		var testtupleKeys []string
		testtupleKeys, err = getKeysFromComposite(stub, "testtuple~traintuple~certified~key", []string{"testtuple", traintupleKey, "true"})
		if err != nil {
			return
		}
		if len(testtupleKeys) == 1 {
			// get testtuple and serialize it
			testtupleKey := testtupleKeys[0]
			testtuple := make(map[string]interface{})
			if err = getElementStruct(stub, testtupleKey, &testtuple); err != nil {
				return
			}
			testtuple["key"] = testtupleKey
			element["testtuple"] = testtuple
		}
		elements = append(elements, element)
	}
	return
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
	outputTraintuple.Fill(stub, traintuple, traintupleKey)
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

// checkNewStatus verifies that the new status is consistent with the tuple current status
func (traintuple *Traintuple) checkNewStatus(stub shim.ChaincodeStubInterface, status string) error {
	// get worker
	worker, err := getDataManagerOwner(stub, traintuple.Dataset.DataManagerKey)
	if err != nil {
		return err
	}
	// check validity of worker and change of status
	if err := checkUpdateTuple(stub, worker, traintuple.Status, status); err != nil {
		return err
	}

	return nil
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
func updateWaitingTraintuples(stub shim.ChaincodeStubInterface, modelTraintupleKey string, parentStatus string) error {

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

		// remove associated composite key
		if err := traintuple.removeModelCompositeKey(stub, modelTraintupleKey); err != nil {
			return err
		}

		// traintuple is already failed, don't update it
		if traintuple.Status == "failed" {
			continue
		}
		if traintuple.Status != "waiting" {
			return fmt.Errorf("traintuple %s has invalid status : '%s' instead of waiting", traintupleKey, traintuple.Status)
		}
		// get traintuple new status
		// TODO use checkNewStatus
		var newStatus string
		if parentStatus == "failed" {
			newStatus = parentStatus
		} else if parentStatus == "done" {
			ready, err := traintuple.isReady(stub, modelTraintupleKey)
			if err != nil {
				return err
			}
			if ready {
				newStatus = "todo"
			}
		}

		// commit new status
		if newStatus != "" {
			if err := traintuple.commitUpdate(stub, traintupleKey, newStatus); err != nil {
				return err
			}
			if newStatus == "todo" {
				out := outputTraintuple{}
				err = out.Fill(stub, traintuple, traintupleKey)
				if err != nil {
					return err
				}
				err = SetEvent(stub, "traintuple-ready", out)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// isReady checks if inModels of a traintuple have been trained, except the newDoneTraintupleKey (since the transaction is not commited)
// and updates the traintuple status if necessary
func (traintuple *Traintuple) isReady(stub shim.ChaincodeStubInterface, newDoneTraintupleKey string) (ready bool, err error) {
	for _, key := range traintuple.InModelKeys {
		// don't check newly done traintuple
		if key == newDoneTraintupleKey {
			continue
		}
		tt := Traintuple{}
		if err := getElementStruct(stub, key, &tt); err != nil {
			return false, err
		}
		if tt.Status != "done" {
			return false, nil
		}
	}
	return true, nil
}

// removeModelCompositeKey remove the Model key state of a traintuple
func (traintuple *Traintuple) removeModelCompositeKey(stub shim.ChaincodeStubInterface, modelKey string) error {
	indexName := "traintuple~inModel~key"
	compositeKey, err := stub.CreateCompositeKey(indexName, []string{"traintuple", modelKey, traintuple.FLTask})

	if err != nil {
		return fmt.Errorf("failed to recreate composite key to update traintuple %s with inModel %s - %s", traintuple.FLTask, modelKey, err.Error())
	}

	if err := stub.DelState(compositeKey); err != nil {
		return fmt.Errorf("failed to delete associated composite key to update traintuple %s with inModel %s - %s", traintuple.FLTask, modelKey, err.Error())
	}
	return nil
}

// commitUpdate update the traintuple after a status update in the ledger
func (traintuple *Traintuple) commitUpdate(stub shim.ChaincodeStubInterface, traintupleKey string, newStatus string) error {
	if traintuple.Status == newStatus {
		return fmt.Errorf("cannot update traintuple %s - status already %s", traintupleKey, newStatus)
	}
	oldStatus := traintuple.Status
	traintuple.Status = newStatus
	traintupleBytes, _ := json.Marshal(traintuple)
	if err := stub.PutState(traintupleKey, traintupleBytes); err != nil {
		return fmt.Errorf("failed to update traintuple %s - %s", traintupleKey, err.Error())
	}

	// get worker
	worker, err := getDataManagerOwner(stub, traintuple.Dataset.DataManagerKey)
	if err != nil {
		return err
	}

	// update associated composite keys
	indexName := "traintuple~worker~status~key"
	oldAttributes := []string{"traintuple", worker, oldStatus, traintupleKey}
	newAttributes := []string{"traintuple", worker, traintuple.Status, traintupleKey}
	if err = updateCompositeKey(stub, indexName, oldAttributes, newAttributes); err != nil {
		return err
	}
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
	for _, testtupleKey := range testtupleKeys {
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

		if testtuple.Status == "todo" {
			out := outputTesttuple{}
			out.Fill(testtupleKey, testtuple)
			err = SetEvent(stub, "testtuple-ready", out)
			if err != nil {
				return err
			}
		}

	}
	return nil
}

// getTraintuplesPayload takes as input a list of keys and returns a paylaod containing a list of associated retrieved elements
func getTraintuplesPayload(stub shim.ChaincodeStubInterface, traintupleKeys []string) ([]map[string]interface{}, error) {

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
	return elements, nil
}

// HashForKey to generate key for an asset
func HashForKey(stub shim.ChaincodeStubInterface, objectType string, hashElements ...string) (key string, err error) {
	toHash := objectType
	sort.Strings(hashElements)
	for _, element := range hashElements {
		toHash += "," + element
	}
	sum := sha256.Sum256([]byte(toHash))
	key = hex.EncodeToString(sum[:])
	if bytes, stubErr := stub.GetState(key); bytes != nil {
		err = fmt.Errorf("this %s already exists (tkey: %s)", objectType, key)
	} else if stubErr != nil {
		return key, stubErr
	}
	return
}
