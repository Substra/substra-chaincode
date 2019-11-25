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
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/////////////////////////////////////////////////////////////
//
// "Regular" tests
// Copied from `traintuple_test.go` and adapted for composite
//
/////////////////////////////////////////////////////////////

func TestTraintupleWithNoTestDatasetComposite(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)
	registerItem(t, *mockStub, "trainDataset")

	objHash := strings.ReplaceAll(objectiveDescriptionHash, "1", "2")
	inpObjective := inputObjective{DescriptionHash: objHash}
	inpObjective.createDefault()
	inpObjective.TestDataset = inputDataset{}
	resp := mockStub.MockInvoke("42", methodAndAssetToByte("registerObjective", inpObjective))
	assert.EqualValues(t, 200, resp.Status, "when adding objective without dataset it should work: ", resp.Message)

	inpAlgo := inputCompositeAlgo{}
	args := inpAlgo.createDefault()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status, "when adding algo it should work: ", resp.Message)

	inpTraintuple := inputCompositeTraintuple{ObjectiveKey: objHash}
	args = inpTraintuple.createDefault()
	resp = mockStub.MockInvoke("42", args)

	assert.EqualValues(t, 200, resp.Status, "when adding traintuple without test dataset it should work: ", resp.Message)

	traintuple := outputCompositeTraintuple{}
	json.Unmarshal(resp.Payload, &traintuple)
	args = [][]byte{[]byte("queryCompositeTraintuple"), keyToJSON(traintuple.Key)}
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status, "It should find the traintuple without error ", resp.Message)
}

func TestTraintupleWithSingleDatasampleComposite(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)
	registerItem(t, *mockStub, "trainDataset")

	objHash := strings.ReplaceAll(objectiveDescriptionHash, "1", "2")
	inpObjective := inputObjective{DescriptionHash: objHash}
	inpObjective.createDefault()
	inpObjective.TestDataset = inputDataset{}
	resp := mockStub.MockInvoke("42", methodAndAssetToByte("registerObjective", inpObjective))
	assert.EqualValues(t, 200, resp.Status, "when adding objective without dataset it should work: ", resp.Message)

	inpAlgo := inputCompositeAlgo{}
	args := inpAlgo.createDefault()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status, "when adding algo it should work: ", resp.Message)

	inpTraintuple := inputCompositeTraintuple{
		ObjectiveKey:   objHash,
		AlgoKey:        compositeAlgoHash,
		DataSampleKeys: []string{trainDataSampleHash1},
	}
	args = inpTraintuple.createDefault()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status, "when adding composite traintuple with a single data samples it should work: ", resp.Message)

	traintuple := outputCompositeTraintuple{}
	err := json.Unmarshal(resp.Payload, &traintuple)
	assert.NoError(t, err, "should be unmarshaled")
	args = [][]byte{[]byte("queryCompositeTraintuple"), keyToJSON(traintuple.Key)}
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status, "It should find the composite traintuple without error ", resp.Message)
}

func TestTraintupleWithDuplicatedDatasamplesComposite(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)
	registerItem(t, *mockStub, "trainDataset")

	objHash := strings.ReplaceAll(objectiveDescriptionHash, "1", "2")
	inpObjective := inputObjective{DescriptionHash: objHash}
	inpObjective.createDefault()
	inpObjective.TestDataset = inputDataset{}
	resp := mockStub.MockInvoke("42", methodAndAssetToByte("registerObjective", inpObjective))
	assert.EqualValues(t, 200, resp.Status, "when adding objective without dataset it should work: ", resp.Message)

	inpAlgo := inputCompositeAlgo{}
	args := inpAlgo.createDefault()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status, "when adding composite algo it should work: ", resp.Message)

	inpTraintuple := inputCompositeTraintuple{
		ObjectiveKey:   objHash,
		DataSampleKeys: []string{trainDataSampleHash1, trainDataSampleHash2, trainDataSampleHash1},
	}
	args = inpTraintuple.createDefault()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 400, resp.Status, "when adding traintuple with a duplicated data samples it should not work: %s", resp.Message)
}

func TestNoPanicWhileQueryingIncompleteTraintupleComposite(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)
	// Add a some dataManager, dataSample and traintuple
	registerItem(t, *mockStub, "traintuple")

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
		getOutputCompositeTraintuple(NewLedgerDB(mockStub), traintupleKey)
	})
}

func TestTraintupleComputePlanCreationComposite(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)

	// Add dataManager, dataSample and algo
	registerItem(t, *mockStub, "compositeAlgo")

	inpTraintuple := inputCompositeTraintuple{ComputePlanID: "someComputePlanID"}
	args := inpTraintuple.createDefault()
	resp := mockStub.MockInvoke("42", args)
	require.EqualValues(t, 400, resp.Status, "should failed for missing rank")
	require.Contains(t, resp.Message, "invalid inputs, a ComputePlan should have a rank", "invalid error message")

	inpTraintuple = inputCompositeTraintuple{Rank: "1"}
	args = inpTraintuple.createDefault()
	resp = mockStub.MockInvoke("42", args)
	require.EqualValues(t, 400, resp.Status, "should failed for invalid rank")
	require.Contains(t, resp.Message, "invalid inputs, a new ComputePlan should have a rank 0")

	inpTraintuple = inputCompositeTraintuple{Rank: "0"}
	args = inpTraintuple.createDefault()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status)
	res := map[string]string{}
	err := json.Unmarshal(resp.Payload, &res)
	assert.NoError(t, err, "should unmarshal without problem")
	assert.Contains(t, res, "key")
	key := res["key"]
	require.EqualValues(t, key, compositeTraintupleKey)

	inpTraintuple = inputCompositeTraintuple{Rank: "0"}
	args = inpTraintuple.createDefault()
	resp = mockStub.MockInvoke("42", args)
	require.EqualValues(t, 409, resp.Status, "should failed for existing ComputePlanID")
	require.Contains(t, resp.Message, "already exists")

	require.EqualValues(t, 409, resp.Status, "should failed for existing FLTask")
	errorPayload := map[string]interface{}{}
	err = json.Unmarshal(resp.Payload, &errorPayload)
	assert.NoError(t, err, "should unmarshal without problem")
	require.Contains(t, errorPayload, "key", "key should be available in payload")
	assert.EqualValues(t, compositeTraintupleKey, errorPayload["key"], "key in error should be compositeTraintupleKey")
}

func TestTraintupleMultipleCommputePlanCreationsComposite(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)

	// Add a some dataManager, dataSample and traintuple
	registerItem(t, *mockStub, "compositeAlgo")

	inpTraintuple := inputCompositeTraintuple{Rank: "0"}
	args := inpTraintuple.createDefault()
	resp := mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status)
	res := map[string]string{}
	err := json.Unmarshal(resp.Payload, &res)
	assert.NoError(t, err, "should unmarshal without problem")
	assert.Contains(t, res, "key")
	key := res["key"]
	// Failed to add a traintuple with the same rank
	inpTraintuple = inputCompositeTraintuple{
		InHeadModelKey:  key,
		InTrunkModelKey: key,
		Rank:            "0",
		ComputePlanID:   key}
	args = inpTraintuple.createDefault()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 400, resp.Status, resp.Message, "should failed to add a traintuple of the same rank")

	// Failed to add a traintuple to an unexisting CommputePlan
	inpTraintuple = inputCompositeTraintuple{
		InHeadModelKey:  key,
		InTrunkModelKey: key,
		Rank:            "1",
		ComputePlanID:   "notarealone"}
	args = inpTraintuple.createDefault()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 400, resp.Status, resp.Message, "should failed to add a traintuple to an unexisting ComputePlanID")

	// Succesfully add a traintuple to the same ComputePlanID
	inpTraintuple = inputCompositeTraintuple{
		InHeadModelKey:  key,
		InTrunkModelKey: key,
		Rank:            "1",
		ComputePlanID:   key}
	args = inpTraintuple.createDefault()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status, resp.Message, "should be able do create a traintuple with the same ComputePlanID")
	err = json.Unmarshal(resp.Payload, &res)
	assert.NoError(t, err, "should unmarshal without problem")
	assert.Contains(t, res, "key")
	ttkey := res["key"]
	// Add new algo to check all ComputePlan algo consistency
	newAlgoHash := strings.Replace(compositeAlgoHash, "a", "b", 1)
	inpAlgo := inputCompositeAlgo{inputAlgo{Hash: newAlgoHash}}
	args = inpAlgo.createDefault()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status)

	inpTraintuple = inputCompositeTraintuple{
		AlgoKey:         newAlgoHash,
		InHeadModelKey:  ttkey,
		InTrunkModelKey: ttkey,
		Rank:            "2",
		ComputePlanID:   key}
	args = inpTraintuple.createDefault()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 400, resp.Status, resp.Message, "should fail for it doesn't have the same composite algo key")
	assert.Contains(t, resp.Message, "does not have the same algo key")
}

func TestTraintupleComposite(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)

	// Add traintuple with invalid field
	inpTraintuple := inputCompositeTraintuple{
		AlgoKey: "aaa",
	}
	args := inpTraintuple.createDefault()
	resp := mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 400, resp.Status, "when adding objective with invalid hash, status %d and message %s", resp.Status, resp.Message)

	// Add traintuple with unexisting algo
	inpTraintuple = inputCompositeTraintuple{}
	args = inpTraintuple.createDefault()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 400, resp.Status, "when adding composite traintuple with unexisting algo, status %d and message %s", resp.Status, resp.Message)

	// Properly add traintuple
	resp, tt := registerItem(t, *mockStub, "compositeTraintuple")

	inpTraintuple = tt.(inputCompositeTraintuple)
	res := map[string]string{}
	err := json.Unmarshal(resp.Payload, &res)
	assert.NoError(t, err, "composite traintuple should unmarshal without problem")
	assert.Contains(t, res, "key")
	traintupleKey := res["key"]
	// Query traintuple from key and check the consistency of returned arguments
	args = [][]byte{[]byte("queryCompositeTraintuple"), keyToJSON(traintupleKey)}
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when querying the composite traintuple - status %d and message %s", resp.Status, resp.Message)
	out := outputCompositeTraintuple{}
	err = json.Unmarshal(resp.Payload, &out)
	assert.NoError(t, err, "when unmarshalling queried composite traintuple")
	expected := outputCompositeTraintuple{
		Key: compositeTraintupleKey,
		Algo: &HashDressName{
			Hash:           compositeAlgoHash,
			Name:           compositeAlgoName,
			StorageAddress: compositeAlgoStorageAddress,
		},
		Creator: worker,
		Dataset: &TtDataset{
			DataSampleKeys: []string{trainDataSampleHash1, trainDataSampleHash2},
			OpenerHash:     dataManagerOpenerHash,
			Perf:           0.0,
			Worker:         worker,
		},
		Objective: &TtObjective{
			Key: objectiveDescriptionHash,
			Metrics: &HashDress{
				Hash:           objectiveMetricsHash,
				StorageAddress: objectiveMetricsStorageAddress,
			},
		},
		OutHeadModel: outModelComposite{
			Permissions: outputPermissions{
				Process: Permission{Public: false, AuthorizedIDs: []string{worker}},
			},
		},
		OutTrunkModel: outModelComposite{
			Permissions: outputPermissions{
				Process: Permission{Public: true, AuthorizedIDs: []string{}},
			},
		},
		Status: StatusTodo,
	}
	assert.Exactly(t, expected, out, "the composite traintuple queried from the ledger differ from expected")

	// Query all traintuples and check consistency
	args = [][]byte{[]byte("queryCompositeTraintuples")}
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when querying composite traintuples - status %d and message %s", resp.Status, resp.Message)
	// TODO add traintuple key to output struct
	// For now we test it as cleanly as its added to the query response
	assert.Contains(t, string(resp.Payload), "key\":\""+compositeTraintupleKey)
	var queryTraintuples []outputCompositeTraintuple
	err = json.Unmarshal(resp.Payload, &queryTraintuples)
	assert.NoError(t, err, "composite traintuples should unmarshal without problem")
	require.NotZero(t, queryTraintuples)
	assert.Exactly(t, out, queryTraintuples[0])

	// Add traintuple with inmodel from the above-submitted traintuple
	inpWaitingTraintuple := inputCompositeTraintuple{
		InHeadModelKey:  compositeTraintupleKey,
		InTrunkModelKey: compositeTraintupleKey}
	args = inpWaitingTraintuple.createDefault()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when adding composite traintuple with status %d and message %s", resp.Status, resp.Message)

	// Query traintuple with status todo and worker as trainworker and check consistency
	filter := inputQueryFilter{
		IndexName:  "compositeTraintuple~worker~status",
		Attributes: worker + ", todo",
	}
	args = [][]byte{[]byte("queryFilter"), assetToJSON(filter)}
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when querying composite traintuple of worker with todo status - status %d and message %s", resp.Status, resp.Message)
	err = json.Unmarshal(resp.Payload, &queryTraintuples)
	assert.NoError(t, err, "composite traintuples should unmarshal without problem")
	assert.Exactly(t, out, queryTraintuples[0])

	// Update status and check consistency
	success := inputLogSuccessCompositeTrain{}
	success.Key = traintupleKey

	argsSlice := [][][]byte{
		[][]byte{[]byte("logStartCompositeTrain"), keyToJSON(traintupleKey)},
		success.createDefault(),
	}
	traintupleStatus := []string{StatusDoing, StatusDone}
	for i := range traintupleStatus {
		resp = mockStub.MockInvoke("42", argsSlice[i])
		require.EqualValuesf(t, 200, resp.Status, "when logging start %s with message %s", traintupleStatus[i], resp.Message)
		filter := inputQueryFilter{
			IndexName:  "compositeTraintuple~worker~status",
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

	// Query CompositeTraintuple From key
	args = [][]byte{[]byte("queryCompositeTraintuple"), keyToJSON(compositeTraintupleKey)}
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when querying composite traintuple with status %d and message %s", resp.Status, resp.Message)
	endTraintuple := outputCompositeTraintuple{}
	assert.NoError(t, json.Unmarshal(resp.Payload, &endTraintuple))
	expected.Dataset.Perf = success.Perf
	expected.Log = success.Log
	expected.OutHeadModel.OutModel = &HashDress{
		Hash:           headModelHash,
		StorageAddress: headModelAddress}
	expected.OutTrunkModel.OutModel = &HashDress{
		Hash:           trunkModelHash,
		StorageAddress: trunkModelAddress}
	expected.Status = traintupleStatus[1]
	assert.Exactly(t, expected, endTraintuple, "retreived CompositeTraintuple does not correspond to what is expected")

	// query all traintuples related to a traintuple with the same algo
	args = [][]byte{[]byte("queryModelDetails"), keyToJSON(traintupleKey)}
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when querying model details with status %d and message %s", resp.Status, resp.Message)
	payload := outputModelDetails{}
	assert.NoError(t, json.Unmarshal(resp.Payload, &payload))
	assert.NotNil(t, payload.CompositeTraintuple, "when querying model tuples, payload should contain one traintuple")

	// query all traintuples related to a traintuple with the same algo
	args = [][]byte{[]byte("queryModels")}
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when querying models with status %d and message %s", resp.Status, resp.Message)
}

func TestQueryTraintupleNotFoundComposite(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)
	registerItem(t, *mockStub, "compositeAlgo")

	inpTraintuple := inputCompositeTraintuple{}
	inpTraintuple.fillDefaults()
	args := inpTraintuple.getArgs()
	resp := mockStub.MockInvoke("42", args)
	var _key struct{ Key string }
	json.Unmarshal(resp.Payload, &_key)

	// queryCompositeTraintuple: normal queryCompositeTraintuple
	args = [][]byte{[]byte("queryCompositeTraintuple"), keyToJSON(_key.Key)}
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when querying the composite traintuple - status %d and message %s", resp.Status, resp.Message)

	// queryCompositeTraintuple: key does not exist
	notFoundKey := "eedbb7c31f62244c0f34461cc168804227115793d01c270021fe3f7935482eed"
	args = [][]byte{[]byte("queryCompositeTraintuple"), keyToJSON(notFoundKey)}
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 404, resp.Status, "when querying the composite traintuple - status %d and message %s", resp.Status, resp.Message)

	// queryCompositeTraintuple: key does not exist and use existing other asset type key
	args = [][]byte{[]byte("queryCompositeTraintuple"), keyToJSON(algoHash)}
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 404, resp.Status, "when querying the composite traintuple - status %d and message %s", resp.Status, resp.Message)
}

func TestInsertTraintupleTwiceComposite(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)
	registerItem(t, *mockStub, "trainDataset")

	inpAlgo := inputCompositeAlgo{}
	args := inpAlgo.createDefault()
	resp := mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status, "when adding algo it should work: ", resp.Message)

	// create a composite traintuple and start a ComplutePlan
	inpTraintuple := inputCompositeTraintuple{
		Rank: "0",
	}
	inpTraintuple.createDefault()
	resp = mockStub.MockInvoke("42", methodAndAssetToByte("createCompositeTraintuple", inpTraintuple))
	assert.EqualValues(t, http.StatusOK, resp.Status)
	var _key struct{ Key string }
	json.Unmarshal(resp.Payload, &_key)

	// create a second composite traintuple in the same ComputePlan
	inpTraintuple.Rank = "1"
	inpTraintuple.ComputePlanID = _key.Key
	inpTraintuple.InHeadModelKey = _key.Key
	inpTraintuple.InTrunkModelKey = _key.Key
	resp = mockStub.MockInvoke("42", methodAndAssetToByte("createCompositeTraintuple", inpTraintuple))
	assert.EqualValues(t, http.StatusOK, resp.Status)

	// re-insert the same composite traintuple and expect a conflict error
	resp = mockStub.MockInvoke("42", methodAndAssetToByte("createCompositeTraintuple", inpTraintuple))
	assert.EqualValues(t, http.StatusConflict, resp.Status)
}

//////////////////////////////////////////////
//
// Composite-specific tests
// Not copied from `traintuple_test.go`
//
/////////////////////////////////////////////

func TestCreateCompositeTraintupleInModels(t *testing.T) {
	testTable := []struct {
		testName         string
		withInHeadModel  bool
		withInTrunkModel bool
		shouldSucceed    bool
		expectedStatus   string
		message          string
	}{
		{
			testName:         "NoHeadAndNoTrunk",
			withInHeadModel:  false,
			withInTrunkModel: false,
			shouldSucceed:    true,
			expectedStatus:   "todo", // no in-models, so we're ready to train
			message:          "One should be able to create a composite traintuple without head or trunk inModels"},
		{
			testName:         "NoHeadAndTrunk",
			withInHeadModel:  true,
			withInTrunkModel: false,
			shouldSucceed:    false,
			message:          "One should NOT be able to create a composite traintuple with a head inModel unless a trunk inModel is also supplied"},
		{
			testName:         "HeadAndNoTrunk",
			withInHeadModel:  false,
			withInTrunkModel: true,
			shouldSucceed:    false,
			message:          "One should NOT be able to create a composite traintuple with a trunk inModel unless a head inModel is also supplied"},
		{
			testName:         "HeadAndTrunk",
			withInHeadModel:  true,
			withInTrunkModel: true,
			shouldSucceed:    true,
			expectedStatus:   "waiting", // waiting for in models to be done before we can start training
			message:          "One should be able to create a composite traintuple with both a head and a trunk inModels"},
	}
	for _, tt := range testTable {
		t.Run(tt.testName, func(t *testing.T) {
			scc := new(SubstraChaincode)
			mockStub := NewMockStubWithRegisterNode("substra", scc)
			registerItem(t, *mockStub, "trainDataset")

			objHash := strings.ReplaceAll(objectiveDescriptionHash, "1", "2")
			inpObjective := inputObjective{DescriptionHash: objHash}
			inpObjective.createDefault()
			inpObjective.TestDataset = inputDataset{}
			resp := mockStub.MockInvoke("42", methodAndAssetToByte("registerObjective", inpObjective))
			assert.EqualValues(t, 200, resp.Status, "when adding objective without dataset it should work: ", resp.Message)

			inpAlgo := inputCompositeAlgo{}
			args := inpAlgo.createDefault()
			resp = mockStub.MockInvoke("42", args)
			assert.EqualValues(t, 200, resp.Status, "when adding algo it should work: ", resp.Message)

			inpTraintuple := inputCompositeTraintuple{ObjectiveKey: objHash}

			if tt.withInHeadModel {
				// create head traintuple
				inpHeadTraintuple := inputCompositeTraintuple{ObjectiveKey: objHash}
				// make the traintuple unique so that it has a unique hash
				inpHeadTraintuple.DataSampleKeys = []string{trainDataSampleHash1}
				args = inpHeadTraintuple.createDefault()
				resp = mockStub.MockInvoke("42", args)
				headTraintuple := outputCompositeTraintuple{}
				json.Unmarshal(resp.Payload, &headTraintuple)

				// make it the head inmodel of inpTraintuple
				inpTraintuple.InHeadModelKey = headTraintuple.Key
			}

			if tt.withInTrunkModel {
				// create trunk traintuple
				inpTrunkTraintuple := inputCompositeTraintuple{ObjectiveKey: objHash}
				// make the traintuple unique so that it has a unique hash
				inpTrunkTraintuple.DataSampleKeys = []string{trainDataSampleHash2}
				args = inpTrunkTraintuple.createDefault()
				resp = mockStub.MockInvoke("42", args)
				trunkTraintuple := outputCompositeTraintuple{}
				json.Unmarshal(resp.Payload, &trunkTraintuple)

				// make it the trunk inmodel of inpTraintuple
				inpTraintuple.InTrunkModelKey = trunkTraintuple.Key
			}

			args = inpTraintuple.createDefault()
			resp = mockStub.MockInvoke("42", args)

			if tt.shouldSucceed {
				assert.EqualValues(t, 200, resp.Status, tt.message+": "+resp.Message)
				traintuple := outputCompositeTraintuple{}
				json.Unmarshal(resp.Payload, &traintuple)
				args = [][]byte{[]byte("queryCompositeTraintuple"), keyToJSON(traintuple.Key)}
				resp = mockStub.MockInvoke("42", args)
				assert.EqualValues(t, 200, resp.Status, "It should find the traintuple without error ", resp.Message)
				traintuple = outputCompositeTraintuple{}
				json.Unmarshal(resp.Payload, &traintuple)
				assert.EqualValues(t, tt.expectedStatus, traintuple.Status, "The traintuple status should be correct")
			} else {
				assert.EqualValues(t, 400, resp.Status, tt.message)
			}
		})
	}
}

func TestCompositeTraintupleInModelTypes(t *testing.T) {
	// Head can only be a composite traintuple's head out model
	allowedHeadTypes := map[AssetType]bool{
		TraintupleType:          false,
		CompositeTraintupleType: true,
		AggregateTupleType:      false,
	}

	// Trunk can be either:
	// - a traintuple's out model
	// - a composite traintuple's head out model
	// - an aggregate tuple's out model
	allowedTrunkTypes := map[AssetType]bool{
		TraintupleType:          true,
		CompositeTraintupleType: true,
		AggregateTupleType:      true,
	}

	for headType, validHeadType := range allowedHeadTypes {
		for trunkType, validTrunkType := range allowedTrunkTypes {
			// Traintuple creation should succeed only if both
			// in-model types are valid
			shouldSucceed := validHeadType && validTrunkType

			successStr := "ShouldSucceed"
			if !shouldSucceed {
				successStr = "ShouldFail"
			}

			testName := fmt.Sprintf("TestTraintuple_%sHeadInModel_%sTrunkInModel_%s", headType, trunkType, successStr)

			t.Run(testName, func(t *testing.T) {
				testCompositeTraintupleInModelTypes(t, headType, trunkType, shouldSucceed)
			})
		}
	}
}

func testCompositeTraintupleInModelTypes(t *testing.T, headType AssetType, trunkType AssetType, shouldSucceed bool) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)
	registerItem(t, *mockStub, "aggregateAlgo")

	inpTraintuple := inputCompositeTraintuple{}

	head, err := registerTraintuple(mockStub, headType)
	assert.NoError(t, err)
	inpTraintuple.InHeadModelKey = head

	trunk, err := registerTraintuple(mockStub, trunkType)
	assert.NoError(t, err)
	inpTraintuple.InTrunkModelKey = trunk

	// create composite traintuple
	inpTraintuple.fillDefaults()
	args := inpTraintuple.getArgs()
	resp := mockStub.MockInvoke("42", args)

	if !shouldSucceed {
		assert.EqualValues(t, 404, resp.Status, "It should NOT be possible to register a traintuple with a %s head and a %s trunk: %s", headType, trunkType, resp.Message)
		return
	}

	assert.EqualValues(t, 200, resp.Status, "It should be possible to register a traintuple with a %s head and a %s trunk: %s", headType, trunkType, resp.Message)
	var keyOnly struct{ Key string }
	json.Unmarshal(resp.Payload, &keyOnly)

	// fetch it back
	args = [][]byte{[]byte("queryCompositeTraintuple"), keyToJSON(keyOnly.Key)}
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status, "It should find the traintuple without error: %s", resp.Message)
	traintuple := outputCompositeTraintuple{}
	json.Unmarshal(resp.Payload, &traintuple)

	require.NotNil(t, traintuple.InHeadModel)
	assert.EqualValues(t, inpTraintuple.InHeadModelKey, traintuple.InHeadModel.TraintupleKey)

	require.NotNil(t, traintuple.InTrunkModel)
	assert.EqualValues(t, inpTraintuple.InTrunkModelKey, traintuple.InTrunkModel.TraintupleKey)
}

func TestCompositeTraintuplePermissions(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)
	registerItem(t, *mockStub, "compositeAlgo")

	inpTraintuple := inputCompositeTraintuple{ObjectiveKey: objectiveDescriptionHash}
	inpTraintuple.fillDefaults()
	// Grant trunk model permissions to no-one
	inpTraintuple.OutTrunkModelPermissions = inputPermissions{Process: inputPermission{Public: false, AuthorizedIDs: []string{}}}
	args := inpTraintuple.getArgs()
	resp := mockStub.MockInvoke("42", args)

	traintuple := outputCompositeTraintuple{}
	json.Unmarshal(resp.Payload, &traintuple)
	args = [][]byte{[]byte("queryCompositeTraintuple"), keyToJSON(traintuple.Key)}
	resp = mockStub.MockInvoke("42", args)
	traintuple = outputCompositeTraintuple{}
	json.Unmarshal(resp.Payload, &traintuple)

	assert.EqualValues(t, false, traintuple.OutHeadModel.Permissions.Process.Public,
		"the head model should not be public")
	assert.EqualValues(t, []string{worker}, traintuple.OutHeadModel.Permissions.Process.AuthorizedIDs,
		"the head model should only be processable by creator")
	assert.EqualValues(t, false, traintuple.OutTrunkModel.Permissions.Process.Public,
		"the trunk model should not be public")
	assert.EqualValues(t, []string{worker}, traintuple.OutHeadModel.Permissions.Process.AuthorizedIDs,
		"if input trunk model permissions are set to 'nobody', this should effectively grant permission to the creator only")
}

func TestCompositeTraintupleLogSuccessFail(t *testing.T) {
	for _, status := range []string{StatusDone, StatusFailed} {
		t.Run("TestCompositeTraintupleLog"+status, func(t *testing.T) {
			scc := new(SubstraChaincode)
			mockStub := NewMockStubWithRegisterNode("substra", scc)
			resp, _ := registerItem(t, *mockStub, "compositeTraintuple")
			var _key struct{ Key string }
			json.Unmarshal(resp.Payload, &_key)
			key := _key.Key

			// start
			resp = mockStub.MockInvoke("42", [][]byte{[]byte("logStartCompositeTrain"), keyToJSON(key)})

			var expectedStatus string

			switch status {
			case StatusDone:
				success := inputLogSuccessCompositeTrain{}
				success.Key = key
				args := success.createDefault()
				resp = mockStub.MockInvoke("42", args)
				require.EqualValuesf(t, 200, resp.Status, "traintuple should be successfully set to 'success': %s", resp.Message)
				expectedStatus = "done"
			case StatusFailed:
				failed := inputLogFailTrain{}
				failed.Key = key
				failed.fillDefaults()
				args := failed.getArgsComposite()
				resp = mockStub.MockInvoke("42", args)
				require.EqualValuesf(t, 200, resp.Status, "traintuple should be successfully set to 'failed': %s", resp.Message)
				expectedStatus = "failed"
			}

			// fetch back
			args := [][]byte{[]byte("queryCompositeTraintuple"), keyToJSON(key)}
			resp = mockStub.MockInvoke("42", args)
			assert.EqualValues(t, 200, resp.Status, "It should find the traintuple without error: %s", resp.Message)
			traintuple := outputCompositeTraintuple{}
			json.Unmarshal(resp.Payload, &traintuple)
			assert.EqualValues(t, expectedStatus, traintuple.Status, "The traintuple status should be set to %s", expectedStatus)
		})
	}
}

// This takes makes sure that, assuming a parent composite traintuple:
// - a child aggregate tuple takes the *trunk* out-model from the parent as its in-model
// - a child composite traintuple takes the *head* out-model from the parent as its head in-model
func TestCorrectParent(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)

	// register parent
	resp, _ := registerItem(t, *mockStub, "compositeTraintuple")
	var _key struct{ Key string }
	json.Unmarshal(resp.Payload, &_key)
	parentKey := _key.Key

	// register aggregate child
	inp1 := inputAggregateTuple{}
	inp1.fillDefaults()
	inp1.InModels = []string{parentKey}
	resp = mockStub.MockInvoke("42", inp1.getArgs())
	json.Unmarshal(resp.Payload, &_key)
	child1Key := _key.Key

	// register composite child
	inp2 := inputCompositeTraintuple{}
	inp2.createDefault()
	inp2.InHeadModelKey = parentKey
	inp2.InTrunkModelKey = traintupleKey
	resp = mockStub.MockInvoke("42", inp2.getArgs())
	json.Unmarshal(resp.Payload, &_key)
	child2Key := _key.Key

	// start
	mockStub.MockInvoke("42", [][]byte{[]byte("logStartCompositeTrain"), keyToJSON(parentKey)})
	// success
	success := inputLogSuccessCompositeTrain{}
	success.Key = parentKey
	args := success.createDefault()
	mockStub.MockInvoke("42", args)

	mockStub.MockTransactionStart("42")
	db := NewLedgerDB(mockStub)

	// fetch aggregate child, and check its in-model is the parent's trunk out-model
	child1, _ := queryAggregateTuple(db, assetToArgs(inputHash{Key: child1Key}))
	assert.Equal(t, trunkModelHash, child1.InModels[0].Hash)

	// fetch composite child, and check its head in-model is the parent's head out-model
	child2, _ := queryCompositeTraintuple(db, assetToArgs(inputHash{Key: child2Key}))
	assert.Equal(t, headModelHash, child2.InHeadModel.Hash)
}
