package main

import (
	"encoding/json"
	"reflect"
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
	args = [][]byte{[]byte("query"), []byte(datasetKey)}
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("when querying the dataset, status %d and message %s", status, resp.Message)
	}
	dataset := make(map[string]interface{})
	if err := bytesToStruct(resp.Payload, &dataset); err != nil {
		t.Errorf("when unmarshalling queried dataset with error %s", err)
	}
	if dataset["name"] != inpDataset.Name {
		t.Errorf("ledger dataset name does not correspond to what was input: %s - %s", dataset["name"], inpDataset.Name)
	}
	if dataset["openerStorageAddress"] != inpDataset.OpenerStorageAddress {
		t.Errorf("ledger dataset opener storage address does not correspond to what was input")
	}
	if dataset["type"] != inpDataset.Type {
		t.Errorf("ledger dataset type does not correspond to what was input")
	}
	if dataset["description"].(map[string]interface{})["hash"] != inpDataset.DescriptionHash {
		t.Errorf("ledger dataset description hash does not correspond to what was input")
	}
	if dataset["description"].(map[string]interface{})["storageAddress"] != inpDataset.DescriptionStorageAddress {
		t.Errorf("ledger dataset description storage address does not correspond to what was input")
	}
	if dataset["challengeKey"] != "" {
		t.Errorf("ledger dataset challenge keys does not correspond to what was input")
	}
	if dataset["permissions"] != inpDataset.Permissions {
		t.Errorf("ledger dataset challenge keys does not correspond to what was input")
	}
	// Query all datasets and check fields match expectations
	args = [][]byte{[]byte("queryDatasets")}
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("when querying datasets - status %d and message %s", status, resp.Message)
	}
	var sPayload []map[string]interface{}
	if err := json.Unmarshal(resp.Payload, &sPayload); err != nil {
		t.Errorf("when unmarshalling queried datasets")
	}
	payload := sPayload[0]
	delete(payload, "key")
	if !reflect.DeepEqual(payload, dataset) {
		t.Errorf("when querying datasets, dataset does not correspond to the input dataset")
	}

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
