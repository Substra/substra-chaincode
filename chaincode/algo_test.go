package main

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func TestAlgo(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := shim.NewMockStub("substra", scc)

	// Add algo with invalid field
	inpAlgo := inputAlgo{
		DescriptionHash: "aaa",
	}
	args := inpAlgo.createSample()
	resp := mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 500 {
		t.Errorf("when adding algo with invalid hash, status %d and message %s", status, resp.Message)
	}

	// Add algo with unexisting challenge
	inpAlgo = inputAlgo{}
	args = inpAlgo.createSample()
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 500 {
		t.Errorf("when adding algo with unexisting challenge, status %d and message %s", status, resp.Message)
	}

	// Properly add algo
	err, resp, tt := registerItem(*mockStub, "algo")
	if err != nil {
		t.Errorf(err.Error())
	}
	inpAlgo = tt.(inputAlgo)
	algoKey := string(resp.Payload)
	if algoKey != inpAlgo.Hash {
		t.Errorf("when adding algo, key does not corresponds to its hash - key: %s and hash %s", algoKey, inpAlgo.Hash)
	}

	// Query algo from key and check the consistency of returned arguments
	args = [][]byte{[]byte("query"), []byte(algoKey)}
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("when querying an algo with status %d and message %s", status, resp.Message)
	}
	algo := make(map[string]interface{})
	if err := bytesToStruct(resp.Payload, &algo); err != nil {
		t.Errorf("when unmarshalling queried algo with error %s", err)
	}
	if algo["name"] != inpAlgo.Name {
		t.Errorf("ledger algo name does not correspond to what was input: %s - %s", algo["name"], inpAlgo.Name)
	}
	if algo["storageAddress"] != inpAlgo.StorageAddress {
		t.Errorf("ledger algo description storage address does not correspond to what was input: %s - %s", algo["storageAddress"], inpAlgo.StorageAddress)
	}
	if algo["challengeKey"] != inpAlgo.ChallengeKey {
		t.Errorf("ledger algo challenge key does not correspond to what was input: %s - %s", algo["challengeKey"], inpAlgo.ChallengeKey)
	}
	if algo["permissions"] != inpAlgo.Permissions {
		t.Errorf("ledger algo permissions does not correspond to what was input: %s - %s", algo["permissions"], inpAlgo.Permissions)
	}
	if algo["description"].(map[string]interface{})["hash"] != inpAlgo.DescriptionHash {
		t.Errorf("ledger algo metrics hash does not correspond to what was input")
	}
	if algo["description"].(map[string]interface{})["storageAddress"] != inpAlgo.DescriptionStorageAddress {
		t.Errorf("ledger algo metrics address does not correspond to what was input")
	}

	// Query all algo and check consistency
	args = [][]byte{[]byte("queryAlgos")}
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("when querying algos - status %d and message %s", status, resp.Message)
	}
	var sPayload []map[string]interface{}
	if err := json.Unmarshal(resp.Payload, &sPayload); err != nil {
		t.Errorf("when unmarshalling queried algos")
	}
	payload := sPayload[0]
	delete(payload, "key")
	if !reflect.DeepEqual(payload, algo) {
		t.Errorf("when querying algos, dataset does not correspond to the input algo")
	}
}
