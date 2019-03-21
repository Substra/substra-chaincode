package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	peer "github.com/hyperledger/fabric/protos/peer"
)

const objectiveDescriptionHash = "5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379"
const objectiveDescriptionStorageAddress = "https://toto/objective/222/description"
const objectiveMetricsHash = "4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379"
const objectiveMetricsStorageAddress = "https://toto/objective/222/metrics"
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
const traintupleKey = "337c5b7d78ffa157471cc790b61caa99b486cbf30b049c8c0550ed40f3fa1d4f"

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
	s := "Smart contract: `" + fnName + "`  \n Inputs: `" + strings.Join(argsNames, "`, `") + "`"
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
	if data.DatasetKeys == "" {
		data.DatasetKeys = datasetOpenerHash
	}
	if data.TestOnly == "" {
		data.TestOnly = "false"
	}
	args, _ := inputStructToBytes(data)
	args = append([][]byte{[]byte("registerData")}, args...)
	return args
}

func (objective *inputObjective) createSample() [][]byte {
	if objective.Name == "" {
		objective.Name = "MSI classification"
	}
	if objective.DescriptionHash == "" {
		objective.DescriptionHash = objectiveDescriptionHash
	}
	if objective.DescriptionStorageAddress == "" {
		objective.DescriptionStorageAddress = "https://toto/objective/222/description"
	}
	if objective.MetricsName == "" {
		objective.MetricsName = "accuracy"
	}
	if objective.MetricsHash == "" {
		objective.MetricsHash = objectiveMetricsHash
	}
	if objective.MetricsStorageAddress == "" {
		objective.MetricsStorageAddress = objectiveMetricsStorageAddress
	}
	if objective.TestData == "" {
		objective.TestData = datasetOpenerHash + ":" + testDataHash1 + ", " + testDataHash2
	}
	objective.Permissions = "all"
	args, _ := inputStructToBytes(objective)
	args = append([][]byte{[]byte("registerObjective")}, args...)
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
	if algo.ObjectiveKey == "" {
		algo.ObjectiveKey = objectiveDescriptionHash
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
	if traintuple.InModels == "" {
		traintuple.InModels = ""
	}
	if traintuple.DatasetKey == "" {
		traintuple.DatasetKey = datasetOpenerHash
	}
	if traintuple.DataKeys == "" {
		traintuple.DataKeys = trainDataHash1 + ", " + trainDataHash2
	}
	args, _ := inputStructToBytes(traintuple)
	args = append([][]byte{[]byte("createTraintuple")}, args...)
	return args
}

func (testtuple *inputTesttuple) createSample() [][]byte {
	if testtuple.TraintupleKey == "" {
		testtuple.TraintupleKey = traintupleKey
	}
	args, _ := inputStructToBytes(testtuple)
	args = append([][]byte{[]byte("createTesttuple")}, args...)
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
		Hashes:      testDataHash1 + ", " + testDataHash2,
		DatasetKeys: datasetOpenerHash,
		TestOnly:    "true",
	}
	args = inpData.createSample()
	resp = mockStub.MockInvoke("42", args)
	if resp.Status != 200 {
		return fmt.Errorf("when adding test data with status %d and message %s", resp.Status, resp.Message), resp, inpData
	} else if itemType == "testData" {
		return nil, resp, inpData
	}
	// 3. add objective
	inpObjective := inputObjective{}
	args = inpObjective.createSample()
	resp = mockStub.MockInvoke("42", args)
	if resp.Status != 200 {
		return fmt.Errorf("when adding objective with status %d and message %s", resp.Status, resp.Message), resp, inpObjective
	} else if itemType == "objective" {
		return nil, resp, inpObjective
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
	printArgsNames("queryDataset", []string{"elementKey"})
	args = [][]byte{[]byte("queryDataset"), []byte(datasetKey)}
	printArgs(args, "queryDataset")
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

	fmt.Println("#### ------------ Add Objective ------------")
	inpObjective := inputObjective{}
	printArgsNames("registerObjective", getFieldNames(&inpObjective))
	args = inpObjective.createSample()
	printArgs(args, "invoke")
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("when adding objective with status %d and message %s", status, resp.Message)
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

	fmt.Println("#### ------------ Query Objectives ------------")
	args = [][]byte{[]byte("queryObjectives")}
	printArgs(args, "query")
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("when querying objective with status %d and message %s", status, resp.Message)
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
	args = [][]byte{[]byte("queryTraintuple"), traintupleKey}
	resp = mockStub.MockInvoke("42", args)
	respTraintuple := resp.Payload
	traintuple := outputTraintuple{}
	if err := bytesToStruct(respTraintuple, &traintuple); err != nil {
		t.Errorf("when unmarshalling queried traintuple with error %s", err)
	}
	trainWorker := traintuple.Data.Worker

	fmt.Println("#### ------------ Add Traintuple with inModel from previous traintuple ------------")
	inpTraintuple = inputTraintuple{
		InModels: string(traintupleKey),
	}
	printArgsNames("createTraintuple", getFieldNames(&inpTraintuple))
	args = inpTraintuple.createSample()
	printArgs(args, "invoke")
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("when adding traintuple with status %d and message %s", status, resp.Message)
	}
	fmt.Printf(">  %s \n\n", string(resp.Payload))
	todoTraintupleKey := string(resp.Payload)

	fmt.Println("#### ------------ Query Traintuples of worker with todo status ------------")
	args = [][]byte{[]byte("queryFilter"), []byte("traintuple~worker~status"), []byte(trainWorker + ", todo")}
	printArgs(args, "invoke")
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("when querying traintuple of worker with todo status - with status %d and message %s", status, resp.Message)
	}
	fmt.Printf(">  %s \n\n", string(resp.Payload))

	fmt.Println("#### ------------ Log Start Training ------------")
	args = [][]byte{[]byte("logStartTrain"), traintupleKey}
	printArgs(args, "invoke")
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("when logging start training with status %d and message %s", status, resp.Message)
	}
	fmt.Printf(">  %s \n\n", string(resp.Payload))

	fmt.Println("#### ------------ Log Success Training ------------")
	perf := "0.9"
	log := "no error, ah ah ah"
	args = [][]byte{[]byte("logSuccessTrain"), traintupleKey, []byte(modelHash + ", " + modelAddress),
		[]byte(perf), []byte(log)}
	printArgs(args, "invoke")
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("when logging successful training with status %d and message %s", status, resp.Message)
	}
	fmt.Printf(">  %s \n\n", string(resp.Payload))

	fmt.Println("#### ------------ Query Traintuple From key ------------")
	args = [][]byte{[]byte("queryTraintuple"), traintupleKey}
	printArgs(args, "queryTraintuple")
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("when querying traintuple with status %d and message %s", status, resp.Message)
	}
	fmt.Printf(">  %s \n\n", string(resp.Payload))

	fmt.Println("#### ------------ Add Non-Certified Testtuple ------------")
	inpTesttuple := inputTesttuple{
		DatasetKey: datasetOpenerHash,
		DataKeys:   trainDataHash1 + ", " + trainDataHash2,
	}
	printArgsNames("createTesttuple", getFieldNames(&inpTesttuple))
	args = inpTesttuple.createSample()
	printArgs(args, "invoke")
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("when adding testtuple with status %d and message %s", status, resp.Message)
	}
	fmt.Printf(">  %s \n\n", string(resp.Payload))

	fmt.Println("#### ------------ Add Certified Testtuple ------------")
	inpTesttuple = inputTesttuple{}
	printArgsNames("createTesttuple", getFieldNames(&inpTesttuple))
	args = inpTesttuple.createSample()
	printArgs(args, "invoke")
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("when adding testtuple with status %d and message %s", status, resp.Message)
	}
	fmt.Printf(">  %s \n\n", string(resp.Payload))
	// Get testtuple key from Payload
	testtupleKey := resp.Payload
	// check not possible to create same testtuple
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 500 {
		t.Errorf("when adding same testtuple with status %d and message %s", status, resp.Message)
	}
	// Get owner of the testtuple
	args = [][]byte{[]byte("queryTesttuple"), testtupleKey}
	resp = mockStub.MockInvoke("42", args)
	respTesttuple := resp.Payload
	testtuple := Testtuple{}
	if err := bytesToStruct(respTesttuple, &testtuple); err != nil {
		t.Errorf("when unmarshalling queried testtuple with error %s", err)
	}
	testWorker := testtuple.Data.Worker

	fmt.Println("#### ------------ Add Testtuple with not done traintuple ------------")
	inpTesttuple = inputTesttuple{
		TraintupleKey: todoTraintupleKey,
	}
	printArgsNames("createTesttuple", getFieldNames(&inpTesttuple))
	args = inpTesttuple.createSample()
	printArgs(args, "invoke")
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("when adding testtuple with status %d and message %s", status, resp.Message)
	}
	fmt.Printf(">  %s \n\n", string(resp.Payload))

	fmt.Println("#### ------------ Query Testtuples of worker with todo status ------------")
	args = [][]byte{[]byte("queryFilter"), []byte("testtuple~worker~status"), []byte(testWorker + ", todo")}
	printArgs(args, "invoke")
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("when querying testtuple of worker with todo status - with status %d and message %s", status, resp.Message)
	}
	fmt.Printf(">  %s \n\n", string(resp.Payload))

	fmt.Println("#### ------------ Log Start Testing ------------")
	args = [][]byte{[]byte("logStartTest"), testtupleKey}
	printArgs(args, "invoke")
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("when logging start testing with status %d and message %s", status, resp.Message)
	}
	fmt.Printf(">  %s \n\n", string(resp.Payload))

	fmt.Println("#### ------------ Log Success Testing ------------")
	perf = "0.89"
	log = "still no error, suprah ah ah"
	args = [][]byte{[]byte("logSuccessTest"), testtupleKey, []byte(perf), []byte(log)}
	printArgs(args, "invoke")
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("when logging successful training with status %d and message %s", status, resp.Message)
	}
	fmt.Printf(">  %s \n\n", string(resp.Payload))

	fmt.Println("#### ------------ Query Testtuple from its key ------------")
	args = [][]byte{[]byte("queryTesttuple"), testtupleKey}
	printArgs(args, "queryTesttuple")
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("when querying testtuple with status %d and message %s", status, resp.Message)
	}
	fmt.Printf(">  %s \n\n", string(resp.Payload))

	fmt.Println("#### ------------ Query all Testtuples ------------")
	args = [][]byte{[]byte("queryTesttuples")}
	printArgs(args, "queryTesttuples")
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("when querying testtuple with status %d and message %s", status, resp.Message)
	}
	fmt.Printf(">  %s \n\n", string(resp.Payload))

	fmt.Println("#### ------------ Query details about a model ------------")
	args = [][]byte{[]byte("queryModelDetails"), []byte(traintupleKey)}
	printArgs(args, "query")
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("when querying model details with status %d and message %s", status, resp.Message)
	}
	fmt.Printf(">  %s \n\n", string(resp.Payload))

	fmt.Println("#### ------------ Query all models ------------")
	args = [][]byte{[]byte("queryModels")}
	printArgs(args, "query")
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("when querying model tuples with status %d and message %s", status, resp.Message)
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
