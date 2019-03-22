package main

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func TestRegisterObjectiveWithDataKeyNotDataManagerKey(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := shim.NewMockStub("substra", scc)

	// Add a dataManager and some data successfuly
	inpDataManager := inputDataManager{}
	args := inpDataManager.createSample()
	mockStub.MockInvoke("42", args)
	inpData := inputData{
		Hashes:      testDataHash1,
		DataManagerKeys: dataManagerOpenerHash,
		TestOnly:    "true",
	}
	args = inpData.createSample()
	mockStub.MockInvoke("42", args)
	inpData = inputData{
		Hashes:      testDataHash2,
		DataManagerKeys: dataManagerOpenerHash,
		TestOnly:    "true",
	}
	args = inpData.createSample()
	r := mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, r.Status)

	// Fail to insert the objective
	inpObjective := inputObjective{TestDataset: testDataHash1 + ":" + testDataHash2}
	args = inpObjective.createSample()
	resp := mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 500, resp.Status, "status should indicate an error since the dataManager key is a data key")
}
func TestObjective(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := shim.NewMockStub("substra", scc)

	// Add objective with invalid field
	inpObjective := inputObjective{
		DescriptionHash: "aaa",
	}
	args := inpObjective.createSample()
	resp := mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 500 {
		t.Errorf("when adding objective with invalid hash, status %d and message %s", status, resp.Message)
	}

	// Add objective with unexisting test data
	inpObjective = inputObjective{}
	args = inpObjective.createSample()
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 500 {
		t.Errorf("when adding objective with unexisting test data, status %d and message %s", status, resp.Message)
	}

	// Properly add objective
	err, resp, tt := registerItem(*mockStub, "objective")
	if err != nil {
		t.Errorf(err.Error())
	}
	inpObjective = tt.(inputObjective)
	objectiveKey := string(resp.Payload)
	if objectiveKey != inpObjective.DescriptionHash {
		t.Errorf("when adding objective: unexpected returned objective key - %s / %s", objectiveKey, inpObjective.DescriptionHash)
	}

	// Query objective from key and check the consistency of returned arguments
	args = [][]byte{[]byte("queryObjective"), []byte(objectiveKey)}
	resp = mockStub.MockInvoke("42", args)
	if status := resp.Status; status != 200 {
		t.Errorf("when querying a dataManager with status %d and message %s", status, resp.Message)
	}
	objective := outputObjective{}
	err = bytesToStruct(resp.Payload, &objective)
	assert.NoError(t, err, "when unmarshalling queried objective")
	expectedObjective := outputObjective{
		Key:   objectiveKey,
		Owner: "bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0",
		TestDataset: &Dataset{
			DataManagerKey: dataManagerOpenerHash,
			DataKeys:   []string{testDataHash1, testDataHash2},
		},
		Name: inpObjective.Name,
		Description: HashDress{
			StorageAddress: inpObjective.DescriptionStorageAddress,
			Hash:           objectiveKey,
		},
		Permissions: inpObjective.Permissions,
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
	if status := resp.Status; status != 200 {
		t.Errorf("when querying objectives - status %d and message %s", status, resp.Message)
	}
	var objectives []outputObjective
	err = json.Unmarshal(resp.Payload, &objectives)
	assert.NoError(t, err, "while unmarshalling objectives")
	assert.Len(t, objectives, 1)
	assert.Exactly(t, expectedObjective, objectives[0], "return objective different from registered one")
}
