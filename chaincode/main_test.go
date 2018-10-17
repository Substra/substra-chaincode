package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	peer "github.com/hyperledger/fabric/protos/peer"
)

const challengeDescriptionHash = "5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379"
const challengeDescriptionStorageAddress = "https://toto/challenge/222/description"
const challengeMetricsHash = "4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379"
const challengeMetricsStorageAddress = "https://toto/challenge/222/metrics"
const datasetOpenerHash = "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
const trainDataHash1 = "aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
const trainDataHash2 = "aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
const testDataHash1 = "bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
const testDataHash2 = "bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
const algoHash = "fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
const algoStorageAddress = "https://toto/algo/222/algo"
const algoName = "hog + svm"
const modelHash = "eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed"
const modelAddress = "https://substrabac/model/toto"
const worker = "bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0"

func TestInit(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := shim.NewMockStub("substra", scc)

	// resp := mockStub.MockInit("42", [][]byte{[]byte("init")})
	resp := mockStub.MockInit("42", [][]byte{[]byte("init")})
	status := resp.Status
	if status != 200 {
		t.Errorf("init failed with status %d and message %s", status, resp.Message)
	}
}

func printArgs(args [][]byte, command string) {
	s := "```\npeer chaincode " + command + " -n mycc -c '{\"Args\":["
	for i, arg := range args {
		s += "\"" + string(arg) + "\""
		if i+1 < len(args) {
			s += ","
		}
	}
	s += "]}' -C myc\n```"
	fmt.Println(s)
}
func printArgsNames(fnName string, argsNames []string) {
	s := "Smart contract: `" + fnName + "`\nInputs: `" + strings.Join(argsNames, "`, `") + "`"
	fmt.Println(s)
}

func (dataset *inputDataset) createSample() [][]byte {
	if dataset.Name == "" {
		dataset.Name = "liver slide"
	}
	if dataset.OpenerHash == "" {
		dataset.OpenerHash = datasetOpenerHash
	}
	if dataset.OpenerStorageAddress == "" {
		dataset.OpenerStorageAddress = "https://toto/dataset/42234/opener"
	}
	if dataset.Type == "" {
		dataset.Type = "images"
	}
	if dataset.DescriptionHash == "" {
		dataset.DescriptionHash = "8d4bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eee"
	}
	if dataset.DescriptionStorageAddress == "" {
		dataset.DescriptionStorageAddress = "https://toto/dataset/42234/description"
	}
	dataset.Permissions = "all"
	args, _ := inputStructToBytes(dataset)
	args = append([][]byte{[]byte("registerDataset")}, args...)
	return args
}

func (data *inputData) createSample() [][]byte {
	if data.Hashes == "" {
		data.Hashes = trainDataHash1 + ", " + trainDataHash2
	}
	if data.DatasetKey == "" {
		data.DatasetKey = datasetOpenerHash
	}
	if data.Size == "" {
		data.Size = "100"
	}
	if data.TestOnly == "" {
		data.TestOnly = "false"
	}
	args, _ := inputStructToBytes(data)
	args = append([][]byte{[]byte("registerData")}, args...)
	return args
}

func (challenge *inputChallenge) createSample() [][]byte {
	if challenge.Name == "" {
		challenge.Name = "MSI classification"
	}
	if challenge.DescriptionHash == "" {
		challenge.DescriptionHash = challengeDescriptionHash
	}
	if challenge.DescriptionStorageAddress == "" {
		challenge.DescriptionStorageAddress = "https://toto/challenge/222/description"
	}
	if challenge.MetricsName == "" {
		challenge.MetricsName = "accuracy"
	}
	if challenge.MetricsHash == "" {
		challenge.MetricsHash = challengeMetricsHash
	}
	if challenge.MetricsStorageAddress == "" {
		challenge.MetricsStorageAddress = challengeMetricsStorageAddress
	}
	if challenge.TestDataKeys == "" {
		challenge.TestDataKeys = testDataHash1 + ", " + testDataHash2
	}
	challenge.Permissions = "all"
	args, _ := inputStructToBytes(challenge)
	args = append([][]byte{[]byte("registerChallenge")}, args...)
	return args
}

func (algo *inputAlgo) createSample() [][]byte {
	if algo.Name == "" {
		algo.Name = algoName
	}
	if algo.Hash == "" {
		algo.Hash = algoHash
	}
	if algo.StorageAddress == "" {
		algo.StorageAddress = algoStorageAddress
	}
	if algo.DescriptionHash == "" {
		algo.DescriptionHash = "e2dbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dca"
	}
	if algo.DescriptionStorageAddress == "" {
		algo.DescriptionStorageAddress = "https://toto/algo/222/description"
	}
	if algo.ChallengeKey == "" {
		algo.ChallengeKey = challengeDescriptionHash
	}
	algo.Permissions = "all"
	args, _ := inputStructToBytes(algo)
	args = append([][]byte{[]byte("registerAlgo")}, args...)
	return args
}

func (traintuple *inputTraintuple) createSample() [][]byte {
	if traintuple.AlgoKey == "" {
		traintuple.AlgoKey = algoHash
	}
	if traintuple.StartModelKey == "" {
		traintuple.StartModelKey = ""
	}
	if traintuple.TrainDataKeys == "" {
		traintuple.TrainDataKeys = trainDataHash1 + ", " + trainDataHash2
	}
	args, _ := inputStructToBytes(traintuple)
	args = append([][]byte{[]byte("createTraintuple")}, args...)
	return args
}

func registerItem(mockStub shim.MockStub, itemType string) (error, peer.Response, interface{}) {
	// 1. add dataset
	inpDataset := inputDataset{}
	args := inpDataset.createSample()
	resp := mockStub.MockInvoke("42", args)
	if resp.Status != 200 {
		return fmt.Errorf("when adding dataset with status %d and message %s", resp.Status, resp.Message), resp, inpDataset
	} else if itemType == "dataset" {
		return nil, resp, inpDataset
	}
	// 2. add test data
	inpData := inputData{
		Hashes:   testDataHash1 + ", " + testDataHash2,
		TestOnly: "true",
	}
	args = inpData.createSample()
	resp = mockStub.MockInvoke("42", args)
	if resp.Status != 200 {
		return fmt.Errorf("when adding test data with status %d and message %s", resp.Status, resp.Message), resp, inpData
	} else if itemType == "testData" {
		return nil, resp, inpData
	}
	// 3. add challenge
	inpChallenge := inputChallenge{}
	args = inpChallenge.createSample()
	resp = mockStub.MockInvoke("42", args)
	if resp.Status != 200 {
		return fmt.Errorf("when adding challenge with status %d and message %s", resp.Status, resp.Message), resp, inpChallenge
	} else if itemType == "challenge" {
		return nil, resp, inpChallenge
	}
	// 4. Add train data
	inpData = inputData{}
	args = inpData.createSample()
	resp = mockStub.MockInvoke("42", args)
	if resp.Status != 200 {
		return fmt.Errorf("when adding train data with status %d and message %s", resp.Status, resp.Message), resp, inpData
	} else if itemType == "trainData" {
		return nil, resp, inpData
	}
	// 5. Add algo
	inpAlgo := inputAlgo{}
	args = inpAlgo.createSample()
	resp = mockStub.MockInvoke("42", args)
	if resp.Status != 200 {
		return fmt.Errorf("when adding algo with status %d and message %s", resp.Status, resp.Message), resp, inpAlgo
	} else if itemType == "algo" {
		return nil, resp, inpAlgo
	}
	// 6. Add traintuple
	inpTraintuple := inputTraintuple{}
	args = inpTraintuple.createSample()
	resp = mockStub.MockInvoke("42", args)
	if resp.Status != 200 {
		return fmt.Errorf("when adding traintuple with status %d and message %s", resp.Status, resp.Message), resp, inpAlgo
	}
	return nil, resp, inpTraintuple
}

func TestDataset(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := shim.NewMockStub("substra", scc)

	// Add dataset with invalid field
	inpDataset := inputDataset{
		OpenerHash: "aaa",
	}
	args := inpDataset.createSample()
	resp := mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 500 {
		t.Errorf("when adding dataset with invalid opener hash, status %d and message %s", status, resp.Message)
	}
	// Properly add dataset
	err, resp, tt := registerItem(*mockStub, "dataset")
	if err != nil {
		t.Errorf(err.Error())
	}
	inpDataset = tt.(inputDataset)
	datasetKey := string(resp.Payload)
	// check returned dataset key corresponds to opener hash
	if datasetKey != datasetOpenerHash {
		t.Errorf("when adding dataset: dataset key does not correspond to dataset opener hash: %s - %s", datasetKey, datasetOpenerHash)
	}
	// Add dataset which already exist
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 500 {
		t.Errorf("when adding dataset which already exists, status %d and message %s", status, resp.Message)
	}
	// Query dataset and check fields match expectations
	args = [][]byte{[]byte("query"), []byte(datasetKey)}
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("when querying the dataset, status %d and message %s", status, resp.Message)
	}
	dataset := make(map[string]interface{})
	if err := bytesToStruct(resp.Payload, &dataset); err != nil {
		t.Errorf("when unmarshalling queried dataset with error %s", err)
	}
	if dataset["name"] != inpDataset.Name {
		t.Errorf("ledger dataset name does not correspond to what was input: %s - %s", dataset["name"], inpDataset.Name)
	}
	if dataset["openerStorageAddress"] != inpDataset.OpenerStorageAddress {
		t.Errorf("ledger dataset opener storage address does not correspond to what was input")
	}
	if dataset["size"] != 0. && dataset["nbData"] != 0. {
		t.Errorf("ledger dataset size is not 0, whereas no data was added")
	}
	if dataset["type"] != inpDataset.Type {
		t.Errorf("ledger dataset type does not correspond to what was input")
	}
	if dataset["description"].(map[string]interface{})["hash"] != inpDataset.DescriptionHash {
		t.Errorf("ledger dataset description hash does not correspond to what was input")
	}
	if dataset["description"].(map[string]interface{})["storageAddress"] != inpDataset.DescriptionStorageAddress {
		t.Errorf("ledger dataset description storage address does not correspond to what was input")
	}
	if reflect.ValueOf(dataset["challengeKeys"]).Len() != 0 {
		t.Errorf("ledger dataset challenge keys does not correspond to what was input")
	}
	if dataset["permissions"] != inpDataset.Permissions {
		t.Errorf("ledger dataset challenge keys does not correspond to what was input")
	}
	// Query all datasets and check fields match expectations
	args = [][]byte{[]byte("queryDatasets")}
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("when querying datasets - status %d and message %s", status, resp.Message)
	}
	var sPayload []map[string]interface{}
	if err := json.Unmarshal(resp.Payload, &sPayload); err != nil {
		t.Errorf("when unmarshalling queried datasets")
	}
	payload := sPayload[0]
	delete(payload, "key")
	if !reflect.DeepEqual(payload, dataset) {
		t.Errorf("when querying datasets, dataset does not correspond to the input dataset")
	}

	args = [][]byte{[]byte("queryDatasetData"), []byte(inpDataset.OpenerHash)}
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("when querying dataset data, status %d and message %s", status, resp.Message)
	}
	if !strings.Contains(string(resp.Payload), "\"trainDataKeys\":[]") {
		t.Errorf("when querying dataset data, trainDataKeys should be []")
	}
}

func TestData(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := shim.NewMockStub("substra", scc)

	// Add data with invalid field
	inpData := inputData{
		Hashes: "aaa",
	}
	args := inpData.createSample()
	resp := mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 500 {
		t.Errorf("when adding data with invalid hash, status %d and message %s", status, resp.Message)
	}

	// Add data with unexiting dataset
	inpData = inputData{}
	args = inpData.createSample()
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 500 {
		t.Errorf("when adding data with unexisting dataset, status %d and message %s", status, resp.Message)
	}
	// TODO Would be nice to check failure when adding data to a dataset owned by a different people

	// Properly add data
	// 1. add associated dataset
	inpDataset := inputDataset{}
	args = inpDataset.createSample()
	mockStub.MockInvoke("42", args)
	// 2. add data
	inpData = inputData{}
	args = inpData.createSample()
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("when adding data, status %d and message %s", status, resp.Message)
	}
	// check payload correspond to input data keys
	dataKeys := string(resp.Payload)
	if dataKeys != inpData.Hashes {
		t.Errorf("when adding data: data keys does not correspond to data hashes: %s - %s", dataKeys, inpData.Hashes)
	}

	// Add data which already exist
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 500 {
		t.Errorf("when adding data which already exist, status %d and message %s", status, resp.Message)
	}

	// Query data and check it corresponds to what was input
	args = [][]byte{[]byte("queryDatasetData"), []byte(inpDataset.OpenerHash)}
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("when querying dataset data with status %d and message %s", status, resp.Message)
	}
	payload := make(map[string]interface{})
	json.Unmarshal(resp.Payload, &payload)
	if _, ok := payload["key"]; !ok {
		t.Errorf("when querying dataset data, payload should contain the dataset key")
	}
	v, ok := payload["trainDataKeys"]
	if !ok {
		t.Errorf("when querying dataset data, payload should contain the train data keys")
	}
	if reflect.DeepEqual(v, strings.Split(strings.Replace(inpData.Hashes, " ", "", -1), ",")) {
		t.Errorf("when querying dataset data, unexpected train keys")
	}
}

func TestChallenge(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := shim.NewMockStub("substra", scc)

	// Add challenge with invalid field
	inpChallenge := inputChallenge{
		DescriptionHash: "aaa",
	}
	args := inpChallenge.createSample()
	resp := mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 500 {
		t.Errorf("when adding challenge with invalid hash, status %d and message %s", status, resp.Message)
	}

	// Add challenge with unexisting test data
	inpChallenge = inputChallenge{}
	args = inpChallenge.createSample()
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 500 {
		t.Errorf("when adding challenge with unexisting test data, status %d and message %s", status, resp.Message)
	}

	// Properly add challenge
	err, resp, tt := registerItem(*mockStub, "challenge")
	if err != nil {
		t.Errorf(err.Error())
	}
	inpChallenge = tt.(inputChallenge)
	challengeKey := string(resp.Payload)
	if challengeKey != inpChallenge.DescriptionHash {
		t.Errorf("when adding challenge: unexpected returned challenge key - %s / %s", challengeKey, inpChallenge.DescriptionHash)
	}

	// Query challenge from key and check the consistency of returned arguments
	args = [][]byte{[]byte("query"), []byte(challengeKey)}
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("when querying a dataset with status %d and message %s", status, resp.Message)
	}
	challenge := make(map[string]interface{})
	if err := bytesToStruct(resp.Payload, &challenge); err != nil {
		t.Errorf("when unmarshalling queried challenge with error %s", err)
	}
	if challenge["name"] != inpChallenge.Name {
		t.Errorf("ledger challenge name does not correspond to what was input: %s - %s", challenge["name"], inpChallenge.Name)
	}
	if challenge["descriptionStorageAddress"] != inpChallenge.DescriptionStorageAddress {
		t.Errorf("ledger challenge description storage address does not correspond to what was input: %s - %s", challenge["descriptionStorageAddress"], inpChallenge.DescriptionStorageAddress)
	}
	if challenge["permissions"] != inpChallenge.Permissions {
		t.Errorf("ledger challenge permissions does not correspond to what was input: %s - %s", challenge["permissions"], inpChallenge.Permissions)
	}
	if challenge["metrics"].(map[string]interface{})["hash"] != inpChallenge.MetricsHash {
		t.Errorf("ledger challenge metrics hash does not correspond to what was input")
	}
	if challenge["metrics"].(map[string]interface{})["name"] != inpChallenge.MetricsName {
		t.Errorf("ledger challenge metrics name does not correspond to what was input")
	}
	if challenge["metrics"].(map[string]interface{})["storageAddress"] != inpChallenge.MetricsStorageAddress {
		t.Errorf("ledger challenge metrics address does not correspond to what was input")
	}
	if reflect.DeepEqual(challenge["testDataKeys"], []string{testDataHash1, testDataHash2}) {
		t.Errorf("ledger challenge test data does not correspond to what was input")
	}

	// Query all challenges and check consistency
	args = [][]byte{[]byte("queryChallenges")}
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("when querying challenges - status %d and message %s", status, resp.Message)
	}
	var sPayload []map[string]interface{}
	if err := json.Unmarshal(resp.Payload, &sPayload); err != nil {
		t.Errorf("when unmarshalling queried challenges")
	}
	payload := sPayload[0]
	delete(payload, "key")
	if !reflect.DeepEqual(payload, challenge) {
		t.Errorf("when querying challenges, dataset does not correspond to the input challenge")
	}
}

func TestAlgo(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := shim.NewMockStub("substra", scc)

	// Add algo with invalid field
	inpAlgo := inputAlgo{
		DescriptionHash: "aaa",
	}
	args := inpAlgo.createSample()
	resp := mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 500 {
		t.Errorf("when adding algo with invalid hash, status %d and message %s", status, resp.Message)
	}

	// Add algo with unexisting challenge
	inpAlgo = inputAlgo{}
	args = inpAlgo.createSample()
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 500 {
		t.Errorf("when adding algo with unexisting challenge, status %d and message %s", status, resp.Message)
	}

	// Properly add algo
	err, resp, tt := registerItem(*mockStub, "algo")
	if err != nil {
		t.Errorf(err.Error())
	}
	inpAlgo = tt.(inputAlgo)
	algoKey := string(resp.Payload)
	if algoKey != inpAlgo.Hash {
		t.Errorf("when adding algo, key does not corresponds to its hash - key: %s and hash %s", algoKey, inpAlgo.Hash)
	}

	// Query algo from key and check the consistency of returned arguments
	args = [][]byte{[]byte("query"), []byte(algoKey)}
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("when querying an algo with status %d and message %s", status, resp.Message)
	}
	algo := make(map[string]interface{})
	if err := bytesToStruct(resp.Payload, &algo); err != nil {
		t.Errorf("when unmarshalling queried algo with error %s", err)
	}
	if algo["name"] != inpAlgo.Name {
		t.Errorf("ledger algo name does not correspond to what was input: %s - %s", algo["name"], inpAlgo.Name)
	}
	if algo["storageAddress"] != inpAlgo.StorageAddress {
		t.Errorf("ledger algo description storage address does not correspond to what was input: %s - %s", algo["storageAddress"], inpAlgo.StorageAddress)
	}
	if algo["challengeKey"] != inpAlgo.ChallengeKey {
		t.Errorf("ledger algo challenge key does not correspond to what was input: %s - %s", algo["challengeKey"], inpAlgo.ChallengeKey)
	}
	if algo["permissions"] != inpAlgo.Permissions {
		t.Errorf("ledger algo permissions does not correspond to what was input: %s - %s", algo["permissions"], inpAlgo.Permissions)
	}
	if algo["description"].(map[string]interface{})["hash"] != inpAlgo.DescriptionHash {
		t.Errorf("ledger algo metrics hash does not correspond to what was input")
	}
	if algo["description"].(map[string]interface{})["storageAddress"] != inpAlgo.DescriptionStorageAddress {
		t.Errorf("ledger algo metrics address does not correspond to what was input")
	}

	// Query all algo and check consistency
	args = [][]byte{[]byte("queryAlgos")}
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("when querying algos - status %d and message %s", status, resp.Message)
	}
	var sPayload []map[string]interface{}
	if err := json.Unmarshal(resp.Payload, &sPayload); err != nil {
		t.Errorf("when unmarshalling queried algos")
	}
	payload := sPayload[0]
	delete(payload, "key")
	if !reflect.DeepEqual(payload, algo) {
		t.Errorf("when querying algos, dataset does not correspond to the input algo")
	}
}

func TestTraintuple(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := shim.NewMockStub("substra", scc)

	// Add traintuple with invalid field
	inpTraintuple := inputTraintuple{
		AlgoKey: "aaa",
	}
	args := inpTraintuple.createSample()
	resp := mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 500 {
		t.Errorf("when adding challenge with invalid hash, status %d and message %s", status, resp.Message)
	}

	// Add traintuple with unexisting algo
	inpTraintuple = inputTraintuple{}
	args = inpTraintuple.createSample()
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 500 {
		t.Errorf("when adding traintuple with unexisting algo, status %d and message %s", status, resp.Message)
	}

	// Properly add traintuple
	err, resp, tt := registerItem(*mockStub, "traintuple")
	if err != nil {
		t.Errorf(err.Error())
	}
	inpTraintuple = tt.(inputTraintuple)
	traintupleKey := string(resp.Payload)

	// Query traintuple from key and check the consistency of returned arguments
	args = [][]byte{[]byte("query"), []byte(traintupleKey)}
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("when querying the traintuple - status %d and message %s", status, resp.Message)
	}
	traintuple := make(map[string]interface{})
	if err := bytesToStruct(resp.Payload, &traintuple); err != nil {
		t.Errorf("when unmarshalling queried traintuple with error %s", err)
	}
	if a, b := traintuple["status"], "todo"; a != b {
		t.Errorf("wrong ledger traintuple status: %s - %s", a, b)
	}
	if a, b := traintuple["permissions"], "all"; a != b {
		t.Errorf("ledger traintuple permissions does not correspond to what was input: %s - %s", a, b)
	}
	if a, b := traintuple["log"], ""; a != b {
		t.Errorf("wrong ledger traintuple log: %s - %s", a, b)
	}
	if a, b := traintuple["challenge"].(map[string]interface{})["hash"], challengeDescriptionHash; a != b {
		t.Errorf("ledger traintuple challenge hash does not corresponds to what was input: %s - %s", a, b)
	}
	if a, b := traintuple["challenge"].(map[string]interface{})["metrics"].(map[string]interface{})["hash"], challengeMetricsHash; a != b {
		t.Errorf("ledger traintuple challenge hash does not corresponds to what was input: %s - %s", a, b)
	}
	if a, b := traintuple["challenge"].(map[string]interface{})["metrics"].(map[string]interface{})["storageAddress"], challengeMetricsStorageAddress; a != b {
		t.Errorf("ledger traintuple challenge hash does not corresponds to what was input: %s - %s", a, b)
	}
	algo := make(map[string]interface{})
	algo["hash"] = algoHash
	algo["storageAddress"] = algoStorageAddress
	algo["name"] = algoName
	if a, b := traintuple["algo"], algo; !reflect.DeepEqual(a, b) {
		t.Errorf("wrong ledger traintuple algo: %s - %s", a, b)
	}
	data := make(map[string]interface{})
	data["worker"] = worker
	data["keys"] = []interface{}{trainDataHash1, trainDataHash2}
	data["openerHash"] = datasetOpenerHash
	data["perf"] = 0.0
	if a, b := traintuple["trainData"], data; !reflect.DeepEqual(a, b) {
		t.Errorf("wrong ledger traintuple train data: %s - %s", a, b)
	}
	data["keys"] = []interface{}{testDataHash1, testDataHash2}
	if a, b := traintuple["testData"], data; !reflect.DeepEqual(a, b) {
		t.Errorf("wrong ledger traintuple test data: %s - %s", a, b)
	}

	// Query all traintuples and check consistency
	args = [][]byte{[]byte("queryTraintuples")}
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("when querying traintuples - status %d and message %s", status, resp.Message)
	}
	var sPayload []map[string]interface{}
	if err := json.Unmarshal(resp.Payload, &sPayload); err != nil {
		t.Errorf("when unmarshalling queried challenges")
	}
	payload := sPayload[0]
	delete(payload, "key")
	if !reflect.DeepEqual(payload, traintuple) {
		t.Errorf("when querying challenges, dataset does not correspond to the input challenge")
	}

	// Query traintuple with status todo and worker as trainworker and check consistency
	args = [][]byte{[]byte("queryFilter"), []byte("traintuple~trainWorker~status"), []byte(worker + ", todo")}
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("when querying traintuple of worker with todo status - status %d and message %s", status, resp.Message)
	}
	if err := bytesToStruct(resp.Payload, &sPayload); err != nil {
		t.Errorf("when unmarshalling queried traintuple with error %s", err)
	}
	if a, b := sPayload[0]["key"], traintupleKey; a != b {
		t.Errorf("wrong retrieved key when querying traintuple of worker with todo status: %s %s", a, b)
	}
	delete(sPayload[0], "key")
	if !reflect.DeepEqual(sPayload[0], traintuple) {
		t.Errorf("unexpected traintuple when querying traintuple with status todo and worker")
	}

	// Update status and check consistency
	trainDataPerf := "0.9"
	logTrain := "no error, ah ah ah"
	testDataPerf := "0.89"
	logTest := "still no error, suprah ah ah"
	argsSlice := [][][]byte{
		[][]byte{[]byte("logStartTrainTest"), []byte(traintupleKey), []byte("training")},
		[][]byte{[]byte("logSuccessTrain"), []byte(traintupleKey), []byte(modelHash + ", " + modelAddress),
			[]byte(trainDataPerf), []byte(logTrain)},
		[][]byte{[]byte("logStartTrainTest"), []byte(traintupleKey), []byte("testing")},
		[][]byte{[]byte("logSuccessTest"), []byte(traintupleKey), []byte(testDataPerf), []byte(logTest)},
	}
	traintupleStatus := []string{"training", "trained", "testing", "done"}
	for i, _ := range traintupleStatus {
		resp = mockStub.MockInvoke("42", argsSlice[i])
		if status := resp.Status; status != 200 {
			t.Errorf("when logging start %s with status %d and message %s",
				traintupleStatus[i], status, resp.Message)
		}
		args = [][]byte{[]byte("queryFilter"), []byte("traintuple~trainWorker~status"), []byte(worker + ", " + traintupleStatus[i])}
		resp = mockStub.MockInvoke("42", args)
		if status := resp.Status; status != 200 {
			t.Errorf("when querying traintuple of worker with %s status - status %d and message %s", traintupleStatus[i], status, resp.Message)
		}
		sPayload = make([]map[string]interface{}, 1)
		if err := bytesToStruct(resp.Payload, &sPayload); err != nil {
			t.Errorf("when unmarshalling queried traintuple with error %s", err)
		}
		if a, b := sPayload[0]["key"], traintupleKey; a != b {
			t.Errorf("wrong retrieved key when querying traintuple of worker with %s status: %s %s", traintupleStatus[i], a, b)
		}
		delete(sPayload[0], "key")
		if a, b := sPayload[0]["status"], traintupleStatus[i]; a != b {
			t.Errorf("wrong retrieved status when querying traintuple of worker with %s status: %s %s", traintupleStatus[i], a, b)
		}
	}

	// Query Traintuple From key
	args = [][]byte{[]byte("query"), []byte(traintupleKey)}
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("when querying traintuple with status %d and message %s",
			status, resp.Message)
	}
	respTraintuple := resp.Payload
	endTraintuple := Traintuple{}
	if err := bytesToStruct(respTraintuple, &endTraintuple); err != nil {
		t.Errorf("when unmarshalling queried traintuple with error %s", err)
	}
	if a := endTraintuple.Log; a != logTrain+logTest {
		t.Errorf("because retrieved log in traintuple does not correspond to what "+
			"was submitted: %s", a)
	}
	endModel := HashDress{
		Hash:           modelHash,
		StorageAddress: modelAddress}
	if endTraintuple.EndModel.Hash != endModel.Hash || endTraintuple.EndModel.StorageAddress != endModel.StorageAddress {
		t.Errorf("because retrieved endModel in traintuple does not correspond to what "+
			"was submitted: %s, %s", endTraintuple.EndModel, endModel)
	}
	// Query Model
	args = [][]byte{[]byte("queryModel"), []byte(modelHash)}
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("when querying model with status %d and message %s", status, resp.Message)
	}
	if !bytes.Equal(resp.Payload, respTraintuple) {
		t.Errorf("when querying model, didn't get the associated traintuple")
	}

	// query traintuples related to a model
	args = [][]byte{[]byte("queryModelTraintuples"), []byte(modelHash)}
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("when querying model traintuples with status %d and message %s",
			status, resp.Message)
	}
	payload = make(map[string]interface{})
	json.Unmarshal(resp.Payload, &payload)
	if l := len(payload); l != 2 {
		t.Errorf("when querying model traintuple, payload should contain an algo "+
			"and a traintuple, but it contains %d elements", l)
	}
}

func TestPipeline(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := shim.NewMockStub("substra", scc)

	fmt.Println("#### ------------ Add Dataset ------------")
	inpDataset := inputDataset{}
	printArgsNames("registerDataset", getFieldNames(&inpDataset))
	args := inpDataset.createSample()
	printArgs(args, "invoke")
	resp := mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("when adding dataset with status %d and message %s", status, resp.Message)
	}
	fmt.Printf(">  %s \n\n", string(resp.Payload))
	// Get dataset key from Payload
	datasetKey := string(resp.Payload)

	fmt.Println("#### ------------ Query Dataset From key ------------")
	printArgsNames("query", []string{"elementKey"})
	args = [][]byte{[]byte("query"), []byte(datasetKey)}
	printArgs(args, "query")
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("when querying a dataset with status %d and message %s", status, resp.Message)
	}
	fmt.Printf(">  %s \n\n", string(resp.Payload))

	fmt.Println("#### ------------ Add test Data ------------")
	inpData := inputData{
		Hashes:   testDataHash1 + ", " + testDataHash2,
		TestOnly: "true",
	}
	printArgsNames("registerData", getFieldNames(&inpData))
	args = inpData.createSample()
	printArgs(args, "invoke")
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("when adding test data with status %d and message %s", status, resp.Message)
	}
	fmt.Printf(">  %s \n\n", string(resp.Payload))

	fmt.Println("#### ------------ Add Challenge ------------")
	inpChallenge := inputChallenge{}
	printArgsNames("registerChallenge", getFieldNames(&inpChallenge))
	args = inpChallenge.createSample()
	printArgs(args, "invoke")
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("when adding challenge with status %d and message %s", status, resp.Message)
	}
	fmt.Printf(">  %s \n\n", string(resp.Payload))

	fmt.Println("#### ------------ Add Algo ------------")
	inpAlgo := inputAlgo{}
	printArgsNames("registerAlgo", getFieldNames(&inpAlgo))
	args = inpAlgo.createSample()
	printArgs(args, "invoke")
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("when adding algo with status %d and message %s", status, resp.Message)
	}
	fmt.Printf(">  %s \n\n", string(resp.Payload))

	fmt.Println("#### ------------ Add Train Data ------------")
	inpData = inputData{}
	printArgsNames("registerData", getFieldNames(&inpData))
	args = inpData.createSample()
	printArgs(args, "invoke")
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("when adding train data with status %d and message %s", status, resp.Message)
	}
	fmt.Printf(">  %s \n\n", string(resp.Payload))

	fmt.Println("#### ------------ Query Datasets ------------")
	args = [][]byte{[]byte("queryDatasets")}
	printArgs(args, "query")
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("when querying dataset with status %d and message %s", status, resp.Message)
	}
	fmt.Printf(">  %s \n\n", string(resp.Payload))

	fmt.Println("#### ------------ Query Challenges ------------")
	args = [][]byte{[]byte("queryChallenges")}
	printArgs(args, "query")
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("when querying challenge with status %d and message %s", status, resp.Message)
	}
	fmt.Printf(">  %s \n\n", string(resp.Payload))

	fmt.Println("#### ------------ Add Traintuple ------------")
	inpTraintuple := inputTraintuple{}
	printArgsNames("createTraintuple", getFieldNames(&inpTraintuple))
	args = inpTraintuple.createSample()
	printArgs(args, "invoke")
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("when adding traintuple with status %d and message %s", status, resp.Message)
	}
	fmt.Printf(">  %s \n\n", string(resp.Payload))
	// Get traintuple key from Payload
	traintupleKey := resp.Payload
	// check not possible to create same traintuple
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 500 {
		t.Errorf("when adding same traintuple with status %d and message %s", status, resp.Message)
	}
	// Get owner of the traintuple
	args = [][]byte{[]byte("query"), traintupleKey}
	resp = mockStub.MockInvoke("42", args)
	respTraintuple := resp.Payload
	traintuple := Traintuple{}
	if err := bytesToStruct(respTraintuple, &traintuple); err != nil {
		t.Errorf("when unmarshalling queried traintuple with error %s", err)
	}
	trainWorker := traintuple.TrainData.Worker

	fmt.Println("#### ------------ Query Traintuples of worker with todo status ------------")
	args = [][]byte{[]byte("queryFilter"), []byte("traintuple~trainWorker~status"), []byte(trainWorker + ", todo")}
	printArgs(args, "invoke")
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("when querying traintuple of worker with todo status - with status %d and message %s", status, resp.Message)
	}
	fmt.Printf(">  %s \n\n", string(resp.Payload))

	fmt.Println("#### ------------ Log Start Training ------------")
	args = [][]byte{[]byte("logStartTrainTest"), traintupleKey, []byte("training")}
	printArgs(args, "invoke")
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("when logging start training with status %d and message %s",
			status, resp.Message)
	}
	fmt.Printf(">  %s \n\n", string(resp.Payload))

	fmt.Println("#### ------------ Log Success Training ------------")
	trainDataPerf := "0.9"
	logTrain := "no error, ah ah ah"
	args = [][]byte{[]byte("logSuccessTrain"), traintupleKey, []byte(modelHash + ", " + modelAddress),
		[]byte(trainDataPerf), []byte(logTrain)}
	printArgs(args, "invoke")
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("when logging successful training with status %d and message %s",
			status, resp.Message)
	}
	fmt.Printf(">  %s \n\n", string(resp.Payload))

	fmt.Println("#### ------------ Log Start Testing ------------")
	args = [][]byte{[]byte("logStartTrainTest"), traintupleKey, []byte("testing")}
	printArgs(args, "invoke")
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("when logging start testing with status %d and message %s",
			status, resp.Message)
	}
	fmt.Printf(">  %s \n\n", string(resp.Payload))

	fmt.Println("#### ------------ Log Success Testing ------------")
	testDataPerf := "0.89"
	logTest := "still no error, suprah ah ah"
	args = [][]byte{[]byte("logSuccessTest"), traintupleKey, []byte(testDataPerf), []byte(logTest)}
	printArgs(args, "invoke")
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("when logging successful training with status %d and message %s",
			status, resp.Message)
	}
	fmt.Printf(">  %s \n\n", string(resp.Payload))

	fmt.Println("#### ------------ Query Traintuple From key ------------")
	args = [][]byte{[]byte("query"), traintupleKey}
	printArgs(args, "query")
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("when querying traintuple with status %d and message %s",
			status, resp.Message)
	}
	fmt.Printf(">  %s \n\n", string(resp.Payload))

	fmt.Println("#### ------------ Query Model ------------")
	args = [][]byte{[]byte("queryModel"), []byte(modelHash)}
	printArgs(args, "query")
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("when querying model with status %d and message %s", status, resp.Message)
	}
	fmt.Printf(">  %s \n\n", string(resp.Payload))

	fmt.Println("#### ------------ Query Model Traintuples ------------")
	args = [][]byte{[]byte("queryModelTraintuples"), []byte(modelHash)}
	printArgs(args, "query")
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("when querying model traintuples with status %d and message %s",
			status, resp.Message)
	}
	fmt.Printf(">  %s \n\n", string(resp.Payload))

	fmt.Println("#### ------------ Query Dataset Data ------------")
	args = [][]byte{[]byte("queryDatasetData"), []byte(datasetOpenerHash)}
	printArgs(args, "query")
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("when querying dataset data with status %d and message %s", status, resp.Message)
	}
	fmt.Printf(">  %s \n\n", string(resp.Payload))

}
