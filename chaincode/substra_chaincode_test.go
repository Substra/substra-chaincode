package main

import (
	"fmt"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type SubmittedData struct {
	Hash        string
	Name        string
	OpenerHash  string
	ProblemKeys string
	TestOnly    string
}

type SubmittedAlgo struct {
	Hash           string
	Name           string
	StorageAddress string
	ProblemKey     string
}

type SubmittedProblem struct {
	DescriptionHash           string
	Name                      string
	DescriptionStorageAddress string
	MetricsStorageAddress     string
	MetricsHash               string
	TestDataKeys              string
}

type SubmittedTrainTuple struct {
	ProblemKey    string
	StartModelKey string
	TrainDataKeys string
}

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

// func createSampleData(dataHash, name, dataOpenerHash, associatedProblems, testOnly string) (args [][]byte) {
func createSampleData(data SubmittedData) (args [][]byte) {
	if data.Hash == "" {
		data.Hash = "aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
	}
	if data.Name == "" {
		data.Name = "liver slide"
	}
	if data.OpenerHash == "" {
		data.OpenerHash = "do4bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
	}
	if data.TestOnly == "" {
		data.TestOnly = "false"
	}
	bFn := []byte("addData")
	bDataHash := []byte(data.Hash)
	bName := []byte(data.Name)
	bDataOpenerHash := []byte(data.OpenerHash)
	bAssociatedProblems := []byte(data.ProblemKeys)
	bTestOnly := []byte(data.TestOnly)
	bPermissions := []byte("all")
	args = [][]byte{bFn, bDataHash, bName, bDataOpenerHash, bAssociatedProblems, bTestOnly, bPermissions}
	return args
}

func createSampleProblem(problem SubmittedProblem) (args [][]byte) {
	if problem.DescriptionHash == "" {
		problem.DescriptionHash = "5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379"
	}
	if problem.Name == "" {
		problem.Name = "MSI classification"
	}
	if problem.DescriptionStorageAddress == "" {
		problem.DescriptionStorageAddress = "https://toto/problem/222/description"
	}
	if problem.MetricsStorageAddress == "" {
		problem.MetricsStorageAddress = "https://toto/problem/222/metrics"
	}
	if problem.MetricsHash == "" {
		problem.MetricsHash = "fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482d8d"
	}
	if problem.TestDataKeys == "" {
		problem.TestDataKeys = "data_da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
	}
	bFn := []byte("addProblem")
	bDescriptionHash := []byte(problem.DescriptionHash)
	bName := []byte(problem.Name)
	bDescriptionStorageAddress := []byte(problem.DescriptionStorageAddress)
	bMetricsStorageAddress := []byte(problem.MetricsStorageAddress)
	bMetricsHash := []byte(problem.MetricsHash)
	bTestData := []byte(problem.TestDataKeys)
	bPermissions := []byte("all")
	args = [][]byte{bFn, bDescriptionHash, bName, bDescriptionStorageAddress,
		bMetricsStorageAddress, bMetricsHash, bTestData, bPermissions}
	return args
}

func createSampleAlgo(algo SubmittedAlgo) (args [][]byte) {
	if algo.Hash == "" {
		algo.Hash = "fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
	}
	if algo.Name == "" {
		algo.Name = "hog + svm"
	}
	if algo.StorageAddress == "" {
		algo.StorageAddress = "https://toto/algo/222/algo"
	}
	if algo.ProblemKey == "" {
		algo.ProblemKey = "problem_5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379"
	}
	bFn := []byte("addAlgo")
	bAlgoHash := []byte(algo.Hash)
	bName := []byte(algo.Name)
	bStorageAddress := []byte(algo.StorageAddress)
	bAssociatedProblem := []byte(algo.ProblemKey)
	bPermissions := []byte("all")
	args = [][]byte{bFn, bAlgoHash, bName, bStorageAddress, bAssociatedProblem, bPermissions}
	return args
}

func createSampleTrainTuple(trainTuple SubmittedTrainTuple) (args [][]byte) {
	if trainTuple.ProblemKey == "" {
		trainTuple.ProblemKey = "problem_5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379"
	}
	if trainTuple.StartModelKey == "" {
		trainTuple.StartModelKey = "algo_fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
	}
	if trainTuple.TrainDataKeys == "" {
		trainTuple.TrainDataKeys = "data_aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
	}
	bFn := []byte("addTrainTuple")
	bProblemKey := []byte(trainTuple.ProblemKey)
	bStartModelKey := []byte(trainTuple.StartModelKey)
	bTrainDataKeys := []byte(trainTuple.TrainDataKeys)
	args = [][]byte{bFn, bProblemKey, bStartModelKey, bTrainDataKeys}
	return args
}

func TestAddData(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := shim.NewMockStub("substra", scc)

	submittedData := SubmittedData{}
	args := createSampleData(submittedData)
	resp := mockStub.MockInvoke("42", args)
	status := resp.Status
	fmt.Println(string(resp.Payload))
	if status != 200 {
		t.Errorf("testAddData failed with status %d and message %s", status, resp.Message)
	}
	// TODO check content of payload := resp.Payload
}

func TestPipeline(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := shim.NewMockStub("substra", scc)
	// Add test data
	submittedData := SubmittedData{
		Hash:     "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc",
		TestOnly: "true",
	}
	args := createSampleData(submittedData)
	resp := mockStub.MockInvoke("42", args)
	status := resp.Status
	if status != 200 {
		t.Errorf("testPipeline failed when adding test data with status %d and message %s", status, resp.Message)
	}
	// Add problem
	submittedProblem := SubmittedProblem{}
	args = createSampleProblem(submittedProblem)
	resp = mockStub.MockInvoke("42", args)
	status = resp.Status
	if status != 200 {
		t.Errorf("testPipeline failed when adding problem with status %d and message %s", status, resp.Message)
	}
	// Add algo
	submittedAlgo := SubmittedAlgo{}
	args = createSampleAlgo(submittedAlgo)
	resp = mockStub.MockInvoke("42", args)
	status = resp.Status
	if status != 200 {
		t.Errorf("testPipeline failed when adding algo with status %d and message %s", status, resp.Message)
	}
	// Add train data
	submittedData = SubmittedData{
		ProblemKeys: "problem_5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379",
	}
	args = createSampleData(submittedData)
	resp = mockStub.MockInvoke("42", args)
	status = resp.Status
	if status != 200 {
		t.Errorf("testPipeline failed when adding train data with status %d and message %s", status, resp.Message)
	}
	// Query all data
	resp = mockStub.MockInvoke("42", [][]byte{[]byte("queryData")})
	status = resp.Status
	if status != 200 {
		t.Errorf("testPipeline failed when adding train data with status %d and message %s", status, resp.Message)
	}

	// Add trainTuple
	trainTuple := SubmittedTrainTuple{}
	args = createSampleTrainTuple(trainTuple)
	resp = mockStub.MockInvoke("42", args)
	status = resp.Status
	if status != 200 {
		t.Errorf("testPipeline failed when adding traintuple with status %d and message %s", status, resp.Message)
	}
	// TODO ADD CHECK content of resp.Payload
}
