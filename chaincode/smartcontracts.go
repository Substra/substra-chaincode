package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

const lenKey int = 64

// registerChallenge stores a new challenge in the ledger.
// If the key exists, it will override the value with the new one
func registerChallenge(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	expectedArgs := getFieldNames(&inputChallenge{})
	if nbArgs := len(expectedArgs); nbArgs != len(args) {
		return nil, fmt.Errorf("incorrect arguments, expecting %d args: %s", nbArgs, strings.Join(expectedArgs, ", "))
	}

	// convert input strings args to input struct inputChallenge
	inpc := inputChallenge{}
	stringToInputStruct(args, &inpc)
	// check validity of input args and convert it to Challenge
	challenge := Challenge{}
	challengeKey, datasetKey, err := challenge.Set(stub, inpc)
	if err != nil {
		return nil, err
	}
	// check challenge is not already in ledger
	if elementBytes, _ := stub.GetState(challengeKey); elementBytes != nil {
		return nil, fmt.Errorf("challenge with this description already exists - %s", string(elementBytes))
	}
	// submit to ledger
	challengeBytes, _ := json.Marshal(challenge)
	if err := stub.PutState(challengeKey, challengeBytes); err != nil {
		return nil, fmt.Errorf("failed to submit to ledger the challenge with key %s, error is %s", challengeKey, err.Error())
	}
	// create composite key
	if err := createCompositeKey(stub, "challenge~owner~key", []string{"challenge", challenge.Owner, challengeKey}); err != nil {
		return nil, err
	}
	// add challenge to dataset
	err = addChallengeDataset(stub, datasetKey, challengeKey)
	// return []byte(challengeKey), err
	return []byte(challengeKey), err
}

// registerDataset stores a new dataset in the ledger.
// If the key exists, it will override the value with the new one
func registerDataset(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	expectedArgs := getFieldNames(&inputDataset{})
	if nbArgs := len(expectedArgs); nbArgs != len(args) {
		return nil, fmt.Errorf("incorrect arguments, expecting %d args: %s", nbArgs, strings.Join(expectedArgs, ", "))
	}
	// convert input strings args to input struct inputDataset
	inp := inputDataset{}
	stringToInputStruct(args, &inp)
	// check validity of input args and convert it to a Dataset
	dataset := Dataset{}
	datasetKey, challengeKeys, err := dataset.Set(stub, inp)
	if err != nil {
		return nil, err
	}
	// check dataset is not already in ledger
	if elementBytes, _ := stub.GetState(datasetKey); elementBytes != nil {
		return nil, fmt.Errorf("dataset with this opener already exists")
	}
	// submit to ledger
	datasetBytes, _ := json.Marshal(dataset)
	err = stub.PutState(datasetKey, datasetBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to add dataset with opener hash %s, error is %s", inp.OpenerHash, err.Error())
	}
	// create composite keys (one for each associated challenge) to find data associated with a challenge
	indexName := "dataset~challenge~key"
	for _, challengeKey := range challengeKeys {
		err = createCompositeKey(stub, indexName, []string{"dataset", challengeKey, datasetKey})
		if err != nil {
			return nil, err
		}
	}
	// create composite key to find dataset associated with a owner
	err = createCompositeKey(stub, "dataset~owner~key", []string{"dataset", dataset.Owner, datasetKey})
	if err != nil {
		return nil, err
	}
	return []byte(datasetKey), nil
}

// registerData stores new data in the ledger (one or more).
// If the key exists, it will override the value with the new one
func registerData(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	expectedArgs := getFieldNames(&inputData{})
	if nbArgs := len(expectedArgs); nbArgs != len(args) {
		return nil, fmt.Errorf("incorrect arguments, expecting %d args: %s", nbArgs, strings.Join(expectedArgs, ", "))
	}

	// convert input strings args to input struct inputData
	inp := inputData{}
	stringToInputStruct(args, &inp)
	// check validity of input args and update dataset associated with data
	dataset := Dataset{}
	datasetKey, dataHashes, testOnly, err := dataset.Update(stub, inp)
	if err != nil {
		return nil, err
	}
	// submit updated dataset to ledger
	datasetBytes, _ := json.Marshal(dataset)
	err = stub.PutState(datasetKey, datasetBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to update %s size and nb data", datasetKey)
	}
	// store each added data in the ledger
	var dataKeys string
	suffix := ", "
	for _, dataHash := range dataHashes {
		dataKeys = dataKeys + dataHash + suffix
		// check data is not already in ledgert
		if elementBytes, _ := stub.GetState(dataHash); elementBytes != nil {
			return nil, fmt.Errorf("data with this hash already exists")
		}
		// create data object
		var data = Data{
			DatasetKey: datasetKey,
			TestOnly:   testOnly}
		dataBytes, _ := json.Marshal(data)
		err = stub.PutState(dataHash, dataBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to add data with hash %s", dataHash)
		}
		// create composite keys to find all data associated with a dataset and only test or train data
		if err = createCompositeKey(stub, "data~dataset~key", []string{"data", datasetKey, dataHash}); err != nil {
			return nil, err
		}
		if err = createCompositeKey(stub, "data~dataset~testOnly~key", []string{"data", datasetKey, strconv.FormatBool(testOnly), dataHash}); err != nil {
			return nil, err
		}
	}
	// return added data keys
	dataKeys = strings.TrimSuffix(dataKeys, suffix)
	return []byte(dataKeys), nil
}

// registerAlgo stores a new algo in the ledger.
// If the key exists, it will override the value with the new one
func registerAlgo(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	expectedArgs := getFieldNames(&inputAlgo{})
	if nbArgs := len(expectedArgs); nbArgs != len(args) {
		return nil, fmt.Errorf("incorrect arguments, expecting %d args: %s", nbArgs, strings.Join(expectedArgs, ", "))
	}

	// convert input strings args to input struct inputAlgo
	inp := inputAlgo{}
	stringToInputStruct(args, &inp)
	// check validity of input args and convert it to Algo
	algo := Algo{}
	algoKey, err := algo.Set(stub, inp)
	if err != nil {
		return nil, err
	}
	// check data is not already in ledgert
	if elementBytes, _ := stub.GetState(algoKey); elementBytes != nil {
		return nil, fmt.Errorf("algo with this hash already exists")
	}
	// submit to ledger
	algoBytes, _ := json.Marshal(algo)
	err = stub.PutState(algoKey, algoBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to add to ledger algo with key %s with error %s", algoKey, err.Error())
	}
	// create composite key
	err = createCompositeKey(stub, "algo~challenge~key", []string{"algo", algo.ChallengeKey, algoKey})
	if err != nil {
		return nil, err
	}
	return []byte(algoKey), nil
}

// createTraintuple add a Traintuple in the ledger
// ....
func createTraintuple(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	expectedArgs := getFieldNames(&inputTraintuple{})
	if nbArgs := len(expectedArgs); nbArgs != len(args) {
		return nil, fmt.Errorf("incorrect arguments, expecting %d args: %s", nbArgs, strings.Join(expectedArgs, ", "))
	}

	// find associated creator and check permissions (TODO later)
	creator, err := getTxCreator(stub)
	if err != nil {
		return nil, err
	}

	algoKey := args[0]
	startModelKey := args[1]
	if len(algoKey) != lenKey || !(len(startModelKey) == lenKey || len(startModelKey) == 0) {
		return nil, fmt.Errorf("invalid algo, or start model key")
	}
	trainDataKeys := strings.Split(strings.Replace(args[2], " ", "", -1), ",")

	// initiate traintuple
	traintuple := Traintuple{
		Status:      "todo",
		Permissions: "all",
		Creator:     creator}

	// check train data are from the same dataset and get the dataset key
	trainDatasetKey, err := checkSameDataset(stub, trainDataKeys)
	if err != nil {
		return nil, err
	}

	// check dataset exists and get trainWorker and trainDataOpener
	trainDataset := Dataset{}
	if err = getElementStruct(stub, trainDatasetKey, &trainDataset); err != nil {
		err = fmt.Errorf("could not retrieve dataset with key %s - %s", trainDatasetKey, err.Error())
		return nil, err
	}
	traintuple.TrainData = &TtData{
		Worker:     trainDataset.Owner,
		Keys:       trainDataKeys,
		OpenerHash: trainDatasetKey,
	}

	// get algo
	algo := Algo{}
	if err = getElementStruct(stub, algoKey, &algo); err != nil {
		err = fmt.Errorf("could not retrieve algo with key %s - %s", algoKey, err.Error())
		return nil, err
	}
	traintuple.Algo = &HashDressName{
		Name:           algo.Name,
		Hash:           algoKey,
		StorageAddress: algo.StorageAddress}
	// get challenge key
	challengeKey := algo.ChallengeKey
	// define for the to-be-traintuple: Challenge, StartModel, TestDataKeys, TestDataOpenerHash, TestWorker
	if startModelKey == "" { // first time algo is trained
		if err = fillTraintupleChallenge(stub, &traintuple, challengeKey); err != nil {
			return nil, err
		}
	} else { // use parent traintuple associated to model
		if err = fillTraintupleFromModel(stub, &traintuple, startModelKey, challengeKey); err != nil {
			return nil, err
		}
	}

	// create key: hash of challenge + algo + start model + train data + creator (keys)
	// certainly not be the most efficient key... but let's make it work and them try to make it better...
	tKey := sha256.Sum256([]byte(challengeKey + algoKey + startModelKey + strings.Join(trainDataKeys, ",") + creator))
	key := hex.EncodeToString(tKey[:])
	// check if traintuple key already exist
	if elementBytes, _ := stub.GetState(key); elementBytes != nil {
		return nil, fmt.Errorf("traintuple with these challenge, algo, start model, and train data already exist")
	}
	traintupleBytes, _ := json.Marshal(traintuple)
	if err = stub.PutState(key, traintupleBytes); err != nil {
		err = fmt.Errorf("could not put in ledger traintuple with startModel %s and challenge %s - %s", startModelKey, challengeKey, err.Error())
		return nil, err
	}
	// create composite keys
	if err = createCompositeKey(stub, "traintuple~algo~key", []string{"traintuple", algoKey, key}); err != nil {
		err = fmt.Errorf("issue creating composite keys - %s", err.Error())
		return nil, err
	}
	if err = createCompositeKey(stub, "traintuple~trainWorker~status~key", []string{"traintuple", traintuple.TrainData.Worker, "todo", key}); err != nil {
		err = fmt.Errorf("issue creating composite keys - %s", err.Error())
		return nil, err
	}
	if err = createCompositeKey(stub, "traintuple~testWorker~status~key", []string{"traintuple", traintuple.TestData.Worker, "todo", key}); err != nil {
		err = fmt.Errorf("issue creating composite keys - %s", err.Error())
		return nil, err
	}
	return []byte(key), nil
}

// logStartTrainTest modifies a traintuple by changing its status from todo to training
func logStartTrainTest(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	expectedArgs := [2]string{"key of the traintuple to update", "new status (training or testing)"}
	if nbArgs := len(expectedArgs); nbArgs != len(args) {
		return nil, fmt.Errorf("incorrect arguments, expecting %d args: %s", nbArgs, strings.Join(expectedArgs[:], ", "))
	}

	traintupleKey := args[0]
	if len(traintupleKey) != lenKey {
		return nil, fmt.Errorf("invalid traintuple key")
	}
	traintupleStatus := args[1]
	traintuple, err := updateStatusTraintuple(stub, traintupleKey, traintupleStatus)
	if err != nil {
		return nil, err
	}
	traintupleBytes, _ := json.Marshal(traintuple)
	return traintupleBytes, nil
}

// logSuccessTrain modifies a traintuple by changing its status from training to trained and
// reports logs and associated performances
func logSuccessTrain(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	expectedArgs := [4]string{"key of the traintuple to update", "end model hash and storage address (endModelHash, endModelStorageAddress)",
		"train perf (float)", "logs"}
	if nbArgs := len(expectedArgs); nbArgs != len(args) {
		return nil, fmt.Errorf("incorrect arguments, expecting %d args: %s", nbArgs, strings.Join(expectedArgs[:], ", "))
	}
	// get traintuple
	traintupleKey := args[0]
	if len(traintupleKey) != lenKey {
		return nil, fmt.Errorf("invalid traintuple key")
	}
	traintuple := Traintuple{}
	if err := getElementStruct(stub, traintupleKey, &traintuple); err != nil {
		return nil, err
	}
	// check validity of worker and change of status
	if err := checkUpdateTraintuple(stub, traintuple.TrainData.Worker, traintuple.Status, "trained"); err != nil {
		return nil, err
	}
	// get end model info and check validity
	endModel := strings.Split(strings.Replace(args[1], " ", "", -1), ",")
	if lenModelHash := len(endModel[0]); lenModelHash != 64 {
		return nil, fmt.Errorf("invalid len of hash of model %d", lenModelHash)
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
	// update traintuple
	traintuple.TrainData.Perf = float32(perf)
	traintuple.Status = "trained"
	traintuple.EndModel = &HashDress{
		Hash:           endModel[0],
		StorageAddress: endModel[1]}
	traintuple.Log += log
	traintupleBytes, _ := json.Marshal(traintuple)
	if err = stub.PutState(traintupleKey, traintupleBytes); err != nil {
		return nil, fmt.Errorf("failed to update traintuple status to trained with key %s", traintupleKey)
	}
	// create composite key with the end model
	indexName := "traintuple~endModel~key"
	attributes := []string{"traintuple", endModel[0], traintupleKey}
	if err = createCompositeKey(stub, indexName, attributes); err != nil {
		return nil, err
	}
	// update associated composite keys
	indexName = "traintuple~trainWorker~status~key"
	oldAttributes := []string{"traintuple", traintuple.TrainData.Worker, "training", traintupleKey}
	newAttributes := []string{"traintuple", traintuple.TrainData.Worker, "trained", traintupleKey}
	if err = updateCompositeKey(stub, indexName, oldAttributes, newAttributes); err != nil {
		return nil, err
	}
	indexName = "traintuple~testWorker~status~key"
	oldAttributes = []string{"traintuple", traintuple.TestData.Worker, "training", traintupleKey}
	newAttributes = []string{"traintuple", traintuple.TestData.Worker, "trained", traintupleKey}
	if err = updateCompositeKey(stub, indexName, oldAttributes, newAttributes); err != nil {
		return nil, err
	}
	return traintupleBytes, nil
}

// logFailTrainTest modifies a traintuple by changing its status to fail and reports associated logs
func logFailTrainTest(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
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
	// get traintuple
	traintuple := Traintuple{}
	if err := getElementStruct(stub, traintupleKey, &traintuple); err != nil {
		return nil, err
	}
	// check validity of worker and change of status
	var worker string
	oldStatus := traintuple.Status
	if stringInSlice(oldStatus, []string{"todo", "training"}) {
		worker = traintuple.TrainData.Worker
	} else if stringInSlice(oldStatus, []string{"trained", "testing"}) {
		worker = traintuple.TestData.Worker
	} else {
		return nil, fmt.Errorf("not possible to change status from %s to failed", oldStatus)
	}
	if err := checkUpdateTraintuple(stub, worker, oldStatus, "failed"); err != nil {
		return nil, err
	}
	// update traintuple
	traintuple.Status = "failed"
	traintuple.Log += log
	traintupleBytes, _ := json.Marshal(traintuple)
	if err := stub.PutState(traintupleKey, traintupleBytes); err != nil {
		return nil, fmt.Errorf("failed to update traintuple status to failed with key %s", traintupleKey)
	}
	// update associated composite keys
	indexName := "traintuple~trainWorker~status~key"
	oldAttributes := []string{"traintuple", traintuple.TrainData.Worker, oldStatus, traintupleKey}
	newAttributes := []string{"traintuple", traintuple.TrainData.Worker, "failed", traintupleKey}
	if err := updateCompositeKey(stub, indexName, oldAttributes, newAttributes); err != nil {
		return nil, err
	}
	indexName = "traintuple~testWorker~status~key"
	oldAttributes = []string{"traintuple", traintuple.TestData.Worker, oldStatus, traintupleKey}
	newAttributes = []string{"traintuple", traintuple.TestData.Worker, "failed", traintupleKey}
	if err := updateCompositeKey(stub, indexName, oldAttributes, newAttributes); err != nil {
		return nil, err
	}
	return traintupleBytes, nil
}

// logSuccessTest modifies a traintuple by changing its status to done, reports perf and logs
func logSuccessTest(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	expectedArgs := [3]string{"key of the traintuple to update", "test perf (float)", "logs"}
	if nbArgs := len(expectedArgs); nbArgs != len(args) {
		return nil, fmt.Errorf("incorrect arguments, expecting %d args: %s", nbArgs, strings.Join(expectedArgs[:], ", "))
	}
	// get traintuple
	traintupleKey := args[0]
	if len(traintupleKey) != lenKey {
		return nil, fmt.Errorf("invalid traintuple key")
	}
	traintuple := Traintuple{}
	if err := getElementStruct(stub, traintupleKey, &traintuple); err != nil {
		return nil, err
	}
	// check validity of worker (transaction requester is test worker) and change of status (currently testing)
	if err := checkUpdateTraintuple(stub, traintuple.TestData.Worker, traintuple.Status, "done"); err != nil {
		return nil, err
	}
	// get test perf, check validity
	perf, err := strconv.ParseFloat(args[1], 32)
	if err != nil {
		return nil, err
	}
	// get logs and check validity
	log := args[2]
	if err = checkLog(log); err != nil {
		return nil, err
	}
	traintuple.Status = "done"
	traintuple.TestData.Perf = float32(perf)
	traintuple.Log += log
	traintupleBytes, _ := json.Marshal(traintuple)
	if err = stub.PutState(traintupleKey, traintupleBytes); err != nil {
		return nil, fmt.Errorf("failed to update traintuple status to trained with key %s", traintupleKey)
	}
	// update associated composite keys
	indexName := "traintuple~trainWorker~status~key"
	oldAttributes := []string{"traintuple", traintuple.TrainData.Worker, "testing", traintupleKey}
	newAttributes := []string{"traintuple", traintuple.TrainData.Worker, "done", traintupleKey}
	if err = updateCompositeKey(stub, indexName, oldAttributes, newAttributes); err != nil {
		return nil, err
	}
	indexName = "traintuple~testWorker~status~key"
	oldAttributes = []string{"traintuple", traintuple.TestData.Worker, "testing", traintupleKey}
	newAttributes = []string{"traintuple", traintuple.TestData.Worker, "done", traintupleKey}
	if err = updateCompositeKey(stub, indexName, oldAttributes, newAttributes); err != nil {
		return nil, err
	}
	return traintupleBytes, nil
}

// query returns an element of the ledger given its key
// For now, ok for everything. Later returns if the requester has permission to see it
func query(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	expectedArgs := [1]string{"element key"}
	if nbArgs := len(expectedArgs); nbArgs != len(args) {
		return nil, fmt.Errorf("incorrect arguments, expecting %d args: %s", nbArgs, strings.Join(expectedArgs[:], ", "))
	}

	key := args[0]
	if len(key) != 64 {
		return nil, fmt.Errorf("invalid key")
	}
	return getElementBytes(stub, key)
}

// queryFilter returns all elements of the ledger matching some filters
// For now, ok for everything. Later returns if the requester has permission to see it
func queryFilter(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	expectedArgs := [2]string{"indexName", "attributes"}
	if nbArgs := len(expectedArgs); nbArgs != len(args) {
		return nil, fmt.Errorf("incorrect arguments, expecting %d args: %s", nbArgs, strings.Join(expectedArgs[:], ", "))
	}
	// check validity of inputs
	indexName := args[0]
	validIndexNames := []string{"traintuple~trainWorker~status", "traintuple~testWorker~status"}
	if !stringInSlice(indexName, validIndexNames) {
		return nil, fmt.Errorf("invalid indexName filter query: %s", indexName)
	}
	indexName = indexName + "~key"
	attributes := strings.Split(strings.Replace(args[1], " ", "", -1), ",")
	attributes = append([]string{strings.Split(indexName, "~")[0]}, attributes...)

	filteredKeys, err := getKeysFromComposite(stub, indexName, attributes)
	if err != nil {
		return nil, fmt.Errorf("issue getting keys from composite key %s - %s", indexName, err.Error())
	}
	//TODO get elements with filtererd keys
	var elements []map[string]interface{}
	for _, key := range filteredKeys {
		var element map[string]interface{}
		if err := getElementStruct(stub, key, &element); err != nil {
			return nil, fmt.Errorf("issue getting element with key %s - %s", key, err.Error())
		}
		element["key"] = key
		elements = append(elements, element)
	}
	payload, err := json.Marshal(elements)
	if err != nil {
		return nil, err
	}
	return payload, nil
}

// queryAll returns all elements of the ledger given its type
// For now, ok for everything. Later returns if the requester has permission to see it
func queryAll(stub shim.ChaincodeStubInterface, args []string, elementType string) ([]byte, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("incorrect number of arguments, expecting nothing")
	}
	var indexName string
	switch elementType {
	case "challenge":
		indexName = "challenge~owner~key"
	case "dataset":
		indexName = "dataset~owner~key"
	case "algo":
		indexName = "algo~challenge~key"
	case "traintuple":
		indexName = "traintuple~algo~key"
	default:
		return nil, fmt.Errorf("no element type %s", elementType)
	}
	elementsKeys, err := getKeysFromComposite(stub, indexName, []string{elementType})
	var elements []map[string]interface{}
	for _, key := range elementsKeys {
		var element map[string]interface{}
		err = getElementStruct(stub, key, &element)
		if err != nil {
			return nil, err
		}
		element["key"] = key
		elements = append(elements, element)

	}
	payload, err := json.Marshal(elements)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

// queryModel returns model's permissions and the associated algo key
func queryModel(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	expectedArgs := [1]string{"model hash"}
	if nbArgs := len(expectedArgs); nbArgs != len(args) {
		return nil, fmt.Errorf("incorrect arguments, expecting %d args: %s", nbArgs, strings.Join(expectedArgs[:], ", "))
	}
	modelHash := args[0]
	return getModel(stub, modelHash)
}

// queryModelTraintuples returns info about the algo related to a model, and all traintuple related to this algo
func queryModelTraintuples(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	expectedArgs := [1]string{"model hash"}
	if nbArgs := len(expectedArgs); nbArgs != len(args) {
		return nil, fmt.Errorf("incorrect arguments, expecting %d args: %s", nbArgs, strings.Join(expectedArgs[:], ", "))
	}
	modelHash := args[0]

	// get associated traintuple
	traintupleBytes, err := getModel(stub, modelHash)
	if err != nil {
		return nil, err
	}
	traintuple := Traintuple{}
	if err = bytesToStruct(traintupleBytes, &traintuple); err != nil {
		return nil, err
	}
	// get associated algo
	algoKey := traintuple.Algo.Hash
	mPayload := make(map[string]interface{})
	var algo map[string]interface{}
	if err := getElementStruct(stub, algoKey, &algo); err != nil {
		return nil, err
	}
	algo["key"] = algoKey
	mPayload["algo"] = algo
	// get traintuples related to algo, whose permissions match the requester
	traintupleKeys, err := getKeysFromComposite(stub, "traintuple~algo~key", []string{"traintuple", algoKey})
	if err != nil {
		return nil, err
	}
	// get all traintuples and serialize them
	var traintuples []map[string]interface{}
	for _, traintupleKey := range traintupleKeys {
		var traintuple map[string]interface{} // Traintuple{}
		if err = getElementStruct(stub, traintupleKey, &traintuple); err != nil {
			return nil, err
		}
		traintuple["key"] = traintupleKey
		traintuples = append(traintuples, traintuple)
	}
	mPayload["traintuples"] = traintuples
	// Marshal payload
	payload, err := json.Marshal(mPayload)
	if err != nil {
		return nil, err
	}
	return payload, nil
}

// queryDatasetData returns info about a dataset and all related data
func queryDatasetData(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	expectedArgs := [1]string{"dataset key"}
	if nbArgs := len(expectedArgs); nbArgs != len(args) {
		return nil, fmt.Errorf("incorrect arguments, expecting %d args: %s", nbArgs, strings.Join(expectedArgs[:], ", "))
	}
	datasetKey := args[0]

	// get dataset info
	var mPayload map[string]interface{}
	if err := getElementStruct(stub, datasetKey, &mPayload); err != nil {
		return nil, err
	}
	mPayload["key"] = datasetKey
	// get related train data
	trainDataKeys, err := getDatasetData(stub, datasetKey, true)
	if err != nil {
		return nil, err
	}
	mPayload["trainDataKeys"] = trainDataKeys
	// Marshal payload
	payload, err := json.Marshal(mPayload)
	if err != nil {
		return nil, err
	}
	return payload, nil
}
