package main

import (
	"chaincode/errors"
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

const (
	StatusDoing   = "doing"
	StatusTodo    = "todo"
	StatusWaiting = "waiting"
	StatusFailed  = "failed"
	StatusDone    = "done"
)

// -------------------------------------------------------------------------------------------
// Methods on receivers traintuple and testuples
// -------------------------------------------------------------------------------------------

// Set is a method of the receiver Traintuple. It checks the validity of inputTraintuple and uses its fields to set the Traintuple
func (traintuple *Traintuple) Set(stub shim.ChaincodeStubInterface, inp inputTraintuple) (traintupleKey string, err error) {

	validate := validator.New()
	if err = validate.Struct(inp); err != nil {
		err = errors.BadRequest(err, "invalid inputs to update data")
		return
	}

	// TODO later: check permissions
	// find associated creator and check permissions (TODO later)
	creator, err := getTxCreator(stub)
	if err != nil {
		return
	}
	traintuple.AssetType = TraintupleType
	traintuple.Creator = creator
	traintuple.Permissions = "all"
	traintuple.Tag = inp.Tag
	// check if algo exists
	if _, err = getElementBytes(stub, inp.AlgoKey); err != nil {
		err = errors.BadRequest(err, "could not retrieve algo with key %s", inp.AlgoKey)
		return
	}
	traintuple.AlgoKey = inp.AlgoKey

	// check objective and add it
	obj := Objective{}
	if err = getElementStruct(stub, inp.ObjectiveKey, &obj); err != nil {
		err = errors.BadRequest(err, "could not retrieve objective with key %s", inp.ObjectiveKey)
		return
	}
	traintuple.ObjectiveKey = inp.ObjectiveKey

	// check if InModels is empty or if mentionned models do exist and fill inModels
	status := StatusTodo
	parentTraintupleKeys := strings.Split(strings.Replace(inp.InModels, " ", "", -1), ",")
	for _, parentTraintupleKey := range parentTraintupleKeys {
		if parentTraintupleKey == "" {
			break
		}
		parentTraintuple := Traintuple{}
		if err = getElementStruct(stub, parentTraintupleKey, &parentTraintuple); err != nil {
			err = errors.BadRequest(err, "could not retrieve parent traintuple with key %s %d", parentTraintupleKeys, len(parentTraintupleKeys))
			return
		}
		// set traintuple to waiting if one of the parent traintuples is not done
		if parentTraintuple.OutModel == nil {
			status = StatusWaiting
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
		err = errors.BadRequest("not possible to create a traintuple with test only data")
		return
	}

	// fill traintuple.Dataset from dataManager and dataSample
	if _, err = getElementBytes(stub, inp.DataManagerKey); err != nil {
		err = errors.BadRequest(err, "could not retrieve dataManager with key %s", inp.DataManagerKey)
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
		err = errors.Conflict(err)
		return
	}

	// check FLTask and Rank and set it when required
	if inp.Rank == "" {
		if inp.FLTask != "" {
			err = errors.BadRequest("invalit inputs, a FLTask should have a rank")
			return
		}
	} else {
		traintuple.Rank, err = strconv.Atoi(inp.Rank)
		if err != nil {
			return
		}
		if inp.FLTask == "" {
			if traintuple.Rank != 0 {
				err = errors.BadRequest("invalid inputs, a new FLTask should have a rank 0")
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
				err = errors.BadRequest("cannot find the FLTask %s", inp.FLTask)
				return
			}
			for _, ttKey := range ttKeys {
				FLTraintuple := Traintuple{}
				err = getElementStruct(stub, ttKey, &FLTraintuple)
				if err != nil {
					return
				} else if FLTraintuple.AlgoKey != inp.AlgoKey {
					err = errors.BadRequest("previous traintuple for FLTask %s does not have the same algo key %s", inp.FLTask, inp.AlgoKey)
					return
				}
			}

			attributes = []string{"traintuple", inp.FLTask, traintuple.Dataset.Worker, inp.Rank}
			ttKeys, err = getKeysFromComposite(stub, "traintuple~fltask~worker~rank~key", attributes)
			if err != nil {
				return
			} else if len(ttKeys) > 0 {
				err = errors.BadRequest("FLTask %s with worker %s rank %d already exists", inp.FLTask, traintuple.Dataset.Worker, traintuple.Rank)
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
		return "", errors.BadRequest(err, "invalid inputs to update data")
	}

	// check associated traintuple
	traintuple := Traintuple{}
	if err = getElementStruct(stub, inp.TraintupleKey, &traintuple); err != nil {
		return testtupleKey, errors.BadRequest(err, "could not retrieve traintuple with key %s", inp.TraintupleKey)
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
	testtuple.AssetType = TesttupleType
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
	case StatusDone:
		testtuple.Status = StatusTodo
	case StatusFailed:
		err = errors.BadRequest(
			"could not register this testtuple, the traintuple %s has a failed status",
			inp.TraintupleKey)
		return
	default:
		testtuple.Status = StatusWaiting
	}

	// Get test dataset from objective
	objective := Objective{}
	if err = getElementStruct(stub, testtuple.Objective.Key, &objective); err != nil {
		return testtupleKey, errors.BadRequest(err, "could not retrieve objective with key %s", testtuple.Objective.Key)
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
		return testtupleKey, errors.BadRequest("invalid input: dataManagerKey and dataSampleKey should be provided together")
	} else if objective.TestDataset != nil {
		dataSampleKeys = objectiveDataSampleKeys
		dataManagerKey = objectiveDataManagerKey
		testtuple.Certified = true
	} else {
		err = errors.BadRequest("can not create a certified testtuple, no data associated with objective %s", testtuple.Objective.Key)
		return
	}
	// retrieve dataManager owner
	dataManager := DataManager{}
	if err = getElementStruct(stub, dataManagerKey, &dataManager); err != nil {
		return testtupleKey, errors.BadRequest(err, "could not retrieve dataManager with key %s", dataManagerKey)
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
	if err != nil {
		err = errors.Conflict(err)
		return
	}

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

	// https://github.com/hyperledger/fabric/blob/release-1.4/core/chaincode/shim/interfaces.go#L339:L343
	// We can only send one event per transaction
	// https://stackoverflow.com/questions/50344232/not-able-to-set-multiple-events-in-chaincode-per-transaction-getting-only-last
	event := TuplesEvent{}
	event.SetTraintuples(out)

	err = SetEvent(stub, "tuples-updated", event)
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

	// https://github.com/hyperledger/fabric/blob/release-1.4/core/chaincode/shim/interfaces.go#L339:L343
	// We can only send one event per transaction
	// https://stackoverflow.com/questions/50344232/not-able-to-set-multiple-events-in-chaincode-per-transaction-getting-only-last
	event := TuplesEvent{}
	event.SetTesttuples(out)

	err = SetEvent(stub, "tuples-updated", event)
	if err != nil {
		return nil, err
	}

	return map[string]string{"key": testtupleKey}, nil
}

// logStartTrain modifies a traintuple by changing its status from todo to doing
func logStartTrain(stub shim.ChaincodeStubInterface, args []string) (outputTraintuple outputTraintuple, err error) {
	inp := inputHashe{}
	err = AssetFromJSON(args[0], &inp)
	if err != nil {
		return
	}

	// get traintuple, check validity of the update
	traintuple := Traintuple{}
	if err = getElementStruct(stub, inp.Key, &traintuple); err != nil {
		return
	}
	if err = traintuple.commitStatusUpdate(stub, inp.Key, StatusDoing); err != nil {
		return
	}
	outputTraintuple.Fill(stub, traintuple, inp.Key)
	return
}

// logStartTest modifies a testtuple by changing its status from todo to doing
func logStartTest(stub shim.ChaincodeStubInterface, args []string) (outputTesttuple outputTesttuple, err error) {
	inp := inputHashe{}
	err = AssetFromJSON(args[0], &inp)
	if err != nil {
		return
	}

	// get testtuple, check validity of the update, and update its status
	testtuple := Testtuple{}
	if err = getElementStruct(stub, inp.Key, &testtuple); err != nil {
		return
	}
	if err = testtuple.commitStatusUpdate(stub, inp.Key, StatusDoing); err != nil {
		return
	}
	outputTesttuple.Fill(inp.Key, testtuple)
	return
}

// logSuccessTrain modifies a traintuple by changing its status from doing to done
// reports logs and associated performances
func logSuccessTrain(stub shim.ChaincodeStubInterface, args []string) (outputTraintuple outputTraintuple, err error) {
	inp := inputLogSuccessTrain{}
	err = AssetFromJSON(args[0], &inp)
	if err != nil {
		return
	}
	traintupleKey := inp.Key

	// get, update and commit traintuple
	traintuple := Traintuple{}
	if err = getElementStruct(stub, traintupleKey, &traintuple); err != nil {
		return
	}
	traintuple.Perf = inp.Perf
	traintuple.OutModel = &HashDress{
		Hash:           inp.OutModel.Hash,
		StorageAddress: inp.OutModel.StorageAddress}
	traintuple.Log += inp.Log

	if err = traintuple.commitStatusUpdate(stub, traintupleKey, StatusDone); err != nil {
		return
	}

	// update depending tuples
	traintuples_event, err := traintuple.updateTraintupleChildren(stub, traintupleKey)
	if err != nil {
		return
	}

	testtuples_event, err := traintuple.updateTesttupleChildren(stub, traintupleKey)
	if err != nil {
		return
	}

	outputTraintuple.Fill(stub, traintuple, inp.Key)

	// https://github.com/hyperledger/fabric/blob/release-1.4/core/chaincode/shim/interfaces.go#L339:L343
	// We can only send one event per transaction
	// https://stackoverflow.com/questions/50344232/not-able-to-set-multiple-events-in-chaincode-per-transaction-getting-only-last
	event := TuplesEvent{}
	event.SetTraintuples(traintuples_event...)
	event.SetTesttuples(testtuples_event...)

	err = SetEvent(stub, "tuples-updated", event)
	if err != nil {
		return
	}

	return
}

// logSuccessTest modifies a testtuple by changing its status to done, reports perf and logs
func logSuccessTest(stub shim.ChaincodeStubInterface, args []string) (outputTesttuple outputTesttuple, err error) {
	inp := inputLogSuccessTest{}
	err = AssetFromJSON(args[0], &inp)
	if err != nil {
		return
	}

	testtuple := Testtuple{}
	if err = getElementStruct(stub, inp.Key, &testtuple); err != nil {
		return
	}

	testtuple.Dataset.Perf = inp.Perf
	testtuple.Log += inp.Log

	if err = testtuple.commitStatusUpdate(stub, inp.Key, StatusDone); err != nil {
		return
	}
	outputTesttuple.Fill(inp.Key, testtuple)
	return
}

// logFailTrain modifies a traintuple by changing its status to fail and reports associated logs
func logFailTrain(stub shim.ChaincodeStubInterface, args []string) (outputTraintuple outputTraintuple, err error) {
	inp := inputLogFailTrain{}
	err = AssetFromJSON(args[0], &inp)
	if err != nil {
		return
	}

	// get, update and commit traintuple
	traintuple := Traintuple{}
	if err = getElementStruct(stub, inp.Key, &traintuple); err != nil {
		return
	}
	traintuple.Log += inp.Log

	if err = traintuple.commitStatusUpdate(stub, inp.Key, StatusFailed); err != nil {
		return
	}

	// update depending tuples
	testtuples_event, err := traintuple.updateTesttupleChildren(stub, inp.Key)
	if err != nil {
		return
	}

	traintuples_event, err := traintuple.updateTraintupleChildren(stub, inp.Key)
	if err != nil {
		return
	}

	// https://github.com/hyperledger/fabric/blob/release-1.4/core/chaincode/shim/interfaces.go#L339:L343
	// We can only send one event per transaction
	// https://stackoverflow.com/questions/50344232/not-able-to-set-multiple-events-in-chaincode-per-transaction-getting-only-last
	event := TuplesEvent{}
	event.SetTraintuples(traintuples_event...)
	event.SetTesttuples(testtuples_event...)

	err = SetEvent(stub, "tuples-updated", event)
	if err != nil {
		return
	}

	return
}

// logFailTest modifies a testtuple by changing its status to fail and reports associated logs
func logFailTest(stub shim.ChaincodeStubInterface, args []string) (outputTesttuple outputTesttuple, err error) {
	inp := inputLogFailTest{}
	err = AssetFromJSON(args[0], &inp)
	if err != nil {
		return
	}

	// get, update and commit testtuple
	testtuple := Testtuple{}
	if err = getElementStruct(stub, inp.Key, &testtuple); err != nil {
		return
	}

	testtuple.Log += inp.Log

	if err = testtuple.commitStatusUpdate(stub, inp.Key, StatusFailed); err != nil {
		return
	}
	outputTesttuple.Fill(inp.Key, testtuple)
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
	if traintuple.AssetType != TraintupleType {
		err = errors.NotFound("no element with key %s", inp.Key)
		return
	}
	outputTraintuple.Fill(stub, traintuple, inp.Key)
	return
}

// queryTraintuples returns all traintuples
func queryTraintuples(stub shim.ChaincodeStubInterface, args []string) (outTraintuples []outputTraintuple, err error) {
	outTraintuples = []outputTraintuple{}

	if len(args) != 0 {
		err = errors.BadRequest("incorrect number of arguments, expecting nothing")
		return
	}
	elementsKeys, err := getKeysFromComposite(stub, "traintuple~algo~key", []string{"traintuple"})
	if err != nil {
		return
	}
	for _, key := range elementsKeys {
		var outputTraintuple outputTraintuple
		outputTraintuple, err = getOutputTraintuple(stub, key)
		if err != nil {
			return
		}
		outTraintuples = append(outTraintuples, outputTraintuple)
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
	if testtuple.AssetType != TesttupleType {
		err = errors.NotFound("no element with key %s", inp.Key)
		return
	}
	out.Fill(inp.Key, testtuple)
	return
}

// queryTesttuples returns all testtuples of the ledger
func queryTesttuples(stub shim.ChaincodeStubInterface, args []string) (outTesttuples []outputTesttuple, err error) {
	outTesttuples = []outputTesttuple{}

	if len(args) != 0 {
		err = errors.BadRequest("incorrect number of arguments, expecting nothing")
		return
	}
	var indexName = "testtuple~traintuple~certified~key"
	elementsKeys, err := getKeysFromComposite(stub, indexName, []string{"testtuple"})
	if err != nil {
		err = fmt.Errorf("issue getting keys from composite key %s - %s", indexName, err.Error())
		return
	}
	for _, key := range elementsKeys {
		var out outputTesttuple
		out, err = getOutputTesttuple(stub, key)
		if err != nil {
			return
		}
		outTesttuples = append(outTesttuples, out)
	}
	return
}

// queryModelDetails returns info about the testtuple and algo related to a traintuple
func queryModelDetails(stub shim.ChaincodeStubInterface, args []string) (outModelDetails outputModelDetails, err error) {
	inp := inputHashe{}
	err = AssetFromJSON(args[0], &inp)
	if err != nil {
		return
	}

	// get associated traintuple
	outModelDetails.Traintuple, err = getOutputTraintuple(stub, inp.Key)
	if err != nil {
		return
	}

	// get certified and non-certified testtuples related to traintuple
	testtupleKeys, err := getKeysFromComposite(stub, "testtuple~traintuple~certified~key", []string{"testtuple", inp.Key})
	if err != nil {
		return
	}
	for _, testtupleKey := range testtupleKeys {
		// get testtuple and serialize it
		var outputTesttuple outputTesttuple
		outputTesttuple, err = getOutputTesttuple(stub, testtupleKey)
		if err != nil {
			return
		}

		if outputTesttuple.Certified {
			outModelDetails.Testtuple = outputTesttuple
		} else {
			outModelDetails.NonCertifiedTesttuples = append(outModelDetails.NonCertifiedTesttuples, outputTesttuple)
		}
	}
	return
}

// queryModels returns all traintuples and associated testuples
func queryModels(stub shim.ChaincodeStubInterface, args []string) (outModels []outputModel, err error) {
	outModels = []outputModel{}

	if len(args) != 0 {
		err = errors.BadRequest("incorrect number of arguments, expecting nothing")
		return
	}

	traintupleKeys, err := getKeysFromComposite(stub, "traintuple~algo~key", []string{"traintuple"})
	if err != nil {
		return
	}
	for _, traintupleKey := range traintupleKeys {
		var outputModel outputModel

		// get traintuple
		outputModel.Traintuple, err = getOutputTraintuple(stub, traintupleKey)
		if err != nil {
			return
		}

		// get associated testtuple
		var testtupleKeys []string
		testtupleKeys, err = getKeysFromComposite(stub, "testtuple~traintuple~certified~key", []string{"testtuple", traintupleKey, "true"})
		if err != nil {
			return
		}
		if len(testtupleKeys) == 1 {
			// get testtuple and serialize it
			testtupleKey := testtupleKeys[0]
			outputModel.Testtuple, err = getOutputTesttuple(stub, testtupleKey)
			if err != nil {
				return
			}
		}
		outModels = append(outModels, outputModel)
	}
	return
}

// --------------------------------------------------------------
// Utils for smartcontracts related to traintuples and testtuples
// --------------------------------------------------------------

// getOutputTraintuple takes as input a traintuple key and returns the outputTraintuple
func getOutputTraintuple(stub shim.ChaincodeStubInterface, traintupleKey string) (outTraintuple outputTraintuple, err error) {
	traintuple := Traintuple{}
	if err = getElementStruct(stub, traintupleKey, &traintuple); err != nil {
		return
	}
	outTraintuple.Fill(stub, traintuple, traintupleKey)
	return
}

// getOutputTraintuples takes as input a list of keys and returns a paylaod containing a list of associated retrieved elements
func getOutputTraintuples(stub shim.ChaincodeStubInterface, traintupleKeys []string) (outTraintuples []outputTraintuple, err error) {
	for _, key := range traintupleKeys {
		var outputTraintuple outputTraintuple
		outputTraintuple, err = getOutputTraintuple(stub, key)
		if err != nil {
			return
		}
		outTraintuples = append(outTraintuples, outputTraintuple)
	}
	return
}

// getOutputTesttuple takes as input a testtuple key and returns the outputTesttuple
func getOutputTesttuple(stub shim.ChaincodeStubInterface, testtupleKey string) (outTesttuple outputTesttuple, err error) {
	testtuple := Testtuple{}
	if err = getElementStruct(stub, testtupleKey, &testtuple); err != nil {
		return
	}
	outTesttuple.Fill(testtupleKey, testtuple)
	return
}

// getOutputTesttuples takes as input a list of keys and returns a paylaod containing a list of associated retrieved elements
func getOutputTesttuples(stub shim.ChaincodeStubInterface, testtupleKeys []string) (outTesttuples []outputTesttuple, err error) {
	for _, key := range testtupleKeys {
		var outputTesttuple outputTesttuple
		outputTesttuple, err = getOutputTesttuple(stub, key)
		if err != nil {
			return
		}
		outTesttuples = append(outTesttuples, outputTesttuple)
	}
	return
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
		return fmt.Errorf("%s is not allowed to update tuple (%s)", txCreator, worker)
	}
	statusPossibilities := map[string]string{
		StatusWaiting: StatusTodo,
		StatusTodo:    StatusDoing,
		StatusDoing:   StatusDone}
	if statusPossibilities[oldStatus] != newStatus && newStatus != StatusFailed {
		return errors.BadRequest("cannot change status from %s to %s", oldStatus, newStatus)
	}
	return nil
}

// validateNewStatus verifies that the new status is consistent with the tuple current status
func (traintuple *Traintuple) validateNewStatus(stub shim.ChaincodeStubInterface, status string) error {
	// check validity of worker and change of status
	if err := checkUpdateTuple(stub, traintuple.Dataset.Worker, traintuple.Status, status); err != nil {
		return err
	}

	return nil
}

// validateNewStatus verifies that the new status is consistent with the tuple current status
func (testtuple *Testtuple) validateNewStatus(stub shim.ChaincodeStubInterface, status string) error {
	// check validity of worker and change of status
	if err := checkUpdateTuple(stub, testtuple.Dataset.Worker, testtuple.Status, status); err != nil {
		return err
	}
	return nil
}

// updateTraintupleChildren updates the status of waiting trainuples  InModels of traintuples once they have been trained (succesfully or failed)
func (parentTraintuple *Traintuple) updateTraintupleChildren(stub shim.ChaincodeStubInterface, parentTraintupleKey string) ([]outputTraintuple, error) {

	// tuples to be sent in event
	otuples := []outputTraintuple{}

	// get traintuples having as inModels the input traintuple
	indexName := "traintuple~inModel~key"
	traintupleKeys, err := getKeysFromComposite(stub, indexName, []string{"traintuple", parentTraintupleKey})
	if err != nil {
		return otuples, fmt.Errorf("error while getting associated traintuples to update their inModel")
	}
	for _, traintupleKey := range traintupleKeys {
		// get and update traintuple
		traintuple := Traintuple{}
		if err := getElementStruct(stub, traintupleKey, &traintuple); err != nil {
			return otuples, err
		}

		// remove associated composite key
		if err := traintuple.removeModelCompositeKey(stub, parentTraintupleKey); err != nil {
			return otuples, err
		}

		// traintuple is already failed, don't update it
		if traintuple.Status == StatusFailed {
			continue
		}

		if traintuple.Status != StatusWaiting {
			return otuples, fmt.Errorf("traintuple %s has invalid status : '%s' instead of waiting", traintupleKey, traintuple.Status)
		}

		// get traintuple new status
		var newStatus string
		if parentTraintuple.Status == StatusFailed {
			newStatus = StatusFailed
		} else if parentTraintuple.Status == StatusDone {
			ready, err := traintuple.isReady(stub, parentTraintupleKey)
			if err != nil {
				return otuples, err
			}
			if ready {
				newStatus = StatusTodo
			}
		}

		// commit new status
		if newStatus == "" {
			continue
		}
		if err := traintuple.commitStatusUpdate(stub, traintupleKey, newStatus); err != nil {
			return otuples, err
		}
		if newStatus == StatusTodo {
			out := outputTraintuple{}
			err = out.Fill(stub, traintuple, traintupleKey)
			if err != nil {
				return otuples, err
			}
			otuples = append(otuples, out)
		}
	}
	return otuples, nil
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
		if tt.Status != StatusDone {
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

// commitStatusUpdate update the traintuple status in the ledger
func (traintuple *Traintuple) commitStatusUpdate(stub shim.ChaincodeStubInterface, traintupleKey string, newStatus string) error {
	if traintuple.Status == newStatus {
		return fmt.Errorf("cannot update traintuple %s - status already %s", traintupleKey, newStatus)
	}

	if err := traintuple.validateNewStatus(stub, newStatus); err != nil {
		return err
	}

	oldStatus := traintuple.Status
	traintuple.Status = newStatus
	traintupleBytes, _ := json.Marshal(traintuple)
	if err := stub.PutState(traintupleKey, traintupleBytes); err != nil {
		return fmt.Errorf("failed to update traintuple %s - %s", traintupleKey, err.Error())
	}

	// update associated composite keys
	indexName := "traintuple~worker~status~key"
	oldAttributes := []string{"traintuple", traintuple.Dataset.Worker, oldStatus, traintupleKey}
	newAttributes := []string{"traintuple", traintuple.Dataset.Worker, traintuple.Status, traintupleKey}
	if err := updateCompositeKey(stub, indexName, oldAttributes, newAttributes); err != nil {
		return err
	}
	return nil
}

// updateTesttupleChildren update testtuples status associated with a done or failed traintuple
func (parentTraintuple *Traintuple) updateTesttupleChildren(stub shim.ChaincodeStubInterface, parentTraintupleKey string) ([]outputTesttuple, error) {

	otuples := []outputTesttuple{}

	var newStatus string
	if parentTraintuple.Status == StatusFailed {
		newStatus = StatusFailed
	} else if parentTraintuple.Status == StatusDone {
		newStatus = StatusTodo
	} else {
		return otuples, nil
	}

	indexName := "testtuple~traintuple~certified~key"
	// get testtuple associated with this traintuple and updates its status
	testtupleKeys, err := getKeysFromComposite(stub, indexName, []string{"testtuple", parentTraintupleKey})
	if err != nil {
		return otuples, err
	}
	for _, testtupleKey := range testtupleKeys {
		// get and update testtuple
		testtuple := Testtuple{}
		if err := getElementStruct(stub, testtupleKey, &testtuple); err != nil {
			return otuples, err
		}
		testtuple.Model = &Model{
			TraintupleKey: parentTraintupleKey,
		}

		if newStatus == StatusTodo {
			testtuple.Model.Hash = parentTraintuple.OutModel.Hash
			testtuple.Model.StorageAddress = parentTraintuple.OutModel.StorageAddress
		}

		if err := testtuple.commitStatusUpdate(stub, testtupleKey, newStatus); err != nil {
			return otuples, err
		}

		if newStatus == StatusTodo {
			out := outputTesttuple{}
			out.Fill(testtupleKey, testtuple)
			otuples = append(otuples, out)
		}
	}
	return otuples, nil
}

// commitStatusUpdate update the testtuple status in the ledger
func (testtuple *Testtuple) commitStatusUpdate(stub shim.ChaincodeStubInterface, testtupleKey string, newStatus string) error {
	if err := testtuple.validateNewStatus(stub, newStatus); err != nil {
		return err
	}

	oldStatus := testtuple.Status
	testtuple.Status = newStatus

	testtupleBytes, _ := json.Marshal(testtuple)
	if err := stub.PutState(testtupleKey, testtupleBytes); err != nil {
		return fmt.Errorf("failed to update testtuple status to %s with key %s", newStatus, testtupleKey)
	}

	// update associated composite key
	indexName := "testtuple~worker~status~key"
	oldAttributes := []string{"testtuple", testtuple.Dataset.Worker, oldStatus, testtupleKey}
	newAttributes := []string{"testtuple", testtuple.Dataset.Worker, testtuple.Status, testtupleKey}
	if err := updateCompositeKey(stub, indexName, oldAttributes, newAttributes); err != nil {
		return err
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
