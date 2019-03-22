package main

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func TestDataManager(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := shim.NewMockStub("substra", scc)

	// Add dataManager with invalid field
	inpDataManager := inputDataManager{
		OpenerHash: "aaa",
	}
	args := inpDataManager.createSample()
	resp := mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 500 {
		t.Errorf("when adding dataManager with invalid opener hash, status %d and message %s", status, resp.Message)
	}
	// Properly add dataManager
	err, resp, tt := registerItem(*mockStub, "dataManager")
	if err != nil {
		t.Errorf(err.Error())
	}
	inpDataManager = tt.(inputDataManager)
	dataManagerKey := string(resp.Payload)
	// check returned dataManager key corresponds to opener hash
	if dataManagerKey != dataManagerOpenerHash {
		t.Errorf("when adding dataManager: dataManager key does not correspond to dataManager opener hash: %s - %s", dataManagerKey, dataManagerOpenerHash)
	}
	// Add dataManager which already exist
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 500 {
		t.Errorf("when adding dataManager which already exists, status %d and message %s", status, resp.Message)
	}
	// Query dataManager and check fields match expectations
	args = [][]byte{[]byte("queryDataManager"), []byte(dataManagerKey)}
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("when querying the dataManager, status %d and message %s", status, resp.Message)
	}
	dataManager := outputDataManager{}
	err = bytesToStruct(resp.Payload, &dataManager)
	assert.NoError(t, err, "when unmarshalling queried dataManager")
	expectedDataManager := outputDataManager{
		ObjectiveKey: inpDataManager.ObjectiveKey,
		Key:          dataManagerKey,
		Owner:        "bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0",
		Name:         inpDataManager.Name,
		Description: &HashDress{
			StorageAddress: inpDataManager.DescriptionStorageAddress,
			Hash:           inpDataManager.DescriptionHash,
		},
		Permissions: inpDataManager.Permissions,
		Opener: HashDress{
			Hash:           dataManagerKey,
			StorageAddress: inpDataManager.OpenerStorageAddress,
		},
		Type: inpDataManager.Type,
	}
	assert.Exactly(t, expectedDataManager, dataManager)

	// Query all dataManagers and check fields match expectations
	args = [][]byte{[]byte("queryDataManagers")}
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("when querying dataManagers - status %d and message %s", status, resp.Message)
	}
	var dataManagers []outputDataManager
	err = json.Unmarshal(resp.Payload, &dataManagers)
	assert.NoError(t, err, "while unmarshalling dataManagers")
	assert.Len(t, dataManagers, 1)
	assert.Exactly(t, expectedDataManager, dataManagers[0], "return objective different from registered one")

	args = [][]byte{[]byte("queryDataset"), []byte(inpDataManager.OpenerHash)}
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("when querying dataManager data, status %d and message %s", status, resp.Message)
	}
	if !strings.Contains(string(resp.Payload), "\"trainDataKeys\":[]") {
		t.Errorf("when querying dataManager data, trainDataKeys should be []")
	}
	if !strings.Contains(string(resp.Payload), "\"testDataKeys\":[]") {
		t.Errorf("when querying dataManager data, testDataKeys should be []")
	}
}

func TestGetTestDataKeys(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := shim.NewMockStub("substra", scc)

	// Input DataManager
	inpDataManager := inputDataManager{}
	args := inpDataManager.createSample()
	mockStub.MockInvoke("42", args)

	// Add both train and test data
	inpData := inputData{Hashes: testDataHash1}
	args = inpData.createSample()
	mockStub.MockInvoke("42", args)
	inpData.Hashes = testDataHash2
	inpData.TestOnly = "true"
	args = inpData.createSample()
	mockStub.MockInvoke("42", args)

	// Querry the DataManager
	args = [][]byte{[]byte("queryDataset"), []byte(inpDataManager.OpenerHash)}
	resp := mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status, "querrying the dataManager should return an ok status")
	payload := map[string]interface{}{}
	err := json.Unmarshal(resp.Payload, &payload)
	assert.NoError(t, err)

	v, ok := payload["testDataKeys"]
	assert.True(t, ok, "payload should contains the test data keys")
	assert.Contains(t, v, testDataHash2, "testDataKeys should contain the test dataHash")
	assert.NotContains(t, v, testDataHash1, "testDataKeys should not contains the train dataHash")
}
func TestData(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := shim.NewMockStub("substra", scc)

	// Add data with invalid field
	inpData := inputData{
		Hashes: "aaa",
	}
	args := inpData.createSample()
	resp := mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 500 {
		t.Errorf("when adding data with invalid hash, status %d and message %s", status, resp.Message)
	}

	// Add data with unexiting dataManager
	inpData = inputData{}
	args = inpData.createSample()
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 500 {
		t.Errorf("when adding data with unexisting dataManager, status %d and message %s", status, resp.Message)
	}
	// TODO Would be nice to check failure when adding data to a dataManager owned by a different people

	// Properly add data
	// 1. add associated dataManager
	inpDataManager := inputDataManager{}
	args = inpDataManager.createSample()
	mockStub.MockInvoke("42", args)
	// 2. add data
	inpData = inputData{}
	args = inpData.createSample()
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("when adding data, status %d and message %s", status, resp.Message)
	}
	// check payload correspond to input data keys
	dataKeys := string(resp.Payload)
	if expectedResp := "{\"keys\": [\"" + strings.Replace(inpData.Hashes, ", ", "\", \"", -1) + "\"]}"; dataKeys != expectedResp {
		t.Errorf("when adding data: data keys does not correspond to data hashes: %s - %s", dataKeys, expectedResp)
	}

	// Add data which already exist
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 500 {
		t.Errorf("when adding data which already exist, status %d and message %s", status, resp.Message)
	}

	/**
	// Query data and check it corresponds to what was input
	args = [][]byte{[]byte("queryDataset"), []byte(inpDataManager.OpenerHash)}
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("when querying dataManager data with status %d and message %s", status, resp.Message)
	}
	payload := make(map[string]interface{})
	json.Unmarshal(resp.Payload, &payload)
	if _, ok := payload["key"]; !ok {
		t.Errorf("when querying dataManager data, payload should contain the dataManager key")
	}
	v, ok := payload["trainDataKeys"]
	if !ok {
		t.Errorf("when querying dataManager data, payload should contain the train data keys")
	}
	if reflect.DeepEqual(v, strings.Split(strings.Replace(inpData.Hashes, " ", "", -1), ",")) {
		t.Errorf("when querying dataManager data, unexpected train keys")
	}
	**/
}
