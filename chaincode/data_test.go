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
)

func TestJsonInputsDataManager(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)

	inpDataManager := inputDataManager{}
	inpDataManager.createDefault()
	payload, err := json.Marshal(inpDataManager)
	assert.NoError(t, err)
	args := [][]byte{[]byte("registerDataManager"), payload}
	resp := mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status)
}
func TestDataManager(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)

	// Add dataManager with invalid field
	inpDataManager := inputDataManager{
		OpenerChecksum: "aaa",
	}
	args := inpDataManager.createDefault()
	resp := mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 400, resp.Status, "when adding dataManager with invalid opener checksum, status %d and message %s", resp.Status, resp.Message)
	// Properly add dataManager
	resp, tt := registerItem(t, *mockStub, "dataManager")

	inpDataManager = tt.(inputDataManager)
	res := outputKey{}
	err := json.Unmarshal(resp.Payload, &res)
	assert.NoError(t, err, "should unmarshal without problem")
	dataManagerKey := res.Key

	// Add dataManager which already exist
	inpDataManager = inputDataManager{}
	args = inpDataManager.createDefault()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 409, resp.Status, "when adding dataManager which already exists, status %d and message %s", resp.Status, resp.Message)
	// Query dataManager and check fields match expectations
	args = [][]byte{[]byte("queryDataManager"), keyToJSON(dataManagerKey)}
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when querying the dataManager, status %d and message %s", resp.Status, resp.Message)
	dataManager := outputDataManager{}
	err = json.Unmarshal(resp.Payload, &dataManager)
	assert.NoError(t, err, "when unmarshalling queried dataManager")
	expectedDataManager := outputDataManager{
		ObjectiveKey: inpDataManager.ObjectiveKey,
		Key:          dataManagerKey,
		Owner:        worker,
		Name:         inpDataManager.Name,
		Description: &ChecksumAddress{
			StorageAddress: inpDataManager.DescriptionStorageAddress,
			Checksum:       inpDataManager.DescriptionChecksum,
		},
		Permissions: outputPermissions{
			Process: Permission{Public: true, AuthorizedIDs: []string{}},
		},
		Opener: &ChecksumAddress{
			Checksum:       inpDataManager.OpenerChecksum,
			StorageAddress: inpDataManager.OpenerStorageAddress,
		},
		Type:     inpDataManager.Type,
		Metadata: map[string]string{},
	}
	assert.Exactly(t, expectedDataManager, dataManager)

	// Query all dataManagers and check fields match expectations
	args = [][]byte{[]byte("queryDataManagers")}
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when querying dataManagers - status %d and message %s", resp.Status, resp.Message)
	var dataManagers []outputDataManager
	err = json.Unmarshal(resp.Payload, &dataManagers)
	assert.NoError(t, err, "while unmarshalling dataManagers")
	assert.Len(t, dataManagers, 1)
	assert.Exactly(t, expectedDataManager, dataManagers[0], "return objective different from registered one")

	args = [][]byte{[]byte("queryDataset"), keyToJSON(inpDataManager.Key)}
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when querying Dataset, status %d and message %s", resp.Status, resp.Message)
	out := outputDataset{}
	err = json.Unmarshal(resp.Payload, &out)
	assert.NoError(t, err, "while unmarshalling dataset")
	assert.Empty(t, out.TrainDataSampleKeys, "when querying Dataset, trainDataSampleKeys should be empty")
	assert.Empty(t, out.TestDataSampleKeys, "when querying Dataset, testDataSampleKeys should be empty")
}

func TestGetTestDatasetKeys(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)

	// Input DataManager
	inpDataManager := inputDataManager{}
	args := inpDataManager.createDefault()
	mockStub.MockInvoke("42", args)

	// Add both train and test dataSample
	inpDataSample := inputDataSample{Keys: []string{testDataSampleKey1}}
	args = inpDataSample.createDefault()
	mockStub.MockInvoke("42", args)
	inpDataSample.Keys = []string{testDataSampleKey2}
	inpDataSample.TestOnly = "true"
	args = inpDataSample.createDefault()
	mockStub.MockInvoke("42", args)

	// Query the DataManager
	args = [][]byte{[]byte("queryDataset"), keyToJSON(inpDataManager.Key)}
	resp := mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status, "querying the dataManager should return an ok status")
	payload := map[string]interface{}{}
	err := json.Unmarshal(resp.Payload, &payload)
	assert.NoError(t, err)

	v, ok := payload["test_data_sample_keys"]
	assert.True(t, ok, "payload should contains the test dataSample keys")
	assert.Contains(t, v, testDataSampleKey2, "testDataSampleKeys should contain the test dataSampleKey")
	assert.NotContains(t, v, testDataSampleKey1, "testDataSampleKeys should not contains the train dataSampleKey")
}
func TestDataset(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)

	// Add dataSample with invalid field
	inpDataSample := inputDataSample{
		Keys: []string{"aaa"},
	}
	args := inpDataSample.createDefault()
	resp := mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 400, resp.Status, "when adding dataSample with invalid key, status %d and message %s", resp.Status, resp.Message)

	// Add dataSample with unexiting dataManager
	inpDataSample = inputDataSample{}
	args = inpDataSample.createDefault()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 400, resp.Status, "when adding dataSample with unexisting dataManager, status %d and message %s", resp.Status, resp.Message)
	// TODO Would be nice to check failure when adding dataSample to a dataManager owned by a different people

	// Properly add dataSample
	// 1. add associated dataManager
	inpDataManager := inputDataManager{}
	args = inpDataManager.createDefault()
	mockStub.MockInvoke("42", args)
	// 2. add dataSample
	inpDataSample = inputDataSample{}
	args = inpDataSample.createDefault()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when adding dataSample, status %d and message %s", resp.Status, resp.Message)
	// check payload correspond to input dataSample keys
	res := map[string]interface{}{}
	err := json.Unmarshal(resp.Payload, &res)
	assert.NoError(t, err, "should unmarshal without problem")
	assert.Contains(t, res, "keys")
	dataSampleKeys := res["keys"]
	expectedResp := inpDataSample.Keys
	assert.ElementsMatch(t, expectedResp, dataSampleKeys, "when adding dataSample: dataSample keys does not correspond to dataSample keys")

	// Add dataSample which already exist
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 409, resp.Status, "when adding dataSample which already exist, status %d and message %s", resp.Status, resp.Message)

	// Query dataSample and check it corresponds to what was input
	args = [][]byte{[]byte("queryDataset"), keyToJSON(inpDataManager.Key)}
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when querying dataManager dataSample with status %d and message %s", resp.Status, resp.Message)
	out := outputDataset{}
	err = json.Unmarshal(resp.Payload, &out)
	assert.NoError(t, err, "while unmarshalling dataset")
	assert.ElementsMatch(t, out.TrainDataSampleKeys, inpDataSample.Keys, "when querying dataManager dataSample, unexpected train keys")

}
