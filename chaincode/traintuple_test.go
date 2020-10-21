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
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTraintupleWithNoTestDataset(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)
	registerItem(t, *mockStub, "trainDataset")

	key := strings.Replace(objectiveKey, "1", "2", 1)
	inpObjective := inputObjective{Key: key}
	inpObjective.createDefault()
	inpObjective.TestDataset = inputDataset{}
	resp := mockStub.MockInvoke("42", methodAndAssetToByte("registerObjective", inpObjective))
	assert.EqualValues(t, 200, resp.Status, "when adding objective without dataset it should work: ", resp.Message)

	inpAlgo := inputAlgo{}
	args := inpAlgo.createDefault()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status, "when adding algo it should work: ", resp.Message)

	inpTraintuple := inputTraintuple{}
	args = inpTraintuple.createDefault()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status, "when adding traintuple without test dataset it should work: ", resp.Message)

	args = [][]byte{[]byte("queryTraintuple"), keyToJSONOld(traintupleKey)}
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status, "It should find the traintuple without error ", resp.Message)
}

func TestTraintupleWithSingleDatasample(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)
	registerItem(t, *mockStub, "trainDataset")

	key := strings.Replace(objectiveKey, "1", "2", 1)
	inpObjective := inputObjective{Key: key}
	inpObjective.createDefault()
	inpObjective.TestDataset = inputDataset{}
	resp := mockStub.MockInvoke("42", methodAndAssetToByte("registerObjective", inpObjective))
	assert.EqualValues(t, 200, resp.Status, "when adding objective without dataset it should work: ", resp.Message)

	inpAlgo := inputAlgo{}
	args := inpAlgo.createDefault()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status, "when adding algo it should work: ", resp.Message)

	inpTraintuple := inputTraintuple{
		DataSampleKeys: []string{trainDataSampleKey1},
	}
	args = inpTraintuple.createDefault()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status, "when adding traintuple with a single data samples it should work: ", resp.Message)

	traintuple := outputKey{}
	err := json.Unmarshal(resp.Payload, &traintuple)
	assert.NoError(t, err, "should be unmarshaled")
	args = [][]byte{[]byte("queryTraintuple"), keyToJSONOld(traintuple.Key)}
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status, "It should find the traintuple without error ", resp.Message)
}

func TestTraintupleWithDuplicatedDatasamples(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)
	registerItem(t, *mockStub, "trainDataset")

	key := strings.Replace(objectiveKey, "1", "2", 1)
	inpObjective := inputObjective{Key: key}
	inpObjective.createDefault()
	inpObjective.TestDataset = inputDataset{}
	resp := mockStub.MockInvoke("42", methodAndAssetToByte("registerObjective", inpObjective))
	assert.EqualValues(t, 200, resp.Status, "when adding objective without dataset it should work: ", resp.Message)

	inpAlgo := inputAlgo{}
	args := inpAlgo.createDefault()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status, "when adding algo it should work: ", resp.Message)

	inpTraintuple := inputTraintuple{
		DataSampleKeys: []string{trainDataSampleKey1, trainDataSampleKey2, trainDataSampleKey1},
	}
	args = inpTraintuple.createDefault()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 400, resp.Status, "when adding traintuple with a duplicated data samples it should not work: %s", resp.Message)
}

func TestNoPanicWhileQueryingIncompleteTraintuple(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)
	// Add a some dataManager, dataSample and traintuple
	registerItem(t, *mockStub, "traintuple")

	// Manually open a ledger transaction
	mockStub.MockTransactionStart("42")
	defer mockStub.MockTransactionEnd("42")

	// Retreive and alter existing objectif to pass Metrics at nil
	db := NewLedgerDB(mockStub)
	objective, err := db.GetObjective(objectiveKey)
	assert.NoError(t, err)
	objective.Metrics = nil
	objBytes, err := json.Marshal(objective)
	assert.NoError(t, err)
	err = mockStub.PutState(objectiveKey, objBytes)
	assert.NoError(t, err)
	// It should not panic
	require.NotPanics(t, func() {
		getOutputTraintuple(NewLedgerDB(mockStub), traintupleKey)
	})
}

func TestTraintupleComputePlanCreation(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)

	// Add dataManager, dataSample and algo
	registerItem(t, *mockStub, "algo")

	inpTraintuple := inputTraintuple{ComputePlanID: "someComputePlanID"}
	args := inpTraintuple.createDefault()
	resp := mockStub.MockInvoke("42", args)
	require.EqualValues(t, 400, resp.Status, "should failed for missing rank")
	require.Contains(t, resp.Message, "invalid inputs, a ComputePlan should have a rank", "invalid error message")

	inpTraintuple = inputTraintuple{Rank: "1"}
	args = inpTraintuple.createDefault()
	resp = mockStub.MockInvoke("42", args)
	require.EqualValues(t, 400, resp.Status, "should failed for invalid rank")
	require.Contains(t, resp.Message, "invalid inputs, a new ComputePlan should have a rank 0")

	inpTraintuple = inputTraintuple{Rank: "0"}
	args = inpTraintuple.createDefault()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status)
	res := outputKey{}
	err := json.Unmarshal(resp.Payload, &res)
	assert.NoError(t, err, "should unmarshal without problem")
	key := res.Key
	require.EqualValues(t, key, traintupleKey)

	inpTraintuple = inputTraintuple{Rank: "0"}
	args = inpTraintuple.createDefault()
	resp = mockStub.MockInvoke("42", args)
	require.EqualValues(t, 409, resp.Status, "should failed for existing ComputePlanID")
	require.Contains(t, resp.Message, "already exists")

	require.EqualValues(t, 409, resp.Status, "should failed for existing FLTask")
	errorPayload := map[string]interface{}{}
	err = json.Unmarshal(resp.Payload, &errorPayload)
	assert.NoError(t, err, "should unmarshal without problem")
	require.Contains(t, errorPayload, "key", "key should be available in payload")
	assert.EqualValues(t, traintupleKey, errorPayload["key"], "key in error should be traintupleKey")
}

func TestTraintupleMultipleCommputePlanCreations(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)

	// Add a some dataManager, dataSample and traintuple
	registerItem(t, *mockStub, "algo")

	inpTraintuple := inputTraintuple{Rank: "0"}
	args := inpTraintuple.createDefault()
	resp := mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status)
	res := outputKey{}
	err := json.Unmarshal(resp.Payload, &res)
	assert.NoError(t, err, "should unmarshal without problem")
	key := res.Key
	db := NewLedgerDB(mockStub)
	tuple, err := db.GetTraintuple(key)
	assert.NoError(t, err)
	// Failed to add a traintuple with the same rank
	inpTraintuple = inputTraintuple{
		InModels:      []string{key},
		Rank:          "0",
		ComputePlanID: tuple.ComputePlanID}
	args = inpTraintuple.createDefault()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 400, resp.Status, resp.Message, "should failed to add a traintuple of the same rank")

	// Failed to add a traintuple to an unexisting CommputePlan
	inpTraintuple = inputTraintuple{
		InModels:      []string{key},
		Rank:          "1",
		ComputePlanID: "notarealone"}
	args = inpTraintuple.createDefault()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 404, resp.Status, resp.Message, "should failed to add a traintuple to an unexisting ComputePlanID")

	// Succesfully add a traintuple to the same ComputePlanID
	inpTraintuple = inputTraintuple{
		InModels:      []string{key},
		Rank:          "1",
		ComputePlanID: tuple.ComputePlanID}
	args = inpTraintuple.createDefault()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status, resp.Message, "should be able do create a traintuple with the same ComputePlanID")
	err = json.Unmarshal(resp.Payload, &res)
	assert.NoError(t, err, "should unmarshal without problem")
	ttkey := res.Key
	// Add new algo to check all ComputePlan algo consistency
	newAlgoHash := strings.Replace(algoHash, "a", "b", 1)
	inpAlgo := inputAlgo{Hash: newAlgoHash}
	args = inpAlgo.createDefault()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status)

	inpTraintuple = inputTraintuple{
		AlgoKey:       newAlgoHash,
		InModels:      []string{ttkey},
		Rank:          "2",
		ComputePlanID: tuple.ComputePlanID}
	args = inpTraintuple.createDefault()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status, resp.Message, "should be able to create a traintuple with the same ComputePlanID and different algo keys")
}

func TestTraintuple(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)

	// Add traintuple with invalid field
	inpTraintuple := inputTraintuple{
		AlgoKey: "aaa",
	}
	args := inpTraintuple.createDefault()
	resp := mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 400, resp.Status, "when adding objective with invalid hash, status %d and message %s", resp.Status, resp.Message)

	// Add traintuple with unexisting algo
	inpTraintuple = inputTraintuple{}
	args = inpTraintuple.createDefault()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 400, resp.Status, "when adding traintuple with unexisting algo, status %d and message %s", resp.Status, resp.Message)

	// Properly add traintuple
	resp, tt := registerItem(t, *mockStub, "traintuple")

	inpTraintuple = tt.(inputTraintuple)
	res := outputKey{}
	err := json.Unmarshal(resp.Payload, &res)
	assert.NoError(t, err, "traintuple should unmarshal without problem")
	traintupleKey := res.Key
	// Query traintuple from key and check the consistency of returned arguments
	args = [][]byte{[]byte("queryTraintuple"), keyToJSONOld(traintupleKey)}
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when querying the traintuple - status %d and message %s", resp.Status, resp.Message)
	out := outputTraintuple{}
	err = json.Unmarshal(resp.Payload, &out)
	assert.NoError(t, err, "when unmarshalling queried traintuple")
	expected := outputTraintuple{
		Key: traintupleKey,
		Algo: &HashDressName{
			Hash:           algoHash,
			Name:           algoName,
			StorageAddress: algoStorageAddress,
		},
		Creator: worker,
		Dataset: &outputTtDataset{
			Key:            dataManagerKey,
			DataSampleKeys: []string{trainDataSampleKey1, trainDataSampleKey2},
			OpenerHash:     dataManagerOpenerHash,
			Worker:         worker,
			Metadata:       map[string]string{},
		},
		Permissions: outputPermissions{
			Process: Permission{Public: true, AuthorizedIDs: []string{}},
		},
		Metadata: map[string]string{},
		Status:   StatusTodo,
	}
	assert.Exactly(t, expected, out, "the traintuple queried from the ledger differ from expected")

	// Query all traintuples and check consistency
	args = [][]byte{[]byte("queryTraintuples")}
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when querying traintuples - status %d and message %s", resp.Status, resp.Message)
	// TODO add traintuple key to output struct
	// For now we test it as cleanly as its added to the query response
	assert.Contains(t, string(resp.Payload), "key\":\""+traintupleKey)
	var queryTraintuples []outputTraintuple
	err = json.Unmarshal(resp.Payload, &queryTraintuples)
	assert.NoError(t, err, "traintuples should unmarshal without problem")
	assert.Exactly(t, out, queryTraintuples[0])

	// Add traintuple with inmodel from the above-submitted traintuple
	inpWaitingTraintuple := inputTraintuple{
		InModels: []string{string(traintupleKey)},
	}
	args = inpWaitingTraintuple.createDefault()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when adding traintuple with status %d and message %s", resp.Status, resp.Message)
	//waitingTraintupleKey := string(resp.Payload)

	// Query traintuple with status todo and worker as trainworker and check consistency
	filter := inputQueryFilter{
		IndexName:  "traintuple~worker~status",
		Attributes: worker + ", todo",
	}
	args = [][]byte{[]byte("queryFilter"), assetToJSON(filter)}
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when querying traintuple of worker with todo status - status %d and message %s", resp.Status, resp.Message)
	err = json.Unmarshal(resp.Payload, &queryTraintuples)
	assert.NoError(t, err, "traintuples should unmarshal without problem")
	assert.Exactly(t, out, queryTraintuples[0])

	// Update status and check consistency
	success := inputLogSuccessTrain{}
	success.Key = traintupleKey

	argsSlice := [][][]byte{
		[][]byte{[]byte("logStartTrain"), keyToJSONOld(traintupleKey)},
		success.createDefault(),
	}
	traintupleStatus := []string{StatusDoing, StatusDone}
	for i := range traintupleStatus {
		resp = mockStub.MockInvoke("42", argsSlice[i])
		require.EqualValuesf(t, 200, resp.Status, "when logging start %s with message %s", traintupleStatus[i], resp.Message)
		filter := inputQueryFilter{
			IndexName:  "traintuple~worker~status",
			Attributes: worker + ", " + traintupleStatus[i],
		}
		args = [][]byte{[]byte("queryFilter"), assetToJSON(filter)}
		resp = mockStub.MockInvoke("42", args)
		assert.EqualValuesf(t, 200, resp.Status, "when querying traintuple of worker with %s status - message %s", traintupleStatus[i], resp.Message)
		sPayload := make([]map[string]interface{}, 1)
		assert.NoError(t, json.Unmarshal(resp.Payload, &sPayload), "when unmarshal queried traintuples")
		assert.EqualValues(t, traintupleKey, sPayload[0]["key"], "wrong retrieved key when querying traintuple of worker with %s status ", traintupleStatus[i])
		assert.EqualValues(t, traintupleStatus[i], sPayload[0]["status"], "wrong retrieved status when querying traintuple of worker with %s status ", traintupleStatus[i])
	}

	// Query Traintuple From key
	args = [][]byte{[]byte("queryTraintuple"), keyToJSONOld(traintupleKey)}
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when querying traintuple with status %d and message %s", resp.Status, resp.Message)
	endTraintuple := outputTraintuple{}
	assert.NoError(t, json.Unmarshal(resp.Payload, &endTraintuple))
	expected.Log = success.Log
	expected.OutModel = &HashDress{
		Hash:           modelHash,
		StorageAddress: modelAddress}
	expected.Status = traintupleStatus[1]
	assert.Exactly(t, expected, endTraintuple, "retreived Traintuple does not correspond to what is expected")

	// query all traintuples related to a traintuple with the same algo
	args = [][]byte{[]byte("queryModelDetails"), keyToJSONOld(traintupleKey)}
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when querying model details with status %d and message %s", resp.Status, resp.Message)
	payload := outputModelDetails{}
	assert.NoError(t, json.Unmarshal(resp.Payload, &payload))
	assert.NotNil(t, payload.Traintuple, "when querying model tuples, payload should contain one traintuple")

	// query all traintuples related to a traintuple with the same algo
	args = [][]byte{[]byte("queryModels")}
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when querying models with status %d and message %s", resp.Status, resp.Message)
}

func TestQueryTraintupleNotFound(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)
	registerItem(t, *mockStub, "traintuple")

	// queryTraintuple: normal case
	args := [][]byte{[]byte("queryTraintuple"), keyToJSONOld(traintupleKey)}
	resp := mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when querying the traintuple - status %d and message %s", resp.Status, resp.Message)

	// queryTraintuple: key does not exist
	notFoundKey := "eedbb7c31f62244c0f34461cc168804227115793d01c270021fe3f7935482eed"
	args = [][]byte{[]byte("queryTraintuple"), keyToJSONOld(notFoundKey)}
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 404, resp.Status, "when querying the traintuple - status %d and message %s", resp.Status, resp.Message)

	// queryTraintuple: key does not exist and use existing other asset type key
	args = [][]byte{[]byte("queryTraintuple"), keyToJSONOld(algoHash)}
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 404, resp.Status, "when querying the traintuple - status %d and message %s", resp.Status, resp.Message)
}

func TestInsertTraintupleTwice(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)
	registerItem(t, *mockStub, "algo")

	// create a traintuple and start a ComplutePlan
	inpTraintuple := inputTraintuple{
		Rank: "0",
	}
	inpTraintuple.createDefault()
	resp := mockStub.MockInvoke("42", methodAndAssetToByte("createTraintuple", inpTraintuple))
	assert.EqualValues(t, http.StatusOK, resp.Status)

	db := NewLedgerDB(mockStub)
	tuple, err := db.GetTraintuple(traintupleKey)
	assert.NoError(t, err)
	// create a second traintuple in the same ComputePlan
	inpTraintuple.Rank = "1"
	inpTraintuple.ComputePlanID = tuple.ComputePlanID
	inpTraintuple.InModels = []string{traintupleKey}
	resp = mockStub.MockInvoke("42", methodAndAssetToByte("createTraintuple", inpTraintuple))
	assert.EqualValues(t, http.StatusOK, resp.Status)

	// re-insert the same traintuple and expect a conflict error
	resp = mockStub.MockInvoke("42", methodAndAssetToByte("createTraintuple", inpTraintuple))
	assert.EqualValues(t, http.StatusConflict, resp.Status)

}
