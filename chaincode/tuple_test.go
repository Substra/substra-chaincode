package main

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func TestTesttupleOnFailedTraintuple(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := shim.NewMockStub("substra", scc)

	// Add a some dataManager, data and traintuple
	err, resp, _ := registerItem(*mockStub, "traintuple")
	assert.NoError(t, err)
	// Mark the traintuple as failed
	args := [][]byte{[]byte("logFailTrain"), resp.Payload, []byte("pas glop")}
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

	// Add a some dataManager, data and traintuple
	err, resp, _ := registerItem(*mockStub, "traintuple")
	assert.NoError(t, err)

	// Add a testtuple that shoulb be certified since it's the same dataManager and
	// data than the objective but explicitly pass as arguments and in disorder
	inpTesttuple := inputTesttuple{
		DataKeys:   testDataHash2 + "," + testDataHash1,
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

	// Add a some dataManager, data and traintuple
	err, resp, _ := registerItem(*mockStub, "traintuple")
	assert.NoError(t, err)

	// Add a certified testtuple
	inpTesttuple1 := inputTesttuple{}
	args := inpTesttuple1.createSample()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status)

	// Fail to add an incomplete uncertified testtuple
	inpTesttuple2 := inputTesttuple{DataKeys: trainDataHash1}
	args = inpTesttuple2.createSample()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 500, resp.Status)
	assert.Contains(t, resp.Message, "invalid input: dataManagerKey and dataKey should be provided together")

	// Add an uncertified testtuple successfully
	inpTesttuple3 := inputTesttuple{
		DataKeys:   trainDataHash1 + "," + trainDataHash2,
		DataManagerKey: dataManagerOpenerHash}
	args = inpTesttuple3.createSample()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status)

	// Fail to add the same testtuple with a different order for dataKeys
	inpTesttuple4 := inputTesttuple{
		DataKeys:   trainDataHash2 + "," + trainDataHash1,
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
	data := make(map[string]interface{})
	data["worker"] = worker
	data["keys"] = []interface{}{trainDataHash1, trainDataHash2}
	data["openerHash"] = dataManagerOpenerHash
	data["perf"] = 0.0
	if a, b := traintuple["data"], data; !reflect.DeepEqual(a, b) {
		t.Errorf("wrong ledger traintuple train data: %s - %s", a, b)
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
	data := make(map[string]interface{})
	data["worker"] = worker
	data["keys"] = []interface{}{trainDataHash1, trainDataHash2}
	data["openerHash"] = dataManagerOpenerHash
	data["perf"] = 0.0
	if a, b := traintuple["data"], data; !reflect.DeepEqual(a, b) {
		t.Errorf("wrong ledger traintuple train data: %s - %s", a, b)
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
