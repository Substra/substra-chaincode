package main

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func TestChallenge(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := shim.NewMockStub("substra", scc)

	// Add challenge with invalid field
	inpChallenge := inputChallenge{
		DescriptionHash: "aaa",
	}
	args := inpChallenge.createSample()
	resp := mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 500 {
		t.Errorf("when adding challenge with invalid hash, status %d and message %s", status, resp.Message)
	}

	// Add challenge with unexisting test data
	inpChallenge = inputChallenge{}
	args = inpChallenge.createSample()
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 500 {
		t.Errorf("when adding challenge with unexisting test data, status %d and message %s", status, resp.Message)
	}

	// Properly add challenge
	err, resp, tt := registerItem(*mockStub, "challenge")
	if err != nil {
		t.Errorf(err.Error())
	}
	inpChallenge = tt.(inputChallenge)
	challengeKey := string(resp.Payload)
	if challengeKey != inpChallenge.DescriptionHash {
		t.Errorf("when adding challenge: unexpected returned challenge key - %s / %s", challengeKey, inpChallenge.DescriptionHash)
	}

	// Query challenge from key and check the consistency of returned arguments
	args = [][]byte{[]byte("query"), []byte(challengeKey)}
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("when querying a dataset with status %d and message %s", status, resp.Message)
	}
	challenge := make(map[string]interface{})
	if err := bytesToStruct(resp.Payload, &challenge); err != nil {
		t.Errorf("when unmarshalling queried challenge with error %s", err)
	}
	if challenge["name"] != inpChallenge.Name {
		t.Errorf("ledger challenge name does not correspond to what was input: %s - %s", challenge["name"], inpChallenge.Name)
	}
	if challenge["descriptionStorageAddress"] != inpChallenge.DescriptionStorageAddress {
		t.Errorf("ledger challenge description storage address does not correspond to what was input: %s - %s", challenge["descriptionStorageAddress"], inpChallenge.DescriptionStorageAddress)
	}
	if challenge["permissions"] != inpChallenge.Permissions {
		t.Errorf("ledger challenge permissions does not correspond to what was input: %s - %s", challenge["permissions"], inpChallenge.Permissions)
	}
	if challenge["metrics"].(map[string]interface{})["hash"] != inpChallenge.MetricsHash {
		t.Errorf("ledger challenge metrics hash does not correspond to what was input")
	}
	if challenge["metrics"].(map[string]interface{})["name"] != inpChallenge.MetricsName {
		t.Errorf("ledger challenge metrics name does not correspond to what was input")
	}
	if challenge["metrics"].(map[string]interface{})["storageAddress"] != inpChallenge.MetricsStorageAddress {
		t.Errorf("ledger challenge metrics address does not correspond to what was input")
	}
	testData := &DatasetData{
		DatasetKey: datasetOpenerHash,
		DataKeys:   []string{testDataHash1, testDataHash2},
	}
	if reflect.DeepEqual(challenge["testData"], testData) {
		t.Errorf("ledger challenge test data does not correspond to what was input")
	}

	// Query all challenges and check consistency
	args = [][]byte{[]byte("queryChallenges")}
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("when querying challenges - status %d and message %s", status, resp.Message)
	}
	var sPayload []map[string]interface{}
	if err := json.Unmarshal(resp.Payload, &sPayload); err != nil {
		t.Errorf("when unmarshalling queried challenges")
	}
	payload := sPayload[0]
	delete(payload, "key")
	if !reflect.DeepEqual(payload, challenge) {
		t.Errorf("when querying challenges, dataset does not correspond to the input challenge")
	}
}
