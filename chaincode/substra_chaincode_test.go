package main

import (
	"fmt"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func TestInit(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := shim.NewMockStub("substra", scc)

	// resp := mockStub.MockInit("42", [][]byte{[]byte("init")})
	resp := mockStub.MockInit("42", [][]byte{[]byte("init")})
	status := resp.Status
	if status != 200 {
		t.Errorf("init failed with status %d", status)
	}
}

func createSampleData(dataHash, name, dataOpenerHash, associatedProblems, testOnly string) (args [][]byte) {
	if dataHash == "" {
		dataHash = "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
	}
	if name == "" {
		name = "liver slide"
	}
	if dataOpenerHash == "" {
		dataOpenerHash = "do4bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
	}
	if testOnly == "" {
		testOnly = "false"
	}
	bFn := []byte("addData")
	bDataHash := []byte(dataHash)
	bName := []byte(name)
	bDataOpenerHash := []byte(dataOpenerHash)
	bAssociatedProblems := []byte(associatedProblems)
	bTestOnly := []byte(testOnly)
	bPermissions := []byte("all")
	args = [][]byte{bFn, bDataHash, bName, bDataOpenerHash, bAssociatedProblems, bTestOnly, bPermissions}
	return args
}

func createSampleProblem(descriptionHash, name, descriptionStorageAddress, metricsStorageAddress, metricsHash, testData string) (args [][]byte) {
	if descriptionHash == "" {
		descriptionHash = "5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379"
	}
	if name == "" {
		name = "MSI classification"
	}
	if descriptionStorageAddress == "" {
		descriptionStorageAddress = "https://toto/problem/222/description"
	}
	if metricsStorageAddress == "" {
		metricsStorageAddress = "https://toto/problem/222/metrics"
	}
	if metricsHash == "" {
		metricsHash = "fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482d8d"
	}
	if testData == "" {
		testData = "data_da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
	}
	bFn := []byte("addProblem")
	bDescriptionHash := []byte(descriptionHash)
	bName := []byte(name)
	bDescriptionStorageAddress := []byte(descriptionStorageAddress)
	bMetricsStorageAddress := []byte(metricsStorageAddress)
	bMetricsHash := []byte(metricsHash)
	bTestData := []byte(testData)
	bPermissions := []byte("all")
	args = [][]byte{bFn, bDescriptionHash, bName, bDescriptionStorageAddress,
		bMetricsStorageAddress, bMetricsHash, bTestData, bPermissions}
	return args
}

func createSampleAlgo(algoHash, name, storageAddress, associatedProblem string) (args [][]byte) {
	if algoHash == "" {
		algoHash = "fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
	}
	if name == "" {
		name = "hog + svm"
	}
	if storageAddress == "" {
		storageAddress = "https://toto/algo/222/algo"
	}
	if associatedProblem == "" {
		associatedProblem = "problem_5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379"
	}
	bFn := []byte("addAlgo")
	bAlgoHash := []byte(algoHash)
	bName := []byte(name)
	bStorageAddress := []byte(storageAddress)
	bAssociatedProblem := []byte(associatedProblem)
	bPermissions := []byte("all")
	args = [][]byte{bFn, bAlgoHash, bName, bStorageAddress, bAssociatedProblem, bPermissions}
	return args
}

func TestAddData(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := shim.NewMockStub("substra", scc)

	args := createSampleData("", "", "", "", "")
	resp := mockStub.MockInvoke("42", args)
	status := resp.Status
	fmt.Println(resp.Message)
	if status != 200 {
		t.Errorf("testAddData failed with status %d", status)
	}
	// TODO ADD CHECK resp.Message
}

func TestAdd(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := shim.NewMockStub("substra", scc)
	// Add test data
	args := createSampleData("", "", "", "", "true")
	resp := mockStub.MockInvoke("42", args)
	status := resp.Status
	if status != 200 {
		t.Errorf("testAdd failed when adding test data with status %d", status)
	}
	// Add problem
	args = createSampleProblem("", "", "", "", "", "")
	resp = mockStub.MockInvoke("42", args)
	fmt.Println(resp.Message)
	status = resp.Status
	if status != 200 {
		t.Errorf("testAdd failed when adding problem with status %d", status)
	}
	// Add algo
	args = createSampleAlgo("", "", "", "")
	resp = mockStub.MockInvoke("42", args)
	fmt.Println(resp.Message)
	status = resp.Status
	if status != 200 {
		t.Errorf("testAdd failed when adding algo with status %d", status)
	}
	// Add train data

	// Add trainTuple
	// TODO ADD CHECK resp.Message
}

func TestAddTrainTuple(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := shim.NewMockStub("substra", scc)

	// Preparation: add problem, data, and algo

	fn := []byte("addTrainTuple")
	problemKey := []byte("problem_5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379")
	startModelKey := []byte("algo_fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc")
	trainData := []byte("data_xx5c1d9cd1c2c1082dde0921b56d110")
	args := [][]byte{fn, problemKey, startModelKey, trainData}

	resp := mockStub.MockInvoke("42", args)
	status := resp.Status
	if status != 200 {
		t.Errorf("addData failed with status %d", status)
	}
	// TODO ADD CHECK resp.Message
}
