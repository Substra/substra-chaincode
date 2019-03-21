package main

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func TestDataset(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := shim.NewMockStub("substra", scc)

	// Add dataset with invalid field
	inpDataset := inputDataset{
		OpenerHash: "aaa",
	}
	args := inpDataset.createSample()
	resp := mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 500 {
		t.Errorf("when adding dataset with invalid opener hash, status %d and message %s", status, resp.Message)
	}
	// Properly add dataset
	err, resp, tt := registerItem(*mockStub, "dataset")
	if err != nil {
		t.Errorf(err.Error())
	}
	inpDataset = tt.(inputDataset)
	datasetKey := string(resp.Payload)
	// check returned dataset key corresponds to opener hash
	if datasetKey != datasetOpenerHash {
		t.Errorf("when adding dataset: dataset key does not correspond to dataset opener hash: %s - %s", datasetKey, datasetOpenerHash)
	}
	// Add dataset which already exist
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 500 {
		t.Errorf("when adding dataset which already exists, status %d and message %s", status, resp.Message)
	}
	// Query dataset and check fields match expectations
	args = [][]byte{[]byte("queryDataset"), []byte(datasetKey)}
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("when querying the dataset, status %d and message %s", status, resp.Message)
	}
	dataset := outputDataset{}
	err = bytesToStruct(resp.Payload, &dataset)
	assert.NoError(t, err, "when unmarshalling queried dataset")
	expectedDataset := outputDataset{
		ChallengeKey: inpDataset.ChallengeKey,
		Key:          datasetKey,
		Owner:        "bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0",
		Name:         inpDataset.Name,
		Description: &HashDress{
			StorageAddress: inpDataset.DescriptionStorageAddress,
			Hash:           inpDataset.DescriptionHash,
		},
		Permissions: inpDataset.Permissions,
		Opener: datasetOpener{
			Hash:           datasetKey,
			StorageAddress: inpDataset.OpenerStorageAddress,
		},
		Type: inpDataset.Type,
	}
	assert.Exactly(t, expectedDataset, dataset)

	// Query all datasets and check fields match expectations
	args = [][]byte{[]byte("queryDatasets")}
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("when querying datasets - status %d and message %s", status, resp.Message)
	}
	var datasets []outputDataset
	err = json.Unmarshal(resp.Payload, &datasets)
	assert.NoError(t, err, "while unmarshalling datasets")
	assert.Len(t, datasets, 1)
	assert.Exactly(t, expectedDataset, datasets[0], "return challenge different from registered one")

	args = [][]byte{[]byte("queryDatasetData"), []byte(inpDataset.OpenerHash)}
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("when querying dataset data, status %d and message %s", status, resp.Message)
	}
	if !strings.Contains(string(resp.Payload), "\"trainDataKeys\":[]") {
		t.Errorf("when querying dataset data, trainDataKeys should be []")
	}
	if !strings.Contains(string(resp.Payload), "\"testDataKeys\":[]") {
		t.Errorf("when querying dataset data, testDataKeys should be []")
	}
}

func TestGetTestDataKeys(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := shim.NewMockStub("substra", scc)

	// Input Dataset
	inpDataset := inputDataset{}
	args := inpDataset.createSample()
	mockStub.MockInvoke("42", args)

	// Add both train and test data
	inpData := inputData{Hashes: testDataHash1}
	args = inpData.createSample()
	mockStub.MockInvoke("42", args)
	inpData.Hashes = testDataHash2
	inpData.TestOnly = "true"
	args = inpData.createSample()
	mockStub.MockInvoke("42", args)

	// Querry the Dataset
	args = [][]byte{[]byte("queryDatasetData"), []byte(inpDataset.OpenerHash)}
	resp := mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status, "querrying the dataset should return an ok status")
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

	// Add data with unexiting dataset
	inpData = inputData{}
	args = inpData.createSample()
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 500 {
		t.Errorf("when adding data with unexisting dataset, status %d and message %s", status, resp.Message)
	}
	// TODO Would be nice to check failure when adding data to a dataset owned by a different people

	// Properly add data
	// 1. add associated dataset
	inpDataset := inputDataset{}
	args = inpDataset.createSample()
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
	args = [][]byte{[]byte("queryDatasetData"), []byte(inpDataset.OpenerHash)}
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("when querying dataset data with status %d and message %s", status, resp.Message)
	}
	payload := make(map[string]interface{})
	json.Unmarshal(resp.Payload, &payload)
	if _, ok := payload["key"]; !ok {
		t.Errorf("when querying dataset data, payload should contain the dataset key")
	}
	v, ok := payload["trainDataKeys"]
	if !ok {
		t.Errorf("when querying dataset data, payload should contain the train data keys")
	}
	if reflect.DeepEqual(v, strings.Split(strings.Replace(inpData.Hashes, " ", "", -1), ",")) {
		t.Errorf("when querying dataset data, unexpected train keys")
	}
	**/
}
