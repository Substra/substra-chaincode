package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	peer "github.com/hyperledger/fabric/protos/peer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const objectiveDescriptionHash = "5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379"
const objectiveDescriptionStorageAddress = "https://toto/objective/222/description"
const objectiveMetricsHash = "4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379"
const objectiveMetricsStorageAddress = "https://toto/objective/222/metrics"
const dataManagerOpenerHash = "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
const trainDataSampleHash1 = "aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
const trainDataSampleHash2 = "aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
const testDataSampleHash1 = "bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
const testDataSampleHash2 = "bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
const algoHash = "fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
const algoStorageAddress = "https://toto/algo/222/algo"
const algoName = "hog + svm"
const modelHash = "eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed"
const modelAddress = "https://substrabac/model/toto"
const worker = "bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0"
const traintupleKey = "8e29bacef1250f8c3bd6ccc72455f764b74ef7e66b9157fd6cd2b0cecef1c687"

var (
	pipeline = flag.Bool("pipeline", false, "Print out the pipeline test output")
	readme   = flag.String("readme", "../README.md", "Pass the path to the README and compare it to the output")
	update   = flag.Bool("update", false, "Update the README.md file")
)

func TestInit(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := shim.NewMockStub("substra", scc)

	// resp := mockStub.MockInit("42", [][]byte{[]byte("init")})
	resp := mockStub.MockInit("42", [][]byte{[]byte("init")})
	assert.EqualValuesf(t, 200, resp.Status, "init failed with status %d and message %s", resp.Status, resp.Message)
}

func assetToJSON(asset interface{}) []byte {
	assetjson, _ := json.Marshal(asset)
	payload, _ := json.Marshal(string(assetjson))
	return payload
}

func keyToJSON(key string) []byte {
	return assetToJSON(inputHashe{Key: key})
}
func printArgs(buf io.Writer, args [][]byte, command string) {
	s := "```\npeer chaincode " + command + " -n mycc -c '{\"Args\":["
	for i, arg := range args {
		s += "\"" + string(arg) + "\""
		if i+1 < len(args) {
			s += ","
		}
	}
	s += "]}' -C myc\n```"
	fmt.Fprintln(buf, s)
}
func printArgsNames(buf io.Writer, fnName string, argsNames []string) {
	s := "Smart contract: `" + fnName + "`  \n Inputs: `" + strings.Join(argsNames, "`, `") + "`"
	fmt.Fprintln(buf, s)
}

func (dataManager *inputDataManager) createSample() [][]byte {
	if dataManager.Name == "" {
		dataManager.Name = "liver slide"
	}
	if dataManager.OpenerHash == "" {
		dataManager.OpenerHash = dataManagerOpenerHash
	}
	if dataManager.OpenerStorageAddress == "" {
		dataManager.OpenerStorageAddress = "https://toto/dataManager/42234/opener"
	}
	if dataManager.Type == "" {
		dataManager.Type = "images"
	}
	if dataManager.DescriptionHash == "" {
		dataManager.DescriptionHash = "8d4bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eee"
	}
	if dataManager.DescriptionStorageAddress == "" {
		dataManager.DescriptionStorageAddress = "https://toto/dataManager/42234/description"
	}
	dataManager.Permissions = "all"
	args := append([][]byte{[]byte("registerDataManager")}, assetToJSON(dataManager))

	return args
}
func (dataSample *inputDataSample) createSample() [][]byte {
	if dataSample.Hashes == "" {
		dataSample.Hashes = trainDataSampleHash1 + ", " + trainDataSampleHash2
	}
	if dataSample.DataManagerKeys == "" {
		dataSample.DataManagerKeys = dataManagerOpenerHash
	}
	if dataSample.TestOnly == "" {
		dataSample.TestOnly = "false"
	}
	args := append([][]byte{[]byte("registerDataSample")}, assetToJSON(dataSample))
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
	if objective.TestDataset == "" {
		objective.TestDataset = dataManagerOpenerHash + ":" + testDataSampleHash1 + ", " + testDataSampleHash2
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
	if traintuple.ObjectiveKey == "" {
		traintuple.ObjectiveKey = objectiveDescriptionHash
	}
	if traintuple.DataManagerKey == "" {
		traintuple.DataManagerKey = dataManagerOpenerHash
	}
	if traintuple.DataSampleKeys == "" {
		traintuple.DataSampleKeys = trainDataSampleHash1 + ", " + trainDataSampleHash2
	}
	args := append([][]byte{[]byte("createTraintuple")}, assetToJSON(traintuple))
	return args
}

func (success *inputSuccessTrain) createSample() [][]byte {
	if success.Key == "" {
		success.Key = traintupleKey
	}
	if success.Log == "" {
		success.Log = "no error, ah ah ah"
	}
	if success.Perf == 0 {
		success.Perf = 0.9
	}
	if success.OutModel.Hash == "" {
		success.OutModel.Hash = modelHash
	}
	if success.OutModel.StorageAddress == "" {
		success.OutModel.StorageAddress = modelAddress
	}
	args := append([][]byte{[]byte("logSuccessTrain")}, assetToJSON(success))
	return args
}
func (testtuple *inputTesttuple) createSample() [][]byte {
	if testtuple.TraintupleKey == "" {
		testtuple.TraintupleKey = traintupleKey
	}
	args, _ := inputStructToBytes(testtuple)
	args = append([][]byte{[]byte("createTesttuple")}, assetToJSON(testtuple))
	return args
}

func registerItem(t *testing.T, mockStub shim.MockStub, itemType string) (peer.Response, interface{}) {
	// 1. add dataManager
	inpDataManager := inputDataManager{}
	args := inpDataManager.createSample()
	resp := mockStub.MockInvoke("42", args)
	require.EqualValuesf(t, 200, resp.Status, "when adding dataManager with status %d and message %s", resp.Status, resp.Message)
	if itemType == "dataManager" {
		return resp, inpDataManager
	}
	// 2. add test dataSample
	inpDataSample := inputDataSample{
		Hashes:          testDataSampleHash1 + ", " + testDataSampleHash2,
		DataManagerKeys: dataManagerOpenerHash,
		TestOnly:        "true",
	}
	args = inpDataSample.createSample()
	resp = mockStub.MockInvoke("42", args)
	require.EqualValuesf(t, 200, resp.Status, "when adding test dataSample with status %d and message %s", resp.Status, resp.Message)
	if itemType == "testDataset" {
		return resp, inpDataSample
	}
	// 3. add objective
	inpObjective := inputObjective{}
	args = inpObjective.createSample()
	resp = mockStub.MockInvoke("42", args)
	require.EqualValuesf(t, 200, resp.Status, "when adding objective with status %d and message %s", resp.Status, resp.Message)
	if itemType == "objective" {
		return resp, inpObjective
	}
	// 4. Add train dataSample
	inpDataSample = inputDataSample{}
	args = inpDataSample.createSample()
	resp = mockStub.MockInvoke("42", args)
	require.EqualValuesf(t, 200, resp.Status, "when adding train dataSample with status %d and message %s", resp.Status, resp.Message)
	if itemType == "trainDataset" {
		return resp, inpDataSample
	}
	// 5. Add algo
	inpAlgo := inputAlgo{}
	args = inpAlgo.createSample()
	resp = mockStub.MockInvoke("42", args)
	require.EqualValuesf(t, 200, resp.Status, "when adding algo with status %d and message %s", resp.Status, resp.Message)
	if itemType == "algo" {
		return resp, inpAlgo
	}
	// 6. Add traintuple
	inpTraintuple := inputTraintuple{}
	args = inpTraintuple.createSample()
	resp = mockStub.MockInvoke("42", args)
	require.EqualValuesf(t, 200, resp.Status, "when adding traintuple with status %d and message %s", resp.Status, resp.Message)
	return resp, inpTraintuple
}

func TestPipeline(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := shim.NewMockStub("substra", scc)
	var out strings.Builder

	fmt.Fprintln(&out, "#### ------------ Add DataManager ------------")
	inpDataManager := inputDataManager{}
	printArgsNames(&out, "registerDataManager", getFieldNames(&inpDataManager))
	args := inpDataManager.createSample()
	printArgs(&out, args, "invoke")
	resp := mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when adding dataManager with status %d and message %s", resp.Status, resp.Message)
	fmt.Fprintf(&out, ">  %s \n\n", string(resp.Payload))
	// Get dataManager key from Payload
	res := map[string]string{}
	err := json.Unmarshal(resp.Payload, &res)
	assert.NoError(t, err, "should unmarshal without problem")
	assert.Contains(t, res, "key")
	dataManagerKey := res["key"]

	fmt.Fprintln(&out, "#### ------------ Query DataManager From key ------------")
	printArgsNames(&out, "queryDataManager", []string{"elementKey"})
	args = [][]byte{[]byte("queryDataManager"), keyToJSON(dataManagerKey)}
	printArgs(&out, args, "queryDataManager")
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when querying a dataManager with status %d and message %s", resp.Status, resp.Message)
	fmt.Fprintf(&out, ">  %s \n\n", string(resp.Payload))

	fmt.Fprintln(&out, "#### ------------ Add test DataSample ------------")
	inpDataSample := inputDataSample{
		Hashes:   testDataSampleHash1 + ", " + testDataSampleHash2,
		TestOnly: "true",
	}
	printArgsNames(&out, "registerDataSample", getFieldNames(&inpDataSample))
	args = inpDataSample.createSample()
	printArgs(&out, args, "invoke")
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when adding test dataSample with status %d and message %s", resp.Status, resp.Message)
	fmt.Fprintf(&out, ">  %s \n\n", string(resp.Payload))

	fmt.Fprintln(&out, "#### ------------ Add Objective ------------")
	inpObjective := inputObjective{}
	printArgsNames(&out, "registerObjective", getFieldNames(&inpObjective))
	args = inpObjective.createSample()
	printArgs(&out, args, "invoke")
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when adding objective with status %d and message %s", resp.Status, resp.Message)
	fmt.Fprintf(&out, ">  %s \n\n", string(resp.Payload))

	fmt.Fprintln(&out, "#### ------------ Add Algo ------------")
	inpAlgo := inputAlgo{}
	printArgsNames(&out, "registerAlgo", getFieldNames(&inpAlgo))
	args = inpAlgo.createSample()
	printArgs(&out, args, "invoke")
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when adding algo with status %d and message %s", resp.Status, resp.Message)
	fmt.Fprintf(&out, ">  %s \n\n", string(resp.Payload))

	fmt.Fprintln(&out, "#### ------------ Add Train DataSample ------------")
	inpDataSample = inputDataSample{}
	printArgsNames(&out, "registerDataSample", getFieldNames(&inpDataSample))
	args = inpDataSample.createSample()
	printArgs(&out, args, "invoke")
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when adding train dataSample with status %d and message %s", resp.Status, resp.Message)
	fmt.Fprintf(&out, ">  %s \n\n", string(resp.Payload))

	fmt.Fprintln(&out, "#### ------------ Query DataManagers ------------")
	args = [][]byte{[]byte("queryDataManagers")}
	printArgs(&out, args, "query")
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when querying dataManager with status %d and message %s", resp.Status, resp.Message)
	fmt.Fprintf(&out, ">  %s \n\n", string(resp.Payload))

	fmt.Fprintln(&out, "#### ------------ Query Objectives ------------")
	args = [][]byte{[]byte("queryObjectives")}
	printArgs(&out, args, "query")
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when querying objective with status %d and message %s", resp.Status, resp.Message)
	fmt.Fprintf(&out, ">  %s \n\n", string(resp.Payload))

	fmt.Fprintln(&out, "#### ------------ Add Traintuple ------------")
	inpTraintuple := inputTraintuple{}
	printArgsNames(&out, "createTraintuple", getFieldNames(&inpTraintuple))
	args = inpTraintuple.createSample()
	printArgs(&out, args, "invoke")
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when adding traintuple with status %d and message %s", resp.Status, resp.Message)
	fmt.Fprintf(&out, ">  %s \n\n", string(resp.Payload))
	// Get traintuple key from Payload
	res = map[string]string{}
	err = json.Unmarshal(resp.Payload, &res)
	assert.NoError(t, err, "should unmarshal without problem")
	assert.Contains(t, res, "key")
	traintupleKey := []byte(res["key"])
	// check not possible to create same traintuple
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 500, resp.Status, "when adding same traintuple with status %d and message %s", resp.Status, resp.Message)
	// Get owner of the traintuple
	args = [][]byte{[]byte("queryTraintuple"), traintupleKey}
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when adding traintuple with status %d and message %s", resp.Status, resp.Message)
	traintuple := outputTraintuple{}
	respTraintuple := resp.Payload
	if err := bytesToStruct(respTraintuple, &traintuple); err != nil {
		t.Errorf("when unmarshalling queried traintuple with error %s", err)
	}
	trainWorker := traintuple.Dataset.Worker

	fmt.Fprintln(&out, "#### ------------ Add Traintuple with inModel from previous traintuple ------------")
	inpTraintuple = inputTraintuple{
		InModels: string(traintupleKey),
	}
	printArgsNames(&out, "createTraintuple", getFieldNames(&inpTraintuple))
	args = inpTraintuple.createSample()
	printArgs(&out, args, "invoke")
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when adding traintuple with status %d and message %s", resp.Status, resp.Message)
	fmt.Fprintf(&out, ">  %s \n\n", string(resp.Payload))
	res = map[string]string{}
	err = json.Unmarshal(resp.Payload, &res)
	assert.NoError(t, err, "should unmarshal without problem")
	assert.Contains(t, res, "key")
	todoTraintupleKey := res["key"]

	fmt.Fprintln(&out, "#### ------------ Query Traintuples of worker with todo status ------------")
	args = [][]byte{[]byte("queryFilter"), []byte("traintuple~worker~status"), []byte(trainWorker + ", todo")}
	printArgs(&out, args, "invoke")
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when querying traintuple of worker with todo status - with status %d and message %s", resp.Status, resp.Message)
	fmt.Fprintf(&out, ">  %s \n\n", string(resp.Payload))

	fmt.Fprintln(&out, "#### ------------ Log Start Training ------------")
	args = [][]byte{[]byte("logStartTrain"), keyToJSON(string(traintupleKey))}
	printArgs(&out, args, "invoke")
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when logging start training with status %d and message %s", resp.Status, resp.Message)
	fmt.Fprintf(&out, ">  %s \n\n", string(resp.Payload))

	fmt.Fprintln(&out, "#### ------------ Log Success Training ------------")
	inp := inputSuccessTrain{Key: string(traintupleKey)}
	args = inp.createSample()
	printArgs(&out, args, "invoke")
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when logging successful training with status %d and message %s", resp.Status, resp.Message)
	fmt.Fprintf(&out, ">  %s \n\n", string(resp.Payload))

	fmt.Fprintln(&out, "#### ------------ Query Traintuple From key ------------")
	args = [][]byte{[]byte("queryTraintuple"), traintupleKey}
	printArgs(&out, args, "queryTraintuple")
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when querying traintuple with status %d and message %s", resp.Status, resp.Message)
	fmt.Fprintf(&out, ">  %s \n\n", string(resp.Payload))

	fmt.Fprintln(&out, "#### ------------ Add Non-Certified Testtuple ------------")
	inpTesttuple := inputTesttuple{
		DataManagerKey: dataManagerOpenerHash,
		DataSampleKeys: trainDataSampleHash1 + ", " + trainDataSampleHash2,
	}
	printArgsNames(&out, "createTesttuple", getFieldNames(&inpTesttuple))
	args = inpTesttuple.createSample()
	printArgs(&out, args, "invoke")
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when adding testtuple with status %d and message %s", resp.Status, resp.Message)
	fmt.Fprintf(&out, ">  %s \n\n", string(resp.Payload))

	fmt.Fprintln(&out, "#### ------------ Add Certified Testtuple ------------")
	inpTesttuple = inputTesttuple{}
	printArgsNames(&out, "createTesttuple", getFieldNames(&inpTesttuple))
	args = inpTesttuple.createSample()
	printArgs(&out, args, "invoke")
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when adding testtuple with status %d and message %s", resp.Status, resp.Message)
	fmt.Fprintf(&out, ">  %s \n\n", string(resp.Payload))
	// Get testtuple key from Payload
	res = map[string]string{}
	err = json.Unmarshal(resp.Payload, &res)
	assert.NoError(t, err, "should unmarshal without problem")
	assert.Contains(t, res, "key")
	testtupleKey := []byte(res["key"])
	// check not possible to create same testtuple
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 500, resp.Status, "when adding same testtuple with status %d and message %s", resp.Status, resp.Message)
	// Get owner of the testtuple
	args = [][]byte{[]byte("queryTesttuple"), testtupleKey}
	resp = mockStub.MockInvoke("42", args)
	respTesttuple := resp.Payload
	testtuple := Testtuple{}
	if err := bytesToStruct(respTesttuple, &testtuple); err != nil {
		t.Errorf("when unmarshalling queried testtuple with error %s", err)
	}
	testWorker := testtuple.Dataset.Worker

	fmt.Fprintln(&out, "#### ------------ Add Testtuple with not done traintuple ------------")
	inpTesttuple = inputTesttuple{
		TraintupleKey: todoTraintupleKey,
	}
	printArgsNames(&out, "createTesttuple", getFieldNames(&inpTesttuple))
	args = inpTesttuple.createSample()
	printArgs(&out, args, "invoke")
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when adding testtuple with status %d and message %s", resp.Status, resp.Message)
	fmt.Fprintf(&out, ">  %s \n\n", string(resp.Payload))

	fmt.Fprintln(&out, "#### ------------ Query Testtuples of worker with todo status ------------")
	args = [][]byte{[]byte("queryFilter"), []byte("testtuple~worker~status"), []byte(testWorker + ", todo")}
	printArgs(&out, args, "invoke")
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when querying testtuple of worker with todo status - with status %d and message %s", resp.Status, resp.Message)
	fmt.Fprintf(&out, ">  %s \n\n", string(resp.Payload))

	fmt.Fprintln(&out, "#### ------------ Log Start Testing ------------")
	args = [][]byte{[]byte("logStartTest"), keyToJSON(string(testtupleKey))}
	printArgs(&out, args, "invoke")
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when logging start testing with status %d and message %s", resp.Status, resp.Message)
	fmt.Fprintf(&out, ">  %s \n\n", string(resp.Payload))

	fmt.Fprintln(&out, "#### ------------ Log Success Testing ------------")
	perf := "0.89"
	log := "still no error, suprah ah ah"
	args = [][]byte{[]byte("logSuccessTest"), testtupleKey, []byte(perf), []byte(log)}
	printArgs(&out, args, "invoke")
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when logging successful training with status %d and message %s", resp.Status, resp.Message)
	fmt.Fprintf(&out, ">  %s \n\n", string(resp.Payload))

	fmt.Fprintln(&out, "#### ------------ Query Testtuple from its key ------------")
	args = [][]byte{[]byte("queryTesttuple"), testtupleKey}
	printArgs(&out, args, "queryTesttuple")
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when querying testtuple with status %d and message %s", resp.Status, resp.Message)
	fmt.Fprintf(&out, ">  %s \n\n", string(resp.Payload))

	fmt.Fprintln(&out, "#### ------------ Query all Testtuples ------------")
	args = [][]byte{[]byte("queryTesttuples")}
	printArgs(&out, args, "queryTesttuples")
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when querying testtuple with status %d and message %s", resp.Status, resp.Message)
	fmt.Fprintf(&out, ">  %s \n\n", string(resp.Payload))

	fmt.Fprintln(&out, "#### ------------ Query details about a model ------------")
	args = [][]byte{[]byte("queryModelDetails"), []byte(traintupleKey)}
	printArgs(&out, args, "query")
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when querying model details with status %d and message %s", resp.Status, resp.Message)
	fmt.Fprintf(&out, ">  %s \n\n", string(resp.Payload))

	fmt.Fprintln(&out, "#### ------------ Query all models ------------")
	args = [][]byte{[]byte("queryModels")}
	printArgs(&out, args, "query")
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when querying model tuples with status %d and message %s", resp.Status, resp.Message)
	fmt.Fprintf(&out, ">  %s \n\n", string(resp.Payload))

	fmt.Fprintln(&out, "#### ------------ Query Dataset ------------")
	args = [][]byte{[]byte("queryDataset"), keyToJSON(dataManagerOpenerHash)}
	printArgs(&out, args, "query")
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when querying dataset with status %d and message %s", resp.Status, resp.Message)
	fmt.Fprintf(&out, ">  %s \n\n", string(resp.Payload))

	// 3. add new data manager and dataSample
	fmt.Fprintln(&out, "#### ------------ Update Data Sample with new data manager ------------")
	newDataManagerKey := "38a320b2a67c8003cc748d6666534f2b01f3f08d175440537a5bf86b7d08d5ee"
	inpDataManager = inputDataManager{OpenerHash: newDataManagerKey}
	args = inpDataManager.createSample()
	mockStub.MockInvoke("42", args)
	// associate a data sample with the old data manager with the updateDataSample
	args = [][]byte{[]byte("updateDataSample"), []byte(trainDataSampleHash1), []byte(newDataManagerKey)}
	printArgs(&out, args, "invoke")
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when updating data sample with new data manager with status %d and message %s", resp.Status, resp.Message)
	fmt.Fprintf(&out, ">  %s \n\n", string(resp.Payload))

	fmt.Fprintln(&out, "#### ------------ Query the new Dataset ------------")
	args = [][]byte{[]byte("queryDataset"), []byte(newDataManagerKey)}
	printArgs(&out, args, "query")
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when querying dataset with status %d and message %s", resp.Status, resp.Message)
	fmt.Fprintf(&out, ">  %s \n\n", string(resp.Payload))

	// Use the output to check the README file and if asked update it
	doc := out.String()
	fromFile, err := ioutil.ReadFile(*readme)
	require.NoErrorf(t, err, "can not read the readme file at the path %s", *readme)
	actualReadme := string(fromFile)
	exampleTitle := "### Examples\n\n"
	index := strings.Index(actualReadme, exampleTitle)
	require.NotEqual(t, -1, index, "README file does not include a Examples section")
	if *update {
		err = ioutil.WriteFile(*readme, []byte(actualReadme[:index+len(exampleTitle)]+doc), 0644)
		assert.NoError(t, err)
	} else {
		testUsage := "The Readme examples are not up to date with the tests"
		testUsage += "\n`-pipeline` to see the output"
		testUsage += "\n`-update` to update the README"
		testUsage += "\n`-readme` to set a different path for the README"
		assert.True(t, strings.Contains(actualReadme, doc), testUsage)
	}
	if *pipeline {
		fmt.Println(doc, index)
	}
}

func TestMain(m *testing.M) {
	//Raise log level to silence it during tests
	logger.SetLevel(shim.LogCritical)
	os.Exit(m.Run())
}
