package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// getFieldName returns a slice containing field names of a struc
func getFieldNames(v interface{}) (fieldNames []string) {
	e := reflect.ValueOf(v).Elem()
	eType := e.Type()
	for i := 0; i < e.NumField(); i++ {
		varName := eType.Field(i).Name
		fieldNames = append(fieldNames, varName)
	}
	return
}

// inputStructToBytes converts fields of a struct (with string fields only, such as input struct defined in ledger.go) to a [][]byte
func inputStructToBytes(v interface{}) (sb [][]byte, err error) {

	e := reflect.ValueOf(v).Elem()
	for i := 0; i < e.NumField(); i++ {
		v := e.Field(i)
		if v.Type().Name() != "string" {
			err = fmt.Errorf("struc should contain only string values")
			return
		}
		varValue := v.String()
		sb = append(sb, []byte(varValue))
	}
	return

}

// stringToInputStruct fills fields of a input struct (such as defined in ledger.go) with elements stored in a slice of string
func stringToInputStruct(args []string, v interface{}) {
	fieldNames := getFieldNames(v)
	e := reflect.ValueOf(v).Elem()
	for i, fn := range fieldNames {
		f := e.FieldByName(fn)
		f.SetString(args[i])
	}
}

// getTxCreator returns the sha256 of the creator of the transaction
func getTxCreator(stub shim.ChaincodeStubInterface) (string, error) {
	// get the agent submitting the transaction
	bCreator, err := stub.GetCreator()
	if err != nil {
		return "", err
	}
	tt := sha256.Sum256(bCreator)
	return hex.EncodeToString(tt[:]), nil
}

// bytesToStruct converts bytes to one a the struct corresponding to elements stored in the ledger
func bytesToStruct(elementBytes []byte, element interface{}) error {
	return json.Unmarshal(elementBytes, &element)
}

// getElementBytes checks if an element is stored in the ledger given its key, and returns associated bytes
func getElementBytes(stub shim.ChaincodeStubInterface, elementKey string) ([]byte, error) {
	elementBytes, err := stub.GetState(elementKey)
	if err != nil {
		return nil, err
	} else if elementBytes == nil {
		return nil, fmt.Errorf("no element with key %s", elementKey)
	}
	return elementBytes, nil
}

// getElementStruct fills an element struct given its key
func getElementStruct(stub shim.ChaincodeStubInterface, elementKey string, element interface{}) error {
	elementBytes, err := getElementBytes(stub, elementKey)
	if err != nil {
		return err
	}
	return bytesToStruct(elementBytes, element)
}

// checkSameDataset checks if data in a slices exist and are from the same dataset. If yes, returns the dataset key
func checkSameDataset(stub shim.ChaincodeStubInterface, dataKeys []string) (string, error) {
	var datasetKey string
	for i, dataKey := range dataKeys {
		data := Data{}
		if err := getElementStruct(stub, dataKey, &data); err != nil {
			return "", fmt.Errorf("issue retrieving %s %s", dataKey, err.Error())
		}
		if i == 0 {
			datasetKey = data.DatasetKey
			continue
		}
		if data.DatasetKey != datasetKey {
			return "", fmt.Errorf("data do not belong to the same dataset")
		}
	}
	return datasetKey, nil
}

// addChallengeDataset adds a challenge key to the list of associated challenges of a dataset
func addChallengeDataset(stub shim.ChaincodeStubInterface, datasetKey string, challengeKey string) error {
	dataset := Dataset{}
	if err := getElementStruct(stub, datasetKey, &dataset); err != nil {
		return nil
	}
	dataset.ChallengeKeys = append(dataset.ChallengeKeys, challengeKey)
	datasetBytes, _ := json.Marshal(dataset)
	return stub.PutState(datasetKey, datasetBytes)
}

// createCompositeKey creates a composite key given an indexName and attributes
// (combination of attributes to form a key)
func createCompositeKey(stub shim.ChaincodeStubInterface, indexName string, attributes []string) error {
	compositeKey, err := stub.CreateCompositeKey(indexName, attributes)
	if err != nil {
		return err
	}
	value := []byte{0x00}
	if err = stub.PutState(compositeKey, value); err != nil {
		return fmt.Errorf("failed to add composite key with index %s to the ledger", indexName)
	}
	return nil
}

// getKeysFromComposite returns element keys associated with a composite key specified by its indexName and attributes
func getKeysFromComposite(stub shim.ChaincodeStubInterface, indexName string,
	attributes []string) ([]string, error) {
	var elementKeys []string
	compositeIterator, err := stub.GetStateByPartialCompositeKey(indexName, attributes)
	if err != nil {
		return elementKeys, err
	}
	defer compositeIterator.Close()
	for i := 0; compositeIterator.HasNext(); i++ {
		compositeKey, err := compositeIterator.Next()
		if err != nil {
			return elementKeys, err
		}
		_, compositeKeyParts, err := stub.SplitCompositeKey(compositeKey.Key)
		if err != nil {
			return elementKeys, err
		}
		elementKeys = append(elementKeys, compositeKeyParts[len(compositeKeyParts)-1])
	}
	return elementKeys, nil
}

// getDatasetData returns all data keys associated to a dataset
func getDatasetData(stub shim.ChaincodeStubInterface, datasetKey string, trainOnly bool) ([]string, error) {
	var indexName string
	var attributes []string
	if trainOnly {
		indexName = "data~dataset~testOnly~key"
		attributes = []string{"data", datasetKey, "false"}
	} else {
		indexName = "data~dataset~key"
		attributes = []string{"data", datasetKey}
	}
	dataKeys, err := getKeysFromComposite(stub, indexName, attributes)
	if err != nil {
		return nil, err
	}
	return dataKeys, nil
}

// strToPerf checks if a string such as dataKey1:perf1, dataKey2:perf2, dataKey3:perf3, ...
// has same data keys as those in dataKeys and returns a []float32 with perf in same order as dataKeys
func strToPerf(stub shim.ChaincodeStubInterface, strDataPerf string, dataKeys []string) ([]float32, error) {
	var perf []float32
	var err error
	m := make(map[string]float32)
	// split each dataKey:perf pair
	sliceDataPerf := strings.Split(strings.Replace(strDataPerf, " ", "", -1), ",")
	// check data is part of the traintuple and fill mapDataPerf
	for _, dataPerf := range sliceDataPerf {
		sdp := strings.Split(dataPerf, ":")
		dataKey := sdp[0]
		// convert perf to float
		perfValue, err := strconv.ParseFloat(sdp[1], 32)
		if err != nil {
			return perf, err
		}
		m[dataKey] = float32(perfValue)
	}
	// check if data keys in maps correspond to dataKeys
	for _, dataKey := range dataKeys {
		value, ok := m[dataKey]
		if !ok {
			err = fmt.Errorf("dataKey %s referenced in slice is not in map", dataKey)
			return perf, err
		}
		perf = append(perf, value)
	}
	return perf, err
}

// checkLog checks the validity of logs
func checkLog(log string) (err error) {
	maxLength := 200
	if length := len(log); length > maxLength {
		err = fmt.Errorf("too long log, is %d and should be %d ", length, maxLength)
	}
	return
}

// stringInSlice checks if a string is in a slice of string
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// updateStatusTraintuple changes the status of a traintuple
func updateStatusTraintuple(stub shim.ChaincodeStubInterface, traintupleKey string, status string) (Traintuple, error) {

	traintuple := Traintuple{}

	// get traintuple
	if err := getElementStruct(stub, traintupleKey, &traintuple); err != nil {
		return Traintuple{}, err
	}
	// check the validity of traintuple update: consistent status and worker
	var worker string
	if status == "training" {
		worker = traintuple.TrainData.Worker
	} else if status == "testing" {
		worker = traintuple.TestData.Worker
	} else {
		return Traintuple{}, fmt.Errorf("status %s is not implemented, expecting training or testing", status)
	}
	if err := checkUpdateTraintuple(stub, worker, traintuple.Status, status); err != nil {
		return Traintuple{}, err
	}
	// update traintuple
	traintuple.Status = status
	traintupleBytes, _ := json.Marshal(traintuple)
	if err := stub.PutState(traintupleKey, traintupleBytes); err != nil {
		return traintuple, fmt.Errorf("failed to update traintuple status to %s with key %s", status, traintupleKey)
	}
	return traintuple, nil
}

// check validity of traintuple update: consistent status and agent submitting the transaction
func checkUpdateTraintuple(stub shim.ChaincodeStubInterface, worker string, oldStatus string, newStatus string) error {
	txCreator, err := getTxCreator(stub)
	if err != nil {
		return err
	}
	if txCreator != worker {
		return fmt.Errorf("%s is not allowed to update traintuple", txCreator)
	}
	statusPossibilities := map[string]string{
		"todo":     "training",
		"training": "trained",
		"trained":  "testing",
		"testing":  "done"}
	if statusPossibilities[oldStatus] != newStatus {
		return fmt.Errorf("cannot change status from %s to %s", oldStatus, newStatus)
	}
	return nil
}

// updateCompositeKey modifies composite keys
func updateCompositeKey(stub shim.ChaincodeStubInterface, indexName string, oldAttributes []string, newAttributes []string) error {
	oldCompositeKey, err := stub.CreateCompositeKey(indexName, oldAttributes)
	if err != nil {
		return err
	}
	err = stub.DelState(oldCompositeKey)
	if err != nil {
		return err
	}
	newCompositeKey, err := stub.CreateCompositeKey(indexName, newAttributes)
	if err != nil {
		return err
	}
	value := []byte{0x00}
	return stub.PutState(newCompositeKey, value)
}

// getModel return the traintuple (as bytes) in which the model is the endModel
func getModel(stub shim.ChaincodeStubInterface, modelHash string) ([]byte, error) {

	traintupleKeys, err := getKeysFromComposite(stub, "traintuple~endModel~key",
		[]string{"traintuple", modelHash})
	if len(traintupleKeys) > 1 {
		return nil, fmt.Errorf("more than one traintuple with endModel hash %s", modelHash)
	}
	traintupleKey := traintupleKeys[0]
	traintupleBytes, err := getElementBytes(stub, traintupleKey)
	if err != nil {
		return nil, fmt.Errorf("error getting associated traintuple %s", traintupleKey)
	}
	return traintupleBytes, err
}

// fillTraintupleChallengeFromModel fills the following fields of the pointed traintuple, given a startModel key:
// Challenge, StartModel, TestDataKeys, TestDataOpenerHash, TestWorker, Rank
func fillTraintupleFromModel(stub shim.ChaincodeStubInterface, traintuple *Traintuple, startModelKey string) error {
	// get parent traintuple
	parentTraintupleKeys, err := getKeysFromComposite(stub, "traintuple~endModel~key",
		[]string{"traintuple", startModelKey})
	if err != nil {
		return err
	}
	if len(parentTraintupleKeys) != 1 {
		return fmt.Errorf("several models associated with start model hash")
	}
	parentTraintupleKey := parentTraintupleKeys[0]
	// model derives from a previous Traintuple
	parentTraintuple := Traintuple{}
	err = getElementStruct(stub, parentTraintupleKey, &parentTraintuple)
	if err != nil {
		return err
	}
	// fill traintuple
	traintuple.Challenge = parentTraintuple.Challenge
	traintuple.StartModel = parentTraintuple.EndModel
	traintuple.TestData = parentTraintuple.TestData
	traintuple.Rank = parentTraintuple.Rank + 1
	return nil
}

// fillTraintupleChallengeFromAlgo fills the following fields of the pointed traintuple, given an algo key:
// Challenge, StartModel, TestDataKeys, TestDataOpenerHash, TestWorker, Rank
func fillTraintupleFromAlgo(stub shim.ChaincodeStubInterface, traintuple *Traintuple, algoKey string, challengeKey string) error {
	// startModel corresponds to the algo itself. Check algo field corresponds to algoKey
	if traintupleAlgoKey := "algo_" + traintuple.Algo.Hash; traintupleAlgoKey != algoKey {
		return fmt.Errorf("input algoKey %s does not correspond to algoKey of traintuple %s",
			algoKey, traintupleAlgoKey)
	}
	traintuple.StartModel = traintuple.Algo
	// get challenge to derive metrics info and test data keys
	retrievedChallenge := Challenge{}
	if err := getElementStruct(stub, challengeKey, &retrievedChallenge); err != nil {
		return err
	}
	metrics := HashDress{
		Hash:           retrievedChallenge.Metrics.Hash,
		StorageAddress: retrievedChallenge.Metrics.StorageAddress,
	}
	traintuple.Challenge = &TtChallenge{
		Hash:    strings.Split(challengeKey, "_")[1],
		Metrics: &metrics,
	}

	// get test worker and test data openerHas from associated dataset
	testData := Data{}
	if err := getElementStruct(stub, retrievedChallenge.TestDataKeys[0], &testData); err != nil {
		return err
	}
	testDatasetKey := testData.DatasetKey
	testDataset := Dataset{}
	if err := getElementStruct(stub, testDatasetKey, &testDataset); err != nil {
		return err
	}
	traintuple.TestData = &TtData{
		Worker:     testDataset.Owner,
		Keys:       retrievedChallenge.TestDataKeys,
		OpenerHash: strings.Split(testDatasetKey, "_")[1],
	}
	// first time algo is trained
	traintuple.Rank = 0
	return nil
}
