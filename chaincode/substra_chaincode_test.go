package main

import (
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

func TestAddProblem(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := shim.NewMockStub("substra", scc)

	fn := []byte("addProblem")
	descriptionHash := []byte("5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379")
	name := []byte("msi classification")
	descriptionStorageAddress := []byte("https://toto/problem/222/description")
	metricsStorageAddress := []byte("https://toto/problem/222/metrics")
	metricsHash := []byte("fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482d8d")
	testData := []byte("data_fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482d8a, data_fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482d8e")
	permissions := []byte("all")
	args := [][]byte{fn, descriptionHash, name, descriptionStorageAddress,
		metricsStorageAddress, metricsHash, testData, permissions}

	resp := mockStub.MockInvoke("42", args)
	status := resp.Status
	if status != 200 {
		t.Errorf("addProblem failed with status %d", status)
	}
	// TODO ADD CHECK resp.Message
}

func TestAddAlgo(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := shim.NewMockStub("substra", scc)

	fn := []byte("addAlgo")
	algoHash := []byte("fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc")
	name := []byte("hog + svm")
	storageAddress := []byte("https://toto/algo/222/algo")
	associatedProblem := []byte("problem_5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379")
	permissions := []byte("all")
	args := [][]byte{fn, algoHash, name, storageAddress, associatedProblem, permissions}

	resp := mockStub.MockInvoke("42", args)
	status := resp.Status
	if status != 200 {
		t.Errorf("addAlgo failed with status %d", status)
	}
	// TODO ADD CHECK resp.Message
}

func TestAddData(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := shim.NewMockStub("substra", scc)

	fn := []byte("addData")
	dataHash := []byte("da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc")
	name := []byte("liver slide")
	dataOpenerHash := []byte("do4bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc")
	associatedProblems := []byte("problem_5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379")
	permissions := []byte("all")
	args := [][]byte{fn, dataHash, name, dataOpenerHash, associatedProblems, permissions}

	resp := mockStub.MockInvoke("42", args)
	status := resp.Status
	if status != 200 {
		t.Errorf("addData failed with status %d", status)
	}
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
