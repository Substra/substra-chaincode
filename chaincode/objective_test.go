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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLeaderBoard(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)
	db := NewLedgerDB(mockStub)
	registerItem(t, *mockStub, "")
	mockStub.MockTransactionStart("42")

	// Add a certified testtuple
	inputTest := inputTesttuple{
		TraintupleKey: traintupleKey,
	}
	keyMap, err := createTesttuple(db, assetToArgs(inputTest))
	assert.NoError(t, err)

	inpLeaderboard := inputLeaderboard{
		ObjectiveKey:   objectiveDescriptionHash,
		AscendingOrder: true,
	}
	// leaderboard should be empty since there is no testtuple done
	leaderboard, err := queryObjectiveLeaderboard(db, assetToArgs(inpLeaderboard))
	assert.NoError(t, err)
	assert.Len(t, leaderboard.Testtuples, 0)

	// Update testtuple status directly
	testtuple, err := db.GetTesttuple(keyMap["key"])
	assert.NoError(t, err)
	testtuple.Status = StatusDone
	testtuple.Dataset.Perf = 0.9
	err = db.Put(keyMap["key"], testtuple)
	assert.NoError(t, err)

	leaderboard, err = queryObjectiveLeaderboard(db, assetToArgs(inpLeaderboard))
	assert.NoError(t, err)
	assert.Equal(t, objectiveDescriptionHash, leaderboard.Objective.Key)
	require.Len(t, leaderboard.Testtuples, 1)
	assert.Equal(t, keyMap["key"], leaderboard.Testtuples[0].Key)
	assert.Equal(t, traintupleKey, leaderboard.Testtuples[0].TraintupleKey)
}
func TestRegisterObjectiveWhitoutDataset(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)

	inpObjective := inputObjective{}
	inpObjective.createDefault()
	inpObjective.TestDataset = inputDataset{}
	resp := mockStub.MockInvoke("42", methodAndAssetToByte("registerObjective", inpObjective))
	assert.EqualValues(t, 200, resp.Status, "when adding objective without dataset it should work: ", resp.Message)
}
func TestRegisterObjectiveWithDataSampleKeyNotDataManagerKey(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)
	mockStub.MockInvoke("42", [][]byte{[]byte("registerNode")})

	// Add a dataManager and some dataSample successfuly
	inpDataManager := inputDataManager{}
	args := inpDataManager.createDefault()
	mockStub.MockInvoke("42", args)
	inpDataSample := inputDataSample{
		Hashes:          []string{testDataSampleHash1},
		DataManagerKeys: []string{dataManagerOpenerHash},
		TestOnly:        "true",
	}
	args = inpDataSample.createDefault()
	mockStub.MockInvoke("42", args)
	inpDataSample = inputDataSample{
		Hashes:          []string{testDataSampleHash2},
		DataManagerKeys: []string{dataManagerOpenerHash},
		TestOnly:        "true",
	}
	args = inpDataSample.createDefault()
	r := mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, r.Status)

	// Fail to insert the objective
	inpObjective := inputObjective{
		TestDataset: inputDataset{
			DataManagerKey: testDataSampleHash1,
			DataSampleKeys: []string{testDataSampleHash2}}}
	args = inpObjective.createDefault()
	resp := mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 400, resp.Status, "status should indicate an error since the dataManager key is a dataSample key")
}
func TestObjective(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)

	// Add objective with invalid field
	inpObjective := inputObjective{
		DescriptionHash: "aaa",
	}
	args := inpObjective.createDefault()
	resp := mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 400, resp.Status, "when adding objective with invalid hash, status %d and message %s", resp.Status, resp.Message)

	// Add objective with unexisting test dataSample
	inpObjective = inputObjective{}
	args = inpObjective.createDefault()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 400, resp.Status, "when adding objective with unexisting test dataSample, status %d and message %s", resp.Status, resp.Message)

	// Properly add objective
	resp, tt := registerItem(t, *mockStub, "objective")

	inpObjective = tt.(inputObjective)
	res := map[string]string{}
	err := json.Unmarshal(resp.Payload, &res)
	assert.NoError(t, err, "should unmarshal without problem")
	assert.Contains(t, res, "key")
	objectiveKey := res["key"]
	assert.EqualValuesf(
		t,
		inpObjective.DescriptionHash,
		objectiveKey,
		"when adding objective: unexpected returned objective key - %s / %s",
		objectiveKey,
		inpObjective.DescriptionHash)

	// Query objective from key and check the consistency of returned arguments
	args = [][]byte{[]byte("queryObjective"), keyToJSON(objectiveKey)}
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when querying a dataManager with status %d and message %s", resp.Status, resp.Message)
	objective := outputObjective{}
	err = bytesToStruct(resp.Payload, &objective)
	assert.NoError(t, err, "when unmarshalling queried objective")
	expectedObjective := outputObjective{
		Key:   objectiveKey,
		Owner: worker,
		TestDataset: &Dataset{
			DataManagerKey: dataManagerOpenerHash,
			DataSampleKeys: []string{testDataSampleHash1, testDataSampleHash2},
		},
		Name: inpObjective.Name,
		Description: HashDress{
			StorageAddress: inpObjective.DescriptionStorageAddress,
			Hash:           objectiveKey,
		},
		Permissions: outputPermissions{
			Process: Permission{Public: true, AuthorizedIDs: []string{}},
		},
		Metrics: &HashDressName{
			Hash:           inpObjective.MetricsHash,
			Name:           inpObjective.MetricsName,
			StorageAddress: inpObjective.MetricsStorageAddress,
		},
	}
	assert.Exactly(t, expectedObjective, objective)

	// Query all objectives and check consistency
	args = [][]byte{[]byte("queryObjectives")}
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when querying objectives - status %d and message %s", resp.Status, resp.Message)
	var objectives []outputObjective
	err = json.Unmarshal(resp.Payload, &objectives)
	assert.NoError(t, err, "while unmarshalling objectives")
	assert.Len(t, objectives, 1)
	assert.Exactly(t, expectedObjective, objectives[0], "return objective different from registered one")
}
