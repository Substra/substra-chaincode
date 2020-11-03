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

	key := strings.Replace(objectiveKey, "1", "2", 1)
	inpObjective := inputObjective{Key: key}
	inpObjective.createDefault()
	inpObjective.TestDataset = inputDataset{}
	resp := mockStub.MockInvoke("42", methodAndAssetToByte("registerObjective", inpObjective))
	assert.EqualValues(t, 200, resp.Status, "when adding objective without dataset it should work: ", resp.Message)

	inpAlgo := inputAggregateAlgo{}
	args := inpAlgo.createDefault()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status, "when adding aggregate algo it should work: ", resp.Message)

	inpTraintuple := inputAggregatetuple{}
	args = inpTraintuple.createDefault()
	resp = mockStub.MockInvoke("42", args)

	assert.EqualValues(t, 200, resp.Status, "when adding aggregate tuple without test dataset it should work: ", resp.Message)

	traintuple := outputAggregatetuple{}
	json.Unmarshal(resp.Payload, &traintuple)
	args = [][]byte{[]byte("queryAggregatetuple"), keyToJSON(traintuple.Key)}
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status, "It should find the aggregate tuple without error ", resp.Message)
}

func TestTraintupleWithSingleDatasampleAggregate(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)
	registerItem(t, *mockStub, "trainDataset")

	key := strings.Replace(objectiveKey, "1", "2", 1)
	inpObjective := inputObjective{Key: key}
	inpObjective.createDefault()
	inpObjective.TestDataset = inputDataset{}
	resp := mockStub.MockInvoke("42", methodAndAssetToByte("registerObjective", inpObjective))
	assert.EqualValues(t, 200, resp.Status, "when adding objective without dataset it should work: ", resp.Message)

	inpAlgo := inputAggregateAlgo{}
	args := inpAlgo.createDefault()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status, "when adding aggregate algo it should work: ", resp.Message)

	inpTraintuple := inputAggregatetuple{
		AlgoKey: aggregateAlgoKey,
	}
	args = inpTraintuple.createDefault()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status, "when adding aggregate tuple with a single data samples it should work: ", resp.Message)

	traintuple := outputKey{}
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
	registerItem(t, *mockStub, "aggregatetuple")

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
		getOutputAggregatetuple(NewLedgerDB(mockStub), traintupleKey)
	})
}

func TestTraintupleComputePlanCreationAggregate(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)

	// Add dataManager, dataSample and algo
	registerItem(t, *mockStub, "aggregateAlgo")

	inpTraintuple := inputAggregatetuple{ComputePlanKey: "someComputePlanKey"}
	args := inpTraintuple.createDefault()
	resp := mockStub.MockInvoke("42", args)
	require.EqualValues(t, 400, resp.Status, "should failed for missing rank")
	require.Contains(t, resp.Message, "invalid inputs, a ComputePlan should have a rank", "invalid error message")

	cpKey := RandomUUID()
	inCP := inputComputePlan{Key: cpKey}
	resp = mockStub.MockInvoke("42", inCP.getArgs())
	require.EqualValues(t, 200, resp.Status)

	inpTraintuple = inputAggregatetuple{Rank: "1"}
	args = inpTraintuple.createDefault()
	resp = mockStub.MockInvoke("42", args)
	require.EqualValues(t, 400, resp.Status, "should failed for invalid rank")
	require.Contains(t, resp.Message, "Field validation for 'ComputePlanKey' failed on the 'required_with' tag")

	inpTraintuple = inputAggregatetuple{Rank: "0", ComputePlanKey: cpKey}
	args = inpTraintuple.createDefault()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status)
	res := outputKey{}
	err := json.Unmarshal(resp.Payload, &res)
	assert.NoError(t, err, "should unmarshal without problem")
	key := res.Key
	require.EqualValues(t, aggregatetupleKey, key)

	inpTraintuple = inputAggregatetuple{}
	args = inpTraintuple.createDefault()
	resp = mockStub.MockInvoke("42", args)
	require.EqualValues(t, 409, resp.Status, "should failed for existing aggregatetuple key")
	require.Contains(t, resp.Message, "already exists")

	require.EqualValues(t, 409, resp.Status, "should failed for existing FLTask")
	errorPayload := map[string]interface{}{}
	err = json.Unmarshal(resp.Payload, &errorPayload)
	assert.NoError(t, err, "should unmarshal without problem")
	require.Contains(t, errorPayload, "key", "key should be available in payload")
	assert.EqualValues(t, aggregatetupleKey, errorPayload["key"], "key in error should be aggregatetupleKey")
}

func TestTraintupleMultipleCommputePlanCreationsAggregate(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)

	// Add a some dataManager, dataSample and traintuple
	registerItem(t, *mockStub, "aggregateAlgo")
	db := NewLedgerDB(mockStub)

	cpKey := RandomUUID()
	inCP := inputComputePlan{Key: cpKey}
	resp := mockStub.MockInvoke("42", inCP.getArgs())
	require.EqualValues(t, 200, resp.Status)

	inpTraintuple := inputAggregatetuple{Rank: "0", ComputePlanKey: cpKey}
	args := inpTraintuple.createDefault()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status)
	res := outputKey{}
	err := json.Unmarshal(resp.Payload, &res)
	assert.NoError(t, err, "should unmarshal without problem")
	key := res.Key
	_, err = db.GetAggregatetuple(key)
	assert.NoError(t, err)

	// Failed to add a traintuple with the same rank
	inpTraintuple = inputAggregatetuple{
		Key:            RandomUUID(),
		InModels:       []string{key},
		Rank:           "0",
		ComputePlanKey: cpKey}
	args = inpTraintuple.createDefault()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 400, resp.Status, resp.Message, "should failed to add an aggregate tuple of the same rank")

	// Failed to add a traintuple to an unexisting CommputePlan
	inpTraintuple = inputAggregatetuple{
		Key:            RandomUUID(),
		InModels:       []string{key},
		Rank:           "1",
		ComputePlanKey: "notarealone"}
	args = inpTraintuple.createDefault()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 404, resp.Status, resp.Message, "should failed to add an aggregate tuple to an unexisting ComputePlanKey")

	// Succesfully add a traintuple to the same ComputePlanKey
	inpTraintuple = inputAggregatetuple{
		Key:            RandomUUID(),
		InModels:       []string{key},
		Rank:           "1",
		ComputePlanKey: cpKey}
	args = inpTraintuple.createDefault()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status, resp.Message, "should be able do create an aggregate tuple with the same ComputePlanKey")
	err = json.Unmarshal(resp.Payload, &res)
	assert.NoError(t, err, "should unmarshal without problem")
}

func TestTraintupleAggregate(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)

	// Add traintuple with invalid field
	inpTraintuple := inputAggregatetuple{
		AlgoKey: "aaa",
	}
	args := inpTraintuple.createDefault()
	resp := mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 400, resp.Status, "when adding objective with invalid key, status %d and message %s", resp.Status, resp.Message)

	// Add traintuple with unexisting algo
	inpTraintuple = inputAggregatetuple{}
	args = inpTraintuple.createDefault()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 400, resp.Status, "when adding aggregate tuple with unexisting algo, status %d and message %s", resp.Status, resp.Message)

	// Properly add traintuple
	resp, tt := registerItem(t, *mockStub, "aggregatetuple")

	inpTraintuple = tt.(inputAggregatetuple)
	res := outputKey{}
	err := json.Unmarshal(resp.Payload, &res)
	assert.NoError(t, err, "aggregate tuple should unmarshal without problem")
	traintupleKey := res.Key
	// Query traintuple from key and check the consistency of returned arguments
	args = [][]byte{[]byte("queryAggregatetuple"), keyToJSON(traintupleKey)}
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when querying the aggregate tuple - status %d and message %s", resp.Status, resp.Message)
	out := outputAggregatetuple{}
	err = json.Unmarshal(resp.Payload, &out)
	assert.NoError(t, err, "when unmarshalling queried aggregate tuple")
	expected := outputAggregatetuple{
		Key: aggregatetupleKey,
		Algo: &KeyChecksumAddressName{
			Key:            aggregateAlgoKey,
			Checksum:       aggregateAlgoChecksum,
			Name:           aggregateAlgoName,
			StorageAddress: aggregateAlgoStorageAddress,
		},
		Creator: worker,
		Worker:  worker,
		Status:  StatusTodo,
		Permissions: outputPermissions{
			Process: Permission{
				Public:        true,
				AuthorizedIDs: []string{},
			},
		},
		Metadata: map[string]string{},
	}
	assert.Exactly(t, expected, out, "the aggregate tuple queried from the ledger differ from expected")

	// Query all traintuples and check consistency
	args = [][]byte{[]byte("queryAggregatetuples")}
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when querying aggregate tuples - status %d and message %s", resp.Status, resp.Message)
	// TODO add traintuple key to output struct
	// For now we test it as cleanly as its added to the query response
	assert.Contains(t, string(resp.Payload), "key\":\""+aggregatetupleKey)
	var queryTraintuples []outputAggregatetuple
	err = json.Unmarshal(resp.Payload, &queryTraintuples)
	assert.NoError(t, err, "aggregate tuples should unmarshal without problem")
	require.NotZero(t, queryTraintuples)
	assert.Exactly(t, out, queryTraintuples[0])

	// Add traintuple with inmodel from the above-submitted traintuple
	inpWaitingTraintuple := inputAggregatetuple{
		InModels: []string{aggregatetupleKey},
		Key:      RandomUUID(),
	}
	args = inpWaitingTraintuple.createDefault()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when adding aggregate tuple with status %d and message %s", resp.Status, resp.Message)

	// Query traintuple with status todo and worker as trainworker and check consistency
	filter := inputQueryFilter{
		IndexName:  "aggregatetuple~worker~status",
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
		[][]byte{[]byte("logStartAggregate"), keyToJSON(traintupleKey)},
		[][]byte{[]byte("logSuccessAggregate"), assetToJSON(success)},
	}
	traintupleStatus := []string{StatusDoing, StatusDone}
	for i := range traintupleStatus {
		resp = mockStub.MockInvoke("42", argsSlice[i])
		require.EqualValuesf(t, 200, resp.Status, "when logging start %s with message %s", traintupleStatus[i], resp.Message)
		filter := inputQueryFilter{
			IndexName:  "aggregatetuple~worker~status",
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

	// Query Aggregatetuple From key
	args = [][]byte{[]byte("queryAggregatetuple"), keyToJSON(aggregatetupleKey)}
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when querying aggregate tuple with status %d and message %s", resp.Status, resp.Message)
	endTraintuple := outputAggregatetuple{}
	assert.NoError(t, json.Unmarshal(resp.Payload, &endTraintuple))
	expected.Log = success.Log
	expected.OutModel = &KeyChecksumAddress{
		Key:            modelKey,
		Checksum:       modelChecksum,
		StorageAddress: modelAddress}
	expected.Status = traintupleStatus[1]
	assert.Exactly(t, expected, endTraintuple, "retreived Aggregatetuple does not correspond to what is expected")

	// query all traintuples related to a traintuple with the same algo
	args = [][]byte{[]byte("queryModelDetails"), keyToJSON(traintupleKey)}
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when querying model details with status %d and message %s", resp.Status, resp.Message)
	payload := outputModelDetails{}
	assert.NoError(t, json.Unmarshal(resp.Payload, &payload))
	assert.NotNil(t, payload.Aggregatetuple, "when querying model tuples, payload should contain one traintuple")

	// query all traintuples related to a traintuple with the same algo
	args = [][]byte{[]byte("queryModels")}
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when querying models with status %d and message %s", resp.Status, resp.Message)
}

func TestQueryTraintupleNotFoundAggregate(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)
	registerItem(t, *mockStub, "aggregateAlgo")

	inpTraintuple := inputAggregatetuple{}
	inpTraintuple.fillDefaults()
	args := inpTraintuple.getArgs()
	resp := mockStub.MockInvoke("42", args)
	var _key struct{ Key string }
	json.Unmarshal(resp.Payload, &_key)

	// queryAggregatetuple: normal queryAggregatetuple
	args = [][]byte{[]byte("queryAggregatetuple"), keyToJSON(_key.Key)}
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when querying the aggregate tuple - status %d and message %s", resp.Status, resp.Message)

	// queryAggregatetuple: key does not exist
	notFoundKey := "eedbb7c3-1f62-244c-0f34-461cc1688042"
	args = [][]byte{[]byte("queryAggregatetuple"), keyToJSON(notFoundKey)}
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 404, resp.Status, "when querying the aggregate tuple - status %d and message %s", resp.Status, resp.Message)

	// queryAggregatetuple: key does not exist and use existing other asset type key
	args = [][]byte{[]byte("queryAggregatetuple"), keyToJSON(algoKey)}
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
	cpKey := RandomUUID()
	inCP := inputComputePlan{Key: cpKey}
	resp = mockStub.MockInvoke("42", inCP.getArgs())
	require.EqualValues(t, 200, resp.Status)
	inpTraintuple := inputAggregatetuple{
		Rank:           "0",
		ComputePlanKey: cpKey,
	}
	inpTraintuple.createDefault()
	resp = mockStub.MockInvoke("42", methodAndAssetToByte("createAggregatetuple", inpTraintuple))
	assert.EqualValues(t, http.StatusOK, resp.Status)
	var _key struct{ Key string }
	json.Unmarshal(resp.Payload, &_key)
	db := NewLedgerDB(mockStub)
	tuple, err := db.GetAggregatetuple(_key.Key)
	assert.NoError(t, err)
	// create a second aggregate tuple in the same ComputePlan
	inpTraintuple.Key = RandomUUID()
	inpTraintuple.Rank = "1"
	inpTraintuple.ComputePlanKey = tuple.ComputePlanKey
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

func TestAggregatetuplePermissions(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)
	registerItem(t, *mockStub, "aggregateAlgo")

	// register nodes
	registerNode := func(nodeName string) {
		initialCreator := mockStub.Creator
		mockStub.Creator = nodeName
		mockStub.MockInvoke("42", [][]byte{[]byte("registerNode")})
		mockStub.Creator = initialCreator
	}
	registerNode("nodeA")
	registerNode("nodeB")
	registerNode("nodeC")
	registerNode("nodeD")

	// register 3 algos
	algo1, err := registerRandomCompositeAlgo(t, mockStub)
	assert.Nil(t, err)
	algo2, err := registerRandomCompositeAlgo(t, mockStub)
	assert.Nil(t, err)
	algo3, err := registerRandomCompositeAlgo(t, mockStub)
	assert.Nil(t, err)

	// register 3 composite traintuples, with various permissions
	registerCompositeTraintuple := func(algoKey string, authorizedIds []string) string {
		inp := inputCompositeTraintuple{Key: RandomUUID(), AlgoKey: algoKey}
		inp.fillDefaults()
		inp.OutTrunkModelPermissions.Process.Public = false
		inp.OutTrunkModelPermissions.Process.AuthorizedIDs = authorizedIds
		resp := mockStub.MockInvoke("42", inp.getArgs())
		assert.EqualValues(t, 200, resp.Status, resp.Message)
		var _key struct{ Key string }
		json.Unmarshal(resp.Payload, &_key)
		return _key.Key
	}
	traintuple1 := registerCompositeTraintuple(algo1, []string{"nodeA", "nodeC"})
	traintuple2 := registerCompositeTraintuple(algo2, []string{"nodeB", "nodeC"})
	traintuple3 := registerCompositeTraintuple(algo3, []string{"nodeA", "nodeC", "nodeD"})

	// create an aggregate tuple with the 3 composite as in-models
	inpAgg := inputAggregatetuple{}
	inpAgg.fillDefaults()
	inpAgg.InModels = []string{traintuple1, traintuple2, traintuple3}
	resp := mockStub.MockInvoke("42", inpAgg.getArgs())
	assert.EqualValues(t, 200, resp.Status, resp.Message)
	var _key struct{ Key string }
	json.Unmarshal(resp.Payload, &_key)
	aggrKey := _key.Key

	// fetch the aggregate tuple back
	aggr := outputAggregatetuple{}
	args := [][]byte{[]byte("queryAggregatetuple"), keyToJSON(aggrKey)}
	resp = mockStub.MockInvoke("42", args)
	aggr = outputAggregatetuple{}
	json.Unmarshal(resp.Payload, &aggr)

	// verify permissions
	assert.EqualValues(t, false, aggr.Permissions.Process.Public,
		"the aggregate tuple should not be public")
	assert.EqualValues(t, []string{worker, "nodeC"}, aggr.Permissions.Process.AuthorizedIDs,
		"the aggregate tuple permissions should be the intersect of the in-model permissions")
}

func TestAggregatetupleLogSuccessFail(t *testing.T) {
	for _, status := range []string{StatusDone, StatusFailed} {
		t.Run("TestAggregatetupleLog"+status, func(t *testing.T) {
			scc := new(SubstraChaincode)
			mockStub := NewMockStubWithRegisterNode("substra", scc)
			resp, _ := registerItem(t, *mockStub, "aggregatetuple")
			var _key struct{ Key string }
			json.Unmarshal(resp.Payload, &_key)
			key := _key.Key

			// start
			resp = mockStub.MockInvoke("42", [][]byte{[]byte("logStartAggregate"), keyToJSON(key)})

			var expectedStatus string

			switch status {
			case StatusDone:
				success := inputLogSuccessTrain{}
				success.Key = key
				success.createDefault()
				success.fillDefaults()
				resp = mockStub.MockInvoke("42", [][]byte{[]byte("logSuccessAggregate"), assetToJSON(success)})
				require.EqualValuesf(t, 200, resp.Status, "traintuple should be successfully set to 'success': %s", resp.Message)
				expectedStatus = "done"
			case StatusFailed:
				failed := inputLogFailTrain{}
				failed.Key = key
				failed.fillDefaults()
				resp = mockStub.MockInvoke("42", [][]byte{[]byte("logFailAggregate"), assetToJSON(failed)})
				require.EqualValuesf(t, 200, resp.Status, "traintuple should be successfully set to 'failed': %s", resp.Message)
				expectedStatus = "failed"
			}

			// fetch back
			args := [][]byte{[]byte("queryAggregatetuple"), keyToJSON(key)}
			resp = mockStub.MockInvoke("42", args)
			assert.EqualValues(t, 200, resp.Status, "It should find the traintuple without error: %s", resp.Message)
			traintuple := outputAggregatetuple{}
			json.Unmarshal(resp.Payload, &traintuple)
			assert.EqualValues(t, expectedStatus, traintuple.Status, "The traintuple status should be set to %s", expectedStatus)
		})
	}
}

func TestQueryAggregatetuple(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)
	_, _ = registerItem(t, *mockStub, "compositeTraintuple")

	in := inputAggregatetuple{}
	in.InModels = []string{traintupleKey, compositeTraintupleKey}
	args := in.createDefault()
	resp := mockStub.MockInvoke("42", args)
	require.EqualValuesf(t, 200, resp.Status, "when adding aggregate tuple with status %d and message %s", resp.Status, resp.Message)

	var keyOnly struct{ Key string }
	json.Unmarshal(resp.Payload, &keyOnly)

	args = [][]byte{[]byte("queryAggregatetuple"), keyToJSON(keyOnly.Key)}
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status, "It should find the aggregate tuple: %s", resp.Message)
	out := outputAggregatetuple{}
	json.Unmarshal(resp.Payload, &out)

	assert.NotEmpty(t, out.Key)
	assert.Equal(t, in.Worker, out.Worker)
	assert.Equal(t, worker, out.Creator)
	assert.Equal(t, in.Tag, out.Tag)
	assert.Len(t, out.InModels, 2)
	assert.Equal(t, traintupleKey, out.InModels[0].TraintupleKey)
	assert.Equal(t, compositeTraintupleKey, out.InModels[1].TraintupleKey)
	assert.Equal(t, aggregateAlgoName, out.Algo.Name)
	assert.Equal(t, aggregateAlgoChecksum, out.Algo.Checksum)
	assert.Equal(t, aggregateAlgoStorageAddress, out.Algo.StorageAddress)
	assert.Equal(t, StatusWaiting, out.Status)
}

func TestCreateFailedAggregate(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)
	registerItem(t, *mockStub, "compositeTraintuple")
	mockStub.MockTransactionStart("42")
	db := NewLedgerDB(mockStub)

	_, err := logStartCompositeTrain(db, assetToArgs(inputKey{Key: compositeTraintupleKey}))
	assert.NoError(t, err)

	_, err = logFailCompositeTrain(db, assetToArgs(inputLogFailTrain{inputLog{Key: compositeTraintupleKey}}))
	assert.NoError(t, err)

	in := inputAggregatetuple{}
	in.fillDefaults()
	in.InModels = []string{compositeTraintupleKey, traintupleKey}
	key, err := createAggregatetupleInternal(db, in, true)
	assert.NoError(t, err)

	out, err := queryAggregatetuple(db, assetToArgs(inputKey{Key: key}))
	assert.NoError(t, err)
	assert.Equal(t, StatusFailed, out.Status)
}
