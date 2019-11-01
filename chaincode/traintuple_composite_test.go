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
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var createInModelTests = []struct {
	withInModelHead  bool
	withInModelTrunk bool
	shouldSucceed    bool
	message          string
}{
	{
		withInModelHead:  false,
		withInModelTrunk: false,
		shouldSucceed:    true,
		message:          "One should be able to create a composite traintuple without head or trunk inModels"},
	{
		withInModelHead:  true,
		withInModelTrunk: false,
		shouldSucceed:    false,
		message:          "One should NOT be able to create a composite traintuple with a head inModel unless a trunk inModel is also supplied"},
	{
		withInModelHead:  false,
		withInModelTrunk: true,
		shouldSucceed:    false,
		message:          "One should NOT be able to create a composite traintuple with a trunk inModel unless a head inModel is also supplied"},
	{
		withInModelHead:  true,
		withInModelTrunk: true,
		shouldSucceed:    true,
		message:          "One should be able to create a composite traintuple with both a head and a trunk inModels"}}

// TODO: give this function a more accurate name
func TestTraintupleWithNoTestDatasetComposite(t *testing.T) {
	for _, tt := range createInModelTests {
		scc := new(SubstraChaincode)
		mockStub := NewMockStubWithRegisterNode("substra", scc)
		registerItem(t, *mockStub, "trainDataset")

		objHash := strings.ReplaceAll(objectiveDescriptionHash, "1", "2")
		inpObjective := inputObjective{DescriptionHash: objHash}
		inpObjective.createDefault()
		inpObjective.TestDataset = inputDataset{}
		resp := mockStub.MockInvoke("42", methodAndAssetToByte("registerObjective", inpObjective))
		assert.EqualValues(t, 200, resp.Status, "when adding objective without dataset it should work: ", resp.Message)

		inpAlgo := inputAlgo{}
		args := inpAlgo.createDefault()
		resp = mockStub.MockInvoke("42", args)
		assert.EqualValues(t, 200, resp.Status, "when adding algo it should work: ", resp.Message)

		inpTraintuple := inputCompositeTraintuple{ObjectiveKey: objHash}

		if tt.withInModelHead {
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

		if tt.withInModelTrunk {
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
		} else {
			assert.EqualValues(t, 400, resp.Status, tt.message)
		}
	}
}

func registerTraintuple(mockStub *MockStub, assetType AssetType, dataSampleKeys []string) (key string, err error) {
	switch assetType {
	case CompositeTraintupleType:
		inpTraintuple := inputCompositeTraintuple{}
		inpTraintuple.DataSampleKeys = dataSampleKeys
		inpTraintuple.fillDefaults()
		args := inpTraintuple.getArgs()
		resp := mockStub.MockInvoke("42", args)
		if resp.Status != 200 {
			err = fmt.Errorf("Failed to register traintuple: %s", resp.Message)
			return
		}
		var _key struct{ Key string }
		json.Unmarshal(resp.Payload, &_key)
		return _key.Key, nil
	case TraintupleType:
		inpTraintuple := inputTraintuple{}
		inpTraintuple.DataSampleKeys = dataSampleKeys
		args := inpTraintuple.createDefault()
		resp := mockStub.MockInvoke("42", args)
		if resp.Status != 200 {
			err = fmt.Errorf("Failed to register traintuple: %s", resp.Message)
			return
		}
		var _key struct{ Key string }
		json.Unmarshal(resp.Payload, &_key)
		return _key.Key, nil
	default:
		err = fmt.Errorf("Invalid asset type: %v", assetType)
		return
	}
}

// TODO: add more tests:
// head:
// - regular head
// - composite head
// trunk:
// - regular trunk
// - composite trunk
// - aggregate trunk
func TestTraintupleRegularHead(t *testing.T) {
	message := "It should be possible to registering a traintuple with a regular traintuple as a head"

	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)
	registerItem(t, *mockStub, "trainDataset")

	objHash := strings.ReplaceAll(objectiveDescriptionHash, "1", "2")
	inpObjective := inputObjective{DescriptionHash: objHash}
	inpObjective.createDefault()
	inpObjective.TestDataset = inputDataset{}
	resp := mockStub.MockInvoke("42", methodAndAssetToByte("registerObjective", inpObjective))

	inpAlgo := inputAlgo{}
	args := inpAlgo.createDefault()
	resp = mockStub.MockInvoke("42", args)

	inpTraintuple := inputCompositeTraintuple{ObjectiveKey: objHash}

	head, err := registerTraintuple(mockStub, TraintupleType, []string{trainDataSampleHash1})
	assert.NoError(t, err)
	inpTraintuple.InHeadModelKey = head

	trunk, err := registerTraintuple(mockStub, TraintupleType, []string{trainDataSampleHash2})
	assert.NoError(t, err)
	inpTraintuple.InTrunkModelKey = trunk

	// create composite traintuple
	inpTraintuple.fillDefaults()
	args = inpTraintuple.getArgs()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status, message+": ", resp.Message)
	var keyOnly struct{ Key string }
	json.Unmarshal(resp.Payload, &keyOnly)

	// fetch it back
	args = [][]byte{[]byte("queryCompositeTraintuple"), keyToJSON(keyOnly.Key)}
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status, "It should find the traintuple without error: %s", resp.Message)
	traintuple := outputCompositeTraintuple{}
	json.Unmarshal(resp.Payload, &traintuple)

	assert.EqualValues(t, inpTraintuple.InHeadModelKey, traintuple.InModelHead.TraintupleKey)
	// TODO: test that fields are populated once parent traintuple is "done"
	// assert.NotEqual(t, "", traintuple.InModelHead.Hash)
	// assert.NotEqual(t, "", traintuple.InModelHead.StorageAddress)

	assert.EqualValues(t, inpTraintuple.InTrunkModelKey, traintuple.InModelTrunk.TraintupleKey)
	// TODO: test that fields are populated once parent traintuple is "done"
	// assert.NotEqual(t, "", traintuple.InModelTrunk.Hash)
	// assert.NotEqual(t, "", traintuple.InModelTrunk.StorageAddress)
}

func TestTraintuplePermissions(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)
	registerItem(t, *mockStub, "trainDataset")

	objHash := strings.ReplaceAll(objectiveDescriptionHash, "1", "2")
	inpObjective := inputObjective{DescriptionHash: objHash}
	inpObjective.createDefault()
	inpObjective.TestDataset = inputDataset{}
	resp := mockStub.MockInvoke("42", methodAndAssetToByte("registerObjective", inpObjective))

	inpAlgo := inputAlgo{}
	args := inpAlgo.createDefault()
	resp = mockStub.MockInvoke("42", args)

	inpTraintuple := inputCompositeTraintuple{ObjectiveKey: objHash}
	inpTraintuple.fillDefaults()
	// Grant trunk model permissions to no-one
	inpTraintuple.InTrunkModelPermission = inputPermissions{Process: inputPermission{Public: false, AuthorizedIDs: []string{}}}
	args = inpTraintuple.getArgs()
	resp = mockStub.MockInvoke("42", args)

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

	inpAlgo := inputAlgo{}
	args := inpAlgo.createDefault()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status, "when adding algo it should work: ", resp.Message)

	inpTraintuple := inputCompositeTraintuple{
		ObjectiveKey:   objHash,
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

// func TestTraintupleWithDuplicatedDatasamplesComposite(t *testing.T) {
// 	scc := new(SubstraChaincode)
// 	mockStub := NewMockStub("substra", scc)
// 	registerItem(t, *mockStub, "trainDataset")

// 	objHash := strings.ReplaceAll(objectiveDescriptionHash, "1", "2")
// 	inpObjective := inputObjective{DescriptionHash: objHash}
// 	inpObjective.createDefault()
// 	inpObjective.TestDataset = inputDataset{}
// 	resp := mockStub.MockInvoke("42", methodAndAssetToByte("registerObjective", inpObjective))
// 	assert.EqualValues(t, 200, resp.Status, "when adding objective without dataset it should work: ", resp.Message)

// 	inpAlgo := inputAlgo{}
// 	args := inpAlgo.createDefault()
// 	resp = mockStub.MockInvoke("42", args)
// 	assert.EqualValues(t, 200, resp.Status, "when adding algo it should work: ", resp.Message)

// 	inpTraintuple := inputCompositeTraintuple{
// 		ObjectiveKey:   objHash,
// 		DataSampleKeys: []string{trainDataSampleHash1, trainDataSampleHash2, trainDataSampleHash1},
// 	}
// 	args = inpTraintuple.createDefault()
// 	resp = mockStub.MockInvoke("42", args)
// 	assert.EqualValuesf(t, 400, resp.Status, "when adding traintuple with a duplicated data samples it should not work: %s", resp.Message)
// }

// // func TestTagTuple(t *testing.T) {
// // 	scc := new(SubstraChaincode)
// // 	mockStub := NewMockStubWithRegisterNode("substra", scc)

// // 	registerItem(t, *mockStub, "algo")

// // 	noTag := "This is not a tag because it's waaaaaaaaaaaaaaaayyyyyyyyyyyyyyyyyyyyyyy too long."

// // 	inpTraintuple := inputCompositeTraintuple{Tag: noTag}
// // 	args := inpTraintuple.createDefault()
// // 	resp := mockStub.MockInvoke("42", args)
// // 	assert.EqualValues(t, 400, resp.Status, resp.Message)

// // 	tag := "This is a tag"

// // 	inpTraintuple = inputCompositeTraintuple{Tag: tag}
// // 	args = inpTraintuple.createDefault()
// // 	resp = mockStub.MockInvoke("42", args)
// // 	assert.EqualValues(t, 200, resp.Status, resp.Message)

// // 	args = [][]byte{[]byte("queryTraintuples")}
// // 	resp = mockStub.MockInvoke("42", args)

// // 	traintuples := []outputCompositeTraintuple{}
// // 	err := json.Unmarshal(resp.Payload, &traintuples)

// // 	assert.NoError(t, err, "should be unmarshaled")
// // 	assert.Len(t, traintuples, 1, "there should be one traintuple")
// // 	assert.EqualValues(t, tag, traintuples[0].Tag)

// // 	inpTesttuple := inputTesttuple{Tag: tag}
// // 	args = inpTesttuple.createDefault()
// // 	resp = mockStub.MockInvoke("42", args)
// // 	assert.EqualValues(t, 200, resp.Status, resp.Message)

// // 	args = [][]byte{[]byte("queryTesttuples")}
// // 	resp = mockStub.MockInvoke("42", args)
// // 	testtuples := []outputTesttuple{}
// // 	err = json.Unmarshal(resp.Payload, &testtuples)
// // 	assert.NoError(t, err, "should be unmarshaled")
// // 	assert.Len(t, testtuples, 1, "there should be one traintuple")
// // 	assert.EqualValues(t, tag, testtuples[0].Tag)

// // 	filter := inputQueryFilter{
// // 		IndexName:  "testtuple~tag",
// // 		Attributes: tag,
// // 	}
// // 	args = [][]byte{[]byte("queryFilter"), assetToJSON(filter)}
// // 	resp = mockStub.MockInvoke("42", args)
// // 	assert.EqualValues(t, 200, resp.Status, resp.Message)
// // 	filtertuples := []outputTesttuple{}
// // 	err = json.Unmarshal(resp.Payload, &filtertuples)
// // 	assert.NoError(t, err, "should be unmarshaled")
// // 	assert.Len(t, testtuples, 1, "there should be one traintuple")
// // 	assert.EqualValues(t, tag, testtuples[0].Tag)

// // }
// func TestNoPanicWhileQueryingIncompleteCompositeTraintuple(t *testing.T) {
// 	scc := new(SubstraChaincode)
// 	mockStub := NewMockStubWithRegisterNode("substra", scc)
// 	// Add a some dataManager, dataSample and traintuple
// 	registerItem(t, *mockStub, "traintuple")

// 	// Manually open a ledger transaction
// 	mockStub.MockTransactionStart("42")
// 	defer mockStub.MockTransactionEnd("42")

// 	// Retreive and alter existing objectif to pass Metrics at nil
// 	db := NewLedgerDB(mockStub)
// 	objective, err := db.GetObjective(objectiveDescriptionHash)
// 	assert.NoError(t, err)
// 	objective.Metrics = nil
// 	objBytes, err := json.Marshal(objective)
// 	assert.NoError(t, err)
// 	err = mockStub.PutState(objectiveDescriptionHash, objBytes)
// 	assert.NoError(t, err)
// 	// It should not panic
// 	require.NotPanics(t, func() {
// 		getOutputTraintuple(NewLedgerDB(mockStub), traintupleKey)
// 	})
// }
// func TestTraintupleComputePlanCreationComposite(t *testing.T) {
// 	scc := new(SubstraChaincode)
// 	mockStub := NewMockStubWithRegisterNode("substra", scc)

// 	// Add dataManager, dataSample and algo
// 	registerItem(t, *mockStub, "algo")

// 	inpTraintuple := inputCompositeTraintuple{ComputePlanID: "someComputePlanID"}
// 	args := inpTraintuple.createDefault()
// 	resp := mockStub.MockInvoke("42", args)
// 	require.EqualValues(t, 400, resp.Status, "should failed for missing rank")
// 	require.Contains(t, resp.Message, "invalid inputs, a ComputePlan should have a rank", "invalid error message")

// 	inpTraintuple = inputCompositeTraintuple{Rank: "1"}
// 	args = inpTraintuple.createDefault()
// 	resp = mockStub.MockInvoke("42", args)
// 	require.EqualValues(t, 400, resp.Status, "should failed for invalid rank")
// 	require.Contains(t, resp.Message, "invalid inputs, a new ComputePlan should have a rank 0")

// 	inpTraintuple = inputCompositeTraintuple{Rank: "0"}
// 	args = inpTraintuple.createDefault()
// 	resp = mockStub.MockInvoke("42", args)
// 	assert.EqualValues(t, 200, resp.Status)
// 	res := map[string]string{}
// 	err := json.Unmarshal(resp.Payload, &res)
// 	assert.NoError(t, err, "should unmarshal without problem")
// 	assert.Contains(t, res, "key")
// 	key := res["key"]
// 	require.EqualValues(t, key, traintupleKey)

// 	inpTraintuple = inputCompositeTraintuple{Rank: "0"}
// 	args = inpTraintuple.createDefault()
// 	resp = mockStub.MockInvoke("42", args)
// 	require.EqualValues(t, 409, resp.Status, "should failed for existing ComputePlanID")
// 	require.Contains(t, resp.Message, "already exists")

// 	require.EqualValues(t, 409, resp.Status, "should failed for existing FLTask")
// 	errorPayload := map[string]interface{}{}
// 	err = json.Unmarshal(resp.Payload, &errorPayload)
// 	assert.NoError(t, err, "should unmarshal without problem")
// 	require.Contains(t, errorPayload, "key", "key should be available in payload")
// 	assert.EqualValues(t, traintupleKey, errorPayload["key"], "key in error should be traintupleKey")
// }

// func TestTraintupleMultipleCommputePlanCreationsComposite(t *testing.T) {
// 	scc := new(SubstraChaincode)
// 	mockStub := NewMockStubWithRegisterNode("substra", scc)

// 	// Add a some dataManager, dataSample and traintuple
// 	registerItem(t, *mockStub, "algo")

// 	inpTraintuple := inputCompositeTraintuple{Rank: "0"}
// 	args := inpTraintuple.createDefault()
// 	resp := mockStub.MockInvoke("42", args)
// 	assert.EqualValues(t, 200, resp.Status)
// 	res := map[string]string{}
// 	err := json.Unmarshal(resp.Payload, &res)
// 	assert.NoError(t, err, "should unmarshal without problem")
// 	assert.Contains(t, res, "key")
// 	key := res["key"]
// 	// Failed to add a traintuple with the same rank
// 	inpTraintuple = inputCompositeTraintuple{
// 		InModels:      []string{key},
// 		Rank:          "0",
// 		ComputePlanID: key}
// 	args = inpTraintuple.createDefault()
// 	resp = mockStub.MockInvoke("42", args)
// 	assert.EqualValues(t, 400, resp.Status, resp.Message, "should failed to add a traintuple of the same rank")

// 	// Failed to add a traintuple to an unexisting CommputePlan
// 	inpTraintuple = inputCompositeTraintuple{
// 		InModels:      []string{key},
// 		Rank:          "1",
// 		ComputePlanID: "notarealone"}
// 	args = inpTraintuple.createDefault()
// 	resp = mockStub.MockInvoke("42", args)
// 	assert.EqualValues(t, 400, resp.Status, resp.Message, "should failed to add a traintuple to an unexisting ComputePlanID")

// 	// Succesfully add a traintuple to the same ComputePlanID
// 	inpTraintuple = inputCompositeTraintuple{
// 		InModels:      []string{key},
// 		Rank:          "1",
// 		ComputePlanID: key}
// 	args = inpTraintuple.createDefault()
// 	resp = mockStub.MockInvoke("42", args)
// 	assert.EqualValues(t, 200, resp.Status, resp.Message, "should be able do create a traintuple with the same ComputePlanID")
// 	err = json.Unmarshal(resp.Payload, &res)
// 	assert.NoError(t, err, "should unmarshal without problem")
// 	assert.Contains(t, res, "key")
// 	ttkey := res["key"]
// 	// Add new algo to check all ComputePlan algo consistency
// 	newAlgoHash := strings.Replace(algoHash, "a", "b", 1)
// 	inpAlgo := inputAlgo{Hash: newAlgoHash}
// 	args = inpAlgo.createDefault()
// 	resp = mockStub.MockInvoke("42", args)
// 	assert.EqualValues(t, 200, resp.Status)

// 	inpTraintuple = inputCompositeTraintuple{
// 		AlgoKey:       newAlgoHash,
// 		InModels:      []string{ttkey},
// 		Rank:          "2",
// 		ComputePlanID: key}
// 	args = inpTraintuple.createDefault()
// 	resp = mockStub.MockInvoke("42", args)
// 	assert.EqualValues(t, 400, resp.Status, resp.Message, "sould fail for it doesn't have the same algo key")
// 	assert.Contains(t, resp.Message, "does not have the same algo key")
// }

// // func TestTesttupleOnFailedTraintuple(t *testing.T) {
// // 	scc := new(SubstraChaincode)
// // 	mockStub := NewMockStubWithRegisterNode("substra", scc)

// // 	// Add a some dataManager, dataSample and traintuple
// // 	resp, _ := registerItem(t, *mockStub, "traintuple")

// // 	res := map[string]string{}
// // 	err := json.Unmarshal(resp.Payload, &res)
// // 	assert.NoError(t, err, "should unmarshal without problem")
// // 	assert.Contains(t, res, "key")
// // 	traintupleKey := res["key"]

// // 	// Mark the traintuple as failed
// // 	fail := inputLogFailTrain{}
// // 	fail.Key = traintupleKey
// // 	args := fail.createDefault()
// // 	resp = mockStub.MockInvoke("42", args)
// // 	assert.EqualValues(t, 200, resp.Status, "should be able to log traintuple as failed")

// // 	// Fail to add a testtuple to this failed traintuple
// // 	inpTesttuple := inputTesttuple{}
// // 	args = inpTesttuple.createDefault()
// // 	resp = mockStub.MockInvoke("42", args)
// // 	assert.EqualValues(t, 400, resp.Status, "status should show an error since the traintuple is failed")
// // 	assert.Contains(t, resp.Message, "could not register this testtuple")
// // }

// // func TestCertifiedExplicitTesttuple(t *testing.T) {
// // 	scc := new(SubstraChaincode)
// // 	mockStub := NewMockStubWithRegisterNode("substra", scc)

// // 	// Add a some dataManager, dataSample and traintuple
// // 	registerItem(t, *mockStub, "traintuple")

// // 	// Add a testtuple that shoulb be certified since it's the same dataManager and
// // 	// dataSample than the objective but explicitly pass as arguments and in disorder
// // 	inpTesttuple := inputTesttuple{
// // 		DataSampleKeys: []string{testDataSampleHash2, testDataSampleHash1},
// // 		DataManagerKey: dataManagerOpenerHash}
// // 	args := inpTesttuple.createDefault()
// // 	resp := mockStub.MockInvoke("42", args)
// // 	assert.EqualValues(t, 200, resp.Status)

// // 	args = [][]byte{[]byte("queryTesttuples")}
// // 	resp = mockStub.MockInvoke("42", args)
// // 	testtuples := [](map[string]interface{}){}
// // 	err := json.Unmarshal(resp.Payload, &testtuples)
// // 	assert.NoError(t, err, "should be unmarshaled")
// // 	assert.Len(t, testtuples, 1, "there should be only one testtuple...")
// // 	assert.True(t, testtuples[0]["certified"].(bool), "... and it should be certified")

// // }
// // func TestConflictCertifiedNonCertifiedTesttuple(t *testing.T) {
// // 	scc := new(SubstraChaincode)
// // 	mockStub := NewMockStubWithRegisterNode("substra", scc)

// // 	// Add a some dataManager, dataSample and traintuple
// // 	registerItem(t, *mockStub, "traintuple")

// // 	// Add a certified testtuple
// // 	inpTesttuple1 := inputTesttuple{}
// // 	args := inpTesttuple1.createDefault()
// // 	resp := mockStub.MockInvoke("42", args)
// // 	assert.EqualValues(t, 200, resp.Status)

// // 	// Fail to add an incomplete uncertified testtuple
// // 	inpTesttuple2 := inputTesttuple{DataSampleKeys: []string{trainDataSampleHash1}}
// // 	args = inpTesttuple2.createDefault()
// // 	resp = mockStub.MockInvoke("42", args)
// // 	assert.EqualValues(t, 400, resp.Status)
// // 	assert.Contains(t, resp.Message, "invalid input: dataManagerKey and dataSampleKey should be provided together")

// // 	// Add an uncertified testtuple successfully
// // 	inpTesttuple3 := inputTesttuple{
// // 		DataSampleKeys: []string{trainDataSampleHash1, trainDataSampleHash2},
// // 		DataManagerKey: dataManagerOpenerHash}
// // 	args = inpTesttuple3.createDefault()
// // 	resp = mockStub.MockInvoke("42", args)
// // 	assert.EqualValues(t, 200, resp.Status)

// // 	// Fail to add the same testtuple with a different order for dataSampleKeys
// // 	inpTesttuple4 := inputTesttuple{
// // 		DataSampleKeys: []string{trainDataSampleHash2, trainDataSampleHash1},
// // 		DataManagerKey: dataManagerOpenerHash}
// // 	args = inpTesttuple4.createDefault()
// // 	resp = mockStub.MockInvoke("42", args)
// // 	assert.EqualValues(t, 409, resp.Status)
// // 	assert.Contains(t, resp.Message, "already exists")
// // }

// func TestCompositeTraintuple(t *testing.T) {
// 	scc := new(SubstraChaincode)
// 	mockStub := NewMockStubWithRegisterNode("substra", scc)

// 	// Add traintuple with invalid field
// 	inpTraintuple := inputCompositeTraintuple{
// 		AlgoKey: "aaa",
// 	}
// 	args := inpTraintuple.createDefault()
// 	resp := mockStub.MockInvoke("42", args)
// 	assert.EqualValuesf(t, 400, resp.Status, "when adding objective with invalid hash, status %d and message %s", resp.Status, resp.Message)

// 	// Add traintuple with unexisting algo
// 	inpTraintuple = inputCompositeTraintuple{}
// 	args = inpTraintuple.createDefault()
// 	resp = mockStub.MockInvoke("42", args)
// 	assert.EqualValuesf(t, 400, resp.Status, "when adding traintuple with unexisting algo, status %d and message %s", resp.Status, resp.Message)

// 	// Properly add traintuple
// 	resp, tt := registerItem(t, *mockStub, "traintuple")

// 	inpTraintuple = tt.(inputCompositeTraintuple)
// 	res := map[string]string{}
// 	err := json.Unmarshal(resp.Payload, &res)
// 	assert.NoError(t, err, "traintuple should unmarshal without problem")
// 	assert.Contains(t, res, "key")
// 	traintupleKey := res["key"]
// 	// Query traintuple from key and check the consistency of returned arguments
// 	args = [][]byte{[]byte("queryTraintuple"), keyToJSON(traintupleKey)}
// 	resp = mockStub.MockInvoke("42", args)
// 	assert.EqualValuesf(t, 200, resp.Status, "when querying the traintuple - status %d and message %s", resp.Status, resp.Message)
// 	out := outputCompositeTraintuple{}
// 	err = json.Unmarshal(resp.Payload, &out)
// 	assert.NoError(t, err, "when unmarshalling queried traintuple")
// 	expected := outputCompositeTraintuple{
// 		Key: traintupleKey,
// 		Algo: &HashDressName{
// 			Hash:           algoHash,
// 			Name:           algoName,
// 			StorageAddress: algoStorageAddress,
// 		},
// 		Creator: worker,
// 		Dataset: &TtDataset{
// 			DataSampleKeys: []string{trainDataSampleHash1, trainDataSampleHash2},
// 			OpenerHash:     dataManagerOpenerHash,
// 			Perf:           0.0,
// 			Worker:         worker,
// 		},
// 		Objective: &TtObjective{
// 			Key: objectiveDescriptionHash,
// 			Metrics: &HashDress{
// 				Hash:           objectiveMetricsHash,
// 				StorageAddress: objectiveMetricsStorageAddress,
// 			},
// 		},
// 		Permissions: outputPermissions{
// 			Process: Permission{Public: true, AuthorizedIDs: []string{}},
// 		},
// 		Status: StatusTodo,
// 	}
// 	assert.Exactly(t, expected, out, "the traintuple queried from the ledger differ from expected")

// 	// Query all traintuples and check consistency
// 	args = [][]byte{[]byte("queryTraintuples")}
// 	resp = mockStub.MockInvoke("42", args)
// 	assert.EqualValuesf(t, 200, resp.Status, "when querying traintuples - status %d and message %s", resp.Status, resp.Message)
// 	// TODO add traintuple key to output struct
// 	// For now we test it as cleanly as its added to the query response
// 	assert.Contains(t, string(resp.Payload), "key\":\""+traintupleKey)
// 	var queryTraintuples []outputCompositeTraintuple
// 	err = json.Unmarshal(resp.Payload, &queryTraintuples)
// 	assert.NoError(t, err, "traintuples should unmarshal without problem")
// 	assert.Exactly(t, out, queryTraintuples[0])

// 	// Add traintuple with inmodel from the above-submitted traintuple
// 	inpWaitingTraintuple := inputCompositeTraintuple{
// 		InModels: []string{string(traintupleKey)},
// 	}
// 	args = inpWaitingTraintuple.createDefault()
// 	resp = mockStub.MockInvoke("42", args)
// 	assert.EqualValuesf(t, 200, resp.Status, "when adding traintuple with status %d and message %s", resp.Status, resp.Message)
// 	//waitingTraintupleKey := string(resp.Payload)

// 	// Query traintuple with status todo and worker as trainworker and check consistency
// 	filter := inputQueryFilter{
// 		IndexName:  "traintuple~worker~status",
// 		Attributes: worker + ", todo",
// 	}
// 	args = [][]byte{[]byte("queryFilter"), assetToJSON(filter)}
// 	resp = mockStub.MockInvoke("42", args)
// 	assert.EqualValuesf(t, 200, resp.Status, "when querying traintuple of worker with todo status - status %d and message %s", resp.Status, resp.Message)
// 	err = json.Unmarshal(resp.Payload, &queryTraintuples)
// 	assert.NoError(t, err, "traintuples should unmarshal without problem")
// 	assert.Exactly(t, out, queryTraintuples[0])

// 	// Update status and check consistency
// 	success := inputLogSuccessTrain{}
// 	success.Key = traintupleKey

// 	argsSlice := [][][]byte{
// 		[][]byte{[]byte("logStartTrain"), keyToJSON(traintupleKey)},
// 		success.createDefault(),
// 	}
// 	traintupleStatus := []string{StatusDoing, StatusDone}
// 	for i := range traintupleStatus {
// 		resp = mockStub.MockInvoke("42", argsSlice[i])
// 		require.EqualValuesf(t, 200, resp.Status, "when logging start %s with message %s", traintupleStatus[i], resp.Message)
// 		filter := inputQueryFilter{
// 			IndexName:  "traintuple~worker~status",
// 			Attributes: worker + ", " + traintupleStatus[i],
// 		}
// 		args = [][]byte{[]byte("queryFilter"), assetToJSON(filter)}
// 		resp = mockStub.MockInvoke("42", args)
// 		assert.EqualValuesf(t, 200, resp.Status, "when querying traintuple of worker with %s status - message %s", traintupleStatus[i], resp.Message)
// 		sPayload := make([]map[string]interface{}, 1)
// 		assert.NoError(t, json.Unmarshal(resp.Payload, &sPayload), "when unmarshal queried traintuples")
// 		assert.EqualValues(t, traintupleKey, sPayload[0]["key"], "wrong retrieved key when querying traintuple of worker with %s status ", traintupleStatus[i])
// 		assert.EqualValues(t, traintupleStatus[i], sPayload[0]["status"], "wrong retrieved status when querying traintuple of worker with %s status ", traintupleStatus[i])
// 	}

// 	// Query CompositeTraintuple From key
// 	args = [][]byte{[]byte("queryTraintuple"), keyToJSON(traintupleKey)}
// 	resp = mockStub.MockInvoke("42", args)
// 	assert.EqualValuesf(t, 200, resp.Status, "when querying traintuple with status %d and message %s", resp.Status, resp.Message)
// 	endTraintuple := outputCompositeTraintuple{}
// 	assert.NoError(t, json.Unmarshal(resp.Payload, &endTraintuple))
// 	expected.Dataset.Perf = success.Perf
// 	expected.Log = success.Log
// 	expected.OutModel = &HashDress{
// 		Hash:           modelHash,
// 		StorageAddress: modelAddress}
// 	expected.Status = traintupleStatus[1]
// 	assert.Exactly(t, expected, endTraintuple, "retreived CompositeTraintuple does not correspond to what is expected")

// 	// query all traintuples related to a traintuple with the same algo
// 	args = [][]byte{[]byte("queryModelDetails"), keyToJSON(traintupleKey)}
// 	resp = mockStub.MockInvoke("42", args)
// 	assert.EqualValuesf(t, 200, resp.Status, "when querying model details with status %d and message %s", resp.Status, resp.Message)
// 	payload := outputModelDetails{}
// 	assert.NoError(t, json.Unmarshal(resp.Payload, &payload))
// 	assert.NotNil(t, payload.CompositeTraintuple, "when querying model tuples, payload should contain one traintuple")

// 	// query all traintuples related to a traintuple with the same algo
// 	args = [][]byte{[]byte("queryModels")}
// 	resp = mockStub.MockInvoke("42", args)
// 	assert.EqualValuesf(t, 200, resp.Status, "when querying models with status %d and message %s", resp.Status, resp.Message)
// }

// func TestQueryTraintupleNotFoundComposite(t *testing.T) {
// 	scc := new(SubstraChaincode)
// 	mockStub := NewMockStubWithRegisterNode("substra", scc)
// 	registerItem(t, *mockStub, "traintuple")

// 	// queryTraintuple: normal case
// 	args := [][]byte{[]byte("queryTraintuple"), keyToJSON(traintupleKey)}
// 	resp := mockStub.MockInvoke("42", args)
// 	assert.EqualValuesf(t, 200, resp.Status, "when querying the traintuple - status %d and message %s", resp.Status, resp.Message)

// 	// queryTraintuple: key does not exist
// 	notFoundKey := "eedbb7c31f62244c0f34461cc168804227115793d01c270021fe3f7935482eed"
// 	args = [][]byte{[]byte("queryTraintuple"), keyToJSON(notFoundKey)}
// 	resp = mockStub.MockInvoke("42", args)
// 	assert.EqualValuesf(t, 404, resp.Status, "when querying the traintuple - status %d and message %s", resp.Status, resp.Message)

// 	// queryTraintuple: key does not exist and use existing other asset type key
// 	args = [][]byte{[]byte("queryTraintuple"), keyToJSON(algoHash)}
// 	resp = mockStub.MockInvoke("42", args)
// 	assert.EqualValuesf(t, 404, resp.Status, "when querying the traintuple - status %d and message %s", resp.Status, resp.Message)
// }

// func TestInsertTraintupleTwiceComposite(t *testing.T) {
// 	scc := new(SubstraChaincode)
// 	mockStub := NewMockStubWithRegisterNode("substra", scc)
// 	registerItem(t, *mockStub, "algo")

// 	// create a traintuple and start a ComplutePlan
// 	inpTraintuple := inputCompositeTraintuple{
// 		Rank: "0",
// 	}
// 	inpTraintuple.createDefault()
// 	resp := mockStub.MockInvoke("42", methodAndAssetToByte("createTraintuple", inpTraintuple))
// 	assert.EqualValues(t, http.StatusOK, resp.Status)

// 	// create a second traintuple in the same ComputePlan
// 	inpTraintuple.Rank = "1"
// 	inpTraintuple.ComputePlanID = traintupleKey
// 	inpTraintuple.InModels = []string{traintupleKey}
// 	resp = mockStub.MockInvoke("42", methodAndAssetToByte("createTraintuple", inpTraintuple))
// 	assert.EqualValues(t, http.StatusOK, resp.Status)

// 	// re-insert the same traintuple and expect a conflict error
// 	resp = mockStub.MockInvoke("42", methodAndAssetToByte("createTraintuple", inpTraintuple))
// 	assert.EqualValues(t, http.StatusConflict, resp.Status)

// }

// func TestRecursiveLogFailedComposite(t *testing.T) {
// 	scc := new(SubstraChaincode)
// 	mockStub := NewMockStub("substra", scc)
// 	mockStub.MockTransactionStart("42")
// 	registerItem(t, *mockStub, "traintuple")
// 	db := NewLedgerDB(mockStub)

// 	childtraintuple := inputTraintuple{}
// 	childtraintuple.createDefault()
// 	childtraintuple.InModels = []string{traintupleKey}
// 	childResp, err := createTraintuple(db, assetToArgs(childtraintuple))
// 	assert.NoError(t, err)

// 	grandChildtraintuple := inputTraintuple{}
// 	grandChildtraintuple.createDefault()
// 	grandChildtraintuple.InModels = []string{childResp["key"]}
// 	grandChildresp, err := createTraintuple(db, assetToArgs(grandChildtraintuple))
// 	assert.NoError(t, err)

// 	grandChildtesttuple := inputTesttuple{TraintupleKey: traintupleKey}
// 	testResp, err := createTesttuple(db, assetToArgs(grandChildtesttuple))
// 	assert.NoError(t, err)

// 	_, err = logStartTrain(db, assetToArgs(inputHash{Key: traintupleKey}))
// 	assert.NoError(t, err)
// 	_, err = logFailTrain(db, assetToArgs(inputHash{Key: traintupleKey}))
// 	assert.NoError(t, err)

// 	train2, err := db.GetTraintuple(grandChildresp["key"])
// 	assert.NoError(t, err)
// 	assert.Equal(t, StatusFailed, train2.Status)

// 	test, err := db.GetTesttuple(testResp["key"])
// 	assert.NoError(t, err)
// 	assert.Equal(t, StatusFailed, test.Status)
// }
