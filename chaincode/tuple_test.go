package main

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func TestNoPanicWhileQueryingIncompleteTraintuple(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := shim.NewMockStub("substra", scc)
	// Add a some dataManager, dataSample and traintuple
	err, _, _ := registerItem(*mockStub, "traintuple")
	assert.NoError(t, err)

	// Manually open a ledger transaction
	mockStub.MockTransactionStart("42")
	defer mockStub.MockTransactionEnd("42")

	// Retreive and alter existing objectif to pass Metrics at nil
	objective := Objective{}
	getElementStruct(mockStub, objectiveDescriptionHash, &objective)
	assert.NoError(t, err)
	objective.Metrics = nil
	objBytes, err := json.Marshal(objective)
	assert.NoError(t, err)
	err = mockStub.PutState(objectiveDescriptionHash, objBytes)
	assert.NoError(t, err)

	// It should not panic
	require.NotPanics(t, func() {
		getOutputTraintuple(mockStub, traintupleKey)
	})
}
func TestTraintupleFLtaskCreation(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := shim.NewMockStub("substra", scc)

	// Add dataManager, dataSample and algo
	err, resp, _ := registerItem(*mockStub, "algo")
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	inpTraintuple := inputTraintuple{FLtask: "someFLtask"}
	args := inpTraintuple.createSample()
	resp = mockStub.MockInvoke("42", args)
	require.EqualValues(t, 500, resp.Status, "should failed for missing rank")
	require.Contains(t, resp.Message, "invalit inputs, a FLtask should have a rank", "invalid error message")

	inpTraintuple = inputTraintuple{Rank: "1"}
	args = inpTraintuple.createSample()
	resp = mockStub.MockInvoke("42", args)
	require.EqualValues(t, 500, resp.Status, "should failed for invalid rank")
	require.Contains(t, resp.Message, "invalid inputs, a new FLtask should have a rank 0")

	inpTraintuple = inputTraintuple{Rank: "0"}
	args = inpTraintuple.createSample()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status)
	key := string(resp.Payload)
	require.EqualValues(t, key, traintupleKey)

	inpTraintuple = inputTraintuple{Rank: "0"}
	args = inpTraintuple.createSample()
	resp = mockStub.MockInvoke("42", args)
	require.EqualValues(t, 500, resp.Status, "should failed for existing FLtask")
	require.Contains(t, resp.Message, "traintuple with these algo, in models, and train dataSample already exist")
}

func TestTraintupleMultipleFLtaskCreations(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := shim.NewMockStub("substra", scc)

	// Add a some dataManager, dataSample and traintuple
	err, resp, _ := registerItem(*mockStub, "algo")
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	inpTraintuple := inputTraintuple{Rank: "0"}
	args := inpTraintuple.createSample()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status)
	key := string(resp.Payload)

	// Failed to add a traintuple with the same rank
	inpTraintuple = inputTraintuple{
		InModels: key,
		Rank:     "0",
		FLtask:   key}
	args = inpTraintuple.createSample()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 500, resp.Status, resp.Message, "should failed to add a traintuple of the same rank")

	// Failed to add a traintuple to an unexisting Fltask
	inpTraintuple = inputTraintuple{
		InModels: key,
		Rank:     "1",
		FLtask:   "notarealone"}
	args = inpTraintuple.createSample()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 500, resp.Status, resp.Message, "should failed to add a traintuple to an unexisting FLtask")

	// Succesfully add a traintuple to the same FLtask
	inpTraintuple = inputTraintuple{
		InModels: key,
		Rank:     "1",
		FLtask:   key}
	args = inpTraintuple.createSample()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status, resp.Message, "should be able do create a traintuple with the same FLtask")
	ttkey := string(resp.Payload)

	// Add new algo to check all fltask algo consistency
	newAlgoHash := strings.Replace(algoHash, "a", "b", 1)
	inpAlgo := inputAlgo{Hash: newAlgoHash}
	args = inpAlgo.createSample()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status)

	inpTraintuple = inputTraintuple{
		AlgoKey:  newAlgoHash,
		InModels: ttkey,
		Rank:     "2",
		FLtask:   key}
	args = inpTraintuple.createSample()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 500, resp.Status, resp.Message, "sould fail for it doesn't have the same algo key")
	assert.Contains(t, resp.Message, "does not have the same algo key")
}

func TestTesttupleOnFailedTraintuple(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := shim.NewMockStub("substra", scc)

	// Add a some dataManager, dataSample and traintuple
	err, resp, _ := registerItem(*mockStub, "traintuple")
	assert.NoError(t, err)
	traintupleKey := resp.Payload

	// Mark the traintuple as failed
	args := [][]byte{[]byte("logFailTrain"), traintupleKey, []byte("pas glop")}
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status, "should be able to log traintuple as failed")

	// Fail to add a testtuple to this failed traintuple
	inpTesttuple := inputTesttuple{}
	args = inpTesttuple.createSample()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 500, resp.Status, "status should show an error since the traintuple is failed")
	assert.Contains(t, resp.Message, "could not register this testtuple")
}

func TestCertifiedExplicitTesttuple(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := shim.NewMockStub("substra", scc)

	// Add a some dataManager, dataSample and traintuple
	err, resp, _ := registerItem(*mockStub, "traintuple")
	assert.NoError(t, err)

	// Add a testtuple that shoulb be certified since it's the same dataManager and
	// dataSample than the objective but explicitly pass as arguments and in disorder
	inpTesttuple := inputTesttuple{
		DataSampleKeys: testDataSampleHash2 + "," + testDataSampleHash1,
		DataManagerKey: dataManagerOpenerHash}
	args := inpTesttuple.createSample()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status)

	args = [][]byte{[]byte("queryTesttuples")}
	resp = mockStub.MockInvoke("42", args)
	testtuples := [](map[string]interface{}){}
	err = json.Unmarshal(resp.Payload, &testtuples)
	assert.NoError(t, err, "should be unmarshaled")
	assert.Len(t, testtuples, 1, "there should be only one testtuple...")
	assert.True(t, testtuples[0]["certified"].(bool), "... and it should be certified")

}
func TestConflictCertifiedNonCertifiedTesttuple(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := shim.NewMockStub("substra", scc)

	// Add a some dataManager, dataSample and traintuple
	err, resp, _ := registerItem(*mockStub, "traintuple")
	assert.NoError(t, err)

	// Add a certified testtuple
	inpTesttuple1 := inputTesttuple{}
	args := inpTesttuple1.createSample()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status)

	// Fail to add an incomplete uncertified testtuple
	inpTesttuple2 := inputTesttuple{DataSampleKeys: trainDataSampleHash1}
	args = inpTesttuple2.createSample()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 500, resp.Status)
	assert.Contains(t, resp.Message, "invalid input: dataManagerKey and dataSampleKey should be provided together")

	// Add an uncertified testtuple successfully
	inpTesttuple3 := inputTesttuple{
		DataSampleKeys: trainDataSampleHash1 + "," + trainDataSampleHash2,
		DataManagerKey: dataManagerOpenerHash}
	args = inpTesttuple3.createSample()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status)

	// Fail to add the same testtuple with a different order for dataSampleKeys
	inpTesttuple4 := inputTesttuple{
		DataSampleKeys: trainDataSampleHash2 + "," + trainDataSampleHash1,
		DataManagerKey: dataManagerOpenerHash}
	args = inpTesttuple4.createSample()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 500, resp.Status)
	assert.Contains(t, resp.Message, "this testtuple already exists")
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
		t.Errorf("when adding objective with invalid hash, status %d and message %s", status, resp.Message)
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
	args = [][]byte{[]byte("queryTraintuple"), []byte(traintupleKey)}
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
	if a, b := traintuple["objective"].(map[string]interface{})["hash"], objectiveDescriptionHash; a != b {
		t.Errorf("ledger traintuple objective hash does not corresponds to what was input: %s - %s", a, b)
	}
	if a, b := traintuple["objective"].(map[string]interface{})["metrics"].(map[string]interface{})["hash"], objectiveMetricsHash; a != b {
		t.Errorf("ledger traintuple objective hash does not corresponds to what was input: %s - %s", a, b)
	}
	if a, b := traintuple["objective"].(map[string]interface{})["metrics"].(map[string]interface{})["storageAddress"], objectiveMetricsStorageAddress; a != b {
		t.Errorf("ledger traintuple objective hash does not corresponds to what was input: %s - %s", a, b)
	}
	algo := make(map[string]interface{})
	algo["hash"] = algoHash
	algo["storageAddress"] = algoStorageAddress
	algo["name"] = algoName
	if a, b := traintuple["algo"], algo; !reflect.DeepEqual(a, b) {
		t.Errorf("wrong ledger traintuple algo: %s - %s", a, b)
	}
	dataset := make(map[string]interface{})
	dataset["worker"] = worker
	dataset["keys"] = []interface{}{trainDataSampleHash1, trainDataSampleHash2}
	dataset["openerHash"] = dataManagerOpenerHash
	dataset["perf"] = 0.0
	if a, b := traintuple["dataset"], dataset; !reflect.DeepEqual(a, b) {
		t.Errorf("wrong ledger traintuple train dataset: %s - %s", a, b)
	}

	// Query all traintuples and check consistency
	args = [][]byte{[]byte("queryTraintuples")}
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("when querying traintuples - status %d and message %s", status, resp.Message)
	}
	var sPayload []map[string]interface{}
	if err := json.Unmarshal(resp.Payload, &sPayload); err != nil {
		t.Errorf("when unmarshalling queried objectives")
	}
	payload := sPayload[0]
	delete(payload, "key")
	if !reflect.DeepEqual(payload, traintuple) {
		t.Errorf("when querying traintuples, does not correspond to what was input")
	}

	// Add traintuple with inmodel from the above-submitted traintuple
	inpWaitingTraintuple := inputTraintuple{
		InModels: string(traintupleKey),
	}
	args = inpWaitingTraintuple.createSample()
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("when adding traintuple with status %d and message %s", status, resp.Message)
	}
	//waitingTraintupleKey := string(resp.Payload)

	// Query traintuple with status todo and worker as trainworker and check consistency
	args = [][]byte{[]byte("queryFilter"), []byte("traintuple~worker~status"), []byte(worker + ", todo")}
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
	perf := "0.9"
	log := "no error, ah ah ah"
	argsSlice := [][][]byte{
		[][]byte{[]byte("logStartTrain"), []byte(traintupleKey)},
		[][]byte{[]byte("logSuccessTrain"), []byte(traintupleKey), []byte(modelHash + ", " + modelAddress),
			[]byte(perf), []byte(log)},
	}
	traintupleStatus := []string{"doing", "done"}
	for i := range traintupleStatus {
		resp = mockStub.MockInvoke("42", argsSlice[i])
		if status := resp.Status; status != 200 {
			t.Errorf("when logging start %s with status %d and message %s",
				traintupleStatus[i], status, resp.Message)
		}
		args = [][]byte{[]byte("queryFilter"), []byte("traintuple~worker~status"), []byte(worker + ", " + traintupleStatus[i])}
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
	args = [][]byte{[]byte("queryTraintuple"), []byte(traintupleKey)}
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
	if a := endTraintuple.Log; a != log {
		t.Errorf("because retrieved log in traintuple does not correspond to what "+
			"was submitted: %s", a)
	}
	outModel := HashDress{
		Hash:           modelHash,
		StorageAddress: modelAddress}
	if endTraintuple.OutModel.Hash != outModel.Hash || endTraintuple.OutModel.StorageAddress != outModel.StorageAddress {
		t.Errorf("because retrieved endModel in traintuple does not correspond to what "+
			"was submitted: %s, %s", endTraintuple.OutModel, outModel)
	}

	// query all traintuples related to a traintuple with the same algo
	args = [][]byte{[]byte("queryModelDetails"), []byte(traintupleKey)}
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("when querying model details with status %d and message %s", status, resp.Message)
	}
	payload = make(map[string]interface{})
	json.Unmarshal(resp.Payload, &payload)
	if l := len(payload); l != 1 {
		t.Errorf("when querying model tuples, payload should contain at this stage only one traintuple, but it contains %d elements", l)
	}

	// query all traintuples related to a traintuple with the same algo
	args = [][]byte{[]byte("queryModels")}
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("when querying models with status %d and message %s", status, resp.Message)
	}
	// var mPayload []map[string]interface{}
	// json.Unmarshal(resp.Payload, &mPayload)
	// fmt.Println(mPayload)
}

/**
}

/**
func TestTesttuple(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := shim.NewMockStub("substra", scc)

	// Add traintuple with invalid field
	inpTraintuple := inputTraintuple{
		AlgoKey: "aaa",
	}
	args := inpTesttuple.createSample()
	resp := mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 500 {
		t.Errorf("when adding objective with invalid hash, status %d and message %s", status, resp.Message)
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
	if a, b := traintuple["objective"].(map[string]interface{})["hash"], objectiveDescriptionHash; a != b {
		t.Errorf("ledger traintuple objective hash does not corresponds to what was input: %s - %s", a, b)
	}
	if a, b := traintuple["objective"].(map[string]interface{})["metrics"].(map[string]interface{})["hash"], objectiveMetricsHash; a != b {
		t.Errorf("ledger traintuple objective hash does not corresponds to what was input: %s - %s", a, b)
	}
	if a, b := traintuple["objective"].(map[string]interface{})["metrics"].(map[string]interface{})["storageAddress"], objectiveMetricsStorageAddress; a != b {
		t.Errorf("ledger traintuple objective hash does not corresponds to what was input: %s - %s", a, b)
	}
	algo := make(map[string]interface{})
	algo["hash"] = algoHash
	algo["storageAddress"] = algoStorageAddress
	algo["name"] = algoName
	if a, b := traintuple["algo"], algo; !reflect.DeepEqual(a, b) {
		t.Errorf("wrong ledger traintuple algo: %s - %s", a, b)
	}
	dataset := make(map[string]interface{})
	dataset["worker"] = worker
	dataset["keys"] = []interface{}{trainDataSampleHash1, trainDataSampleHash2}
	dataset["openerHash"] = dataManagerOpenerHash
	dataset["perf"] = 0.0
	if a, b := traintuple["dataset"], dataset; !reflect.DeepEqual(a, b) {
		t.Errorf("wrong ledger traintuple train dataset: %s - %s", a, b)
	}

	// Query all traintuples and check consistency
	args = [][]byte{[]byte("queryTraintuples")}
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("when querying traintuples - status %d and message %s", status, resp.Message)
	}
	var sPayload []map[string]interface{}
	if err := json.Unmarshal(resp.Payload, &sPayload); err != nil {
		t.Errorf("when unmarshalling queried objectives")
	}
	payload := sPayload[0]
	delete(payload, "key")
	if !reflect.DeepEqual(payload, traintuple) {
		t.Errorf("when querying objectives, dataManager does not correspond to the input objective")
	}

	// Query traintuple with status todo and worker as trainworker and check consistency
	args = [][]byte{[]byte("queryFilter"), []byte("traintuple~worker~status"), []byte(worker + ", todo")}
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
	perf := "0.9"
	log := "no error, ah ah ah"
	argsSlice := [][][]byte{
		[][]byte{[]byte("logStartTrain"), []byte(traintupleKey)},
		[][]byte{[]byte("logSuccessTrain"), []byte(traintupleKey), []byte(modelHash + ", " + modelAddress),
			[]byte(perf), []byte(log)},
	}
	traintupleStatus := []string{"doing", "done"}
	for i, _ := range traintupleStatus {
		resp = mockStub.MockInvoke("42", argsSlice[i])
		if status := resp.Status; status != 200 {
			t.Errorf("when logging start %s with status %d and message %s",
				traintupleStatus[i], status, resp.Message)
		}
		args = [][]byte{[]byte("queryFilter"), []byte("traintuple~worker~status"), []byte(worker + ", " + traintupleStatus[i])}
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
	if a := endTraintuple.Log; a != log {
		t.Errorf("because retrieved log in traintuple does not correspond to what "+
			"was submitted: %s", a)
	}
	outModel := HashDress{
		Hash:           modelHash,
		StorageAddress: modelAddress}
	if endTraintuple.OutModel.Hash != outModel.Hash || endTraintuple.OutModel.StorageAddress != outModel.StorageAddress {
		t.Errorf("because retrieved endModel in traintuple does not correspond to what "+
			"was submitted: %s, %s", endTraintuple.OutModel, outModel)
	}

	// query all traintuples related to a traintuple with the same algo
	args = [][]byte{[]byte("queryTraintuplesAlgo"), []byte(traintupleKey)}
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("when querying algo traintuples with status %d and message %s", status, resp.Message)
	}
	payload = make(map[string]interface{})
	json.Unmarshal(resp.Payload, &payload)
	if l := len(payload); l != 2 {
		t.Errorf("when querying algo traintuples, payload should contain an algo "+"and a traintuple, but it contains %d elements", l)
	}
}
**/
