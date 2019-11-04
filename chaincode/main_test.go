// Copyright 2018 Owkin, inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"reflect"
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
const compositeAlgoHash = "fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcd"
const compositeAlgoStorageAddress = "https://toto/compositealgo/222/algo"
const compositeAlgoName = "hog + svm composite"
const modelHash = "eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed"
const modelAddress = "https://substrabac/model/toto"
const worker = "SampleOrg"
const traintupleKey = "9da043ddc233996d2e62c196471290de4726fc59d65dbbd2b32a920326e8adf3"
const compositeTraintupleKey = "43849d898536c772bd385f431893d05f0f05a3ffd73852e332a5be22b30fc8f2"

var (
	pipeline = flag.Bool("pipeline", false, "Print out the pipeline test output")
	readme   = flag.String("readme", "../README.md", "Pass the path to the README and compare it to the output")
	update   = flag.Bool("update", false, "Update the README.md file")
)

func TestInit(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)

	// resp := mockStub.MockInit("42", [][]byte{[]byte("init")})
	resp := mockStub.MockInit("42", [][]byte{[]byte("init")})
	assert.EqualValuesf(t, 200, resp.Status, "init failed with status %d and message %s", resp.Status, resp.Message)
}

func methodToByte(methodName string) [][]byte {
	return [][]byte{[]byte(methodName)}
}

func methodAndAssetToByte(methodName string, asset interface{}) [][]byte {
	return [][]byte{[]byte(methodName), assetToJSON(asset)}
}

func methodAndKeyToByte(key string, asset interface{}) [][]byte {
	return [][]byte{[]byte("queryAlgo"), keyToJSON(key)}
}

func assetToArgs(asset interface{}) []string {
	return []string{string(assetToJSON(asset))}
}

func assetToJSON(asset interface{}) []byte {
	assetjson, _ := json.Marshal(asset)
	return assetjson
}

func keyToJSON(key string) []byte {
	return assetToJSON(inputHash{Key: key})
}
func printArgs(buf io.Writer, args [][]byte, command string) {
	fmt.Fprintln(buf, "##### Command peer example:")
	fmt.Fprintf(buf, "```bash\npeer chaincode %s -n mycc -c '{\"Args\":[\"%s\"", command, args[0])
	if len(args) == 2 {
		escapedJSON, _ := json.Marshal(string(args[1]))
		fmt.Fprintf(buf, ",%s", escapedJSON)
	}
	fmt.Fprint(buf, "]}' -C myc\n```\n")
}

func prettyPrintStruct(buf io.Writer, margin string, strucType reflect.Type) {
	fmt.Fprintln(buf, "{")
	prettyPrintStructElements(buf, margin+" ", strucType)
	fmt.Fprint(buf, "}")
}
func prettyPrintStructElements(buf io.Writer, margin string, strucType reflect.Type) {
	for i := 0; i < strucType.NumField(); i++ {
		f := strucType.Field(i)
		fieldType := f.Type.Kind()
		fieldStr := ""
		switch fieldType {
		case reflect.Struct:
			if f.Anonymous {
				prettyPrintStructElements(buf, margin, f.Type)
			} else {
				fmt.Fprintf(buf, "%s\"%s\": (%s)", margin, f.Tag.Get("json"), f.Tag.Get("validate"))
				prettyPrintStruct(buf, margin+" ", f.Type)
				fmt.Fprint(buf, ",\n")
			}
			continue
		case reflect.Bool:
			jsonTag := strings.Split(f.Tag.Get("json"), ",")
			fmt.Fprintf(buf, "%s\"%s\": %s (%s),\n", margin, jsonTag[0], fieldType, jsonTag[1])
			continue
		case reflect.Slice:
			if f.Type.Elem().Kind() == reflect.Struct {
				fmt.Fprintf(buf, "%s\"%s\": (%s) [", margin, f.Tag.Get("json"), f.Tag.Get("validate"))
				prettyPrintStruct(buf, margin+" ", f.Type.Elem())
				fmt.Fprint(buf, "],\n")
				continue
			}
			fieldStr = fmt.Sprintf("[%s]", f.Type.Elem().Kind())
		default:
			fieldStr = fmt.Sprint(fieldType)
		}
		fmt.Fprintf(buf, "%s\"%s\": %s (%s),\n", margin, f.Tag.Get("json"), fieldStr, f.Tag.Get("validate"))
	}
	l := len(margin) - 2
	if l > 0 {
		fmt.Fprint(buf, margin[:l])
	}
}

func printInputStuct(buf io.Writer, fnName string, inputType reflect.Type) {
	fmt.Fprintf(buf, "Smart contract: `%s`\n\n##### JSON Inputs:\n```go\n", fnName) // ", fnName)
	prettyPrintStruct(buf, "", inputType)
	fmt.Fprint(buf, "\n```\n")
}

func registerItem(t *testing.T, mockStub MockStub, itemType string) (peer.Response, interface{}) {
	// 1. add dataManager
	inpDataManager := inputDataManager{}
	args := inpDataManager.createDefault()
	resp := mockStub.MockInvoke("42", args)
	require.EqualValuesf(t, 200, resp.Status, "when adding dataManager with status %d and message %s", resp.Status, resp.Message)
	if itemType == "dataManager" {
		return resp, inpDataManager
	}
	// 2. add test dataSample
	inpDataSample := inputDataSample{
		Hashes:          []string{testDataSampleHash1, testDataSampleHash2},
		DataManagerKeys: []string{dataManagerOpenerHash},
		TestOnly:        "true",
	}
	args = inpDataSample.createDefault()
	resp = mockStub.MockInvoke("42", args)
	require.EqualValuesf(t, 200, resp.Status, "when adding test dataSample with status %d and message %s", resp.Status, resp.Message)
	if itemType == "testDataset" {
		return resp, inpDataSample
	}
	// 3. add objective
	inpObjective := inputObjective{}
	args = inpObjective.createDefault()
	resp = mockStub.MockInvoke("42", args)
	require.EqualValuesf(t, 200, resp.Status, "when adding objective with status %d and message %s", resp.Status, resp.Message)
	if itemType == "objective" {
		return resp, inpObjective
	}
	// 4. Add train dataSample
	inpDataSample = inputDataSample{}
	args = inpDataSample.createDefault()
	resp = mockStub.MockInvoke("42", args)
	require.EqualValuesf(t, 200, resp.Status, "when adding train dataSample with status %d and message %s", resp.Status, resp.Message)
	if itemType == "trainDataset" {
		return resp, inpDataSample
	}
	// 5. Add algo
	inpAlgo := inputAlgo{}
	args = inpAlgo.createDefault()
	resp = mockStub.MockInvoke("42", args)
	require.EqualValuesf(t, 200, resp.Status, "when adding algo with status %d and message %s", resp.Status, resp.Message)
	if itemType == "algo" {
		return resp, inpAlgo
	}
	// 6. Add composite algo
	inpCompositeAlgo := inputCompositeAlgo{}
	args = inpCompositeAlgo.createDefault()
	resp = mockStub.MockInvoke("42", args)
	require.EqualValuesf(t, 200, resp.Status, "when adding algo with status %d and message %s", resp.Status, resp.Message)
	if itemType == "compositealgo" {
		return resp, inpCompositeAlgo
	}
	// 7. Add traintuple
	inpTraintuple := inputTraintuple{}
	args = inpTraintuple.createDefault()
	resp = mockStub.MockInvoke("42", args)
	require.EqualValuesf(t, 200, resp.Status, "when adding traintuple with status %d and message %s", resp.Status, resp.Message)
	return resp, inpTraintuple
}

func printResp(buf io.Writer, payload []byte) {
	var toPrint []byte
	if strings.HasPrefix(string(payload), "{") {
		obj := map[string]interface{}{}
		json.Unmarshal(payload, &obj)
		toPrint, _ = json.MarshalIndent(obj, "", " ")
	} else if strings.HasPrefix(string(payload), "[") {
		obj := []map[string]interface{}{}
		json.Unmarshal(payload, &obj)
		toPrint, _ = json.MarshalIndent(obj, "", " ")
	} else {
		toPrint = payload
	}
	fmt.Fprintf(buf, "##### Command output:\n```json\n%s\n```\n", toPrint)
}

func TestPipeline(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStub("substra", scc)
	var out strings.Builder
	callAssertAndPrint := func(peerCmd, smartContract string, inputAsset interface{}) peer.Response {
		var args [][]byte
		if inputAsset != nil {
			printInputStuct(&out, smartContract, reflect.TypeOf(inputAsset))
			args = methodAndAssetToByte(smartContract, inputAsset)
		} else {
			args = methodToByte(smartContract)
		}
		printArgs(&out, args, peerCmd)
		resp := mockStub.MockInvoke("42", args)
		require.EqualValuesf(t, 200, resp.Status, "problem when calling %s, return status %d and message %s", smartContract, resp.Status, resp.Message)
		printResp(&out, resp.Payload)
		return resp
	}

	fmt.Fprintln(&out, "#### ------------ Add Node ------------")
	callAssertAndPrint("invoke", "registerNode", nil)

	fmt.Fprintln(&out, "#### ------------ Add DataManager ------------")
	inpDataManager := inputDataManager{}
	inpDataManager.createDefault()
	resp := callAssertAndPrint("invoke", "registerDataManager", inpDataManager)
	// Get dataManager key from Payload
	res := map[string]string{}
	err := json.Unmarshal(resp.Payload, &res)
	assert.NoError(t, err, "should unmarshal without problem")
	assert.Contains(t, res, "key")
	dataManagerKey := res["key"]

	fmt.Fprintln(&out, "#### ------------ Query DataManager From key ------------")
	callAssertAndPrint("invoke", "queryDataManager", inputHash{dataManagerKey})

	fmt.Fprintln(&out, "#### ------------ Add test DataSample ------------")
	inpDataSample := inputDataSample{
		Hashes:   []string{testDataSampleHash1, testDataSampleHash2},
		TestOnly: "true",
	}
	inpDataSample.createDefault()
	callAssertAndPrint("invoke", "registerDataSample", inpDataSample)

	fmt.Fprintln(&out, "#### ------------ Add Objective ------------")
	inpObjective := inputObjective{}
	inpObjective.createDefault()
	callAssertAndPrint("invoke", "registerObjective", inpObjective)

	fmt.Fprintln(&out, "#### ------------ Add Algo ------------")
	inpAlgo := inputAlgo{}
	inpAlgo.createDefault()
	callAssertAndPrint("invoke", "registerAlgo", inpAlgo)

	fmt.Fprintln(&out, "#### ------------ Add Train DataSample ------------")
	inpDataSample = inputDataSample{}
	inpDataSample.createDefault()
	callAssertAndPrint("invoke", "registerDataSample", inpDataSample)

	fmt.Fprintln(&out, "#### ------------ Query DataManagers ------------")
	callAssertAndPrint("query", "queryDataManagers", nil)

	fmt.Fprintln(&out, "#### ------------ Query DataSamples ------------")
	callAssertAndPrint("query", "queryDataSamples", nil)

	fmt.Fprintln(&out, "#### ------------ Query Objectives ------------")
	callAssertAndPrint("query", "queryObjectives", nil)

	fmt.Fprintln(&out, "#### ------------ Add Traintuple ------------")
	inpTraintuple := inputTraintuple{}
	args := inpTraintuple.createDefault()
	resp = callAssertAndPrint("invoke", "createTraintuple", inpTraintuple)
	// Get traintuple key from Payload
	res = map[string]string{}
	err = json.Unmarshal(resp.Payload, &res)
	assert.NoError(t, err, "should unmarshal without problem")
	assert.Contains(t, res, "key")
	traintupleKey := res["key"]
	// check not possible to create same traintuple
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 409, resp.Status, "when adding same traintuple with status %d and message %s", resp.Status, resp.Message)
	// Get owner of the traintuple
	args = [][]byte{[]byte("queryTraintuple"), keyToJSON(traintupleKey)}
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when adding traintuple with status %d and message %s", resp.Status, resp.Message)
	traintuple := outputTraintuple{}
	respTraintuple := resp.Payload
	if err := bytesToStruct(respTraintuple, &traintuple); err != nil {
		t.Errorf("when unmarshalling queried traintuple with error %s", err)
	}
	trainWorker := traintuple.Dataset.Worker

	fmt.Fprintln(&out, "#### ------------ Add Traintuple with inModel from previous traintuple ------------")
	inpTraintuple = inputTraintuple{}
	inpTraintuple.InModels = []string{traintupleKey}
	inpTraintuple.createDefault()
	resp = callAssertAndPrint("invoke", "createTraintuple", inpTraintuple)
	printResp(&out, resp.Payload)
	res = map[string]string{}
	err = json.Unmarshal(resp.Payload, &res)
	assert.NoError(t, err, "should unmarshal without problem")
	assert.Contains(t, res, "key")
	todoTraintupleKey := res["key"]

	fmt.Fprintln(&out, "#### ------------ Query Traintuples of worker with todo status ------------")
	filter := inputQueryFilter{
		IndexName:  "traintuple~worker~status",
		Attributes: trainWorker + ", todo",
	}
	callAssertAndPrint("invoke", "queryFilter", filter)

	fmt.Fprintln(&out, "#### ------------ Log Start Training ------------")
	callAssertAndPrint("invoke", "logStartTrain", inputHash{traintupleKey})

	fmt.Fprintln(&out, "#### ------------ Log Success Training ------------")
	inp := inputLogSuccessTrain{}
	inp.Key = string(traintupleKey)
	inp.createDefault()
	callAssertAndPrint("invoke", "logSuccessTrain", inp)

	fmt.Fprintln(&out, "#### ------------ Query Traintuple From key ------------")
	callAssertAndPrint("invoke", "queryTraintuple", inputHash{traintupleKey})

	fmt.Fprintln(&out, "#### ------------ Add Non-Certified Testtuple ------------")
	inpTesttuple := inputTesttuple{
		DataManagerKey: dataManagerOpenerHash,
		DataSampleKeys: []string{trainDataSampleHash1, trainDataSampleHash2},
	}
	inpTesttuple.createDefault()
	callAssertAndPrint("invoke", "createTesttuple", inpTesttuple)

	fmt.Fprintln(&out, "#### ------------ Add Certified Testtuple ------------")
	inpTesttuple = inputTesttuple{}
	args = inpTesttuple.createDefault()
	resp = callAssertAndPrint("invoke", "createTesttuple", inpTesttuple)
	// Get testtuple key from Payload
	res = map[string]string{}
	err = json.Unmarshal(resp.Payload, &res)
	assert.NoError(t, err, "should unmarshal without problem")
	assert.Contains(t, res, "key")
	testtupleKey := res["key"]
	// check not possible to create same testtuple
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 409, resp.Status, "when adding same testtuple with status %d and message %s", resp.Status, resp.Message)
	// Get owner of the testtuple
	args = [][]byte{[]byte("queryTesttuple"), keyToJSON(testtupleKey)}
	resp = mockStub.MockInvoke("42", args)
	respTesttuple := resp.Payload
	testtuple := outputTesttuple{}
	if err := bytesToStruct(respTesttuple, &testtuple); err != nil {
		t.Errorf("when unmarshalling queried testtuple with error %s", err)
	}
	testWorker := testtuple.Dataset.Worker

	fmt.Fprintln(&out, "#### ------------ Add Testtuple with not done traintuple ------------")
	inpTesttuple = inputTesttuple{}
	inpTesttuple.TraintupleKey = todoTraintupleKey
	inpTesttuple.createDefault()
	callAssertAndPrint("invoke", "createTesttuple", inpTesttuple)

	fmt.Fprintln(&out, "#### ------------ Query Testtuples of worker with todo status ------------")
	filter = inputQueryFilter{
		IndexName:  "testtuple~worker~status",
		Attributes: testWorker + ", todo",
	}
	callAssertAndPrint("invoke", "queryFilter", filter)

	fmt.Fprintln(&out, "#### ------------ Log Start Testing ------------")
	callAssertAndPrint("invoke", "logStartTest", inputHash{testtupleKey})

	fmt.Fprintln(&out, "#### ------------ Log Success Testing ------------")
	success := inputLogSuccessTest{}
	success.Key = testtupleKey
	args = success.createDefault()
	callAssertAndPrint("invoke", "logSuccessTest", success)

	fmt.Fprintln(&out, "#### ------------ Query Testtuple from its key ------------")
	callAssertAndPrint("query", "queryTesttuple", inputHash{testtupleKey})

	fmt.Fprintln(&out, "#### ------------ Query all Testtuples ------------")
	callAssertAndPrint("query", "queryTesttuples", nil)

	fmt.Fprintln(&out, "#### ------------ Query details about a model ------------")
	callAssertAndPrint("query", "queryModelDetails", inputHash{traintupleKey})

	fmt.Fprintln(&out, "#### ------------ Query all models ------------")
	callAssertAndPrint("query", "queryModels", nil)

	fmt.Fprintln(&out, "#### ------------ Query Dataset ------------")
	callAssertAndPrint("query", "queryDataset", inputHash{dataManagerOpenerHash})

	fmt.Fprintln(&out, "#### ------------ Query nodes ------------")
	callAssertAndPrint("query", "queryNodes", nil)

	// 3. add new data manager and dataSample
	fmt.Fprintln(&out, "#### ------------ Update Data Sample with new data manager ------------")
	newDataManagerKey := "38a320b2a67c8003cc748d6666534f2b01f3f08d175440537a5bf86b7d08d5ee"
	inpDataManager = inputDataManager{OpenerHash: newDataManagerKey}
	args = inpDataManager.createDefault()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when adding dataManager with status %d and message %s", resp.Status, resp.Message)
	// associate a data sample with the old data manager with the updateDataSample
	updateData := inputUpdateDataSample{
		DataManagerKeys: []string{newDataManagerKey},
		Hashes:          []string{trainDataSampleHash1},
	}
	callAssertAndPrint("invoke", "updateDataSample", updateData)

	fmt.Fprintln(&out, "#### ------------ Query the new Dataset ------------")
	callAssertAndPrint("query", "queryDataset", inputHash{newDataManagerKey})

	fmt.Fprintln(&out, "#### ------------ Create a ComputePlan ------------")
	callAssertAndPrint("invoke", "createComputePlan", defaultComputePlan)

	fmt.Fprintln(&out, "#### ------------ Query an ObjectiveLeaderboard ------------")
	inpLeaderboard := inputLeaderboard{
		ObjectiveKey:   objectiveDescriptionHash,
		AscendingOrder: true,
	}
	callAssertAndPrint("invoke", "queryObjectiveLeaderboard", inpLeaderboard)

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

func initializeMockStateDB(t *testing.T, stub *MockStub) {
	stub.MockTransactionStart("42")
	stub.PutState("key", []byte("value"))
}

func TestQueryEmptyResponse(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)
	initializeMockStateDB(t, mockStub)

	smartContracts := []string{
		"queryAlgos",
		"queryDataSamples",
		"queryObjectives",
		"queryDataManagers",
		"queryTraintuples",
		"queryTesttuples",
		"queryModels",
	}

	for _, contractName := range smartContracts {
		t.Run(contractName, func(t *testing.T) {
			args := [][]byte{[]byte(contractName)}
			resp := mockStub.MockInvoke("42", args)

			expectedPayload, _ := json.Marshal(make([]string, 0))
			assert.Equal(t, expectedPayload, resp.Payload, "payload is not an empty list")
		})
	}
}
