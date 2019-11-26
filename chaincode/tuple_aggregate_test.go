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

/////////////////////////////////////////////////////////////
//
// "Regular" tests
// Copied from `traintuple_test.go` and adapted for aggregate
//
/////////////////////////////////////////////////////////////

func TestTraintupleWithNoTestDatasetAggregate(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)
	registerItem(t, *mockStub, "trainDataset")

	objHash := strings.ReplaceAll(objectiveDescriptionHash, "1", "2")
	inpObjective := inputObjective{DescriptionHash: objHash}
	inpObjective.createDefault()
	inpObjective.TestDataset = inputDataset{}
	resp := mockStub.MockInvoke("42", methodAndAssetToByte("registerObjective", inpObjective))
	assert.EqualValues(t, 200, resp.Status, "when adding objective without dataset it should work: ", resp.Message)

	inpAlgo := inputAggregateAlgo{}
	args := inpAlgo.createDefault()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status, "when adding aggregate algo it should work: ", resp.Message)

	inpTraintuple := inputAggregateTuple{ObjectiveKey: objHash}
	args = inpTraintuple.createDefault()
	resp = mockStub.MockInvoke("42", args)

	assert.EqualValues(t, 200, resp.Status, "when adding aggregate tuple without test dataset it should work: ", resp.Message)

	traintuple := outputAggregateTuple{}
	json.Unmarshal(resp.Payload, &traintuple)
	args = [][]byte{[]byte("queryAggregatetuple"), keyToJSON(traintuple.Key)}
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status, "It should find the aggregate tuple without error ", resp.Message)
}

func TestTraintupleWithSingleDatasampleAggregate(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)
	registerItem(t, *mockStub, "trainDataset")

	objHash := strings.ReplaceAll(objectiveDescriptionHash, "1", "2")
	inpObjective := inputObjective{DescriptionHash: objHash}
	inpObjective.createDefault()
	inpObjective.TestDataset = inputDataset{}
	resp := mockStub.MockInvoke("42", methodAndAssetToByte("registerObjective", inpObjective))
	assert.EqualValues(t, 200, resp.Status, "when adding objective without dataset it should work: ", resp.Message)

	inpAlgo := inputAggregateAlgo{}
	args := inpAlgo.createDefault()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status, "when adding aggregate algo it should work: ", resp.Message)

	inpTraintuple := inputAggregateTuple{
		ObjectiveKey: objHash,
		AlgoKey:      aggregateAlgoHash,
	}
	args = inpTraintuple.createDefault()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status, "when adding aggregate tuple with a single data samples it should work: ", resp.Message)

	traintuple := outputAggregateTuple{}
	err := json.Unmarshal(resp.Payload, &traintuple)
	assert.NoError(t, err, "should be unmarshaled")
	args = [][]byte{[]byte("queryAggregatetuple"), keyToJSON(traintuple.Key)}
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status, "It should find the aggregate tuple without error ", resp.Message)
}

func TestNoPanicWhileQueryingIncompleteTraintupleAggregate(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)
	// Add a some dataManager, dataSample and traintuple
	registerItem(t, *mockStub, "aggregateTuple")

	// Manually open a ledger transaction
	mockStub.MockTransactionStart("42")
	defer mockStub.MockTransactionEnd("42")

	// Retreive and alter existing objectif to pass Metrics at nil
	db := NewLedgerDB(mockStub)
	objective, err := db.GetObjective(objectiveDescriptionHash)
	assert.NoError(t, err)
	objective.Metrics = nil
	objBytes, err := json.Marshal(objective)
	assert.NoError(t, err)
	err = mockStub.PutState(objectiveDescriptionHash, objBytes)
	assert.NoError(t, err)
	// It should not panic
	require.NotPanics(t, func() {
		getOutputAggregateTuple(NewLedgerDB(mockStub), traintupleKey)
	})
}

func TestTraintupleComputePlanCreationAggregate(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)

	// Add dataManager, dataSample and algo
	registerItem(t, *mockStub, "aggregateAlgo")

	inpTraintuple := inputAggregateTuple{ComputePlanID: "someComputePlanID"}
	args := inpTraintuple.createDefault()
	resp := mockStub.MockInvoke("42", args)
	require.EqualValues(t, 400, resp.Status, "should failed for missing rank")
	require.Contains(t, resp.Message, "invalid inputs, a ComputePlan should have a rank", "invalid error message")

	inpTraintuple = inputAggregateTuple{Rank: "1"}
	args = inpTraintuple.createDefault()
	resp = mockStub.MockInvoke("42", args)
	require.EqualValues(t, 400, resp.Status, "should failed for invalid rank")
	require.Contains(t, resp.Message, "invalid inputs, a new ComputePlan should have a rank 0")

	inpTraintuple = inputAggregateTuple{Rank: "0"}
	args = inpTraintuple.createDefault()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status)
	res := map[string]string{}
	err := json.Unmarshal(resp.Payload, &res)
	assert.NoError(t, err, "should unmarshal without problem")
	assert.Contains(t, res, "key")
	key := res["key"]
	require.EqualValues(t, aggregateTupleKey, key)

	inpTraintuple = inputAggregateTuple{Rank: "0"}
	args = inpTraintuple.createDefault()
	resp = mockStub.MockInvoke("42", args)
	require.EqualValues(t, 409, resp.Status, "should failed for existing ComputePlanID")
	require.Contains(t, resp.Message, "already exists")

	require.EqualValues(t, 409, resp.Status, "should failed for existing FLTask")
	errorPayload := map[string]interface{}{}
	err = json.Unmarshal(resp.Payload, &errorPayload)
	assert.NoError(t, err, "should unmarshal without problem")
	require.Contains(t, errorPayload, "key", "key should be available in payload")
	assert.EqualValues(t, aggregateTupleKey, errorPayload["key"], "key in error should be aggregateTupleKey")
}

func TestTraintupleMultipleCommputePlanCreationsAggregate(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)

	// Add a some dataManager, dataSample and traintuple
	registerItem(t, *mockStub, "aggregateAlgo")

	inpTraintuple := inputAggregateTuple{Rank: "0"}
	args := inpTraintuple.createDefault()
	resp := mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status)
	res := map[string]string{}
	err := json.Unmarshal(resp.Payload, &res)
	assert.NoError(t, err, "should unmarshal without problem")
	assert.Contains(t, res, "key")
	key := res["key"]
	// Failed to add a traintuple with the same rank
	inpTraintuple = inputAggregateTuple{
		InModels:      []string{key},
		Rank:          "0",
		ComputePlanID: key}
	args = inpTraintuple.createDefault()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 400, resp.Status, resp.Message, "should failed to add an aggregate tuple of the same rank")

	// Failed to add a traintuple to an unexisting CommputePlan
	inpTraintuple = inputAggregateTuple{
		InModels:      []string{key},
		Rank:          "1",
		ComputePlanID: "notarealone"}
	args = inpTraintuple.createDefault()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 400, resp.Status, resp.Message, "should failed to add an aggregate tuple to an unexisting ComputePlanID")

	// Succesfully add a traintuple to the same ComputePlanID
	inpTraintuple = inputAggregateTuple{
		InModels:      []string{key},
		Rank:          "1",
		ComputePlanID: key}
	args = inpTraintuple.createDefault()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status, resp.Message, "should be able do create an aggregate tuple with the same ComputePlanID")
	err = json.Unmarshal(resp.Payload, &res)
	assert.NoError(t, err, "should unmarshal without problem")
	assert.Contains(t, res, "key")
	ttkey := res["key"]
	// Add new algo to check all ComputePlan algo consistency
	newAlgoHash := strings.Replace(aggregateAlgoHash, "a", "b", 1)
	inpAlgo := inputAggregateAlgo{inputAlgo{Hash: newAlgoHash}}
	args = inpAlgo.createDefault()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status)

	inpTraintuple = inputAggregateTuple{
		AlgoKey:       newAlgoHash,
		InModels:      []string{ttkey},
		Rank:          "2",
		ComputePlanID: key}
	args = inpTraintuple.createDefault()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 400, resp.Status, resp.Message, "should fail for it doesn't have the same aggregate algo key")
	assert.Contains(t, resp.Message, "does not have the same algo key")
}

func TestTraintupleAggregate(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)

	// Add traintuple with invalid field
	inpTraintuple := inputAggregateTuple{
		AlgoKey: "aaa",
	}
	args := inpTraintuple.createDefault()
	resp := mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 400, resp.Status, "when adding objective with invalid hash, status %d and message %s", resp.Status, resp.Message)

	// Add traintuple with unexisting algo
	inpTraintuple = inputAggregateTuple{}
	args = inpTraintuple.createDefault()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 400, resp.Status, "when adding aggregate tuple with unexisting algo, status %d and message %s", resp.Status, resp.Message)

	// Properly add traintuple
	resp, tt := registerItem(t, *mockStub, "aggregateTuple")

	inpTraintuple = tt.(inputAggregateTuple)
	res := map[string]string{}
	err := json.Unmarshal(resp.Payload, &res)
	assert.NoError(t, err, "aggregate tuple should unmarshal without problem")
	assert.Contains(t, res, "key")
	traintupleKey := res["key"]
	// Query traintuple from key and check the consistency of returned arguments
	args = [][]byte{[]byte("queryAggregatetuple"), keyToJSON(traintupleKey)}
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when querying the aggregate tuple - status %d and message %s", resp.Status, resp.Message)
	out := outputAggregateTuple{}
	err = json.Unmarshal(resp.Payload, &out)
	assert.NoError(t, err, "when unmarshalling queried aggregate tuple")
	expected := outputAggregateTuple{
		Key: aggregateTupleKey,
		Algo: &HashDressName{
			Hash:           aggregateAlgoHash,
			Name:           aggregateAlgoName,
			StorageAddress: aggregateAlgoStorageAddress,
		},
		Creator: worker,
		Worker:  worker,
		Objective: &TtObjective{
			Key: objectiveDescriptionHash,
			Metrics: &HashDress{
				Hash:           objectiveMetricsHash,
				StorageAddress: objectiveMetricsStorageAddress,
			},
		},
		Status: StatusTodo,
	}
	assert.Exactly(t, expected, out, "the aggregate tuple queried from the ledger differ from expected")

	// Query all traintuples and check consistency
	args = [][]byte{[]byte("queryAggregatetuples")}
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when querying aggregate tuples - status %d and message %s", resp.Status, resp.Message)
	// TODO add traintuple key to output struct
	// For now we test it as cleanly as its added to the query response
	assert.Contains(t, string(resp.Payload), "key\":\""+aggregateTupleKey)
	var queryTraintuples []outputAggregateTuple
	err = json.Unmarshal(resp.Payload, &queryTraintuples)
	assert.NoError(t, err, "aggregate tuples should unmarshal without problem")
	require.NotZero(t, queryTraintuples)
	assert.Exactly(t, out, queryTraintuples[0])

	// Add traintuple with inmodel from the above-submitted traintuple
	inpWaitingTraintuple := inputAggregateTuple{
		InModels: []string{aggregateTupleKey},
	}
	args = inpWaitingTraintuple.createDefault()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when adding aggregate tuple with status %d and message %s", resp.Status, resp.Message)

	// Query traintuple with status todo and worker as trainworker and check consistency
	filter := inputQueryFilter{
		IndexName:  "aggregateTuple~worker~status",
		Attributes: worker + ", todo",
	}
	args = [][]byte{[]byte("queryFilter"), assetToJSON(filter)}
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when querying aggregate tuple of worker with todo status - status %d and message %s", resp.Status, resp.Message)
	err = json.Unmarshal(resp.Payload, &queryTraintuples)
	assert.NoError(t, err, "aggregate tuples should unmarshal without problem")
	assert.Exactly(t, out, queryTraintuples[0])

	// Update status and check consistency
	success := inputLogSuccessTrain{}
	success.Key = traintupleKey
	success.fillDefaults()

	argsSlice := [][][]byte{
		[][]byte{[]byte("logStartAggregateTrain"), keyToJSON(traintupleKey)},
		[][]byte{[]byte("logSuccessAggregateTrain"), assetToJSON(success)},
	}
	traintupleStatus := []string{StatusDoing, StatusDone}
	for i := range traintupleStatus {
		resp = mockStub.MockInvoke("42", argsSlice[i])
		require.EqualValuesf(t, 200, resp.Status, "when logging start %s with message %s", traintupleStatus[i], resp.Message)
		filter := inputQueryFilter{
			IndexName:  "aggregateTuple~worker~status",
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

	// Query AggregateTuple From key
	args = [][]byte{[]byte("queryAggregatetuple"), keyToJSON(aggregateTupleKey)}
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when querying aggregate tuple with status %d and message %s", resp.Status, resp.Message)
	endTraintuple := outputAggregateTuple{}
	assert.NoError(t, json.Unmarshal(resp.Payload, &endTraintuple))
	expected.Log = success.Log
	expected.OutModel = &HashDress{
		Hash:           modelHash,
		StorageAddress: modelAddress}
	expected.Status = traintupleStatus[1]
	assert.Exactly(t, expected, endTraintuple, "retreived AggregateTuple does not correspond to what is expected")

	// query all traintuples related to a traintuple with the same algo
	args = [][]byte{[]byte("queryModelDetails"), keyToJSON(traintupleKey)}
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when querying model details with status %d and message %s", resp.Status, resp.Message)
	payload := outputModelDetails{}
	assert.NoError(t, json.Unmarshal(resp.Payload, &payload))
	assert.NotNil(t, payload.AggregateTuple, "when querying model tuples, payload should contain one traintuple")

	// query all traintuples related to a traintuple with the same algo
	args = [][]byte{[]byte("queryModels")}
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when querying models with status %d and message %s", resp.Status, resp.Message)
}

func TestQueryTraintupleNotFoundAggregate(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)
	registerItem(t, *mockStub, "aggregateAlgo")

	inpTraintuple := inputAggregateTuple{}
	inpTraintuple.fillDefaults()
	args := inpTraintuple.getArgs()
	resp := mockStub.MockInvoke("42", args)
	var _key struct{ Key string }
	json.Unmarshal(resp.Payload, &_key)

	// queryAggregateTuple: normal queryAggregateTuple
	args = [][]byte{[]byte("queryAggregatetuple"), keyToJSON(_key.Key)}
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when querying the aggregate tuple - status %d and message %s", resp.Status, resp.Message)

	// queryAggregateTuple: key does not exist
	notFoundKey := "eedbb7c31f62244c0f34461cc168804227115793d01c270021fe3f7935482eed"
	args = [][]byte{[]byte("queryAggregatetuple"), keyToJSON(notFoundKey)}
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 404, resp.Status, "when querying the aggregate tuple - status %d and message %s", resp.Status, resp.Message)

	// queryAggregateTuple: key does not exist and use existing other asset type key
	args = [][]byte{[]byte("queryAggregatetuple"), keyToJSON(algoHash)}
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 404, resp.Status, "when querying the aggregate tuple - status %d and message %s", resp.Status, resp.Message)
}

func TestInsertTraintupleTwiceAggregate(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)
	registerItem(t, *mockStub, "trainDataset")

	inpAlgo := inputAggregateAlgo{}
	args := inpAlgo.createDefault()
	resp := mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status, "when adding algo it should work: ", resp.Message)

	// create a aggregate tuple and start a ComplutePlan
	inpTraintuple := inputAggregateTuple{
		Rank: "0",
	}
	inpTraintuple.createDefault()
	resp = mockStub.MockInvoke("42", methodAndAssetToByte("createAggregatetuple", inpTraintuple))
	assert.EqualValues(t, http.StatusOK, resp.Status)
	var _key struct{ Key string }
	json.Unmarshal(resp.Payload, &_key)

	// create a second aggregate tuple in the same ComputePlan
	inpTraintuple.Rank = "1"
	inpTraintuple.ComputePlanID = _key.Key
	inpTraintuple.InModels = []string{_key.Key}
	resp = mockStub.MockInvoke("42", methodAndAssetToByte("createAggregatetuple", inpTraintuple))
	assert.EqualValues(t, http.StatusOK, resp.Status)

	// re-insert the same aggregate tuple and expect a conflict error
	resp = mockStub.MockInvoke("42", methodAndAssetToByte("createAggregatetuple", inpTraintuple))
	assert.EqualValues(t, http.StatusConflict, resp.Status)
}

//////////////////////////////////////////////
//
// Aggregate-specific tests
// Not copied from `traintuple_test.go`
//
/////////////////////////////////////////////

func TestAggregateTuplePermissions(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)
	registerItem(t, *mockStub, "trainDataset")

	objHash := strings.ReplaceAll(objectiveDescriptionHash, "1", "2")
	inpObjective := inputObjective{DescriptionHash: objHash}
	inpObjective.createDefault()
	inpObjective.TestDataset = inputDataset{}
	resp := mockStub.MockInvoke("42", methodAndAssetToByte("registerObjective", inpObjective))

	inpAlgo := inputAggregateAlgo{}
	args := inpAlgo.createDefault()
	resp = mockStub.MockInvoke("42", args)

	inpTraintuple := inputAggregateTuple{ObjectiveKey: objHash}
	inpTraintuple.fillDefaults()
	args = inpTraintuple.getArgs()
	resp = mockStub.MockInvoke("42", args)

	traintuple := outputAggregateTuple{}
	json.Unmarshal(resp.Payload, &traintuple)
	args = [][]byte{[]byte("queryAggregatetuple"), keyToJSON(traintuple.Key)}
	resp = mockStub.MockInvoke("42", args)
	traintuple = outputAggregateTuple{}
	json.Unmarshal(resp.Payload, &traintuple)

	// TODO (aggregate): check permissions
	// assert.EqualValues(t, false, traintuple.OutHeadModel.Permissions.Process.Public,
	// 	"the head model should not be public")
	// assert.EqualValues(t, []string{worker}, traintuple.OutHeadModel.Permissions.Process.AuthorizedIDs,
	// 	"the head model should only be processable by creator")
}

func TestAggregateTupleLogSuccessFail(t *testing.T) {
	for _, status := range []string{StatusDone, StatusFailed} {
		t.Run("TestAggregateTupleLog"+status, func(t *testing.T) {
			scc := new(SubstraChaincode)
			mockStub := NewMockStubWithRegisterNode("substra", scc)
			resp, _ := registerItem(t, *mockStub, "aggregateTuple")
			var _key struct{ Key string }
			json.Unmarshal(resp.Payload, &_key)
			key := _key.Key

			// start
			resp = mockStub.MockInvoke("42", [][]byte{[]byte("logStartAggregateTrain"), keyToJSON(key)})

			var expectedStatus string

			switch status {
			case StatusDone:
				success := inputLogSuccessTrain{}
				success.Key = key
				success.createDefault()
				success.fillDefaults()
				resp = mockStub.MockInvoke("42", [][]byte{[]byte("logSuccessAggregateTrain"), assetToJSON(success)})
				require.EqualValuesf(t, 200, resp.Status, "traintuple should be successfully set to 'success': %s", resp.Message)
				expectedStatus = "done"
			case StatusFailed:
				failed := inputLogFailTrain{}
				failed.Key = key
				failed.fillDefaults()
				resp = mockStub.MockInvoke("42", [][]byte{[]byte("logFailAggregateTrain"), assetToJSON(failed)})
				require.EqualValuesf(t, 200, resp.Status, "traintuple should be successfully set to 'failed': %s", resp.Message)
				expectedStatus = "failed"
			}

			// fetch back
			args := [][]byte{[]byte("queryAggregatetuple"), keyToJSON(key)}
			resp = mockStub.MockInvoke("42", args)
			assert.EqualValues(t, 200, resp.Status, "It should find the traintuple without error: %s", resp.Message)
			traintuple := outputAggregateTuple{}
			json.Unmarshal(resp.Payload, &traintuple)
			assert.EqualValues(t, expectedStatus, traintuple.Status, "The traintuple status should be set to %s", expectedStatus)
		})
	}
}

func TestQueryAggregateTuple(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)
	_, _ = registerItem(t, *mockStub, "compositeTraintuple")

	in := inputAggregateTuple{}
	in.InModels = []string{traintupleKey, compositeTraintupleKey}
	args := in.createDefault()
	resp := mockStub.MockInvoke("42", args)
	require.EqualValuesf(t, 200, resp.Status, "when adding aggregate tuple with status %d and message %s", resp.Status, resp.Message)

	var keyOnly struct{ Key string }
	json.Unmarshal(resp.Payload, &keyOnly)

	args = [][]byte{[]byte("queryAggregatetuple"), keyToJSON(keyOnly.Key)}
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status, "It should find the aggregate tuple: %s", resp.Message)
	out := outputAggregateTuple{}
	json.Unmarshal(resp.Payload, &out)

	assert.NotEmpty(t, out.Key)
	assert.Equal(t, in.Worker, out.Worker)
	assert.Equal(t, worker, out.Creator)
	assert.Equal(t, in.Tag, out.Tag)
	assert.Len(t, out.InModels, 2)
	assert.Equal(t, traintupleKey, out.InModels[0].TraintupleKey)
	assert.Equal(t, compositeTraintupleKey, out.InModels[1].TraintupleKey)
	assert.Equal(t, aggregateAlgoName, out.Algo.Name)
	assert.Equal(t, in.AlgoKey, out.Algo.Hash)
	assert.Equal(t, aggregateAlgoStorageAddress, out.Algo.StorageAddress)
	assert.Equal(t, StatusWaiting, out.Status)
	assert.Equal(t, objectiveDescriptionHash, out.Objective.Key)
	assert.Equal(t, objectiveMetricsHash, out.Objective.Metrics.Hash)
	assert.Equal(t, objectiveMetricsStorageAddress, out.Objective.Metrics.StorageAddress)
}
