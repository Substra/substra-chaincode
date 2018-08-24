package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

const challengeDescriptionHash = "5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379"
const datasetOpenerHash = "do1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
const trainDataHash1 = "aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
const trainDataHash2 = "aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
const testDataHash1 = "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
const testDataHash2 = "da2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
const algoHash = "fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
const modelHash = "modbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482mod"
const modelAddress = "https://substrabac/model/toto"

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

func createSampleDataset(dataset inputDataset) [][]byte {
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
	args, _ := inputStructToBytes(&dataset)
	args = append([][]byte{[]byte("registerDataset")}, args...)
	printArgs(args, "invoke")
	return args
}

func createSampleData(data inputData) [][]byte {
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
	args, _ := inputStructToBytes(&data)
	args = append([][]byte{[]byte("registerData")}, args...)
	printArgs(args, "invoke")
	return args
}

func createSampleChallenge(challenge inputChallenge) [][]byte {
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
		challenge.MetricsHash = "fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482d8d"
	}
	if challenge.MetricsStorageAddress == "" {
		challenge.MetricsStorageAddress = "https://toto/challenge/222/metrics"
	}
	if challenge.TestDataKeys == "" {
		challenge.TestDataKeys = "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
	}
	challenge.Permissions = "all"
	args, _ := inputStructToBytes(&challenge)
	args = append([][]byte{[]byte("registerChallenge")}, args...)
	printArgs(args, "invoke")
	return args
}

func createSampleAlgo(algo inputAlgo) [][]byte {
	if algo.Name == "" {
		algo.Name = "hog + svm"
	}
	if algo.Hash == "" {
		algo.Hash = algoHash
	}
	if algo.StorageAddress == "" {
		algo.StorageAddress = "https://toto/algo/222/algo"
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
	args, _ := inputStructToBytes(&algo)
	args = append([][]byte{[]byte("registerAlgo")}, args...)
	printArgs(args, "invoke")
	return args
}

func createSampleTraintuple(traintuple inputTraintuple) [][]byte {
	if traintuple.ChallengeKey == "" {
		traintuple.ChallengeKey = challengeDescriptionHash
	}
	if traintuple.AlgoKey == "" {
		traintuple.AlgoKey = algoHash
	}
	if traintuple.StartModelKey == "" {
		traintuple.StartModelKey = algoHash
	}
	if traintuple.TrainDataKeys == "" {
		traintuple.TrainDataKeys = trainDataHash1 + ", " + trainDataHash2
	}
	args, _ := inputStructToBytes(&traintuple)
	args = append([][]byte{[]byte("createTraintuple")}, args...)
	printArgs(args, "invoke")
	return args
}

func createSamplePerf(dataHashes []string) []byte {
	dataPerf := ""
	for i, dataHash := range dataHashes {
		dataPerf += dataHash + ":0.9" + strconv.Itoa(i)
		if i < len(dataHashes)-1 {
			dataPerf += ", "
		}
	}
	return []byte(dataPerf)
}

func TestPipeline(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := shim.NewMockStub("substra", scc)
	// Add dataset
	fmt.Println("#### ------------ Add Dataset ------------")
	inpDataset := inputDataset{}
	args := createSampleDataset(inpDataset)
	resp := mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("testPipeline failed when adding dataset with status %d and message %s", status, resp.Message)
	}
	fmt.Println(">  " + string(resp.Payload))
	// Get dataset key from Payload
	datasetKey := string(resp.Payload)
	if datasetKey != datasetOpenerHash {
		t.Errorf("dataset key does not correspond to dataset opener hash: %s - %s", datasetKey, datasetOpenerHash)
	}
	// Query Dataset
	fmt.Println("#### ------------ Query Dataset From key ------------")
	args = [][]byte{[]byte("query"), []byte(datasetKey)}
	printArgs(args, "query")
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("testPipeline failed when querying a dataset with status %d and message %s", status, resp.Message)
	}
	fmt.Println(">  " + string(resp.Payload))
	// Add test data
	fmt.Println("#### ------------ Add test Data ------------")
	inpData := inputData{
		Hashes:   testDataHash1 + ", " + testDataHash2,
		TestOnly: "true",
	}
	args = createSampleData(inpData)
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("testPipeline failed when adding test data with status %d and message %s", status, resp.Message)
	}
	fmt.Println(">  " + string(resp.Payload))
	// Add challenge
	fmt.Println("#### ------------ Add Challenge ------------")
	inpChallenge := inputChallenge{}
	args = createSampleChallenge(inpChallenge)
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("testPipeline failed when adding challenge with status %d and message %s", status, resp.Message)
	}
	fmt.Println(">  " + string(resp.Payload))
	// Add algo
	fmt.Println("#### ------------ Add Algo ------------")
	inpAlgo := inputAlgo{}
	args = createSampleAlgo(inpAlgo)
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("testPipeline failed when adding algo with status %d and message %s", status, resp.Message)
	}
	fmt.Println(">  " + string(resp.Payload))
	// Add train data
	fmt.Println("#### ------------ Add Train Data ------------")
	inpData = inputData{}
	args = createSampleData(inpData)
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("testPipeline failed when adding train data with status %d and message %s", status, resp.Message)
	}
	fmt.Println(">  " + string(resp.Payload))
	// Query all dataset
	fmt.Println("#### ------------ Query Datasets ------------")
	args = [][]byte{[]byte("queryDatasets")}
	printArgs(args, "query")
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("testPipeline failed when querying dataset with status %d and message %s", status, resp.Message)
	}
	fmt.Println(">  " + string(resp.Payload))

	// Query all challenges
	fmt.Println("#### ------------ Query Challenges ------------")
	args = [][]byte{[]byte("queryChallenges")}
	printArgs(args, "query")
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("testPipeline failed when querying challenge with status %d and message %s", status, resp.Message)
	}
	fmt.Println(">  " + string(resp.Payload))
	// Add trainTuple
	fmt.Println("#### ------------ Add Traintuple ------------")
	inpTraintuple := inputTraintuple{}
	args = createSampleTraintuple(inpTraintuple)
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("testPipeline failed when adding traintuple with status %d and message %s", status, resp.Message)
	}
	fmt.Println(">  " + string(resp.Payload))
	// Get traintuple key from Payload
	traintupleKey := resp.Payload
	// Log Start Training
	fmt.Println("#### ------------ Log Start Training ------------")
	args = [][]byte{[]byte("logStartTrainTest"), traintupleKey, []byte("training")}
	printArgs(args, "invoke")
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("testPipeline failed when logging start training with status %d and message %s",
			status, resp.Message)
	}
	fmt.Println(">  " + string(resp.Payload))
	// Log Success Training
	fmt.Println("#### ------------ Log Success Training ------------")
	trainDataPerf := createSamplePerf([]string{trainDataHash1, trainDataHash2})
	logTrain := "no error, ah ah ah"
	args = [][]byte{[]byte("logSuccessTrain"), traintupleKey, []byte(modelHash + ", " + modelAddress),
		trainDataPerf, []byte(logTrain)}
	printArgs(args, "invoke")
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("testPipeline failed when logging successful training with status %d and message %s",
			status, resp.Message)
	}
	fmt.Println(">  " + string(resp.Payload))
	// Log Start Testing
	fmt.Println("#### ------------ Log Start Testing ------------")
	args = [][]byte{[]byte("logStartTrainTest"), traintupleKey, []byte("testing")}
	printArgs(args, "invoke")
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("testPipeline failed when logging start testing with status %d and message %s",
			status, resp.Message)
	}
	fmt.Println(">  " + string(resp.Payload))
	// Log Success Testing
	fmt.Println("#### ------------ Log Success Testing ------------")
	testDataPerf := createSamplePerf([]string{testDataHash1, testDataHash2})
	perf := []byte("0.99")
	logTest := "still no error, suprah ah ah"
	args = [][]byte{[]byte("logSuccessTest"), traintupleKey, testDataPerf, perf, []byte(logTest)}
	printArgs(args, "invoke")
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("testPipeline failed when logging successful training with status %d and message %s",
			status, resp.Message)
	}
	fmt.Println(">  " + string(resp.Payload))
	// Query traintuple and check elements
	fmt.Println("#### ------------ Query Traintuple From key ------------")
	args = [][]byte{[]byte("query"), traintupleKey}
	printArgs(args, "query")
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("testPipeline failed when querying traintuple with status %d and message %s",
			status, resp.Message)
	}
	respTraintuple := resp.Payload
	traintuple := Traintuple{}
	err := bytesToStruct(respTraintuple, &traintuple)
	if err != nil {
		t.Errorf("testPipeline failed when unmarshalling queried traintuple with error %s", err)
	}
	if traintuple.Log != logTrain+logTest {
		t.Errorf("testPipeline failed because retrieved log in traintuple does not correspond to what "+
			"was submitted: %s", traintuple.Log)
	}
	endModel := HashDress{
		Hash:           modelHash,
		StorageAddress: modelAddress}
	if traintuple.EndModel.Hash != endModel.Hash || traintuple.EndModel.StorageAddress != endModel.StorageAddress {
		t.Errorf("testPipeline failed because retrieved endModel in traintuple does not correspond to what "+
			"was submitted: %s, %s", traintuple.EndModel, endModel)
	}
	fmt.Println(">  " + string(resp.Payload))
	// Query model and checks it returns the associated traintuple
	fmt.Println("#### ------------ Query Model ------------")
	args = [][]byte{[]byte("queryModel"), []byte(modelHash)}
	printArgs(args, "query")
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("testPipeline failed when querying model with status %d and message %s", status, resp.Message)
	}
	if !bytes.Equal(resp.Payload, respTraintuple) {
		t.Errorf("testPipeline failed when querying model, didn't get the associated traintuple")
	}
	fmt.Println(">  " + string(resp.Payload))
	// Query model traintuples
	fmt.Println("#### ------------ Query Model Traintuples ------------")
	args = [][]byte{[]byte("queryModelTraintuples"), []byte(modelHash)}
	printArgs(args, "query")
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("testPipeline failed when querying model traintuples with status %d and message %s",
			status, resp.Message)
	}
	payload := make(map[string]interface{})
	json.Unmarshal(resp.Payload, &payload)
	if l := len(payload); l != 2 {
		t.Errorf("testPipeline failed when querying model traintuple, payload should contain an algo "+
			"and a traintuple, but it contains %d elements", l)
	}
	fmt.Println(">  " + string(resp.Payload))
	// Query dataset data
	fmt.Println("#### ------------ Query Dataset Data ------------")
	args = [][]byte{[]byte("queryDatasetData"), []byte(datasetOpenerHash)}
	printArgs(args, "query")
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("testPipeline failed when querying dataset data with status %d and message %s", status, resp.Message)
	}
	payload = make(map[string]interface{})
	json.Unmarshal(resp.Payload, &payload)
	if _, ok := payload[datasetOpenerHash]; !ok {
		t.Errorf("testPipeline failed when querying dataset data, payload should contain the dataset info")
	}
	if _, ok := payload["trainDataKeys"]; !ok {
		t.Errorf("testPipeline failed when querying dataset data, payload should contain related data")
	}

	fmt.Println(">  " + string(resp.Payload))
	// TODO ADD CHECK content of resp.Payload

}

func TestFails(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := shim.NewMockStub("substra", scc)
	fmt.Println("#### ------------ TO FAIL - Query Model Traintuples with unexisting model hash ------------")
	unexistingModelHash := "modbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482tot"
	args := [][]byte{[]byte("queryModelTraintuples"), []byte(unexistingModelHash)}
	resp := mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 500 {
		t.Errorf("testFails did not fail when querying model traintuples, got status %d", status)
	}
}
