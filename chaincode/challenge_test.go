package main

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func TestRegisterChallengeWithDataKeyNotDatasetKey(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := shim.NewMockStub("substra", scc)

	// Add a dataset and some data successfuly
	inpDataset := inputDataset{}
	args := inpDataset.createSample()
	mockStub.MockInvoke("42", args)
	inpData := inputData{
		Hashes:      testDataHash1,
		DatasetKeys: datasetOpenerHash,
		TestOnly:    "true",
	}
	args = inpData.createSample()
	mockStub.MockInvoke("42", args)
	inpData = inputData{
		Hashes:      testDataHash2,
		DatasetKeys: datasetOpenerHash,
		TestOnly:    "true",
	}
	args = inpData.createSample()
	r := mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, r.Status)

	// Fail to insert the challenge
	inpChallenge := inputChallenge{TestData: testDataHash1 + ":" + testDataHash2}
	args = inpChallenge.createSample()
	resp := mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 500, resp.Status, "status should indicate an error since the dataset key is a data key")
}
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
	args = [][]byte{[]byte("queryChallenge"), []byte(challengeKey)}
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("when querying a dataset with status %d and message %s", status, resp.Message)
	}
	challenge := outputChallenge{}
	err = bytesToStruct(resp.Payload, &challenge)
	assert.NoError(t, err, "when unmarshalling queried challenge")
	expectedChallenge := outputChallenge{
		Key:   challengeKey,
		Owner: "bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0",
		TestData: &DatasetData{
			DatasetKey: datasetOpenerHash,
			DataKeys:   []string{testDataHash1, testDataHash2},
		},
		Name: inpChallenge.Name,
		Description: challengeDescription{
			StorageAddress: inpChallenge.DescriptionStorageAddress,
			Hash:           challengeKey,
		},
		Permissions: inpChallenge.Permissions,
		Metrics: &HashDressName{
			Hash:           inpChallenge.MetricsHash,
			Name:           inpChallenge.MetricsName,
			StorageAddress: inpChallenge.MetricsStorageAddress,
		},
	}
	assert.Exactly(t, expectedChallenge, challenge)

	// Query all challenges and check consistency
	args = [][]byte{[]byte("queryChallenges")}
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("when querying challenges - status %d and message %s", status, resp.Message)
	}
	var challenges []outputChallenge
	err = json.Unmarshal(resp.Payload, &challenges)
	assert.NoError(t, err, "while unmarshalling challenges")
	assert.Len(t, challenges, 1)
	assert.Exactly(t, expectedChallenge, challenges[0], "return challenge different from registered one")
}
